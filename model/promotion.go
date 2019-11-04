package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type PromotionType string

const (
	PromotionTypePromotion PromotionType = "Promotion"
	PromotionTypeCoupon    PromotionType = "Coupon"
	PromotionTypeVoucher   PromotionType = "Voucher"
)

// Promotion discount price service.
type Promotion struct {
	orm.ModelBase

	Title         string         `form:"title" json:"title"`
	TypePromotion string         `form:"type_promotion" json:"type_promotion" gorm:"type:varchar(25)"`
	Discount      int            `form:"discount" json:"discount"`
	Amount        int            `form:"amount" json:"amount"`
	Code          string         `form:"code" json:"code" gorm:"type:varchar(25)"`
	Name          string         `form:"name" json:"name" gorm:"type:varchar(25)"`
	StartDate     time.Time      `from:"start_time" gorm:"column:start_time" json:"start_time"`
	EndDate       time.Time      `from:"end_time" gorm:"column:end_time" json:"end_time"`
	Condition     string         `form:"condition" json:"condition"`
	Image         string         `form:"image" json:"image" gorm:"type:varchar(255)"`
	Customers     []*Customer    `gorm:"many2many:promotion_customer" json:"customer"`
	ChatChannels  []*ChatChannel `json:"chat_channels" gorm:"many2many:chat_channel_promotion"`
	AccountID     uint           `json:"account_id"`
	Account       Account        `gorm:"ForeignKey:AccountID"`
	Settings      []*Setting     `json:"settings" gorm:"many2many:promotion_setting"`
	services      []*Service     `json:"services" gorm:"many2many:service_promotion"`
}

// SavePromotion is function create Promotion.
func (promotion *Promotion) SavePromotion() *Promotion {
	db := DB()
	if err := db.Set("gorm:association_autoupdate", false).Create(&promotion).Error; err != nil {
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
