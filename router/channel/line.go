package channel

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
	fmt.Println("=================++++++++=======")
	name := c.Param("account")
	ChannelID := c.Param("ChannelID")
	account := model.Account{}
	ChatChannel := model.ChatChannel{}

	if err := model.DB().Where("name = ?", name).Find(&account).Error; err != nil {
		fmt.Println("tet2")
		return c.NoContent(http.StatusNotFound)
	}

	if err := model.DB().Where("channel_id = ?", ChannelID).Find(&ChatChannel).Error; err != nil {
		fmt.Println("tet3")
		return c.NoContent(http.StatusNotFound)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	bot, err := lib.ConnectLineBot(ChatChannel.ChannelSecret, ChatChannel.ChannelAccessToken)

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
		switch eventType := event.Type; eventType {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "calendar" {
					fmt.Println("+++++++++")
					// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("content")).Do(); err != nil {
					// 	log.Print(err)
					// }
					m := lib.MakeCalenda()
					m = fmt.Sprintf(`{"type": "bubble",
					"styles": {
					  "footer": {
						"separator": true
					  }
					},
					"body": {
					  "type": "box",
					  "layout": "vertical",
					  "contents": [%s]}}`, m)

					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
					if err != nil {
						log.Println(err)
					}
					flexMessage := linebot.NewFlexMessage("ตาราง", flexContainer)
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
