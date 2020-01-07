package channel

import (
	"fmt"

	. "github.com/IntouchOpec/base-go-echo/conf"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/hb-go/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
)

func promotionCardTemplate(promotion *model.Promotion) string {
	return ""
	// fmt.Sprintf(`{
	// 	"type": "bubble",
	// 	"hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s" },
	// 	"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
	// 		{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl"
	// 		},
	// 		{ "type": "box", "layout": "baseline", "flex": 1, "contents": [
	// 			{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 } ] } ]
	// 	},
	// 	"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
	// 		{ "type": "button", "flex": 2, "style": "primary", "color": "#aaaaaa", "action":
	// 		{ "type": "uri", "label": "Add to Cart", "uri": "https://linecorp.com" } } ]
	// 	}
	//   },`, promotion.PromImage, promotion.PromTitle, promotion.PromCondition)
}

func PromotionHandler(c *Context) (linebot.SendingMessage, error) {
	var promotions []*model.Promotion
	var promotionCards string
	model.DB().Where("promotion_type = ?", "Promotion").Find(&promotions)
	for _, promotion := range promotions {
		promotionCards = promotionCards + promotionCardTemplate(promotion)
	}
	var promotionsTemplate string = fmt.Sprintf(`{
		"type": "carousel",
		"contents": [%s]
	  }`, promotionCards[:len(promotionCards)-1])
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(promotionsTemplate))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), nil
}

func VoucherHandler(c *Context) (linebot.SendingMessage, error) {
	customer := model.Customer{}
	model.DB().Preload("Promotions", "prom_type = ?", "voucher", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Settings")
	}).Where("cus_line_id = ?", c.Massage.ID).Find(&customer)
	var template string
	// for _, voucher := range customer.VoucherCustomers {
	// template = template + VoucherTemplate(voucher) + ","
	// }
	template = fmt.Sprintf(`{ "type": "carousel", "contents": [%s]}`, template[:len(template)-1])
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("your voucher", flexContainer), nil
}

func VoucherTemplate(promotion *model.Promotion, voucher *model.Voucher) string {
	StartDateStr := voucher.PromStartDate.Format("02-12-2006")
	EndDateStr := voucher.PromEndDate.Format("02-12-2006")
	temp := fmt.Sprintf(`{
		"type": "bubble",
		"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
		"body": { "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "%s", "weight": "bold", "size": "xl"},
			{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
				{ "type": "box", "layout": "baseline", "contents": [
					{ "type": "text", "text": "Date", "color": "#aaaaaa", "size": "sm", "flex": 1 },
					{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 } ]
				},
				{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
					{ "type": "text", "text": "exp", "color": "#aaaaaa", "flex": 1 },
					{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 } ]
				},
				{"type": "text", "margin": "lg", "text": "%s", "align": "center"},
				{"type": "button", "style": "secondary", "action": { "type": "uri", "label": "%s", "uri": "https://web.linecorp.com" }
				}
			  ]
			}
		  ]
		},
		"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
			{ "type": "button", "style": "link", "height": "sm",
			  "action": { "type": "uri", "label": "เงื่อนไขการใช้", "uri": "https://web.linecorp.com" } }
		  ],
		  "flex": 0
		}
	  }`, "https://web."+Conf.Server.DomainWeb+promotion.PromImage, promotion.PromTitle, StartDateStr, EndDateStr, voucher.PromCondition, promotion.PromCode)
	return temp
}
