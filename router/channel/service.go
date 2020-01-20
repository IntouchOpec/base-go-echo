package channel

import (
	"fmt"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/model"

	"github.com/line/line-bot-sdk-go/linebot"
)

func ServiceListHandler(c *Context) (linebot.SendingMessage, error) {
	var flexContainerStr string
	var packageModels []*model.Package
	now := time.Now()
	var timeStart time.Time
	var timeEnd time.Time
	var timeStartStr string
	var timeEndStr string
	var timemillisecon time.Duration
	timeStart = now.Add(30 * time.Hour)
	db := c.DB
	if err := db.Order("pac_order").Where("account_id = ? and pac_is_active = ?", c.Account.ID, true).Find(&packageModels).Error; err != nil {
		return nil, err
	}

	for _, packageModel := range packageModels {
		timemillisecon = time.Duration(packageModel.PacTime.Unix()) * time.Millisecond
		timeEnd = timeStart.Add(timemillisecon)
		timeStartStr = timeStart.Format("01:02")
		timeEndStr = timeEnd.Format("01:02")
		flexContainerStr += fmt.Sprintf(cardTemplate, packageModel.PacName, fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, packageModel.PacImage), timeStartStr, timeEndStr, timeStartStr, timeEndStr) + ","
	}

	flexContainerStr = fmt.Sprintf(carouselTemplate, flexContainerStr[:len(flexContainerStr)-1])

	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexContainerStr))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), err
}

var cardTemplate string = `
{
	"type": "bubble",
	"header": {
	  "type": "box",
	  "layout": "horizontal",
	  "contents": [
		{
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "text",
			  "text": "%s",
			  "size": "sm",
			  "weight": "bold",
			  "color": "#AAAAAA"
			},
			{
			  "type": "box",
			  "layout": "horizontal",
			  "flex": 1,
			  "contents": [
				{
				  "type": "image",
				  "url": "%s",
				  "size": "5xl",
				  "gravity": "center",
				  "flex": 1
				}
			  ]
			}
		  ]
		}
	  ]
	},
	"body": {
	  "type": "box",
	  "layout": "vertical",
	  "spacing": "sm",
	  "contents": [
		{
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "text",
			  "text": "เวลา %s - %s น.",
			  "align": "center"
			}
		  ]
		},
		{
		  "type": "button",
		  "style": "primary",
		  "action": {
			"type": "postback",
			"label": "จอง",
			"text": "man",
			"data": "action=booking&start=%s&end=%s"
		  }
		}
	  ]
	}
  }`
