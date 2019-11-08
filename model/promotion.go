package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type PromotionType string

const (
	PromotionPromotionType PromotionType = "Promotion"
	PromotionTypeCoupon    PromotionType = "Coupon"
	PromotionTypeVoucher   PromotionType = "Voucher"
)

// Promotion discount price service.
type Promotion struct {
	orm.ModelBase

	PromTitle         string             `form:"title" json:"prom_title" gorm:"type:varchar(50)"`
	PromPromotionType string             `form:"promotion_type" json:"prom_promotion_type" gorm:"type:varchar(25)"`
	PromDiscount      int                `form:"discount" json:"prom_discount"`
	PromAmount        int                `form:"amount" json:"prom_amount"`
	PromCode          string             `form:"code" json:"prom_code" gorm:"type:varchar(25)"`
	PromName          string             `form:"name" json:"prom_name" gorm:"type:varchar(25)"`
	PromStartDate     time.Time          `from:"start_time" gorm:"column:start_time" json:"prom_start_time"`
	PromEndDate       time.Time          `from:"end_time" gorm:"column:end_time" json:"prom_end_time"`
	PromCondition     string             `form:"condition" json:"prom_condition"`
	PromImage         string             `form:"image" json:"prom_image" gorm:"type:varchar(255)"`
	PromAccountID     uint               `json:"prom_account_id"`
	Customers         []*Customer        `gorm:"many2many:promotion_customer" json:"customer"`
	ChatChannels      []*ChatChannel     `json:"chat_channels" gorm:"many2many:chat_channel_promotion"`
	Account           Account            `gorm:"ForeignKey:PromAccountID"`
	Settings          []*Setting         `json:"settings" gorm:"many2many:promotion_setting"`
	ProviderServices  []*ProviderService `json:"provider_service" gorm:"many2many:promotion_provider_service"`
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
