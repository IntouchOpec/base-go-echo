package channel

import "github.com/line/line-bot-sdk-go/linebot"

func PaymentHandler(c *Context) (linebot.SendingMessage, error) {
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(cardPatmentTemplate))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), nil
}

var cardPatmentTemplate string = `
{
	"type": "bubble",
	"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
	"body": { "type": "box", "layout": "vertical", "spacing": "md", "contents": [
		{ "type": "text", "text": "จองสำเร็จ", "wrap": true, "weight": "bold", "gravity": "center", "size": "xl" },
		{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm",
		  "contents": [
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "Date", "color": "#aaaaaa", "size": "sm", "flex": 1 },
				{ "type": "text", "text": "%s", "wrap": true, "size": "sm", "color": "#666666", "flex": 4 } ]
			},
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "Place", "color": "#aaaaaa", "size": "sm", "flex": 1 },
				{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 4 }
			  ] } ] },
		{ "type": "box", "layout": "vertical", "margin": "xxl", "contents": [
			{ "type": "spacer" },
			{ "type": "image", "url": "%s", "aspectMode": "cover", "size": "xl" },
			{ "type": "image", "url": "https://cdn4.iconfinder.com/data/icons/user-interface-131/32/reload-512.png", "size": "xxs", "margin": "sm" },
			{ "type": "text", "color": "#aaaaaa", "wrap": true, "margin": "xxl", "size": "xs", "text": "%s" }
		  ]
		}
	  ]
	}
  }`
