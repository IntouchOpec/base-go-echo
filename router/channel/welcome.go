package channel

import (
	"fmt"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

func WelcomeHandle(c *Context) (linebot.SendingMessage, error) {
	customer := model.Customer{}
	// setting := model.Setting{Name: model.NameLIFFregister}
	setting := c.ChatChannel.GetSetting([]string{model.NameLIFFregister})
	// if err := c.DB.Model(&c.ChatChannel).Association("Settings").Find(&setting, "name = ?", model.NameLIFFregister).Error; err != nil {

	// }
	fmt.Println(setting, c.ChatChannel.Settings, model.NameLIFFregister)
	if err := model.DB().FirstOrCreate(&customer, model.Customer{
		CusLineID: c.Event.Source.UserID,
		AccountID: c.ChatChannel.AccountID}).Error; err != nil {
		return nil, err
	}
	jsonFlexMessage := fmt.Sprintf(FollowTemplate, c.ChatChannel.ChaImage, c.ChatChannel.ChaName, c.ChatChannel.ChaWelcomeMessage, setting["LIFFregister"], c.ChatChannel.ChaWebSite)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(jsonFlexMessage))
	if err != nil {
		return nil, err
	}
	fmt.Println("jsonFlexMessage", jsonFlexMessage)
	FlexMessage := linebot.NewFlexMessage(c.ChatChannel.ChaWelcomeMessage, flexContainer)
	return FlexMessage, nil
}

var FollowTemplate string = `{
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
  }`
