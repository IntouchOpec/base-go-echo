package channel

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/labstack/echo"
)

// HandleWebHookLineAPI webhook for connent line api.
func HandleWebHookLineAPI(c echo.Context) error {
	// client := &lib.ClientLine{}
	name := c.Param("account")
	ChannelID := c.Param("ChannelID")
	account := model.Account{}
	chatChannel := model.ChatChannel{}
	db := model.DB()
	if err := db.Where("acc_name = ?", name).Find(&account).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := db.Where("cha_channel_id = ?", ChannelID).Find(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	events, err := bot.ParseRequest(c.Request())

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return c.String(400, linebot.ErrInvalidSignature.Error())
		}
		return c.String(500, "internal")
	}
	customer := model.Customer{}
	eventLogs := []model.EventLog{}
	// actionLogs := []model.ActionLog{}
	// actionLog := model.ActionLog{}
	eventLog := model.EventLog{}
	for _, event := range events {
		var keyWord string

		db.Where("cus_line_id = ?", event.Source.UserID).Find(&customer)
		chatAnswer := model.ChatAnswer{}
		eventType := event.Type
		// eventLog
		chatAnswer.AnsInputType = string(eventType)
		var messageReply linebot.SendingMessage

		switch eventType {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				messageText := message.Text
				if len(messageText) >= 8 {
					keyWord = messageText[0:8]
				}
				switch messageText {

				case "service":

				case "calendar", "booking":

				case "Service":

				case "booking ":

				case "promotio":

				case "location":

				default:
					db.Find(&chatAnswer)

					messageReply = linebot.NewTextMessage(chatAnswer.AnsReply)
					eventLog.EvenType = string(event.Type)
				}

				if keyWord == "calendar" || messageText == "booking" {
					var m string
					if len(message.Text) > 8 {
						m = lib.MakeCalenda(messageText[9:19])
					} else {
						m = lib.MakeCalenda("")
					}
					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {
						log.Println(err)
					}
					flexMessage := linebot.NewFlexMessage("ตาราง", flexContainer)
					_, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do()

					if err != nil {
						act := model.ActionLog{
							ActName:       "calendar",
							ActStatus:     model.StatusFail,
							ActUserID:     event.Source.UserID,
							ChatChannelID: chatChannel.ID,
							CustomerID:    customer.ID}
						act.CreateAction()
						return err
					}
					act := model.ActionLog{
						ActName:       "calendar",
						ActChannel:    model.ActionChannelLine,
						ActStatus:     model.StatusSuccess,
						ActUserID:     event.Source.UserID,
						ChatChannelID: chatChannel.ID,
						CustomerID:    customer.ID}
					db.Create(&act)
				} else if keyWord == "Service" {
					// t, _ := time.Parse("2006-01-02 15:04", messageText[8:]+" 01:00")

					// serviceSlot := []model.ServiceSlot{}
					// if err := model.DB().Preload("Bookings", "booked_date = ?", t).Preload("service").Where("Day = ?", int(t.Weekday())).Find(&serviceSlot).Error; err != nil {
					// 	fmt.Println(err)
					// 	return nil
					// }

					// m := serviceListLineTemplate(serviceSlot, messageText[8:])

					// flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {
						log.Println(err)
					}
					// flexMessage := linebot.NewFlexMessage("ตาราง", flexContainer)

					// res, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
					if err != nil {
						return err
					}
					act := model.ActionLog{
						ActName:       "service",
						ActStatus:     model.StatusSuccess,
						ActUserID:     event.Source.UserID,
						ChatChannelID: chatChannel.ID,
						CustomerID:    customer.ID}
					act.CreateAction()
				} else if keyWord == "booking " {
					t, _ := time.Parse("2006-01-02 15:04", messageText[8:18]+" 01:00")
					service := model.Service{}
					if err := model.DB().Preload("ServiceSlots", "start = ? and day = ?", messageText[19:24], int(t.Weekday())).Where("Name = ?", messageText[31:]).Find(&service).Error; err != nil {
						fmt.Println(err, "err")
						return nil
					}

					custo := model.Customer{}
					if err := model.DB().Where("cus_line_id = ?", event.Source.UserID).Find(&custo).Error; err != nil {
						fmt.Println(err, "err")
						return nil
					}
					// serviceSlot := service.ServiceSlots[0]
					// booking := model.Booking{ServiceSlotID: service.ServiceSlots[0].ID, BookedDate: t, ChatChannelID: chatChannel.ID, CustomerID: custo.ID}
					// _, err := booking.SaveBooking()
					var m string
					if err != nil {
						m = ` { "type": "bubble", "size": "nano", "header": { "type": "box", "layout": "vertical", "contents": [
								{ "type": "text", "text": "คิวเต็ม", "color": "#ffffff", "align": "start", "size": "md", "gravity": "center" }
							  ], "backgroundColor": "#FF6B6E", "paddingTop": "19px", "paddingAll": "12px", "paddingBottom": "16px"},
							"body": { "type": "box", "layout": "vertical", "contents": [{ "type": "box", "layout": "horizontal", "contents": [
									{ "type": "text", "text": "กรุณาเลือกคิวใหม่", "color": "#8C8C8C", "size": "sm", "wrap": true }],
								  "flex": 1
								}
							  ],
							  "spacing": "md",  "paddingAll": "12px"
							},
							"styles": {
							  "footer": { "separator": false } }
						  }`
					} else if keyWord == "my voucher" {
						// customer := model.Customer{}
						// model.DB().Preload("Promotions", "type = ?", "voucher", func(db *gorm.DB) *gorm.DB {
						// 	return db.Preload("Settings")
						// }).Where("line_id = ?", message.ID).Find(&customer)
						// var template string
						// for index := 0; index < len(customer.Promotions); index++ {
						// 	// promotion := customer.Promotions[index]
						// 	// template = template + web.VoucherTemplate(customer.Promotions[index]) + ","
						// }
						// template = fmt.Sprintf(`{ "type": "carousel", "contents": [%s]}`, template[:len(template)-1])
						// flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
						// if err != nil {
						// 	return err
						// }
						// flexMessage := linebot.NewFlexMessage("your voucher", flexContainer)
						// res, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
						// if err != nil {
						// 	act := model.ActionLog{Name: "user voucher", Status: model.StatusFail, Type: model.TypeActionLine,
						// 		UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
						// 	act.CreateAction()
						// 	return err
						// }
						// fmt.Println(res)
						// act := model.ActionLog{Name: "user voucher", Status: model.StatusSuccess,
						// 	Type: model.TypeActionLine, UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
						// act.CreateAction()
					} else if keyWord == "comment" {

					} else {
						// m = ThankyouTemplate(chatChannel, booking, serviceSlot)
					}
					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {
						return err
					}
					flexMessage := linebot.NewFlexMessage("ขอบคุณ", flexContainer)
					// flexMessage := linebot.NewPostbackAction("label", "data", "text", "displayText")
					bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
				} else if keyWord == "promotio" {
					// promotions := []*model.Promotion{}
					// model.DB().Where("promotion_type = ?", "Promotion").Find(&promotions)
					// m := PromotionsTemplate(promotions)
					// flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					// if err != nil {

					// 	return err
					// }
					// flexMessage := linebot.NewFlexMessage("ตาราง", flexContainer)
					// res, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
					// fmt.Println(res)
					// if err != nil {
					// 	return err
					// }
					// act := model.ActionLog{Name: "promotion", Status: model.StatusSuccess, Type: model.TypeActionLine, UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
					// act.CreateAction()
				} else if keyWord == "location" {
					// position := chatChannel.GetSetting([]string{"Latitude", "Longitude"})
					// Latitude, _ := strconv.ParseFloat(position["Latitude"], 64)
					// Longitude, _ := strconv.ParseFloat(position["Longitude"], 64)

					// location := linebot.NewLocationMessage(chatChannel.Name, chatChannel.Address, Latitude, Longitude)

					// _, err := bot.ReplyMessage(event.ReplyToken, location).Do()
					// if err != nil {
					// 	return err
					// }
					// act := model.ActionLog{Name: "location", Status: model.StatusSuccess, Type: model.TypeActionLine, UserID: event.Source.UserID,
					// 	ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
					// act.CreateAction()
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

		case linebot.EventTypeFollow:
			welcomeHandle(&c, event, &chatChannel, bot)
		case linebot.EventTypeUnfollow:
			evenLog := model.EventLog{
				ChatChannelID:  chatChannel.ID,
				EvenReplyToken: event.ReplyToken,
				EvenType:       string(event.Type),
				EvenLineID:     event.Source.UserID,
				CustomerID:     customer.ID}
			model.DB().Create(&evenLog)
		}
		_, err = bot.ReplyMessage(event.ReplyToken, messageReply).Do()
		fmt.Println(err)

		eventLog.EvenReplyToken = event.ReplyToken
		eventLog.CustomerID = customer.ID
		eventLog.ChatChannelID = chatChannel.ID
		eventLog.ChatChannelID = chatChannel.ID
		eventLog.EvenLineID = event.Source.UserID

		// actionLog.CustomerID = customer.ID
		// actionLog.ChatChannelID = chatChannel.ID

		eventLogs = append(eventLogs, eventLog)
		// actionLogs = append(actionLogs, actionLog)

	}
	if err := db.Create(&eventLogs).Error; err != nil {
		fmt.Println(err)
	}
	// if err := db.Create(&actionLogs).Error; err != nil {
	// 	fmt.Println(err)
	// }

	return c.JSON(200, "")

}

func welcomeHandle(c *echo.Context, event *linebot.Event, chatChannel *model.ChatChannel, bot *lib.ClientLine) {
	customer := model.Customer{
		CusLineID: event.Source.UserID,
		AccountID: chatChannel.AccountID}
	settingNames := []string{"LIFFregister"}
	setting := chatChannel.GetSetting(settingNames)
	if err := model.DB().FirstOrCreate(&customer, model.Customer{
		CusLineID: event.Source.UserID,
		AccountID: chatChannel.AccountID}).Error; err != nil {
		// return c.JSON(http.StatusBadRequest, err)
	}
	jsonFlexMessage := FollowTemplate(chatChannel, setting)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(jsonFlexMessage))
	if err != nil {
		// return c.JSON(http.StatusBadRequest, err)
	}
	flexMessage := linebot.NewFlexMessage(chatChannel.ChaWelcomeMessage, flexContainer)
	_, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
	fmt.Println(err)
	evenLog := model.EventLog{
		ChatChannelID:  chatChannel.ID,
		EvenReplyToken: event.ReplyToken,
		EvenType:       string(event.Type),
		EvenLineID:     event.Source.UserID,
		CustomerID:     customer.ID}
	model.DB().Create(&evenLog)
	act := model.ActionLog{
		ActName:       "follow",
		ActStatus:     model.StatusSuccess,
		ActUserID:     event.Source.UserID,
		ChatChannelID: chatChannel.ID,
		CustomerID:    customer.ID}
	if err := model.DB().Create(&act).Error; err != nil {
		// return c.JSON(http.StatusBadRequest, err)
	}
	// return c.JSON(200, "")
}

// serviceListLineTemplate
// func serviceListLineTemplate(serviceSlot []model.ServiceSlot, dateTime string) string {
// 	var slotTime string
// 	var buttonTime string
// 	var serviceList string
// 	var count int
// 	count = 0
// 	for t := 0; t < len(serviceSlot); t++ {
// 		if count == 2 {
// 			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
// 			buttonTime = ""
// 			count = 0
// 		}

// 		if len(serviceSlot[t].Bookings) > 0 {
// 			if serviceSlot[t].Bookings[0].Queue < serviceSlot[t].Amount {
// 				buttonTime = buttonTime + fmt.Sprintf(`{"type": "button", "style": "secondary", "margin": "sm", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`,
// 					serviceSlot[t].Start, serviceSlot[t].End, "เต็มแล้ว")
// 			} else {
// 				buttonTime = buttonTime + fmt.Sprintf(`{"type": "button","style": "primary", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`,
// 					serviceSlot[t].Start, serviceSlot[t].End, "booking"+" "+dateTime+" "+serviceSlot[t].Start+"-"+serviceSlot[t].End+" "+serviceSlot[t].Service.Name)
// 			}
// 		} else {
// 			buttonTime = buttonTime + fmt.Sprintf(`{"type": "button","style": "primary", "margin": "sm", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`,
// 				serviceSlot[t].Start, serviceSlot[t].End, "booking"+" "+dateTime+" "+serviceSlot[t].Start+"-"+serviceSlot[t].End+" "+serviceSlot[t].Service.Name)
// 		}

// 		count = count + 1
// 		if t == len(serviceSlot)-1 {
// 			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
// 			serviceList += fmt.Sprintf(`{"type": "bubble", "hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s"},
// 			"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
// 				{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl" },
// 				{ "type": "box", "layout": "baseline", "contents": [
// 					{ "type": "text", "text": "฿%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
// 				] }
// 				%s]
// 			}}`, serviceSlot[t].Service.Image, serviceSlot[t].Service.Name, strconv.FormatInt(int64(serviceSlot[t].Service.Price), 10), slotTime)
// 		} else if serviceSlot[t].Service.ID != serviceSlot[t+1].Service.ID {
// 			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
// 			serviceList = serviceList + fmt.Sprintf(`{"type": "bubble", "hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s"},
// 			"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
// 				{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl" },
// 				{ "type": "box", "layout": "baseline", "contents": [
// 					{ "type": "text", "text": "฿%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
// 				] }
// 				%s
// 			]
// 			}},`, serviceSlot[t].Service.Image, serviceSlot[t].Service.Name, strconv.FormatInt(int64(serviceSlot[t].Service.Price), 10), slotTime)
// 			slotTime = ""
// 			count = 0
// 			buttonTime = ""
// 		}

// 	}
// 	var nextPage string = `{ "type": "bubble", "body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
// 		{ "type": "button", "flex": 1, "gravity": "center", "action": { "type": "uri", "label": "See more", "uri": "https://linecorp.com" } }] }}`
// 	var serviceTamplate string = fmt.Sprintf(`{ "type": "carousel", "contents": [%s, %s]}`, serviceList, nextPage)
// 	return serviceTamplate
// }

// FollowTemplate
func FollowTemplate(chatChannel *model.ChatChannel, settings map[string]string) string {
	template := fmt.Sprintf(`{
		"type": "bubble",
		"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover"},
		"body": { "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "%s ", "weight": "bold", "size": "xl" },
			{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
				{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
					{ "type": "text", "text": "%s ", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
				]}
			]}]
		},
		"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
			{ "type": "button", "style": "link", "height": "sm", "action": { "type": "uri", "label": "REGISTER", "uri": "line://app/%s"} },
			{ "type": "button", "style": "link", "height": "sm", "action": { "type": "uri", "label": "WEBSITE", "uri": "https://%s"}},
			{ "type": "spacer", "size": "sm" }
		],
		"flex": 0
		}
	  }`, chatChannel.ChaImage, chatChannel.ChaName, chatChannel.ChaWelcomeMessage, settings["LIFFregister"], chatChannel.ChaWebSite)
	return template
}

// // ThankyouTemplate
// func ThankyouTemplate(ChatChannel model.ChatChannel, booking model.Booking, serviceSlot *model.ServiceSlot) string {
// 	var serviceTamplate string = fmt.Sprintf(`{
// 		"type": "bubble",
// 		"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
// 		"body": {
// 		  "type": "box",
// 		  "layout": "vertical",
// 		  "contents": [
// 			{ "type": "text", "text": "จองสำเร็จ", "weight": "bold", "size": "xl" },
// 			{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
// 				{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
// 					{ "type": "text", "text": "Place", "color": "#aaaaaa", "size": "sm", "flex": 1 },
// 					{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
// 				] },
// 				{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
// 					{ "type": "text", "text": "Time", "color": "#aaaaaa", "size": "sm", "flex": 1 },
// 					{ "type": "text", "text": "%s - %s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
// 				] }
// 			  ] }
// 		  	]
// 		},
// 		"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
// 			{ "type": "button", "style": "link", "height": "sm", "action": { "type": "uri", "label": "CALL", "uri": "https://linecorp.com" }
// 			},
// 			{ "type": "spacer", "size": "sm" }
// 		],
// 		"flex": 0
// 		}
// 	  }`, ChatChannel.Image, ChatChannel.Address, serviceSlot.Start, serviceSlot.End)
// 	return serviceTamplate
// }

// func PromotionsTemplate(promotions []*model.Promotion) string {
// 	var promotionCards string
// 	for _, promotion := range promotions {
// 		promotionCards = promotionCards + PromotionCardTemplate(promotion)
// 	}
// 	var promotionsTemplate string = fmt.Sprintf(`{
// 		"type": "carousel",
// 		"contents": [%s]
// 	  }`, promotionCards[:len(promotionCards)-1])
// 	return promotionsTemplate
// }

// func PromotionCardTemplate(promotion *model.Promotion) string {
// 	return fmt.Sprintf(`{
// 		"type": "bubble",
// 		"hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s" },
// 		"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
// 			{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl"
// 			},
// 			{ "type": "box", "layout": "baseline", "flex": 1, "contents": [
// 				{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 } ] } ]
// 		},
// 		"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
// 			{ "type": "button", "flex": 2, "style": "primary", "color": "#aaaaaa", "action":
// 			{ "type": "uri", "label": "Add to Cart", "uri": "https://linecorp.com" } } ]
// 		}
// 	  },`, promotion.Image, promotion.Title, promotion.Condition)
// }
