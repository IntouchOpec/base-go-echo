package model

import (
	"time"

	"github.com/hb-go/gorm"
)

type PromotionType string

const (
	PromotionTypePromotion PromotionType = "Promotion"
	PromotionTypeCoupon    PromotionType = "Coupon"
	PromotionTypeVoucher   PromotionType = "Voucher"
)

// Promotion discount price product.
type Promotion struct {
	gorm.Model
	// ID            uint        `gorm:"primary_key" json:"id"`
	Title         string      `json:"title"`
	TypePromotion string      `json:"type_promotion" gorm:"type:varchar(25)"`
	Discount      int         `json:"discount"`
	Amount        int         `json:"amount"`
	Code          string      `json:"code" gorm:"type:varchar(25)"`
	Name          string      `json:"name" gorm:"type:varchar(25)"`
	StartDate     time.Time   `gorm:"column:start_time" json:"start_time"`
	EndDate       time.Time   `gorm:"column:end_time" json:"end_time"`
	Condition     string      `json:"condition"`
	Image         string      `json:"image" gorm:"type:varchar(255)"`
	ChatChannelID uint        `json:"chat_channel_id"`
	ChatChannel   ChatChannel `json:"chat_channel"`
	AccountID     uint        `json:"account_id" gorm:"not null;"`
	Account       Account     `gorm:"ForeignKey:id"`
	Settings      []*Setting  `json:"settings" gorm:"many2many:promotion_setting"`
	Customers     []*Customer `json:"customers" gorm:"many2many:customer_promotion"`
	Products      []*Product  `json:"products" gorm:"many2many:product_promotion"`
}

// SavePromotion is function create Promotion.
func (promotion *Promotion) SavePromotion() *Promotion {
	if err := DB().Create(&promotion).Error; err != nil {
		return nil
	}
	return promotion
}

func (promotion *Promotion) UpdatePromotion(id int) *Promotion {
	if err := DB().Find(&promotion, id).Error; err != nil {
		return nil
	}
	if err := DB().Save(&promotion).Error; err != nil {
		return nil
	}
	return promotion
}

func GetPromotionList(chatChannelID int) *[]Promotion {
	promotions := []Promotion{}
	if err := DB().Where("chatChannelID = ?").Find(&promotions).Error; err != nil {
		return nil
	}
	return &promotions
}

func GetPromotion(id int) *Promotion {
	promotion := Promotion{}
	if err := DB().Find(&promotion, id).Error; err != nil {
		return nil
	}
	return &promotion
}

func DeletePromotion(id int) *Promotion {
	promotion := Promotion{}
	if err := DB().Find(&promotion).Error; err != nil {
		return nil
	}

	if err := DB().Delete(&promotion).Error; err != nil {
		return nil
	}

	return &promotion
}

func (promo *Promotion) GetSettingPromotion(settingNames []string) map[string]string {
	if err := DB().Preload("Settings", "name in (?)", settingNames).Find(&promo).Error; err != nil {
		return nil
	}

	var m map[string]string
	m = make(map[string]string)
	for key := range promo.Settings {
		m[promo.Settings[key].Name] = promo.Settings[key].Value
	}
	return m
}
