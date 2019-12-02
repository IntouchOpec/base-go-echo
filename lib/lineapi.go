package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ClientLine interface for pointer
type ClientLine struct {
	linebot.Client
}

const (
	APIEndpointBase = "https://api.line.me"

	APIEndpointInsightFollowers = "/v2/bot/insight/followers?date=%s"
)

// ConnectLineBot init token connent line.
func ConnectLineBot(ChannelSecret string, ChannelAccsssToken string) (*ClientLine, error) {
	bot, err := linebot.New(
		ChannelSecret,
		ChannelAccsssToken,
	)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	line := ClientLine{*bot}
	return &line, nil
}

type ResponseInsightFollowrs struct {
	Status          string
	Followers       int
	TargetedReaches int
	Blocks          int
}

func DateStringFormantLine(time time.Time) (string, error) {
	return fmt.Sprintf("%d%d%d", time.Year(), time.Month(), time.Day()-1), nil
}

func InsightFollowers(channelAccsssToken string) (*ResponseInsightFollowrs, error) {
	Authorization := fmt.Sprintf("Bearer %s", channelAccsssToken)
	timeNow := time.Now()
	timeFormat, _ := DateStringFormantLine(timeNow)
	url := fmt.Sprintf("%s%s", APIEndpointBase, fmt.Sprintf(APIEndpointInsightFollowers, timeFormat))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", Authorization)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	result := ResponseInsightFollowrs{}

	if err := decoder.Decode(&result); err != nil {
		if err == io.EOF {
			return &result, nil
		}
		return nil, err
	}

	return &result, nil
}

// ReplyLineMessage send message type other.
func (client *ClientLine) ReplyLineMessage(chatAws model.ChatAnswer, replyToken string) {
	// chatAws.
	// message := ""
	switch replyType := chatAws.AnsReplyType; replyType {
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
		flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(chatAws.AnsSource))
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
func MakeCalenda(date string) string {
	var contents string
	var calendar string
	var color string
	year, month, _ := time.Now().Date()
	t := time.Now()
	if len(date) != 0 {
		time2, _ := time.Parse("2006-01-02", date)
		year, month, _ = time2.Date()
	}

	color = "#000000"

	endOfMonth := time.Date(year, month+1, 1, 0, 0, 0, -1, time.UTC)

	Weekday := int(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday())

	for day := 0; day < Weekday; day++ {
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
		if day == t.Day() && month == t.Month() {
			color = "#1db446"
		} else {
			color = "#000000"
		}
		dayStr := strconv.FormatInt(int64(day), 10)
		monthStr := strconv.FormatInt(int64(month), 10)
		if len(dayStr) == 1 {
			dayStr = fmt.Sprintf("0%s", dayStr)
		}
		if len(monthStr) == 1 {
			monthStr = fmt.Sprintf("0%s", monthStr)
		}
		contents = contents + fmt.Sprintf(`{ "type": "text", "text": "%s", "size": "sm", "color": "%s", "align": "center", "gravity": "center",
					"action": { "type": "message", "label": "%s", "text": "Service%d-%s-%s"}},`, dayStr, color, day, year, monthStr, dayStr)
		contents = contents + `{"type": "separator"},`
		Weekday = int(time.Date(year, month, day, 0, 0, 0, -1, time.UTC).Weekday())
		if endOfMonth.Day() == day {
			for dw := int(endOfMonth.Weekday()); dw < 6; dw++ {
				contents = contents + fmt.Sprintf(`{ "type": "text", "text": " ", "size": "sm", "color": "#000000", "align": "center", "gravity": "center"},`)
			}
		}

		// 6 == saturday
		if (int(Weekday) == 5) || endOfMonth.Day() == day {

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

		weekdaysStr = weekdaysStr + fmt.Sprintf(`{ "type": "text", "text": "%s", "size": "sm", "color": "#000000", "align": "center" },`, weekdays[weekday])
	}
	weekdaysStr = fmt.Sprintf(`{"type": "box","layout": "horizontal","margin": "md","contents": [%s]},`, weekdaysStr[:len(weekdaysStr)-1])
	HeaderCalendat := fmt.Sprintf("%s %s", month, strconv.FormatInt(int64(year), 10))
	var nextMonth string = strconv.FormatInt(int64(month+1), 10)
	var nextYear int = year
	if nextMonth == "13" {
		nextMonth = "01"
		nextYear = year + 1
	}

	if len(nextMonth) == 1 {
		nextMonth = "0" + nextMonth
	}
	actionNextMonth := fmt.Sprintf("%d-%s-01", nextYear, nextMonth)
	m := fmt.Sprintf(`{"type": "bubble","styles": {"footer": {"separator": true}},
	"body": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "%s", "size": "sm", "weight": "bold", "color": "#1db446", "flex": 0 },
				{ "type": "text", "text": "ถัดไป", "size": "sm", "color": "#111111", "align": "end", "action": { "type": "message", "label": " ", "text": "calendar %s"} }
			]
		}, %s]}}`, HeaderCalendat, actionNextMonth, weekdaysStr+`{"type": "separator"},`+calendar[:len(calendar)-1])
	return m
}

// func

// CreateLIIF request url and size
// func (client *ClientLine) CreateLIIF(view linebot.View) {
// res, err := client.AddLIFF(view).Do()
// if err != nil {
// return nil, err
// }
// return res, nil
// }
