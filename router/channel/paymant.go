package channel

import "github.com/line/line-bot-sdk-go/linebot"

func PaymentHandler(c *Context) (linebot.SendingMessage, error) {
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(cardPatmentTemplate))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("", flexContainer), nil
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

var receiptTemplate string = `{
	"type": "bubble",
	"body": {
	  "type": "box",
	  "layout": "vertical",
	  "contents": [
		{ "type": "text", "text": "RECEIPT", "weight": "bold", "color": "#1DB446", "size": "sm" },
		{ "type": "text", "text": "%s", "weight": "bold", "size": "xxl", "margin": "md" },
		{ "type": "text", "text": "%s", "size": "xs", "color": "#aaaaaa", "wrap": true },
		{
		  "type": "separator",
		  "margin": "xxl"
		},
		{
		  "type": "box",
		  "layout": "vertical",
		  "margin": "xxl",
		  "spacing": "sm",
		  "contents": [
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "Energy Drink", "size": "sm", "color": "#555555", "flex": 0 },
				{ "type": "text", "text": "$2.99", "size": "sm", "color": "#111111", "align": "end" }
			] },
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "Chewing Gum", "size": "sm", "color": "#555555", "flex": 0 },
				{ "type": "text", "text": "$0.99", "size": "sm", "color": "#111111", "align": "end" }
			  ] },
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "Bottled Water", "size": "sm", "color": "#555555", "flex": 0 },
				{ "type": "text", "text": "$3.33", "size": "sm", "color": "#111111", "align": "end" }
			  ] },
			{ "type": "separator", "margin": "xxl" },
			{ "type": "box", "layout": "horizontal", "margin": "xxl", "contents": [
				{ "type": "text", "text": "ITEMS", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "3", "size": "sm", "color": "#111111", "align": "end" }
			  ] },
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "TOTAL", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "$7.31", "size": "sm", "color": "#111111", "align": "end" }
			  ] },
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "CASH", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "$8.0", "size": "sm", "color": "#111111", "align": "end" }
			  ]
			},
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "CHANGE", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "$0.69", "size": "sm", "color": "#111111", "align": "end" }
			  ] }
		  ]
		},
		{ "type": "separator", "margin": "xxl" },
		{ "type": "box", "layout": "horizontal", "margin": "md", "contents": [
			{ "type": "text", "text": "PAYMENT ID", "size": "xs", "color": "#aaaaaa", "flex": 0 },
			{ "type": "text", "text": "#%s", "color": "#aaaaaa", "size": "xs", "align": "end" }
		  ] }
	  ]
	},
	"styles": {
	  "footer": {
		"separator": true
	  }
	}
  }`
