package channel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/labstack/echo"
)

// var dp lib.DialogFlowProcessor
// dp.Init(account.AccProjectID, account.AccAuthJSONFilePath, account.AccLang, account.AccTimeZone)
// replyDialogflow := dp.ProcessNLP(messageText, customer.CusDisplayName)
// fmt.Println(replyDialogflow)
// HandleWebHookLineAPI webhook for connent line api.
func HandleWebHookLineAPI(c echo.Context) error {
	db := model.DB()
	name := c.Param("account")
	ChannelID := c.Param("ChannelID")
	var account model.Account
	var chatChannel model.ChatChannel
	// var customer model.Customer
	// var eventLog model.EventLog
	var con Context
	con.DB = db

	if err := db.Where("acc_name = ?", name).Find(&account).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := db.Preload("Settings", "name = ?", model.NameLIFFIDPayment).Where("cha_channel_id = ?", ChannelID).Find(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		fmt.Println(err, "bot", bot)
		return c.JSON(http.StatusBadRequest, err)
	}

	con.Account = account
	con.ChatChannel = chatChannel
	events, err := bot.ParseRequest(c.Request())
	if err != nil {
		fmt.Println(err, "events", events)
		if err == linebot.ErrInvalidSignature {
			return c.String(400, linebot.ErrInvalidSignature.Error())
		}
		return c.String(500, "internal")
	}

	for _, event := range events {
		con.Event = event
		db.Where("cus_line_id = ? and account_id = ?", event.Source.UserID, account.ID).Find(&con.Customer)
		// con.Customer = customer
		chatAnswer := model.ChatAnswer{}
		eventType := event.Type
		chatAnswer.AnsInputType = string(eventType)
		var messageReply linebot.SendingMessage
		fmt.Println(eventType, "eventType")
		switch eventType {
		case linebot.EventTypePostback:
			u, _ := url.Parse(fmt.Sprintf("?%s", event.Postback.Data))
			var postBackActionStr string
			q, _ := url.ParseQuery(u.RawQuery)
			for key, value := range q {
				postBackActionStr += fmt.Sprintf(`"%s": "%s",`, key, value[0])
			}
			postBackActionStr = fmt.Sprintf(fmt.Sprintf("{%s}", postBackActionStr[:len(postBackActionStr)-1]))
			postBackAction := PostbackAction{}
			if err := json.Unmarshal([]byte(postBackActionStr), &postBackAction); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			con.PostbackAction = &postBackAction
			switch postBackAction.Action {
			case "location":
				messageReply, err = LocationHandler(&con)
			case "promotiom":
				messageReply, err = PromotionHandler(&con)
			case "voucher":
				messageReply, err = VoucherListHandler(&con)
			case "service":
				messageReply, err = ChooseService(&con)
			case "report":
				fmt.Println("report")
			case "content":
				fmt.Println("content")
			case "choive_man":
				fmt.Println("choive_man")
				messageReply, err = CalandarHandler(&con, postBackAction.DateStr)
			case "calendar_next":
				messageReply, err = CalandarHandler(&con, postBackAction.DateStr)
			case "calendar":
				fmt.Println("calendar")
				messageReply, err = ServiceList(&con)
			case "choose_timeslot":
				fmt.Println("choose_timeslot")
				messageReply, err = ServiceListLineHandler(&con)
			case "booking_timeslot":
				fmt.Println("booking_timeslot")
				messageReply, err = BookingTimeSlotHandler(&con)
			case "booking_now":
				fmt.Println("booking_now")
				messageReply, err = ServiceNowListHandler(&con)
			case "booking_appointment":
				fmt.Println("booking_appointment")
				messageReply, err = ServiceDateListHandler(&con, event.Postback.Params.Datetime)
			case "booking":
				fmt.Println("booking")
				fmt.Println(postBackAction)
				messageReply, err = BookingServiceHandler(&con)
			}
			_, err = bot.ReplyMessage(event.ReplyToken, messageReply).Do()
			if err != nil {
				fmt.Println(err)
			}
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := db.Where("account_id = ? and chat_input = ?", account.ID, message.Text).Find(&chatAnswer).Error; err != nil {
					db.Where("account_id = ? and chat_input = 'error'", account.ID).Find(&chatAnswer)
				}
				messageReply, err := bot.ReplyLineMessage(chatAnswer)
				if err != nil {
					messageReply = linebot.NewTextMessage("error")
				}
				_, err = bot.ReplyMessage(event.ReplyToken, messageReply).Do()
			case *linebot.ImageMessage:
			case *linebot.VideoMessage:
			case *linebot.AudioMessage:
			case *linebot.FileMessage:
			case *linebot.LocationMessage:
				messageReply := linebot.NewTextMessage(fmt.Sprintf("%v", message))
				_, err = bot.ReplyMessage(event.ReplyToken, messageReply).Do()
				// _, err := bot.ReplyMessage(event.ReplyToken, textMessage).Do()
				// if err != nil {
				// 	act := model.ActionLog{
				// 		ActName:       "LocationMessage",
				// 		ActStatus:     model.StatusFail,
				// 		ChatChannelID: chatChannel.ID,
				// 		CustomerID:    customer.ID}
				// 	act.CreateAction()
				// 	return err
				// }
				// act := model.ActionLog{
				// 	ActName:       "LocationMessage",
				// 	ActStatus:     model.StatusSuccess,
				// 	ChatChannelID: chatChannel.ID,
				// 	CustomerID:    customer.ID}
				// if err := model.DB().Create(&act).Error; err != nil {
				// 	// return c.JSON(http.StatusBadRequest, err)
				// }
			case *linebot.StickerMessage:
			case *linebot.TemplateMessage:
			case *linebot.ImagemapMessage:
			case *linebot.FlexMessage:
			}
			_, err = bot.ReplyMessage(event.ReplyToken, messageReply).Do()
			if err != nil {
				fmt.Println("err", err, "ReplyMessage")
			}
			return c.JSON(200, "")

		case linebot.EventTypeFollow:
			messageReply, err = welcomeHandle(&c, event, &chatChannel)
			if err != nil {
				fmt.Println("err", err)
				return c.JSON(http.StatusBadRequest, err)
			}
			_, err = bot.ReplyMessage(event.ReplyToken, messageReply).Do()
			if err != nil {
				fmt.Println("err", err)
				return c.JSON(http.StatusBadRequest, err)
			}
		case linebot.EventTypeUnfollow:
		case linebot.EventTypeJoin:
			fmt.Println(linebot.EventTypeJoin)

		}
	}
	return c.JSON(200, "")
}
