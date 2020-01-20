package channel

import (
	"fmt"

	"github.com/IntouchOpec/base-go-echo/model"

	"github.com/line/line-bot-sdk-go/linebot"
)

var cardTemplate string = `{
	"type": "bubble",
	"header": { "type": "box", "layout": "horizontal", "contents": [
		{ "type": "box", "layout": "vertical",
		  "contents": [
			{ "type": "text", "text": "VOUCHERS", "size": "sm", "weight": "bold", "color": "#AAAAAA" },
			{ "type": "box", "layout": "horizontal", "flex": 1, "contents": [
				{ "type": "image", "url": "https://scdn.line-apps.com/n/channel_devcenter/img/flexsnapshot/clip/clip7.jpg", "size": "5xl", "aspectMode": "cover", "aspectRatio": "150:196", "gravity": "center", "flex": 1 },
				{ "type": "box", "layout": "vertical", "contents": [
					{ "type": "image", "url": "https://scdn.line-apps.com/n/channel_devcenter/img/flexsnapshot/clip/clip8.jpg", "size": "full", "aspectMode": "cover", "aspectRatio": "150:98", "gravity": "center" },
					{ "type": "image", "url": "https://scdn.line-apps.com/n/channel_devcenter/img/flexsnapshot/clip/clip9.jpg", "size": "full", "aspectMode": "cover", "aspectRatio": "150:98", "gravity": "center" }]}]}]}]},
	"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
		{ "type": "button", "style": "primary", "action": { "type": "datetimepicker", "label": "ทำนัดอัตโนมัติ", "data": "choive auto", "mode": "datetime", "initial": "2020-01-17T14:48", "max": "2021-01-17T14:48", "min": "2019-01-17T14:48" },	  },
		{ "type": "button", "style": "primary", "action": { "type": "postback", "label": "ทำนัดเอง", "text": "man", "data": "abc=abc" }}]}}`

func ServiceListHandler(c *Context) (linebot.SendingMessage, error) {
	var flexContainerStr string
	var packageModels []*model.Package
	db := c.DB
	if err := db.Where("account_id = ?", c.Account.ID).Find(&packageModels).Error; err != nil {
		return nil, err
	}

	for _, packageModel := range packageModels {
		flexContainerStr += fmt.Sprintf(cardTemplate, packageModel) + ","
	}
	flexContainerStr = fmt.Sprintf(carouselTemplate, flexContainerStr)

	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexContainerStr))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), err
}
