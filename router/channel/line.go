package channel

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/hb-go/gorm"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/labstack/echo"
)

// HandleWebHookLineAPI webhook for connent line api.
func HandleWebHookLineAPI(c echo.Context) error {
	client := &lib.ClientLine{}
	name := c.Param("account")
	ChannelID := c.Param("ChannelID")
	account := model.Account{}
	chatChannel := model.ChatChannel{}

	if err := model.DB().Where("name = ?", name).Find(&account).Error; err != nil {
		fmt.Println(err, 1)
		return c.NoContent(http.StatusNotFound)
	}

	if err := model.DB().Where("Channel_ID = ?", ChannelID).Find(&chatChannel).Error; err != nil {
		fmt.Println(err, 2)

		return c.NoContent(http.StatusNotFound)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)

	if err != nil {
		fmt.Println(err, 3)
		log.Print(err)
	}

	events, err := bot.ParseRequest(c.Request())

	if err != nil {
		fmt.Println(err, 4)

		if err == linebot.ErrInvalidSignature {
			c.String(400, linebot.ErrInvalidSignature.Error())
		} else {
			c.String(500, "internal")
		}
	}
	customer := model.Customer{}
	for _, event := range events {
		chatAnswer := model.ChatAnswer{}
		var keyWord string
		model.DB().Where("Line_id = ? and chat_channel_id = ?", event.Source.UserID, chatChannel.ID).Find(&customer)
		switch eventType := event.Type; eventType {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				messageText := message.Text
				if len(messageText) >= 8 {
					keyWord = messageText[0:8]
				}
				fmt.Println(keyWord)
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
					res, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
					fmt.Println(res)
					if err != nil {
						act := model.ActionLog{Name: "calendar", Status: model.StatusFail, Type: model.TypeActionLine, UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
						act.CreateAction()
						return err
					}
					act := model.ActionLog{Name: "calendar", Status: model.StatusSuccess, Type: model.TypeActionLine, UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
					act.CreateAction()
				} else if keyWord == "product " {
					t, _ := time.Parse("2006-01-02 15:04", messageText[8:]+" 01:00")

					subProduct := []model.SubProduct{}
					if err := model.DB().Preload("Bookings", "booked_date = ?", t).Preload("Product").Where("Day = ?", int(t.Weekday())).Find(&subProduct).Error; err != nil {
						fmt.Println(err)
						return nil
					}

					m := ProductListLineTemplate(subProduct, messageText[8:])

					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {
						log.Println(err)
					}
					flexMessage := linebot.NewFlexMessage("ตาราง", flexContainer)

					res, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
					if err != nil {
						return err
					}
					fmt.Println(res, "=====__")
					act := model.ActionLog{Name: "product", Status: model.StatusSuccess, Type: model.TypeActionLine, UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
					act.CreateAction()
				} else if keyWord == "booking " {
					t, _ := time.Parse("2006-01-02 15:04", messageText[8:18]+" 01:00")
					product := model.Product{}
					if err := model.DB().Preload("SubProducts", "start = ? and day = ?", messageText[19:24], int(t.Weekday())).Where("Name = ?", messageText[31:]).Find(&product).Error; err != nil {
						fmt.Println(err, "err")
						return nil
					}

					custo := model.Customer{}
					if err := model.DB().Where("line_id = ? and chat_channel_id = ?", event.Source.UserID, chatChannel.ID).Find(&custo).Error; err != nil {
						fmt.Println(err, "err")
						return nil
					}
					subProduct := product.SubProducts[0]
					booking := model.Booking{SubProductID: product.SubProducts[0].ID, BookedDate: t, ChatChannelID: chatChannel.ID, CustomerID: custo.ID}
					_, err := booking.SaveBooking()
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
						customer := model.Customer{}
						model.DB().Preload("Promotions", "type = ?", "voucher", func(db *gorm.DB) *gorm.DB {
							return db.Preload("Settings")
						}).Where("line_id = ?", message.ID).Find(&customer)
						var template string
						for index := 0; index < len(customer.Promotions); index++ {
							// promotion := customer.Promotions[index]
							// template = template + web.VoucherTemplate(customer.Promotions[index]) + ","
						}
						template = fmt.Sprintf(`{ "type": "carousel", "contents": [%s]}`, template[:len(template)-1])
						flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
						if err != nil {
							return err
						}
						flexMessage := linebot.NewFlexMessage("your voucher", flexContainer)
						res, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
						if err != nil {
							act := model.ActionLog{Name: "user voucher", Status: model.StatusFail, Type: model.TypeActionLine,
								UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
							act.CreateAction()
							return err
						}
						fmt.Println(res)
						act := model.ActionLog{Name: "user voucher", Status: model.StatusSuccess,
							Type: model.TypeActionLine, UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
						act.CreateAction()
					} else if keyWord == "comment" {

					} else {
						m = ThankyouTemplate(chatChannel, booking, subProduct)
					}
					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {
						return err
					}
					flexMessage := linebot.NewFlexMessage("ขอบคุณ", flexContainer)
					// flexMessage := linebot.NewPostbackAction("label", "data", "text", "displayText")
					bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
				} else if keyWord == "promotio" {
					promotions := []*model.Promotion{}
					model.DB().Where("type_promotion = ?", "Promotion").Find(&promotions)
					fmt.Println(len(promotions))
					m := PromotionsTemplate(promotions)
					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {

						return err
					}
					flexMessage := linebot.NewFlexMessage("ตาราง", flexContainer)
					res, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
					fmt.Println(res)
					if err != nil {
						return err
					}
					act := model.ActionLog{Name: "promotion", Status: model.StatusSuccess, Type: model.TypeActionLine, UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
					act.CreateAction()
				} else if keyWord == "location" {
					position := chatChannel.GetSetting([]string{"Latitude", "Longitude"})
					Latitude, _ := strconv.ParseFloat(position["Latitude"], 64)
					Longitude, _ := strconv.ParseFloat(position["Longitude"], 64)

					location := linebot.NewLocationMessage(chatChannel.Name, chatChannel.Address, Latitude, Longitude)

					res, err := bot.ReplyMessage(event.ReplyToken, location).Do()
					if err != nil {
						return err
					}
					fmt.Println(res)
					act := model.ActionLog{Name: "location", Status: model.StatusSuccess, Type: model.TypeActionLine, UserID: event.Source.UserID,
						ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
					act.CreateAction()
				} else {
					model.DB().Where("Input = ?", message.Text).Find(&chatAnswer)
					if err := model.DB().Where("Input = ?", message.Text).Find(&chatAnswer).Error; err != nil {
						fmt.Println(err)
					}
					client.ReplyLineMessage(chatAnswer, event.ReplyToken)
				}
				evenLog := model.EventLog{ChatChannelID: chatChannel.ID, ReplyToken: event.ReplyToken, Type: string(event.Type),
					LineID: event.Source.UserID, Text: messageText, CustomerID: customer.ID}
				model.DB().Create(&evenLog)
			case *linebot.ImageMessage:

			case *linebot.VideoMessage:

			case *linebot.AudioMessage:

			case *linebot.FileMessage:

			case *linebot.LocationMessage:
				textMessage := linebot.NewTextMessage(fmt.Sprintf("%v", message))
				res, err := bot.ReplyMessage(event.ReplyToken, textMessage).Do()
				if err != nil {
					act := model.ActionLog{Name: "LocationMessage", Status: model.StatusFail, Type: model.TypeActionLine, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
					act.CreateAction()
					return err
				}
				fmt.Println(res)
				act := model.ActionLog{Name: "LocationMessage", Status: model.StatusSuccess, Type: model.TypeActionLine, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
				act.CreateAction()
			case *linebot.StickerMessage:

			case *linebot.TemplateMessage:

			case *linebot.ImagemapMessage:

			case *linebot.FlexMessage:

			}

		case linebot.EventTypeFollow:
			customer := model.Customer{LineID: event.Source.UserID, ChatChannelID: chatChannel.ID}
			settingNames := []string{"LIFFregister"}
			setting := chatChannel.GetSetting(settingNames)
			if err := model.DB().FirstOrCreate(&customer, model.Customer{LineID: event.Source.UserID, ChatChannelID: chatChannel.ID}).Error; err != nil {
				fmt.Println(err)
			}
			jsonFlexMessage := FollowTemplate(chatChannel, setting)
			flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(jsonFlexMessage))
			if err != nil {
				log.Print(err)
				return err
			}
			flexMessage := linebot.NewFlexMessage(chatChannel.WelcomeMessage, flexContainer)
			bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
			evenLog := model.EventLog{ChatChannelID: chatChannel.ID, ReplyToken: event.ReplyToken, Type: string(event.Type), LineID: event.Source.UserID, CustomerID: customer.ID}
			model.DB().Create(&evenLog)
			act := model.ActionLog{Name: "follow", Status: model.StatusSuccess, Type: model.TypeActionLine, UserID: event.Source.UserID, ChatChannelID: chatChannel.ID, CustomerID: customer.ID}
			if err := model.DB().Create(&act).Error; err != nil {
				fmt.Println(err, "====")
			}
		case linebot.EventTypeUnfollow:
			evenLog := model.EventLog{ChatChannelID: chatChannel.ID, ReplyToken: event.ReplyToken, Type: string(event.Type), LineID: event.Source.UserID, CustomerID: customer.ID}
			model.DB().Create(&evenLog)
		}

	}

	return c.JSON(200, "")

}

// ProductListLineTemplate
func ProductListLineTemplate(subProduct []model.SubProduct, dateTime string) string {
	var slotTime string
	var buttonTime string
	var productList string
	var count int
	count = 0
	for t := 0; t < len(subProduct); t++ {
		if count == 2 {
			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
			buttonTime = ""
			count = 0
		}

		if len(subProduct[t].Bookings) > 0 {
			fmt.Println(subProduct[t].Bookings[0].Queue < subProduct[t].Amount)
			if subProduct[t].Bookings[0].Queue < subProduct[t].Amount {
				buttonTime = buttonTime + fmt.Sprintf(`{"type": "button", "style": "secondary", "margin": "sm", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`,
					subProduct[t].Start, subProduct[t].End, "เต็มแล้ว")
			} else {
				buttonTime = buttonTime + fmt.Sprintf(`{"type": "button","style": "primary", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`,
					subProduct[t].Start, subProduct[t].End, "booking"+" "+dateTime+" "+subProduct[t].Start+"-"+subProduct[t].End+" "+subProduct[t].Product.Name)
			}
		} else {
			buttonTime = buttonTime + fmt.Sprintf(`{"type": "button","style": "primary", "margin": "sm", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`,
				subProduct[t].Start, subProduct[t].End, "booking"+" "+dateTime+" "+subProduct[t].Start+"-"+subProduct[t].End+" "+subProduct[t].Product.Name)
		}

		count = count + 1
		if t == len(subProduct)-1 {
			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
			productList += fmt.Sprintf(`{"type": "bubble", "hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s"},
			"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
				{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl" },
				{ "type": "box", "layout": "baseline", "contents": [
					{ "type": "text", "text": "฿%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
				] }
				%s]
			}}`, subProduct[t].Product.Image, subProduct[t].Product.Name, strconv.FormatInt(int64(subProduct[t].Product.Price), 10), slotTime)
		} else if subProduct[t].Product.ID != subProduct[t+1].Product.ID {
			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
			productList = productList + fmt.Sprintf(`{"type": "bubble", "hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s"},
			"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
				{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl" },
				{ "type": "box", "layout": "baseline", "contents": [
					{ "type": "text", "text": "฿%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
				] }
				%s
			]
			}},`, subProduct[t].Product.Image, subProduct[t].Product.Name, strconv.FormatInt(int64(subProduct[t].Product.Price), 10), slotTime)
			slotTime = ""
			count = 0
			buttonTime = ""
		}

	}
	var nextPage string = `{ "type": "bubble", "body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
		{ "type": "button", "flex": 1, "gravity": "center", "action": { "type": "uri", "label": "See more", "uri": "https://linecorp.com" } }] }}`
	var productTamplate string = fmt.Sprintf(`{ "type": "carousel", "contents": [%s, %s]}`, productList, nextPage)
	return productTamplate
}

// FollowTemplate
func FollowTemplate(chatChannel model.ChatChannel, settings map[string]string) string {
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
	  }`, chatChannel.Image, chatChannel.Name, chatChannel.WelcomeMessage, settings["LIFFregister"], chatChannel.WebSite)
	return template
}

// ThankyouTemplate
func ThankyouTemplate(ChatChannel model.ChatChannel, booking model.Booking, subProduct *model.SubProduct) string {
	var productTamplate string = fmt.Sprintf(`{
		"type": "bubble",
		"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{ "type": "text", "text": "จองสำเร็จ", "weight": "bold", "size": "xl" },
			{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
				{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
					{ "type": "text", "text": "Place", "color": "#aaaaaa", "size": "sm", "flex": 1 },
					{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
				] },
				{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
					{ "type": "text", "text": "Time", "color": "#aaaaaa", "size": "sm", "flex": 1 },
					{ "type": "text", "text": "%s - %s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
				] }
			  ] }
		  	]
		},
		"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
			{ "type": "button", "style": "link", "height": "sm", "action": { "type": "uri", "label": "CALL", "uri": "https://linecorp.com" }
			},
			{ "type": "spacer", "size": "sm" }
		],
		"flex": 0
		}
	  }`, ChatChannel.Image, ChatChannel.Address, subProduct.Start, subProduct.End)
	return productTamplate
}

func PromotionsTemplate(promotions []*model.Promotion) string {
	var promotionCards string
	for _, promotion := range promotions {
		promotionCards = promotionCards + PromotionCardTemplate(promotion)
	}
	var promotionsTemplate string = fmt.Sprintf(`{
		"type": "carousel",
		"contents": [%s]
	  }`, promotionCards[:len(promotionCards)-1])
	return promotionsTemplate
}

func PromotionCardTemplate(promotion *model.Promotion) string {
	return fmt.Sprintf(`{
		"type": "bubble",
		"hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s" },
		"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
			{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl"
			},
			{ "type": "box", "layout": "baseline", "flex": 1, "contents": [
				{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 } ] } ]
		},
		"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
			{ "type": "button", "flex": 2, "style": "primary", "color": "#aaaaaa", "action": 
			{ "type": "uri", "label": "Add to Cart", "uri": "https://linecorp.com" } } ]
		}
	  },`, promotion.Image, promotion.Title, promotion.Condition)
}
