package channel

import (
	"fmt"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/router/web"
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
	for _, promotion := range customer.Promotions {
		template = template + web.VoucherTemplate(promotion) + ","
	}
	template = fmt.Sprintf(`{ "type": "carousel", "contents": [%s]}`, template[:len(template)-1])
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("your voucher", flexContainer), nil
}
