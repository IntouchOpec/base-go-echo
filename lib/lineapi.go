package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	return time.Format("20060102"), nil
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
func (client *ClientLine) ReplyLineMessage(chatAws model.ChatAnswer) (linebot.SendingMessage, error) {
	var message linebot.SendingMessage
	switch replyType := chatAws.AnsReplyType; replyType {
	case linebot.MessageTypeText:
		message = linebot.NewTextMessage(chatAws.AnsReply)
	case linebot.MessageTypeImage:
		message = linebot.NewTextMessage(chatAws.AnsReply)
	case linebot.MessageTypeVideo:
		videoMessage := linebot.VideoMessage{}
		if err := json.Unmarshal([]byte(chatAws.AnsSource), videoMessage); err != nil {
			return nil, err
		}
		message = &videoMessage
	case linebot.MessageTypeAudio:
		audieoMessage := linebot.AudioMessage{}
		if err := json.Unmarshal([]byte(chatAws.AnsSource), audieoMessage); err != nil {
			return nil, err
		}
		message = &audieoMessage
	case linebot.MessageTypeFile:
	case linebot.MessageTypeLocation:
		locationMessage := linebot.LocationMessage{}
		if err := json.Unmarshal([]byte(chatAws.AnsSource), locationMessage); err != nil {
			return nil, err
		}
		message = &locationMessage
	case linebot.MessageTypeSticker:
	case linebot.MessageTypeTemplate:
	case linebot.MessageTypeImagemap:
	case linebot.MessageTypeFlex:
		flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(chatAws.AnsSource))
		if err != nil {
			return nil, err
		}
		message = linebot.NewFlexMessage(chatAws.AnsInput, flexContainer)
	default:

	}
	return message, nil
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
