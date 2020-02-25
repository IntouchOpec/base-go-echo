package channel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/IntouchOpec/base-go-echo/lib/dialogflow"
	"github.com/IntouchOpec/base-go-echo/lib/lineapi"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/labstack/echo"
)

// HandleWebHookLineAPI webhook for connent line api.
func HandleWebHookLineAPI(c echo.Context) error {
	name := c.Param("account")
	ChannelID := c.Param("ChannelID")
	var con Context
	con.DB = model.DB()

	if err := con.DB.Where("acc_name = ?", name).Find(&con.Account).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	fmt.Println(name, ChannelID)
	if err := con.DB.Preload("Settings").Where("cha_channel_id = ?", ChannelID).Find(&con.ChatChannel).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	bot, err := lineapi.ConnectLineBot(con.ChatChannel.ChaChannelSecret, con.ChatChannel.ChaChannelAccessToken)

	if err != nil {
		fmt.Println(err, "bot", bot)
		return c.JSON(http.StatusBadRequest, err)
	}

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
		con.DB.Where("cus_line_id = ? and account_id = ?", event.Source.UserID, con.Account.ID).Find(&con.Customer)
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
				fmt.Println(err)
				return c.JSON(http.StatusBadRequest, err)
			}
			con.PostbackAction = &postBackAction
			switch postBackAction.Action {
			case "location":
				messageReply, err = LocationHandler(&con)
			case "promotion":
				messageReply, err = PromotionHandler(&con)
			case "voucher":
				messageReply, err = VoucherListHandler(&con)
			case "service":
				ChooseService(&con)
			case "report":
				fmt.Println("report")
				messageReply, err = ReportListHandler(&con)
			case "comment":
				fmt.Println("comment")
				messageReply, err = ReportListHandler(&con)
			case "content":
				fmt.Println("content")
				messageReply, err = ContentListHandler(&con)
			case "choive_man":
				fmt.Println("choive_man")
			case "checkout":
				fmt.Println("checkout")
				messageReply, err = ChackOutHandler(&con)
			case "status":
				// messageReply, err = CalandarHandler(&con, postBackAction.DateStr)
			case "calendar_next":
				// messageReply, err = CalandarHandler(&con, postBackAction.DateStr)
			// case "calendar":
			// 	fmt.Println("calendar")
			// messageReply, err = ServiceList(&con)
			// case "choose_timeslot":
			// fmt.Println("choose_timeslot")
			// messageReply, err = ServiceListLineHandler(&con)
			case "booking":
				fmt.Println("booking")
				BookingHandler(&con)
			}
			_, err = bot.ReplyMessage(event.ReplyToken, messageReply).Do()
			if err != nil {
				fmt.Println(err)
			}
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var dp dialogflow.DialogFlowProcessor
				dp.Init(con.Account.AccProjectID, con.Account.AccAuthJSONFilePath, con.Account.AccLang, con.Account.AccTimeZone)
				replyDialogflow := dp.ProcessNLP(message.Text, con.Customer.CusDisplayName)
				fmt.Println(replyDialogflow)
				if err := con.DB.Where("account_id = ? and chat_input = ?", con.Account.ID, message.Text).Find(&chatAnswer).Error; err != nil {
					con.DB.Where("account_id = ? and chat_input = 'error'", con.Account.ID).Find(&chatAnswer)
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
			messageReply, err = WelcomeHandle(&con)
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
