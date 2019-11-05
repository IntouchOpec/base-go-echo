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

// ChatChannelListHandler
func ChatChannelListHandler(c *Context) error {
	chatChannels := []*model.ChatChannel{}
	a := auth.Default(c)
	model.DB().Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannels)
	return c.Render(http.StatusOK, "chat-channel-list", echo.Map{
		"title": "chat_channel",
		"list":  chatChannels,
	})
}

func ChatChannelGetChannelAccessTokenHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	db.Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id)
	bot, err := linebot.New(chatChannel.ChannelID, chatChannel.ChannelSecret)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.IssueAccessToken(chatChannel.ChannelID, chatChannel.ChannelSecret).Do()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	chatChannel.ChannelAccessToken = res.AccessToken
	if err := db.Save(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	db.Model(&chatChannel).Association("Settings").Append(&model.Setting{Name: "statusAccessToken", Value: "success"}, model.Setting{Name: "dateStatusToken", Value: fmt.Sprintf("%s", time.Now())})
	return c.JSON(http.StatusOK, "res")
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
	}).Preload("EventLogs").Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id)
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
		ChannelID:          chatChannel.ChannelID,
		Name:               chatChannel.Name,
		LineID:             chatChannel.LineID,
		ChannelSecret:      chatChannel.ChannelSecret,
		ChannelAccessToken: chatChannel.ChannelAccessToken,
		Type:               chatChannel.Type,
		PhoneNumber:        chatChannel.PhoneNumber,
		AccountID:          a.User.GetAccountID(),
		Image:              chatChannel.Image,
		WebSite:            chatChannel.WebSite,
		WelcomeMessage:     chatChannel.WelcomeMessage,
		Address:            chatChannel.Address,
		Settings:           *settingsModel,
	}
	chatChannelModel.SaveChatChannel()
	if chatChannel.Type == "Line" {
		bot, err := linebot.New(chatChannel.ChannelID, chatChannel.ChannelSecret)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		URLRegister := fmt.Sprintf("https://%s/register/%s", Conf.Server.DomainLineChannel, chatChannel.LineID)
		view := linebot.View{Type: "full", URL: URLRegister}
		var status string = "success"
		res, err := bot.AddLIFF(view).Do()
		if err != nil {
			status = "error"
		}
		model.DB().Model(&chatChannelModel).Association("Settings").Append(
			&model.Setting{Name: "LIFFregister", Value: res.LIFFID},
			&model.Setting{Name: "statusLIFFregister", Value: status},
			&model.Setting{Name: "statusAccessToken", Value: status},
			&model.Setting{Name: "dateStatusToken", Value: status},
		)
	}
	return c.JSON(http.StatusCreated, chatChannelModel)
}

// ChatChannelEditHandler
func ChatChannelEditHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	db.Preload("Services").Preload("Customers").Preload("ActionLogs").Preload("EventLogs").Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id)
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
