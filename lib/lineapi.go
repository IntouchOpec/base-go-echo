package lib

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ClientLine interface for pointer
type ClientLine struct {
	*linebot.Client
}

// ConnectLineBot init token connent line.
func ConnectLineBot(ChannelSecret string, ChannelAccsssToken string) (*linebot.Client, error) {
	bot, err := linebot.New(
		ChannelSecret,
		ChannelAccsssToken,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return bot, nil
}

// ReplyLineMessage send message type other.
func (client *ClientLine) ReplyLineMessage(chatAws model.ChatAnswer, replyToken string) {
	// chatAws.
	// message := ""
	switch typeReply := chatAws.TypeReply; typeReply {
	case linebot.MessageTypeText:
		textMessage := linebot.NewTextMessage("My name is John Wick")
		client.ReplyMessage(replyToken, textMessage).Do()
	case linebot.MessageTypeImage:
		// if actions == nil {
		// 	linebot.NewImagemapMessage(lineURL, textAction, linebot.ImagemapBaseSize{Width: Width, Height: Height})
		// }
		// linebot.NewImagemapMessage(lineURL, textAction, linebot.ImagemapBaseSize{Width: Width, Height: Height}, *actions...)
	case linebot.MessageTypeVideo:
		// message := FlexMessage(chatAws.Source)
	case linebot.MessageTypeAudio:
		// message := FlexMessage(chatAws.Source)
	case linebot.MessageTypeFile:
		// message := FlexMessage(chatAws.Source)
	case linebot.MessageTypeLocation:
		// message := FlexMessage(chatAws.Source)
	case linebot.MessageTypeSticker:
		// message := FlexMessage(chatAws.Source)
	case linebot.MessageTypeTemplate:
		// message := FlexMessage(chatAws.Source)
	case linebot.MessageTypeImagemap:
		// message := FlexMessage(chatAws.Source)
	case linebot.MessageTypeFlex:
		flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(chatAws.Source))
		if err != nil {
			log.Println(err)
		}
		flexMessage := linebot.NewFlexMessage("FlexWithJSON", flexContainer)
		client.ReplyMessage(replyToken, flexMessage).Do()
	default:
		// message := ""
	}

	// client.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do()

}

// FlexMessage
// func FlexMessage(JSON string) *linebot.FlexMessage {
// 	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(JSON))
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	flexMessage := linebot.NewFlexMessage("FlexWithJSON", flexContainer)
// 	return flexMessage
// }

// ImageMessage create interface for send image. ** if not action send actions = nil
// func ImageMessage(lineURL string, textAction string, Width int, Height int, actions *[]linebot.ImagemapAction) *linebot.ImagemapMessage {
// 	if actions == nil {
// 		return linebot.NewImagemapMessage(lineURL, textAction, linebot.ImagemapBaseSize{Width: Width, Height: Height})
// 	}
// 	return linebot.NewImagemapMessage(lineURL, textAction, linebot.ImagemapBaseSize{Width: Width, Height: Height}, *actions...)
// }

// ImageMessageAction action when click image.
// func ImageMessageAction(lineURL string, textAction string) *[]linebot.ImagemapAction {
// 	uriAction := linebot.URIImagemapAction{
// 		LinkURL: lineURL,
// 		Area: linebot.ImagemapArea{
// 			Height: 0,
// 			Width:  0,
// 			X:      0,
// 			Y:      0,
// 		},
// 	}
// 	messageAction := linebot.MessageImagemapAction{
// 		Text: textAction,
// 		Area: linebot.ImagemapArea{
// 			Height: 0,
// 			Width:  0,
// 			X:      0,
// 			Y:      0,
// 		},
// 	}
// 	actions := &[]linebot.ImagemapAction{
// 		&uriAction,
// 		&messageAction,
// 	}
// 	return actions
// }

// GetTokenChannelAccessToken get token lineAPI webhook Expires in 30 day
// func GetTokenChannelAccessToken() string {
// linebot.
// return
// }

func monthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}

func WeekdayHeader() string {
	weekdays := []string{"อา", "จ", "อ", "พ", "พฤ", "ศ", "ส"}
	var weekdaysStr string
	for weekday := 0; weekday <= len(weekdays); weekday++ {
		weekdaysStr = weekdaysStr + fmt.Sprintf(`{
        "type": "text",
        "text": %s,
        "size": "sm",
        "color": "#000000",
        "align": "center"
      },`, weekdays[weekday])
	}
	return weekdaysStr[:len(weekdaysStr)-1]
}

func MonthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}

//
func MakeCalenda() string {
	var contents string
	var calendar string

	t := time.Now()

	year, month, _ := time.Now().Date()

	var color string = "#000000"

	endOfMonth := time.Date(year, month+1, 1, 0, 0, 0, -1, time.UTC)
	var Weekday int = int(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday())

	for day := 0; day < Weekday; day++ {
		fmt.Println(day)
		contents = contents + fmt.Sprintf(`{
			"type":    "text",
			"text":    " ",
			"size":    "sm",
			"color":   "#000000",
			"align":   "center",
			"gravity": "center"},`)
	}
	Weekday = int(t.Weekday())

	for day := 1; day <= endOfMonth.Day(); day++ {
		if day == t.Day() {
			color = "#1db446"
		} else {
			color = "#000000"
		}

		contents = contents + fmt.Sprintf(`{
			"type":    "text",
			"text":    "%s",
			"size":    "sm",
			"color":   "%s",
			"align":   "center",
			"gravity": "center",
			"action": { "type": "message",
            "label": "%s",
            "text": "%s"}},`, strconv.FormatInt(int64(day), 10), color, day, day)
		contents = contents + `{"type": "separator"},`
		Weekday = int(time.Date(year, month, day, 0, 0, 0, -1, time.UTC).Weekday())
		if endOfMonth.Day() == day {
			for dw := int(endOfMonth.Weekday()); dw < 6; dw++ {
				contents = contents + fmt.Sprintf(`{
          "type":    "text",
          "text":    " ",
          "size":    "sm",
          "color":   "#000000",
          "align":   "center",
          "gravity": "center"},`)
			}
			// contents = contents + `{"type": "text", "text": " ", "size": "sm", "color": "#000000", "align": "center"},`
		}

		// 6 == saturday
		if (int(Weekday) == 5 && day != 1) || endOfMonth.Day() == day {

			calendar = calendar + fmt.Sprintf(`{
				"type":     "box",
				"layout":   "horizontal",
				"margin":   "md",
				"contents": [%s]
			},`, contents[:len(contents)-1])
			contents = ""
		}
	}
	weekdays := []string{"อา", "จ", "อ", "พ", "พฤ", "ศ", "ส"}
	var weekdaysStr string
	for weekday := 0; weekday < len(weekdays); weekday++ {

		weekdaysStr = weekdaysStr + fmt.Sprintf(`{
        "type": "text",
        "text": "%s",
        "size": "sm",
        "color": "#000000",
        "align": "center"
      },`, weekdays[weekday])
	}
	weekdaysStr = fmt.Sprintf(`{
                "type":     "box",
                "layout":   "horizontal",
                "margin":   "md",
                "contents": [%s]},`, weekdaysStr[:len(weekdaysStr)-1])

	return weekdaysStr + `{"type": "separator"},` + calendar[:len(calendar)-1]
}
