package channel

import (
	"fmt"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/model"

	"github.com/line/line-bot-sdk-go/linebot"
)

func ServiceNowListHandler(c *Context) (linebot.SendingMessage, error) {
	var flexContainerStr string
	var packageModels []*model.Package
	now := time.Now()
	var timeStart time.Time
	var timeEnd time.Time
	var timeStartStr string
	var timeEndStr string
	var duration time.Duration
	var button string
	var pagination Pagination
	var total int
	pagination.ParseQueryUnmarshal(c.Event.Postback.Data)
	pagination.SetPagination()
	timeStart = now.Add(30 * time.Minute)
	db := c.DB
	filter := db.Model(&packageModels).Where("account_id = ? and pac_is_active = ?", c.Account.ID, true).Count(&total)
	pagination.MakePagination(total, 9)
	filter.Order("pac_order").Limit(pagination.Record).Offset(pagination.Offset).Find(&packageModels)
	for _, packageModel := range packageModels {
		duration = time.Duration(packageModel.PacTime.Hour() * int(time.Hour))
		timeEnd = timeStart.Add(duration)
		duration = time.Duration(packageModel.PacTime.Minute() * int(time.Minute))
		timeEnd = timeEnd.Add(duration)
		timeStartStr = timeStart.Format("15:04")
		timeEndStr = timeEnd.Format("15:04")
		flexContainerStr += fmt.Sprintf(cardPackageTemplate, packageModel.PacName, fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, packageModel.PacImage), timeStartStr, timeEndStr, timeStartStr, timeEndStr) + ","
	}
	if len(packageModels) < 9 {
		var services []*model.Service
		filter = db.Model(&services).Where("account_id = ? and ser_active = ?", c.Account.ID, true).Count(&total)
		pagination.SetPagination()
		pagination.MakePagination(total, 9-len(packageModels))
		filter.Limit(pagination.Record).Offset(pagination.Offset).Preload("ServiceItems", "ss_is_active = ?", true).Find(&services)
		for _, service := range services {
			button = ""
			if len(service.ServiceItems) == 0 {
				continue
			}
			for _, item := range service.ServiceItems {
				button += fmt.Sprintf(buttonTemplate, item.SSName, fmt.Sprintf("action=booking&service_item_id=%d&start=%s&end=%s&day=%s", item.ID, timeStartStr, timeEndStr, timeStart.Format("2006-01-02")))
			}
			flexContainerStr += fmt.Sprintf(cardServiceTemplate, service.SerName, fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, service.SerImage), service.SerDetail, button[:len(button)-1]) + ","
		}
	}
	flexContainerStr = fmt.Sprintf(carouselTemplate, flexContainerStr[:len(flexContainerStr)-1])
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexContainerStr))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), err
}

func ServiceDateListHandler(c *Context, date string) (linebot.SendingMessage, error) {
	var flexContainerStr string
	var packageModels []*model.Package
	var services []*model.Service
	timeStart, err := time.Parse("2006-01-02T15:04", date)
	if err != nil {
		return nil, err
	}
	var timeEnd time.Time
	var timeStartStr string
	var timeEndStr string
	var duration time.Duration
	var button string
	db := c.DB
	if err := db.Limit(9).Order("pac_order").Where("account_id = ? and pac_is_active = ?", c.Account.ID, true).Find(&packageModels).Error; err != nil {
		return nil, err
	}

	for _, packageModel := range packageModels {
		duration = time.Duration(packageModel.PacTime.Hour() * int(time.Hour))
		timeEnd = timeStart.Add(duration)
		duration = time.Duration(packageModel.PacTime.Minute() * int(time.Minute))
		timeEnd = timeEnd.Add(duration)
		timeStartStr = timeStart.Format("15:04")
		timeEndStr = timeEnd.Format("15:04")
		flexContainerStr += fmt.Sprintf(cardPackageTemplate,
			packageModel.PacName, fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, packageModel.PacImage), timeStartStr, timeEndStr, timeStartStr, timeEndStr, packageModel.ID) + ","
	}
	if len(packageModels) < 9 {
		if err := db.Preload("ServiceItems", "ss_is_active = ?", true).Where("account_id = ?", c.Account.ID).Find(&services).Error; err != nil {
			return nil, err
		}
		for _, service := range services {
			button = ""
			if len(service.ServiceItems) == 0 {
				continue
			}
			for _, item := range service.ServiceItems {
				duration = time.Duration(item.SSTime.Hour() * int(time.Hour))
				timeEnd = timeStart.Add(duration)
				duration = time.Duration(item.SSTime.Minute() * int(time.Minute))
				timeEnd = timeEnd.Add(duration)
				timeStartStr = timeStart.Format("15:04")
				timeEndStr = timeEnd.Format("15:04")
				button += fmt.Sprintf(buttonTemplate, item.SSName, fmt.Sprintf("action=booking&service_item_id=%d&start=%s&end=%s&day=%s", item.ID, timeStartStr, timeEndStr, timeStart.Format("2006-01-02")))
			}
			flexContainerStr += fmt.Sprintf(cardServiceTemplate, service.SerName, fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, service.SerImage), service.SerDetail, button[:len(button)-1]) + ","
		}
	}
	flexContainerStr = fmt.Sprintf(carouselTemplate, flexContainerStr[:len(flexContainerStr)-1])
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexContainerStr))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), err
}

var buttonTemplate string = `
{ "type": "button", "style": "primary", "action": { "type": "postback", "label": "%s", "data": "%s" }},`

var cardServiceTemplate string = `
{
	"type": "bubble",
	"header": { "type": "box", "layout": "horizontal", "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "%s", "size": "sm", "weight": "bold", "color": "#AAAAAA" },
			{ "type": "box", "layout": "horizontal", "flex": 1, "contents": [
				{ "type": "image", "url": "%s", "size": "5xl", "gravity": "center", "flex": 1 }
			  ]
			}
		  ]
		}
	  ]
	},
	"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "%s", "align": "center" } ]
		},
		%s
	  ]
	}
}`

var cardPackageTemplate string = `
{
	"type": "bubble",
	"header": { "type": "box", "layout": "horizontal", "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "%s", "size": "sm", "weight": "bold", "color": "#AAAAAA" },
			{ "type": "box", "layout": "horizontal", "flex": 1, "contents": [
				{ "type": "image", "url": "%s", "size": "5xl", "gravity": "center", "flex": 1 }
			  ] } ] } ] },
	"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "เวลา %s - %s น.", "align": "center" } ] },
		{ "type": "button", "style": "primary", "action": {
			"type": "postback", "label": "จอง", "data": "action=booking&start=%s&end=%s&package_id=%s"
		  }
		}
	  ]
	}
}`
