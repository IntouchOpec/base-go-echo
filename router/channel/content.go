package channel

import (
	"fmt"
	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

// . "github.com/IntouchOpec/base-go-echo/conf"

func ContentListHandler(c *Context) (linebot.SendingMessage, error) {
	var flexContainerStr string
	var contents []model.Content
	var total int
	var urlLiff string
	var textContent string
	filter := c.DB.Model(&contents).Where("account_id = ? and con_is_active = ?", c.Account.ID, true).Count(&total)
	filter.Find(&contents)
	setting := c.ChatChannel.GetSetting([]string{model.NameLIFFIDContent})
	for _, content := range contents {
		if len(content.ConDetail) > 30 {
			textContent = content.ConDetail[:30]
		} else {
			textContent = content.ConDetail
		}
		urlLiff = fmt.Sprintf("line://app/%s?contentID=%d", setting[model.NameLIFFIDContent], content.ID)
		flexContainerStr += fmt.Sprintf(ContentTemplate,
			fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, content.ConImage),
			urlLiff,
			content.ConTitle,
			textContent,
			urlLiff,
		)
	}
	flexContainerStr = fmt.Sprintf(carouselTemplate, flexContainerStr[:len(flexContainerStr)-1])
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexContainerStr))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), err
}

var ContentTemplate string = `
{
	"type": "bubble",
	"hero": {
	  "type": "image",
	  "url": "%s",
	  "size": "full",
	  "aspectRatio": "20:13",
	  "aspectMode": "cover",
	  "action": {
		"type": "uri",
		"uri": "%s"
	  }
	},
	"body": {
	  "type": "box",
	  "layout": "vertical",
	  "contents": [
		{
		  "type": "text",
		  "text": "%s",
		  "weight": "bold",
		  "size": "xl"
		},
		{
		  "type": "text",
		  "text": "%s"
		}
	  ]
	},
	"footer": {
	  "type": "box",
	  "layout": "vertical",
	  "spacing": "sm",
	  "contents": [
		{
		  "type": "button",
		  "style": "link",
		  "height": "sm",
		  "action": {
			"type": "uri",
			"label": "WEBSITE",
			"uri": "%s"
		  }
		},
		{
		  "type": "spacer",
		  "size": "sm"
		}
	  ],
	  "flex": 0
	}
  },`
