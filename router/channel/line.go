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
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/labstack/echo"
)

// Routers channel.
// func Routers() *echo.Echo {
// 	e := echo.New()

// 	e.POST("/callback/:account", Callback)
// 	return e
// }

// HandleWebHookLineAPI webhook for connent line api.
func HandleWebHookLineAPI(c echo.Context) error {
	client := &lib.ClientLine{}
	name := c.Param("account")
	ChannelID := c.Param("ChannelID")
	account := model.Account{}
	chatChannel := model.ChatChannel{}

	if err := model.DB().Where("name = ?", name).Find(&account).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := model.DB().Where("channel_id = ?", ChannelID).Find(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)

	if err != nil {
		log.Print(err)
	}

	events, err := bot.ParseRequest(c.Request())

	if err != nil {

		if err == linebot.ErrInvalidSignature {
			c.String(400, linebot.ErrInvalidSignature.Error())
		} else {
			c.String(500, "internal")
		}
	}

	for _, event := range events {
		ChatAnswer := model.ChatAnswer{}
		var keyWord string

		switch eventType := event.Type; eventType {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				messageText := message.Text
				if len(messageText) > 8 {
					keyWord = messageText[0:8]
				}
				if keyWord == "calendar" || messageText == "Booking" {
					fmt.Println(keyWord)
					var m string
					if len(message.Text) > 8 {
						m = lib.MakeCalenda("")
					} else {
						m = lib.MakeCalenda("")
					}
					fmt.Println(m)
					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {
						log.Println(err)
					}
					flexMessage := linebot.NewFlexMessage("ตาราง", flexContainer)
					bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
				} else if keyWord == "product " {
					t, _ := time.Parse("2006-01-02 15:04", messageText[8:]+" 04:35")

					subProduct := []model.SubProduct{}
					product := model.Product{}

					if err := model.DB().Find(&product).Where("Day = ?", int(t.Weekday())).Related(&subProduct).Error; err != nil {
						fmt.Println(err)
						return nil
					}

					m := ProductListLineTemplate(product, subProduct, messageText[8:])
					model.DB().Preload("Orders").Where("Day = ?", int(t.Weekday())).Preload("Product").Find(&subProduct)

					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {
						log.Println(err)
					}
					flexMessage := linebot.NewFlexMessage("ตาราง", flexContainer)

					bot.ReplyMessage(event.ReplyToken, flexMessage).Do()
				} else if keyWord == "booking " {
					fmt.Println(messageText[8:18], messageText[19:24])
					t, _ := time.Parse("2006-01-02 15:04", messageText[8:])

					subProduct := model.SubProduct{}
					if err := model.DB().Where("start = ?", messageText[19:24]).Find(&subProduct).Error; err != nil {
						fmt.Println(err)
						return nil
					}

					booking := model.Booking{SubProductID: subProduct.ID, BookingDate: t, AccountID: chatChannel.AccountID}
					booking.SaveBooking()
					m := ThankyouTemplate(chatChannel, booking, subProduct)
					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {
						return err
					}
					flexMessage := linebot.NewFlexMessage("ขอบคุณ", flexContainer)

					bot.ReplyMessage(event.ReplyToken, flexMessage).Do()

				}
				// model.DB().Where("Input = ?", message.Text).Find(&ChatAnswer)
				// if err := model.DB().Where("Input = ?", message.Text).Find(&ChatAnswer).Error; err != nil {
				// 	fmt.Println(err)
				// }
				// client.ReplyLineMessage(ChatAnswer, event.ReplyToken)
			case *linebot.ImageMessage:

			case *linebot.VideoMessage:

			case *linebot.AudioMessage:

			case *linebot.FileMessage:

			case *linebot.LocationMessage:

			case *linebot.StickerMessage:

			case *linebot.TemplateMessage:

			case *linebot.ImagemapMessage:

			case *linebot.FlexMessage:

			}

		case linebot.EventTypeFollow:
			if err := model.DB().Where("Input = followAuto").Find(&ChatAnswer).Error; err != nil {
				fmt.Println(err)
			}
			client.ReplyLineMessage(ChatAnswer, event.ReplyToken)
			// var contents []linebot.FlexComponent
			// text := linebot.TextComponent{
			// 	Type:   linebot.FlexComponentTypeText,
			// 	Text:   "taey line bookking plaform",
			// 	Weight: "bold",
			// 	Size:   linebot.FlexTextSizeTypeXl,
			// 	Action: linebot.NewURIAction("register", "https://15e330d8.ngrok.io/register"),
			// }
			// contents = append(contents, &text)
			// // Make Hero
			// hero := linebot.ImageComponent{
			// 	Type:        linebot.FlexComponentTypeImage,
			// 	URL:         "https://scontent.fbkk2-7.fna.fbcdn.net/v/t1.0-9/55771768_3311003885591805_86103752003551232_o.jpg?_nc_cat=109&_nc_eui2=AeGPFqTgk7ynFe18QHmR-69H6MogRu5OFJXtXwbMnKDQa2IZeLa57IEayXcXzhyzKDfBKx_tYZevLlEoaJ_bJn6Fl9hCv6mhlWYOOV3ltGoR9Q&_nc_oc=AQkpFLS6szBuMWyOhKz-Ope9I4YkWTFea1DFHE9oNPodtflCUt53bb_kjVd7SVx236w&_nc_ht=scontent.fbkk2-7.fna&oh=62d415b199aaa244c8bea5b9e60dd44b&oe=5DD5122F",
			// 	Size:        "full",
			// 	AspectRatio: linebot.FlexImageAspectRatioType1to1,
			// 	AspectMode:  linebot.FlexImageAspectModeTypeCover,
			// 	Action:      linebot.NewURIAction("register", "https://15e330d8.ngrok.io/register"),
			// }
			// // Make Body
			// body := linebot.BoxComponent{
			// 	Type:     linebot.FlexComponentTypeBox,
			// 	Layout:   linebot.FlexBoxLayoutTypeVertical,
			// 	Contents: contents,
			// }
			// // Build Container
			// bubble := linebot.BubbleContainer{
			// 	Type: linebot.FlexContainerTypeBubble,
			// 	Hero: &hero,
			// 	Body: &body,
			// }
			// // New Flex Message
			// flexMessage := linebot.NewFlexMessage("ขอบคุณที่มาเป็นเพื่อนกันนะ", &bubble)
			// if _, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
			// 	log.Print(err)
			// }
		}
	}

	return c.JSON(200, "")

}

// ProductListLineTemplate
func ProductListLineTemplate(product model.Product, subProduct []model.SubProduct, dateTime string) string {
	var slotTime string
	var buttonTime string

	for t := 0; t < len(subProduct); t++ {
		if t%2 == 0 {
			if t != 0 {
				slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
				buttonTime = ""
			}
		}

		buttonTime = buttonTime + fmt.Sprintf(`{
			"type": "button",
			"style": "primary",
			"action": {
			  "type": "message",
			  "label": "%s - %s",
			  "text": "%s"
			}},`, subProduct[t].Start, subProduct[t].End, "booking"+" "+dateTime+" "+subProduct[t].Start+"-"+subProduct[t].End)

		if t == len(subProduct)-1 {
			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
		}
	}
	var productTamplate string = fmt.Sprintf(`{
		"type": "carousel",
		"contents": [
		{"type": "bubble",
		"hero": {
		  "type": "image",
		  "size": "full",
		  "aspectRatio": "20:13",
		  "aspectMode": "cover",
		  "url": "%s"
		},
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "spacing": "sm",
		  "contents": [
			{
			  "type": "text",
			  "text": "%s",
			  "wrap": true,
			  "weight": "bold",
			  "size": "xl"
			},
			{
			  "type": "box",
			  "layout": "baseline",
			  "contents": [
				{
				  "type": "text",
				  "text": "฿%s",
				  "wrap": true,
				  "weight": "bold",
				  "size": "xl",
				  "flex": 0
				}
			]
			}
			%s
		  ]
		}}]
	  }`, product.Image, product.Name, strconv.FormatInt(int64(product.Price), 10), slotTime)

	return productTamplate
}

// ThankyouTemplate
func ThankyouTemplate(ChatChannel model.ChatChannel, booking model.Booking, subProduct model.SubProduct) string {
	var productTamplate string = fmt.Sprintf(`{
		"type": "bubble",
		"hero": {
		  "type": "image",
		  "url": "%s",
		  "size": "full",
		  "aspectRatio": "20:13",
		  "aspectMode": "cover"
		},
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "text",
			  "text": "จองสำเร็จ",
			  "weight": "bold",
			  "size": "xl"
			},
			{
			  "type": "box",
			  "layout": "vertical",
			  "margin": "lg",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "box",
				  "layout": "baseline",
				  "spacing": "sm",
				  "contents": [
					{
					  "type": "text",
					  "text": "Place",
					  "color": "#aaaaaa",
					  "size": "sm",
					  "flex": 1
					},
					{
					  "type": "text",
					  "text": "%s",
					  "wrap": true,
					  "color": "#666666",
					  "size": "sm",
					  "flex": 5
					}
				  ]
				},
				{
				  "type": "box",
				  "layout": "baseline",
				  "spacing": "sm",
				  "contents": [
					{
					  "type": "text",
					  "text": "Time",
					  "color": "#aaaaaa",
					  "size": "sm",
					  "flex": 1
					},
					{
					  "type": "text",
					  "text": "%s - %s",
					  "wrap": true,
					  "color": "#666666",
					  "size": "sm",
					  "flex": 5
					}
				  ]
				}
			  ]
			}
		  ]
		},
		"footer": {
		  "type": "box",
		  "layout": "vertical",
		  "spacing": "sm",
		  "contents": [
			{
			  "type": "button",
			  "style": "link",
			  "height": "sm",
			  "action": {
				"type": "uri",
				"label": "CALL",
				"uri": "https://linecorp.com"
			  }
			},
			{
			  "type": "spacer",
			  "size": "sm"
			}
		  ],
		  "flex": 0
		}
	  }`, ChatChannel.Image, ChatChannel.Address, subProduct.Start, subProduct.End)
	return productTamplate
}
