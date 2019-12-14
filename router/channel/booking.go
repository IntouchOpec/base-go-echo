package channel

import (
	"fmt"
	"log"
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
			{ "type": "button", "style": "primary",
			  "action": { "type": "message", "label": "ทำนัดอัตโนมัติ", "text": "service auto" }
			},	
			{
			  "type": "button",
			  "style": "primary",
			  "action": { "type": "message", "label": "ทำนัดเอง", "text": "service choose" }
			}
		  ]
		}
	  }
	]
}`

func ChooseService(c *Context) linebot.SendingMessage {
	m := fmt.Sprintf(serviceMassage, c.ChatChannel.ChaImage)
	flexContainer, _ := linebot.UnmarshalFlexMessageJSON([]byte(m))
	return linebot.NewFlexMessage("service", flexContainer)
}

func CalandarHandler(c *Context) linebot.SendingMessage {
	var m string
	if len(c.Massage) > 8 {
		m = lib.MakeCalenda(c.Massage[9:19])
	} else {
		m = lib.MakeCalenda("")
	}
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
	if err != nil {
		log.Println(err)
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer)
}

// linebot.SendingMessage
func serviceListLineHanlder(timeSlot []model.TimeSlot, dateTime string) linebot.SendingMessage {
	var slotTime string
	var buttonTime string
	var serviceList string
	var count int
	count = 0
	for t := 0; t < len(timeSlot); t++ {
		if count == 2 {
			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
			buttonTime = ""
			count = 0
		}

		if len(timeSlot[t].Bookings) > 0 {
			if timeSlot[t].Bookings[0].BooQueue < timeSlot[t].TimeAmount {
				buttonTime = buttonTime + fmt.Sprintf(`{"type": "button", "style": "secondary", "margin": "sm", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`,
					timeSlot[t].TimeStart, timeSlot[t].TimeEnd, "เต็มแล้ว")
			} else {
				buttonTime = buttonTime + fmt.Sprintf(`{"type": "button","style": "primary", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`,
					timeSlot[t].TimeStart, timeSlot[t].TimeEnd, "booking"+" "+dateTime+" "+timeSlot[t].TimeStart+"-"+timeSlot[t].TimeEnd+" "+timeSlot[t].ProviderService.Service.SerName)
			}
		} else {
			buttonTime = buttonTime + fmt.Sprintf(`{"type": "button","style": "primary", "margin": "sm", "action": { "type": "message", "label": "%s-%s", "text": "%s" }},`,
				timeSlot[t].TimeStart, timeSlot[t].TimeEnd, "booking"+" "+dateTime+" "+timeSlot[t].TimeStart+"-"+timeSlot[t].TimeEnd+" "+timeSlot[t].ProviderService.Service.SerName)
		}

		count = count + 1
		if t == len(timeSlot)-1 {
			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
			serviceList += fmt.Sprintf(`{"type": "bubble", "hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s"},
				"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
					{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl" },
					{ "type": "box", "layout": "baseline", "contents": [
						{ "type": "text", "text": "฿%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
					] }
					%s]
				}}`, "https://"+Conf.Server.DomainLineChannel+timeSlot[t].ProviderService.Provider.ProvImage, timeSlot[t].ProviderService.Service.SerName, strconv.FormatInt(int64(timeSlot[t].ProviderService.PSPrice), 10), slotTime)
		} else if timeSlot[t].ProviderService.ID != timeSlot[t+1].ProviderService.ID {
			slotTime = slotTime + fmt.Sprintf(`,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`, buttonTime[:len(buttonTime)-1])
			serviceList = serviceList + fmt.Sprintf(`{"type": "bubble", "hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s"},
				"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
					{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl" },
					{ "type": "box", "layout": "baseline", "contents": [
						{ "type": "text", "text": "฿%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
					] }
					%s
				]
				}},`, "https://"+Conf.Server.DomainLineChannel+timeSlot[t].ProviderService.Provider.ProvImage, timeSlot[t].ProviderService.Service.SerName, strconv.FormatInt(int64(timeSlot[t].ProviderService.PSPrice), 10), slotTime)
			slotTime = ""
			count = 0
			buttonTime = ""
		}
	}
	nextPage := `{ "type": "bubble", "body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
			{ "type": "button", "flex": 1, "gravity": "center", "action": { "type": "uri", "label": "See more", "uri": "https://linecorp.com" } }] }}`

	serviceTamplate := fmt.Sprintf(`{ "type": "carousel", "contents": [%s, %s]}`, serviceList, nextPage)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(serviceTamplate))
	fmt.Println(serviceTamplate)
	if err != nil {
		log.Println(err, "====>>")
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer)
}
