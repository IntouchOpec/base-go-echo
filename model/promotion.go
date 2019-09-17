package model

import (
	"time"

	"github.com/hb-go/gorm"
)

// Promotion discount price product.
type Promotion struct {
	gorm.Model
	// ID        uint      `gorm:"primary_key" json:"id"`
	Discount      int         `json:"discount"`
	Amount        int         `json:"amount"`
	Name          string      `json:"name" gorm:"unique; type:varchar(25)"`
	StartDate     time.Time   `gorm:"column:start_time" json:"start_time,omitempty"`
	EndDate       time.Time   `gorm:"column:end_time" json:"end_time,omitempty"`
	ChatChannelID uint        `json:"chat_channel_id"`
	ChatChannel   ChatChannel `json:"chat_channels"`
	AccountID     uint        `form:"account_id" json:"account_id" gorm:"not null;"`
	Account       Account     `gorm:"ForeignKey:id"`
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
