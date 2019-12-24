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

	PromTitle          string       `form:"title" json:"prom_title" gorm:"type:varchar(50)"`
	PromType           string       `form:"prom_type" json:"prom_type" gorm:"type:varchar(25)"`
	PromDiscount       int          `form:"discount" json:"prom_discount"`
	PromCode           string       `form:"code" json:"prom_code" gorm:"type:varchar(25)"`
	PromName           string       `form:"name" json:"prom_name" gorm:"type:varchar(25)"`
	PromImage          string       `form:"image" json:"prom_image" gorm:"type:varchar(255)"`
	ProUsed            int          `json:"pro_used" gorm:"default:0"`
	AccountID          uint         `json:"account_id"`
	RegisterPromotions []*Promotion `json:"register_promotions"`
	Account            Account      `gorm:"ForeignKey:AccountID"`
	Settings           []*Setting   `json:"settings" gorm:"many2many:promotion_setting"`
	// Customers          []*Customer  `gorm:"many2many:promotion_customer" json:"customer"`
	// ChatChannels       []*ChatChannel `json:"chat_channels" gorm:"many2many:chat_channel_promotion"`
	// ProviderServices   []*ProviderService `json:"provider_service" gorm:"many2many:promotion_provider_service"`
}

type Voucher struct {
	orm.ModelBase
	PromotionID   uint         `json:"promotion_id"`
	ChatChannelID uint         `json:"chat_channel_id"`
	ChatChannel   *ChatChannel `json:"chat_channel" gorm:"many2many:chat_channel_promotion"`
	Promotion     Promotion    `json:"promotion"`
	PromStartDate time.Time    `from:"start_time" gorm:"column:start_time" json:"prom_start_time"`
	PromEndDate   time.Time    `from:"end_time" gorm:"column:end_time" json:"prom_end_time"`
	PromAmount    int          `form:"amount" json:"prom_amount"`
	PromCondition string       `form:"condition" json:"prom_condition"`
}

type Coupon struct {
	orm.ModelBase
	PromotionID   uint         `json:"promotion_id"`
	Promotion     Promotion    `json:"promotion"`
	ChatChannelID uint         `json:"chat_channel_id"`
	ChatChannel   *ChatChannel `json:"chat_channel" gorm:"many2many:chat_channel_promotion"`
	PromStartDate time.Time    `from:"start_time" gorm:"column:start_time" json:"prom_start_time"`
	PromEndDate   time.Time    `from:"end_time" gorm:"column:end_time" json:"prom_end_time"`
	PromAmount    int          `form:"amount" json:"prom_amount"`
	PromCondition string       `form:"condition" json:"prom_condition"`
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
	if err := DB().Where("chat_channel_id = ?").Find(&promotions).Error; err != nil {
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

func DeletePromotion(id string) (*Promotion, error) {
	promotion := Promotion{}
	if err := DB().Find(&promotion).Error; err != nil {
		return nil, err
	}

	if err := DB().Delete(&promotion).Error; err != nil {
		return nil, err
	}

	return &promotion, nil
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

func (prom *Promotion) RemoveImage() error {
	prom.PromImage = ""
	if err := db.Save(&prom).Error; err != nil {
		return err
	}
	return nil
}
