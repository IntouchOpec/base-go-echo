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

type PostbackAction struct {
	Action  string `json:"action"`
	DateStr string `json:"date"`
}

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
		var keyWord string
		con.Event = event.Source
		db.Where("cus_line_id = ? and chat_channel_id = ?", event.Source.UserID, chatChannel.ID).Find(&customer)
		con.Customer = customer
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
			case "my_voucher":
				messageReply, err = VoucherListHandler(&con)
			case "comment":

			case "choive_man":
				fmt.Println("choive_man")
				messageReply, err = CalandarHandler(&con, postBackAction.DateStr)
			// case "choive_auto":
			// 	fmt.Println("choive_auto")
			// 	messageReply, err = ServiceNowListHandler(&con)
			// 	fmt.Println(err)
			case "booking_now":
				fmt.Println("booking_now")
				messageReply, err = ServiceNowListHandler(&con)
			case "booking_appointment":
				fmt.Println("booking_appointment")
				messageReply, err = ServiceDateListHandler(&con, event.Postback.Params.Datetime)
			case "booking":
				fmt.Println("booking_appointment")
				messageReply, err = ServiceDateListHandler(&con, event.Postback.Params.Datetime)
			}
			_, err = bot.ReplyMessage(event.ReplyToken, messageReply).Do()
			if err != nil {
				fmt.Println(err)
			}
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
					// messageReply, err = CalandarHandler(&con)
				case "Service ":
					messageReply, err = SaveServiceHandler(&con)
				case "booking ":
					messageReply, err = ServiceListLineHandler(&con)
				case "timeslot":
					messageReply, err = ThankyouTemplate(&con)
				case "check":
					messageReply, err = CheckStatusOpen(&con)

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
