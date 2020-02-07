package channel

import (
	"fmt"

	"github.com/jinzhu/gorm"

	. "github.com/IntouchOpec/base-go-echo/conf"

	"github.com/IntouchOpec/base-go-echo/model"
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
