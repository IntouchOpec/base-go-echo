package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/lib"

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
	filterChatChannel := db.Model(&chatChannels).Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterChatChannel.Limit(pagination.Record).Offset(pagination.Offset).Find(&chatChannels)

	return c.Render(http.StatusOK, "chat-channel-list", echo.Map{
		"title":      "chat_channel",
		"list":       chatChannels,
		"pagination": pagination,
	})
}

func ChatChannelGetChannelAccessTokenHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	db.Preload("Settings", "name in (?)", []string{"dateStatusToken", "statusAccessToken"}).Where("account_id = ?", accID).Find(&chatChannel, id)
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
	dateStatusToken := model.Setting{Detail: "", Name: "dateStatusToken", Value: "success"}
	statusAccessToken := model.Setting{Detail: "", Name: "statusAccessToken", Value: time.Now().Format("Mon Jan 2 2006")}
	if len(chatChannel.Settings) == 0 {
		db.Save(&statusAccessToken)
		db.Save(&dateStatusToken)
		db.Model(&chatChannel).Association("Settings").Append(&statusAccessToken)
		db.Model(&chatChannel).Association("Settings").Append(&dateStatusToken)
	} else {
		for _, setting := range chatChannel.Settings {
			if setting.Name == dateStatusToken.Name {
				setting.Value = dateStatusToken.Value
			}
			if statusAccessToken.Name == setting.Name {
				setting.Value = statusAccessToken.Value
			}
			db.Save(&setting)
		}
	}

	return c.JSON(http.StatusOK, chatChannel)
}

func ChatChannelAddRegisterLIFF(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	setting := model.Setting{}
	if err := c.Bind(&setting); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	db.Preload("Settings", "name in (?)", []string{"LIFFregister", "statusLIFFregister"}).Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id)

	bot, err := linebot.New(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	URLRegister := fmt.Sprintf("https://web.%s/register/%s", Conf.Server.Domain, chatChannel.ChaLineID)
	// URLRegister := fmt.Sprintf("https://%s/register/%s", "586f1140.ngrok.io", chatChannel.ChaLineID)
	view := linebot.View{Type: "full", URL: URLRegister}
	res, err := bot.AddLIFF(view).Do()

	LIFFregister := model.Setting{Detail: "", Name: "LIFFregister", Value: res.LIFFID}
	statusLIFFregister := model.Setting{Detail: "", Name: "statusLIFFregister", Value: "success"}
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if len(chatChannel.Settings) == 0 {
		db.Save(&LIFFregister)
		db.Save(&statusLIFFregister)
		db.Model(&chatChannel).Association("Settings").Append(&LIFFregister)
		db.Model(&chatChannel).Association("Settings").Append(&statusLIFFregister)
	} else {
		for _, setting := range chatChannel.Settings {
			if setting.Name == LIFFregister.Name {
				setting.Value = LIFFregister.Value
			}
			if statusLIFFregister.Name == setting.Name {
				setting.Value = statusLIFFregister.Value
			}
			db.Save(&setting)
		}
	}

	return c.JSON(http.StatusOK, res)
}

type SettingResponse struct {
	LIFFregister       string `json:"LIFFregister"`
	StatusLIFFregister string `json:"statusLIFFregister"`
	StatusAccessToken  string `json:"statusAccessToken"`
	DateStatusToken    string `json:"dateStatusToken"`
}

type DeplayDetailChatChannel struct {
	ID    uint
	Name  string
	Value string
}

type DeplayDetailChatChannels []DeplayDetailChatChannel

// ChatChannelDetailHandler
func ChatChannelDetailHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	var totalEvent int
	var totalAction int
	var richMenu string
	var paginationEventLogs Pagination
	var paginationActionLogs Pagination
	var deplayDetailChatChannels []DeplayDetailChatChannel
	setting := model.Setting{}
	actionLogs := []model.ActionLog{}
	db := model.DB()

	eventLogs := []model.EventLog{}
	eventLogsFilter := db.Model(&eventLogs).Where("chat_channel_id = ?", id).Count(&totalEvent)
	paginationEventLogs = MakePagination(totalEvent, 0, 10)
	eventLogsFilter.Preload("Customer").Find(&eventLogs).Limit(10).Offset(0).Order("id")

	filteractionLogs := db.Model(&actionLogs).Where("chat_channel_id = ?", id).Count(&totalAction)
	paginationActionLogs = MakePagination(totalAction, 0, 10)
	filteractionLogs.Preload("Customer").Find(&actionLogs).Limit(10).Offset(0).Order("id")

	db.Preload("Settings").Preload("Account").Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, id)

	insightFollowers, err := lib.InsightFollowers(chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.Render(http.StatusOK, "chat-channel-detail", echo.Map{
			"title":                    "chat_channel",
			"detail":                   chatChannel,
			"actionLogs":               actionLogs,
			"eventLogs":                eventLogs,
			"insightFollowers":         insightFollowers,
			"paginationActionLogs":     paginationActionLogs,
			"paginationEventLogs":      paginationEventLogs,
			"richMenuDefault":          "",
			"urlRichMenu":              setting,
			"deplayDetailChatChannels": deplayDetailChatChannels,
			"list":                     chatChannel.Settings,
		})
	}
	bot, _ := linebot.New(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	timeNow := time.Now()
	day := fmt.Sprintf("%d", timeNow.Day()-1)
	month := fmt.Sprintf("%d", timeNow.Month())

	if len(day) == 1 {
		day = "0" + day
	}
	if len(month) == 1 {
		month = "0" + month
	}

	dateLineFormat := fmt.Sprintf("%d%s%s", timeNow.Year(), month, day)
	MessageQuota, _ := bot.GetMessageQuota().Do()
	MessageQuotaConsumption, err := bot.GetMessageQuotaConsumption().Do()

	if err != nil {
		return c.Render(http.StatusOK, "chat-channel-detail", echo.Map{
			"title":                    "chat_channel",
			"detail":                   chatChannel,
			"actionLogs":               actionLogs,
			"eventLogs":                eventLogs,
			"insightFollowers":         insightFollowers,
			"paginationActionLogs":     paginationActionLogs,
			"paginationEventLogs":      paginationEventLogs,
			"richMenuDefault":          richMenu,
			"urlRichMenu":              setting,
			"deplayDetailChatChannels": deplayDetailChatChannels,
			"list":                     chatChannel.Settings,
		})
	}

	MessageConsumption, _ := bot.GetMessageConsumption().Do()
	NumberReplyMessages, _ := bot.GetNumberReplyMessages(dateLineFormat).Do()
	NumberPushMessages, _ := bot.GetNumberPushMessages(dateLineFormat).Do()
	NumberBroadcastMessages, _ := bot.GetNumberBroadcastMessages(dateLineFormat).Do()
	NumberMulticastMessages, _ := bot.GetNumberMulticastMessages(dateLineFormat).Do()
	richMenuDefault, _ := bot.GetDefaultRichMenu().Do()

	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Quota Total Usage", Value: fmt.Sprintf("%d", MessageQuota.TotalUsage)})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Quota Type", Value: MessageQuota.Type})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Quota Value", Value: fmt.Sprintf("%d", MessageQuota.Value)})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Quota Consumption Type", Value: MessageQuotaConsumption.Type})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Quota Consumption TotalUsage", Value: fmt.Sprintf("%d", MessageQuotaConsumption.TotalUsage)})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Quota Consumption Value", Value: fmt.Sprintf("%d", MessageQuotaConsumption.Value)})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Consumption TotalUsage", Value: fmt.Sprintf("%d", MessageConsumption.TotalUsage)})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Reply Messages Status", Value: NumberReplyMessages.Status})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Reply Messages Success", Value: strconv.FormatInt(NumberReplyMessages.Success, 16)})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Push Messages Status", Value: NumberPushMessages.Status})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Push Messages Success", Value: strconv.FormatInt(NumberPushMessages.Success, 16)})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Broadcast Messages Status", Value: NumberBroadcastMessages.Status})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Broadcast Messages Success", Value: strconv.FormatInt(NumberBroadcastMessages.Success, 16)})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Multicast Messages Status", Value: NumberMulticastMessages.Status})
	deplayDetailChatChannels = append(deplayDetailChatChannels, DeplayDetailChatChannel{Name: "Multicast Messages Success", Value: strconv.FormatInt(NumberMulticastMessages.Success, 16)})

	if richMenuDefault != nil {
		richMenu = richMenuDefault.RichMenuID
		db.Where("name = ?", richMenuDefault.RichMenuID).Find(&setting)
	}

	return c.Render(http.StatusOK, "chat-channel-detail", echo.Map{
		"title":                    "chat_channel",
		"detail":                   chatChannel,
		"actionLogs":               actionLogs,
		"eventLogs":                eventLogs,
		"insightFollowers":         insightFollowers,
		"paginationActionLogs":     paginationActionLogs,
		"paginationEventLogs":      paginationEventLogs,
		"richMenuDefault":          richMenu,
		"urlRichMenu":              setting,
		"deplayDetailChatChannels": deplayDetailChatChannels,
		"list":                     chatChannel.Settings,
	})
}

// ChatChannelCreateHandler
func ChatChannelCreateViewHandler(c *Context) error {
	typeChatChannels := []string{"Facebook", "Line"}
	systemConfirmation := []string{"auto", "man"}
	csrfValue := c.Get("_csrf")
	return c.Render(http.StatusOK, "chat-channel-form", echo.Map{
		"title":              "chat_channel",
		"typeChatChannels":   typeChatChannels,
		"mode":               "Create",
		"_csrf":              csrfValue,
		"method":             "POST",
		"systemConfirmation": systemConfirmation,
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
		AccountID:             a.User.GetAccountID(),
		ChaImage:              chatChannel.Image,
		ChaWebSite:            chatChannel.WebSite,
		ChaWelcomeMessage:     chatChannel.WelcomeMessage,
		ChaAddress:            chatChannel.Address,
		Settings:              *settingsModel,
	}
	if err := model.DB().Create(&chatChannelModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if chatChannel.Type == "Line" {
		bot, err := linebot.New(chatChannel.ChannelID, chatChannel.ChannelSecret)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		URLRegister := fmt.Sprintf("https://web.%s/register/%s", Conf.Server.Domain, chatChannel.LineID)
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
			&model.Setting{Detail: "", Name: "LIFFregister", Value: LIFFID},
			&model.Setting{Detail: "", Name: "statusLIFFregister", Value: status},
			&model.Setting{Detail: "", Name: "statusAccessToken", Value: status},
			&model.Setting{Detail: "", Name: "dateStatusToken", Value: time.Now().Format("Mon Jan 2 2006")},
		).Error; err != nil {
			return c.JSON(http.StatusBadRequest, chatChannelModel)
		}
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": fmt.Sprintf("/admin/chat_channel/%d", chatChannelModel.ID),
		"data":     chatChannel,
	})
}

// ChatChannelEditHandler
func ChatChannelEditHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)

	db := model.DB()
	db.Where("account = ?", a.GetAccountID()).Preload("Settings").Find(&chatChannel, id)
	setting := SetSettingResponse(chatChannel.Settings)

	typeChatChannels := []string{"Facebook", "Line"}
	return c.Render(http.StatusOK, "chat-channel-form", echo.Map{"method": "PUT",
		"title":            "chat_channel",
		"detail":           chatChannel,
		"typeChatChannels": typeChatChannels,
		"mode":             "Edit",
		"setting":          setting,
	})
}

func SetSettingResponse(settings []*model.Setting) map[string]string {
	var m map[string]string
	m = make(map[string]string)
	for key := range settings {
		m[settings[key].Name] = settings[key].Value
	}
	return m
}

func ChatChannelBroadcastMessageViewHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	customerTypes := []model.CustomerType{}
	a := auth.Default(c)
	db := model.DB()
	db.Preload("Customers").Where("account_id = ?",
		a.User.GetAccountID()).Find(&chatChannel, id)
	db.Where("account_id = ?", a.GetAccountID()).Find(&customerTypes)
	return c.Render(http.StatusOK, "chat-channel-broadcast-message", echo.Map{
		"method":        "POST",
		"title":         "chat_channel",
		"detail":        chatChannel,
		"customerTypes": customerTypes,
		"mode":          "Edit",
	})
}

type RequestBroadcastMessage struct {
	SandDate       int       `json:"sand_date"`
	CustomerType   int       `json:"customer_state"`
	Time           time.Time `json:"time"`
	LineName       string    `json:"line_name"`
	CustomerTypeID string    `json:"customer_type_id"`
}

func ChatChannelBroadcastMessageHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()

	sandDate := c.FormValue("sand_date")
	customerState := c.FormValue("customer_state")
	state := c.FormValue("state")

	db.Preload("Customers").Where("account_id = ?",
		a.User.GetAccountID()).Find(&chatChannel, id)

	bot, err := linebot.New(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	var message linebot.SendingMessage
	switch os := state; os {
	case "Massage":
		message = linebot.NewTextMessage(c.FormValue("text"))
	case "Image":
		image := c.FormValue("image")
		filePath, _, err := lib.UploadteImage(image)
		urlFile := "https://web." + Conf.Server.Domain + filePath
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		message = linebot.NewImageMessage(urlFile, urlFile)
	case "Video":
		video := c.FormValue("video")
		fmt.Println(video[:9])
		filePath, _, err := lib.UploadFile(video, ".mp4")
		urlFile := "https://web." + Conf.Server.Domain + filePath
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		message = linebot.NewVideoMessage(urlFile, urlFile)
	case "Audio":
		audio := c.FormValue("audio")
		filePath, _, err := lib.UploadFile(audio, ".mp3")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		i, err := strconv.Atoi(c.FormValue("duration"))
		fmt.Println(c.FormValue("duration"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		urlFile := "https://web." + Conf.Server.Domain + filePath
		message = linebot.NewAudioMessage(urlFile, i)
	case "Line_Bot_Designer":
		flex := c.FormValue("line_bot_designer")
		flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flex))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		message = linebot.NewFlexMessage("test", flexContainer)
	}

	if customerState == "1" {
		customers := []model.Customer{}
		lineNames := c.FormValue("line_name")
		db.Where("cus_display_name = ?", lineNames).Find(&customers)
		var recipient []string
		for _, customer := range customers {
			recipient = append(recipient, customer.CusLineID)
		}
		bot.Multicast(recipient, message).Do()
	}

	if customerState == "2" {
		customers := []model.Customer{}
		customerTypeID := c.FormValue("customer_type_id")
		db.Preload("CustomerType", "id = ?", customerTypeID).Find(&customers)

		var recipient []string
		for _, customer := range customers {
			recipient = append(recipient, customer.CusLineID)
		}
		_, err = bot.Multicast(recipient, message).Do()

	}

	if customerState == "3" {
		_, err = bot.BroadcastMessage(message).Do()
	}

	if customerState == "4" {
		var testers []model.User
		db.Where("tester = ?", true).Find(&testers)
		var recipient []string
		for _, tester := range testers {
			recipient = append(recipient, tester.LineID)
		}
		_, err = bot.Multicast(recipient, message).Do()
	}

	if sandDate == "3" {
		_, err = bot.BroadcastMessage().Do()
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": fmt.Sprintf("/admin/chat_channel/%d", chatChannel.ID),
	})

}

func ChatChannelDeleteHandler(c *Context) error {
	id := c.Param("id")

	chatChannel := model.DeleteChatChannel(id)
	return c.JSON(http.StatusOK, chatChannel)
}