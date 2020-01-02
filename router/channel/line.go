package channel

import (
	"fmt"
	"net/http"

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
	var customer model.Customer
	// var eventLog model.EventLog
	var con Context
	con.DB = db

	if err := db.Where("acc_name = ?", name).Find(&account).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := db.Where("cha_channel_id = ?", ChannelID).Find(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	con.Account = account
	con.ChatChannel = chatChannel
	events, err := bot.ParseRequest(c.Request())

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return c.String(400, linebot.ErrInvalidSignature.Error())
		}
		return c.String(500, "internal")
	}

	for _, event := range events {
		var keyWord string
		con.Event = event.Source
		db.Where("cus_line_id = ? and chat_channel_id = ?", event.Source.UserID, chatChannel.ID).Find(&customer)
		con.Customer = customer
		chatAnswer := model.ChatAnswer{}
		eventType := event.Type
		chatAnswer.AnsInputType = string(eventType)
		var messageReply linebot.SendingMessage

		switch eventType {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				messageText := message.Text
				keyWord = message.Text
				con.Massage = message
				if len(messageText) >= 8 {
					keyWord = messageText[0:8]
				}
				switch keyWord {
				case "service", "service ":
					if len(messageText) > 8 {
						messageReply, err = ServiceList(&con)
					} else {
						messageReply, err = ChooseService(&con)
					}
				case "calendar", "booking":
					messageReply, err = CalandarHandler(&con)
				case "Service ":
					messageReply, err = SaveServiceHandler(&con)
				case "booking ":
					messageReply, err = ServiceListLineHandler(&con)
				case "timeslot":
					messageReply, err = ThankyouTemplate(&con)
				case "promotio":
					messageReply, err = PromotionHandler(&con)
				case "location":
					messageReply, err = LocationHandler(&con)
				case "check":

				case "my voucher":

				case "comment":

				default:
					// if err := db.Find(&chatAnswer).Error; err != nil {
					// 	db.Find(&chatAnswer)
					// }
					// messageReply = linebot.NewTextMessage(chatAnswer.AnsReply)
					// eventLog.EvenType = string(event.Type)
				}
			case *linebot.ImageMessage:
			case *linebot.VideoMessage:
			case *linebot.AudioMessage:
			case *linebot.FileMessage:
			case *linebot.LocationMessage:
				textMessage := linebot.NewTextMessage(fmt.Sprintf("%v", message))
				_, err := bot.ReplyMessage(event.ReplyToken, textMessage).Do()
				if err != nil {
					act := model.ActionLog{
						ActName:       "LocationMessage",
						ActStatus:     model.StatusFail,
						ChatChannelID: chatChannel.ID,
						CustomerID:    customer.ID}
					act.CreateAction()
					return err
				}
				act := model.ActionLog{
					ActName:       "LocationMessage",
					ActStatus:     model.StatusSuccess,
					ChatChannelID: chatChannel.ID,
					CustomerID:    customer.ID}
				if err := model.DB().Create(&act).Error; err != nil {
					// return c.JSON(http.StatusBadRequest, err)
				}
			case *linebot.StickerMessage:
			case *linebot.TemplateMessage:
			case *linebot.ImagemapMessage:
			case *linebot.FlexMessage:
			}
			_, err = bot.ReplyMessage(event.ReplyToken, messageReply).Do()
			return c.JSON(200, "")

		case linebot.EventTypeFollow:
			messageReply = welcomeHandle(&c, event, &chatChannel)
		case linebot.EventTypeUnfollow:
		case linebot.EventTypeJoin:
			fmt.Println(linebot.EventTypeJoin)

		}
	}
	return c.JSON(200, "")
}
