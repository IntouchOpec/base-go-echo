package model

import (
	"time"

	"github.com/hb-go/gorm"
)

// Promotion discount price product.
type Promotion struct {
	gorm.Model
	// ID        uint      `gorm:"primary_key" json:"id"`
	Discount  int       `json:"discount"`
	Amount    int       `json:"amount"`
	Name      string    `json:"name" gorm:"unique; type:varchar(25)"`
	StartDate time.Time `gorm:"column:start_time" json:"start_time,omitempty"`
	EndDate   time.Time `gorm:"column:end_time" json:"end_time,omitempty"`
	AccountID uint      `form:"account_id" json:"account_id" gorm:"not null;"`
	Account   Account   `gorm:"ForeignKey:id"`
}

// SavePromotion is function create Promotion.
func (promotion *Promotion) SavePromotion() *Promotion {
	if err := DB().Create(&promotion).Error; err != nil {
		return nil
	}
	return promotion
}
