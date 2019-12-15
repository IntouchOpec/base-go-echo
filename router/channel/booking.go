package channel

import (
	"fmt"
	"strconv"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

var serviceMassage string = `
{
	"type": "carousel",
	"contents": [
	  { "type": "bubble",
		"hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s" },
		"body": { "type": "box", "layout": "vertical", "spacing": "sm",
		  "contents": [
			{ "type": "button", "style": "primary", "action": { "type": "message", "label": "ทำนัดอัตโนมัติ", "text": "service auto" } },	
			{ "type": "button", "style": "primary", "action": { "type": "message", "label": "ทำนัดเอง", "text": "service choose" } }
		  ]
		}
	  }
	]
}`
var buttonTimeSecondaryTemplate string = `{"type": "button", "style": "secondary", "margin": "sm", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`
var buttonTimePrimaryTemplate string = `{"type": "button","style": "primary", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`
var buttonTimePrimaryLastTemplate string = `{"type": "button","style": "primary", "margin": "sm", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`
var slotTimeTemplate string = `,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`
var serviceListTemplate string = `{"type": "bubble", "hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s"},
"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
	{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl" },
	{ "type": "box", "layout": "baseline", "contents": [
		{ "type": "text", "text": "฿%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
	] }
	%s]
}}`
var nextPageTemplate string = `{ "type": "bubble", "body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
	{ "type": "button", "flex": 1, "gravity": "center", "action": { "type": "uri", "label": "See more", "uri": "https://linecorp.com" } }] }}`
var thankyouTemplate string = `{
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
}`

func ChooseService(c *Context) (linebot.SendingMessage, error) {
	m := fmt.Sprintf(serviceMassage, c.ChatChannel.ChaImage)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), nil
}

func ServiceList(c *Context) (linebot.SendingMessage, error) {
	var services []model.Service
	var template string
	var image string
	var button string
	if err := c.DB.Where("ser_account_id = ?", c.ChatChannel.AccountID).Find(&services).Error; err != nil {
		return nil, err
	}
	for _, service := range services {
		button = fmt.Sprintf(buttonTimePrimaryTemplate, service.SerName, "", service.SerName)
		image = "https://" + Conf.Server.DomainLineChannel + service.SerImage
		template += fmt.Sprintf(serviceListTemplate, image, service.SerName, strconv.FormatInt(int64(service.SerPrice), 10), button)
	}
	template = fmt.Sprintf(slotTimeTemplate, template)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), err
}

func CalandarHandler(c *Context) (linebot.SendingMessage, error) {
	var m string
	text := c.Massage.Text
	if len(text) > 8 {
		m = lib.MakeCalenda(text[9:19])
	} else {
		m = lib.MakeCalenda("")
	}
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), nil
}

// linebot.SendingMessage
func ServiceListLineHanlder(timeSlot []model.TimeSlot, dateTime string) (linebot.SendingMessage, error) {
	var slotTime string
	var buttonTime string
	var serviceList string
	var actionMessge string
	var count int
	count = 0
	for t := 0; t < len(timeSlot); t++ {
		if count == 2 {
			slotTime = slotTime + fmt.Sprintf(slotTimeTemplate, buttonTime[:len(buttonTime)-1])
			buttonTime = ""
			count = 0
		}

		if len(timeSlot[t].Bookings) > 0 {
			if timeSlot[t].Bookings[0].BooQueue < timeSlot[t].TimeAmount {
				buttonTime = buttonTime + fmt.Sprintf(buttonTimeSecondaryTemplate, timeSlot[t].TimeStart, timeSlot[t].TimeEnd, "เต็มแล้ว")
			} else {
				actionMessge = "booking" + " " + dateTime + " " + timeSlot[t].TimeStart + "-" + timeSlot[t].TimeEnd + " " + timeSlot[t].ProviderService.Service.SerName
				buttonTime = buttonTime + fmt.Sprintf(buttonTimePrimaryTemplate, timeSlot[t].TimeStart, timeSlot[t].TimeEnd, actionMessge)
			}
		} else {
			actionMessge = "booking" + " " + dateTime + " " + timeSlot[t].TimeStart + "-" + timeSlot[t].TimeEnd + " " + timeSlot[t].ProviderService.Service.SerName
			buttonTime = buttonTime + fmt.Sprintf(buttonTimePrimaryLastTemplate, timeSlot[t].TimeStart, timeSlot[t].TimeEnd, actionMessge)
		}

		count = count + 1
		if t == len(timeSlot)-1 {
			slotTime = slotTime + fmt.Sprintf(slotTimeTemplate, buttonTime[:len(buttonTime)-1])
			serviceList += fmt.Sprintf(serviceListTemplate, "https://"+Conf.Server.DomainLineChannel+timeSlot[t].ProviderService.Provider.ProvImage, timeSlot[t].ProviderService.Service.SerName, strconv.FormatInt(int64(timeSlot[t].ProviderService.PSPrice), 10), slotTime)
		} else if timeSlot[t].ProviderService.ID != timeSlot[t+1].ProviderService.ID {
			slotTime = slotTime + fmt.Sprintf(slotTimeTemplate, buttonTime[:len(buttonTime)-1])
			serviceList = serviceList + fmt.Sprintf(serviceListTemplate, "https://"+Conf.Server.DomainLineChannel+timeSlot[t].ProviderService.Provider.ProvImage, timeSlot[t].ProviderService.Service.SerName, strconv.FormatInt(int64(timeSlot[t].ProviderService.PSPrice), 10), slotTime)
			slotTime = ""
			count = 0
			buttonTime = ""
		}
	}
	nextPage := nextPageTemplate

	serviceTamplate := fmt.Sprintf(`{ "type": "carousel", "contents": [%s, %s]}`, serviceList, nextPage)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(serviceTamplate))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), err
}

// ThankyouTemplate
func ThankyouTemplate(c *Context) (linebot.SendingMessage, error) {
	var timeSlot model.TimeSlot
	thankyou := fmt.Sprintf(thankyouTemplate, c.ChatChannel.ChaImage, c.ChatChannel.ChaAddress, timeSlot.TimeStart, timeSlot.TimeEnd)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(thankyou))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), nil
}
