package channel

import (
	"fmt"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

func welcomeHandle(c *echo.Context, event *linebot.Event, chatChannel *model.ChatChannel) linebot.SendingMessage {
	customer := model.Customer{}

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
	FlexMessage := linebot.NewFlexMessage(chatChannel.ChaWelcomeMessage, flexContainer)
	return FlexMessage
}

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
