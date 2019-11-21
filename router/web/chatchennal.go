package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
)

type ChatChannelForm struct {
	Name               string `form:"name"`
	PhoneNumber        string `form:"phone_number"`
	LineID             string `form:"line_id"`
	ChannelID          string `form:"channel_id"`
	WebSite            string `form:"website"`
	ChannelSecret      string `form:"channel_secret"`
	WelcomeMessage     string `form:"welcome_message"`
	ChannelAccessToken string `form:"channel_access_token"`
	Type               string `form:"type"`
	Address            string `form:"address"`
	Settings           string `form:"settings"`
	Image              string `form:"Image"`
}

func ChatChannelListHandler(c *Context) error {
	chatChannels := []*model.ChatChannel{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	filterChatChannel := db.Preload("Account").Where("cha_account_id = ?", a.User.GetAccountID()).Find(&chatChannels).Count(&total)
	filterChatChannel.Limit(limit).Offset(page).Find(&chatChannels)
	pagination := MakePagination(total, page, limit)
	return c.Render(http.StatusOK, "chat-channel-list", echo.Map{
		"title":      "chat_channel",
		"list":       chatChannels,
		"pagination": pagination,
	})
}

func ChatChannelGetChannelAccessTokenHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	db.Preload("Account").Where("cha_account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id)
	bot, err := linebot.New(chatChannel.ChaChannelID, chatChannel.ChaChannelSecret)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.IssueAccessToken(chatChannel.ChaChannelID, chatChannel.ChaChannelSecret).Do()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	chatChannel.ChaChannelAccessToken = res.AccessToken
	if err := db.Save(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	db.Model(&chatChannel).Association("Settings").Append(
		&model.Setting{Name: "statusAccessToken", Value: "success"},
		model.Setting{Name: "dateStatusToken", Value: fmt.Sprintf("%s", time.Now())})

	return c.JSON(http.StatusOK, res)
}

func ChatChannelAddRegisterLIFF(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	if err := db.Preload("Account").Where("cha_account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	bot, err := linebot.New(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if bot == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	URLRegister := fmt.Sprintf("https://%s/register/%s", Conf.Server.DomainLineChannel, chatChannel.ChaLineID)
	view := linebot.View{Type: "full", URL: URLRegister}
	var status string = "success"
	var LIFFID string = ""
	res, err := bot.AddLIFF(view).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	} else {
		LIFFID = res.LIFFID
	}
	fmt.Println(err != nil)
	if err := model.DB().Model(&chatChannel).Association("Settings").Append(
		&model.Setting{Name: "LIFFregister", Value: LIFFID},
		&model.Setting{Name: "statusLIFFregister", Value: status},
		&model.Setting{Name: "statusAccessToken", Value: status},
		&model.Setting{Name: "dateStatusToken", Value: status},
	).Error; err != nil {
		return c.JSON(http.StatusBadRequest, chatChannel)
	}
	return c.JSON(http.StatusOK, res)
}

type SettingResponse struct {
	LIFFregister       string `json:"LIFFregister"`
	StatusLIFFregister string `json:"statusLIFFregister"`
	StatusAccessToken  string `json:"statusAccessToken"`
	DateStatusToken    string `json:"dateStatusToken"`
}

// ChatChannelDetailHandler
func ChatChannelDetailHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	model.DB().Preload("Services").Preload("Customers").Preload("ActionLogs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Customers")
	}).Preload("EventLogs").Preload("Account").Where("cha_account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id)
	customerSum := len(chatChannel.Customers)
	serviceSum := len(chatChannel.Services)
	settings := chatChannel.GetSetting([]string{"LIFFregister", "statusLIFFregister", "statusAccessToken", "dateStatusToken"})

	return c.Render(http.StatusOK, "chat-channel-detail", echo.Map{
		"title":       "chat_channel",
		"detail":      chatChannel,
		"customerSum": customerSum,
		"serviceSum":  serviceSum,
		"settings":    settings,
	})
}

// ChatChannelCreateHandler
func ChatChannelCreateViewHandler(c *Context) error {
	typeChatChannels := []string{"Facebook", "Line"}
	csrfValue := c.Get("_csrf")
	return c.Render(http.StatusOK, "chat-channel-form", echo.Map{
		"title":            "chat_channel",
		"typeChatChannels": typeChatChannels,
		"mode":             "Create",
		"_csrf":            csrfValue,
	})
}

// ChatChannelCreatePostHandler
func ChatChannelCreatePostHandler(c *Context) error {
	chatChannel := ChatChannelForm{}
	if err := c.Bind(&chatChannel); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	settingsModel := &[]*model.Setting{}
	err := json.Unmarshal([]byte(chatChannel.Settings), settingsModel)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	a := auth.Default(c)
	chatChannelModel := model.ChatChannel{
		ChaChannelID:          chatChannel.ChannelID,
		ChaName:               chatChannel.Name,
		ChaLineID:             chatChannel.LineID,
		ChaChannelSecret:      chatChannel.ChannelSecret,
		ChaChannelAccessToken: chatChannel.ChannelAccessToken,
		ChaType:               chatChannel.Type,
		ChaPhoneNumber:        chatChannel.PhoneNumber,
		ChaAccountID:          a.User.GetAccountID(),
		ChaImage:              chatChannel.Image,
		ChaWebSite:            chatChannel.WebSite,
		ChaWelcomeMessage:     chatChannel.WelcomeMessage,
		ChaAddress:            chatChannel.Address,
		Settings:              *settingsModel,
	}
	if err := chatChannelModel.SaveChatChannel(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if chatChannel.Type == "Line" {
		bot, err := linebot.New(chatChannel.ChannelID, chatChannel.ChannelSecret)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		URLRegister := fmt.Sprintf("https://%s/register/%s", Conf.Server.DomainLineChannel, chatChannel.LineID)
		view := linebot.View{Type: "full", URL: URLRegister}
		var status string = "success"
		var LIFFID string = ""
		res, err := bot.AddLIFF(view).Do()
		if err != nil {
			status = "error"
		} else {
			LIFFID = res.LIFFID
		}
		if err := model.DB().Model(&chatChannelModel).Association("Settings").Append(
			&model.Setting{Name: "LIFFregister", Value: LIFFID},
			&model.Setting{Name: "statusLIFFregister", Value: status},
			&model.Setting{Name: "statusAccessToken", Value: status},
			&model.Setting{Name: "dateStatusToken", Value: status},
		).Error; err != nil {
			return c.JSON(http.StatusBadRequest, chatChannelModel)
		}
	}
	redirect := fmt.Sprintf("/admin/chat_channel/%d", chatChannelModel.ID)
	return c.JSON(http.StatusCreated, redirect)
}

// ChatChannelEditHandler
func ChatChannelEditHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	if err := db.Preload("Services").Preload("Customers").Preload("ActionLogs").Preload("EventLogs").Preload("Account").Where("cha_account_id = ?",
		a.User.GetAccountID()).Find(&chatChannel, id).Error; err != nil {
		return c.Render(http.StatusOK, "404-page", echo.Map{})
	}
	customerSum := len(chatChannel.Customers)
	serviceSum := len(chatChannel.Services)
	typeChatChannels := []string{"Facebook", "Line"}
	return c.Render(http.StatusOK, "chat-channel-form", echo.Map{
		"title":            "chat_channel",
		"detail":           chatChannel,
		"customerSum":      customerSum,
		"serviceSum":       serviceSum,
		"typeChatChannels": typeChatChannels,
		"mode":             "Edit",
	})
}

func ChatChannelDeleteHandler(c *Context) error {
	id := c.Param("id")

	idInt, _ := strconv.Atoi(id)
	chatChannel := model.DeleteChatChannel(idInt)
	return c.JSON(http.StatusOK, chatChannel)
}
