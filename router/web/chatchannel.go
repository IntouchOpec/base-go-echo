package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

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
	ChaOpenDate        string `form:"cha_open_date" json:"cha_open_date"`
	ChaCloseDate       string `form:"cha_close_date" json:"cha_close_date"`
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
	dateStatusToken := model.Setting{Detail: model.DetailDateStatusToken, Name: model.NameDateStatusToken, Value: "success"}
	statusAccessToken := model.Setting{Detail: model.DetailStatusAccessToken, Name: model.NameStatusAccessToken, Value: time.Now().Format("Mon Jan 2 2006")}
	if chatChannel.Settings == nil {
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

	db.Preload("Settings", "name in (?)",
		[]string{"LIFFregister", "statusLIFFregister", model.NameLIFFIDContent, model.NameLIFFIDReport, model.NameLIFFIDPayment}).Where("account_id = ?",
		a.User.GetAccountID()).Find(&chatChannel, id)

	bot, err := linebot.New(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	URLRegister := fmt.Sprintf("https://web.%s/register/%s", Conf.Server.Domain, chatChannel.ChaLineID)
	URLContent := fmt.Sprintf("https://web.%s/content/%s", Conf.Server.Domain, chatChannel.ChaLineID)
	URLPayment := fmt.Sprintf("https://web.%s/omise", Conf.Server.Domain)
	URLReport := fmt.Sprintf("https://web.%s/report/%s", Conf.Server.Domain, chatChannel.ChaLineID)
	var LIFFIDContent string
	var LIFFIDReport string
	var LIFFIDPayment string
	// URLRegister := fmt.Sprintf("https://%s/register/%s", "586f1140.ngrok.io", chatChannel.ChaLineID)
	view := linebot.View{Type: "full", URL: URLRegister}
	viewURLContent := linebot.View{Type: "full", URL: URLContent}
	viewURLReport := linebot.View{Type: "full", URL: URLReport}
	viewURLPayment := linebot.View{Type: "full", URL: URLPayment}
	res, err := bot.AddLIFF(view).Do()
	res, err = bot.AddLIFF(viewURLContent).Do()
	LIFFIDContent = res.LIFFID
	res, err = bot.AddLIFF(viewURLReport).Do()
	LIFFIDReport = res.LIFFID
	res, err = bot.AddLIFF(viewURLPayment).Do()
	LIFFIDPayment = res.LIFFID
	LIFFregister := model.Setting{Detail: model.DetailLIFFIDRegister, Name: model.NameLIFFregister, Value: res.LIFFID}
	statusLIFFregister := model.Setting{Detail: model.DetailStatusLIFFregister, Name: model.NameStatusLIFFregister, Value: "success"}
	LIFFIDContentSetting := model.Setting{Detail: model.DetailLIFFIDContent, Name: model.NameLIFFIDContent, Value: LIFFIDContent}
	LIFFIDReportSetting := model.Setting{Detail: model.DetailLIFFIDReport, Name: model.NameLIFFIDReport, Value: LIFFIDReport}
	LIFFIDPaymentSetting := model.Setting{Detail: model.DetailLIFFIDPayment, Name: model.NameLIFFIDPayment, Value: LIFFIDPayment}
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if chatChannel.Settings == nil {
		db.Save(&LIFFregister)
		db.Save(&statusLIFFregister)
		db.Model(&chatChannel).Association("Settings").Append(&LIFFregister)
		db.Model(&chatChannel).Association("Settings").Append(&statusLIFFregister, &LIFFIDContentSetting, &LIFFIDReportSetting, &LIFFIDPaymentSetting)
	} else {
		for _, setting := range chatChannel.Settings {
			if setting.Name == LIFFregister.Name {
				setting.Value = LIFFregister.Value
			}
			if statusLIFFregister.Name == setting.Name {
				setting.Value = statusLIFFregister.Value
			}
			if LIFFIDContentSetting.Name == setting.Name {
				setting.Value = LIFFIDContentSetting.Value
			}
			if LIFFIDReportSetting.Name == setting.Name {
				setting.Value = LIFFIDReportSetting.Value
			}
			if LIFFIDPaymentSetting.Name == setting.Name {
				setting.Value = LIFFIDPaymentSetting.Value
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
	var vouchers []*model.Voucher
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

	db.Preload("Settings").Preload("Voucher", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Promotion")
	}).Preload("Account").Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, id)

	db.Preload("Promotion").Where("account_id = ? and is_active = ? and chat_channel_id = ?", a.GetAccountID(), true, chatChannel.ID).Find(&vouchers)

	insightFollowers, err := lib.InsightFollowers(chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.Render(http.StatusOK, "chat-channel-detail", echo.Map{
			"vouchers":                 vouchers,
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
	timeNow := time.Now().AddDate(0, 0, -1)

	dateLineFormat := timeNow.Format("20060102")
	fmt.Println("dateLineFormat", "===")
	fmt.Println("dateLineFormat", dateLineFormat)
	MessageQuota, _ := bot.GetMessageQuota().Do()
	MessageQuotaConsumption, err := bot.GetMessageQuotaConsumption().Do()
	if err != nil {
		return c.Render(http.StatusOK, "chat-channel-detail", echo.Map{
			"vouchers":                 vouchers,
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
		"vouchers":                 vouchers,
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
	csrfValue := c.Get("_csrf")
	return c.Render(http.StatusOK, "chat-channel-form", echo.Map{
		"title":            "chat_channel",
		"mode":             "Create",
		"method":           "POST",
		"typeChatChannels": typeChatChannels,
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
	openDate, err := time.Parse("15:04", chatChannel.ChaOpenDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	closeDate, err := time.Parse("15:04", chatChannel.ChaCloseDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
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
		ChaOpenDate:           openDate,
		ChaCloseDate:          closeDate,
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
		URLContent := fmt.Sprintf("https://web.%s/content/%s", Conf.Server.Domain, chatChannel.LineID)
		URLPayment := fmt.Sprintf("https://web.%s/omise", Conf.Server.Domain, chatChannel.LineID)
		URLReport := fmt.Sprintf("https://web.%s/report/%s", Conf.Server.Domain, chatChannel.LineID)
		viewURLRegister := linebot.View{Type: "full", URL: URLRegister}
		viewURLContent := linebot.View{Type: "full", URL: URLContent}
		viewURLReport := linebot.View{Type: "full", URL: URLReport}
		viewURLPayment := linebot.View{Type: "full", URL: URLPayment}
		var status string = "success"
		var LIFFIDRegister string = ""
		var LIFFIDContent string = ""
		var LIFFIDReport string = ""
		var LIFFIDPayment string = ""
		res, err := bot.AddLIFF(viewURLRegister).Do()
		if err != nil {
			status = "error"
		} else {
			LIFFIDRegister = res.LIFFID
			res, err = bot.AddLIFF(viewURLContent).Do()
			LIFFIDContent = res.LIFFID
			res, err = bot.AddLIFF(viewURLReport).Do()
			LIFFIDReport = res.LIFFID
			res, err = bot.AddLIFF(viewURLPayment).Do()
			LIFFIDPayment = res.LIFFID
		}
		if err := model.DB().Model(&chatChannelModel).Association("Settings").Append(
			&model.Setting{Detail: model.DetailLIFFIDRegister, Name: model.NameLIFFregister, Value: LIFFIDRegister},
			&model.Setting{Detail: model.DetailLIFFIDContent, Name: model.NameLIFFIDContent, Value: LIFFIDContent},
			&model.Setting{Detail: model.DetailLIFFIDReport, Name: model.NameLIFFIDReport, Value: LIFFIDReport},
			&model.Setting{Detail: model.DetailLIFFIDPayment, Name: model.NameLIFFIDPayment, Value: LIFFIDPayment},
			&model.Setting{Detail: model.DetailStatusLIFFregister, Name: model.NameStatusLIFFregister, Value: status},
			&model.Setting{Detail: model.DetailStatusAccessToken, Name: model.NameStatusAccessToken, Value: status},
			&model.Setting{Detail: model.DetailDateStatusToken, Name: model.NameDateStatusToken, Value: time.Now().Format("Mon Jan 2 2006")},
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

func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
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
		ctx := context.Background()
		imagePath, err := lib.UploadGoolgeStorage(ctx, image, "images/broadcast/")
		urlFile := fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, imagePath)
		fmt.Println("urlFile", urlFile)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		message = linebot.NewImageMessage(urlFile, urlFile)
	case "Video":
		video := c.FormValue("video")
		ctx := context.Background()
		videoPath, err := lib.UploadGoolgeStorage(ctx, video, "video/broadcast/")
		urlFile := fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, videoPath)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		message = linebot.NewVideoMessage(urlFile, urlFile)
	case "Audio":
		audio := c.FormValue("audio")
		ctx := context.Background()
		audioPath, err := lib.UploadGoolgeStorage(ctx, audio, "audio/broadcast/")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		i, err := strconv.Atoi(c.FormValue("duration"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		urlFile := fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, audioPath)
		message = linebot.NewAudioMessage(urlFile, i)
	case "Line_Bot_Designer":
		flex := c.FormValue("line_bot_designer")
		flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flex))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		message = linebot.NewFlexMessage("test", flexContainer)
	}
	var multicastCall *linebot.MulticastCall
	var broadcastMessageCall *linebot.BroadcastMessageCall
	var recipient []string
	switch customerState {
	case "1":
		customers := []model.Customer{}
		lineNames := c.FormValue("line_name")
		db.Where("cus_display_name = ?", lineNames).Find(&customers)
		for _, customer := range customers {
			recipient = append(recipient, customer.CusLineID)
		}
		multicastCall = bot.Multicast(recipient, message)
	case "2":
		customers := []model.Customer{}
		customerTypeID := c.FormValue("customer_type_id")
		db.Preload("CustomerType", "id = ?", customerTypeID).Find(&customers)
		var recipient []string
		for _, customer := range customers {
			recipient = append(recipient, customer.CusLineID)
		}
		multicastCall = bot.Multicast(recipient, message)
	case "3":
		broadcastMessageCall = bot.BroadcastMessage(message)
	case "4":
		var testers []model.User
		db.Where("tester = ?", true).Find(&testers)
		var recipient []string
		for _, tester := range testers {
			recipient = append(recipient, tester.LineID)
		}
		multicastCall = bot.Multicast(recipient, message)
	}
	switch sandDate {
	case "1":
		if broadcastMessageCall == nil {
			_, err = multicastCall.Do()
		} else {
			_, err = broadcastMessageCall.Do()
		}
	case "2":
		date := c.FormValue("date")
		timeValue := c.FormValue("time")

		datetime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", date, timeValue))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		now := time.Now()
		then := now.Add(7 * time.Hour)
		seconds := datetime.Sub(then)
		time.AfterFunc(seconds, func() {
			if broadcastMessageCall == nil {
				_, err = multicastCall.Do()
			} else {
				_, err = broadcastMessageCall.Do()
			}
		})
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

func ChatChannelVoucherRegisterHandler(c *Context) error {
	id := c.Param("id")
	db := model.DB()
	accID := auth.Default(c).GetAccountID()
	chatChannel := model.ChatChannel{}
	voucher := model.Voucher{}
	voucherID := c.FormValue("voucher_id")
	if err := db.Where("account_id = ?", accID).Find(&chatChannel, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Where("account_id = ?", accID).Find(&voucher, voucherID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Model(&chatChannel).Association("Voucher").Append(&voucher).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, chatChannel)
}
