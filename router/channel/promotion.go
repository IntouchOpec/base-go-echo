package channel

import (
	"fmt"

	"github.com/jinzhu/gorm"

	. "github.com/IntouchOpec/base-go-echo/conf"

	"github.com/IntouchOpec/base-go-echo/model"
	// "github.com/hb-go/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
)

func VoucherListHandler(c *Context) (linebot.SendingMessage, error) {
	var voucherCustomers []model.VoucherCustomer
	var template string
	var startDateStr string
	var endDateStr string
	c.DB.Preload("Voucher", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Promotion")
	}).Where("account_id = ? and customer_id = ? and vc_used = ?",
		c.Account.ID, c.Customer.ID, false).Find(&voucherCustomers)
	for _, voucherCustomer := range voucherCustomers {
		startDateStr = voucherCustomer.Voucher.PromStartDate.Format("2006-01-02")
		endDateStr = voucherCustomer.Voucher.PromEndDate.Format("2006-01-02")
		template += fmt.Sprintf(VoucherCard,
			fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, voucherCustomer.Voucher.Promotion.PromImage),
			voucherCustomer.Voucher.Promotion.PromTitle,
			startDateStr, endDateStr, voucherCustomer.Voucher.PromCondition, voucherCustomer.Voucher.Promotion.PromCode) + ","
	}
	if len(voucherCustomers) == 0 {
		template = notFoundVoucher
	} else {
		template = fmt.Sprintf(carouselTemplate, template[:len(template)-1])
	}
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
	if err != nil {
		return nil, err
	}

	return linebot.NewFlexMessage("your voucher", flexContainer), nil
}

func PromotionHandler(c *Context) (linebot.SendingMessage, error) {
	var promotions []*model.Promotion
	var startDateStr string
	var endDateStr string
	var template string
	model.DB().Preload("PromotionDetail").Where("prom_type = ?", "Promotion").Find(&promotions)
	for _, promotion := range promotions {
		if promotion.PromotionDetail == nil {
			continue
		}
		startDateStr = promotion.PromotionDetail.PDStartDate.Format("2006-01-02")
		endDateStr = promotion.PromotionDetail.PDEndDate.Format("2006-01-02")
		template = template + fmt.Sprintf(VoucherCard,
			fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, promotion.PromImage),
			promotion.PromTitle,
			startDateStr,
			endDateStr,
			promotion.PromotionDetail.PDCondition, promotion.PromCode) + ","
	}
	template = fmt.Sprintf(`{
		"type": "carousel",
		"contents": [%s]
	  }`, template[:len(template)-1])
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
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
	template = fmt.Sprintf(`{ "type": "carousel", "contents": [%s]}`, template[:len(template)-1])
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("your voucher", flexContainer), nil
}

// func VoucherTemplate(promotion *model.Promotion, voucher *model.Voucher) string {
// 	StartDateStr := voucher.PromStartDate.Format("02-12-2006")
// 	EndDateStr := voucher.PromEndDate.Format("02-12-2006")
// 	temp := fmt.Sprintf(`{
// 		"type": "bubble",
// 		"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
// 		"body": { "type": "box", "layout": "vertical", "contents": [
// 			{ "type": "text", "text": "%s", "weight": "bold", "size": "xl"},
// 			{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
// 				{ "type": "box", "layout": "baseline", "contents": [
// 					{ "type": "text", "text": "Date", "color": "#aaaaaa", "size": "sm", "flex": 1 },
// 					{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 } ]
// 				},
// 				{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
// 					{ "type": "text", "text": "exp", "color": "#aaaaaa", "flex": 1 },
// 					{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 } ]
// 				},
// 				{"type": "text", "margin": "lg", "text": "%s", "align": "center"},
// 				{"type": "button", "style": "secondary", "action": { "type": "uri", "label": "%s", "uri": "https://web.linecorp.com" }
// 				}
// 			  ]
// 			}
// 		  ]
// 		},
// 		"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
// 			{ "type": "button", "style": "link", "height": "sm",
// 			  "action": { "type": "uri", "label": "เงื่อนไขการใช้", "uri": "https://web.linecorp.com" } }
// 		  ],
// 		  "flex": 0
// 		}
// 	  }`, "https://web."+Conf.Server.DomainWeb+promotion.PromImage, promotion.PromTitle, StartDateStr, EndDateStr, voucher.PromCondition, promotion.PromCode)
// 	return temp
// }
