package model

import (
	"fmt"

	"github.com/hb-go/gorm"
)

// Product souce product and service.
type Product struct {
	gorm.Model
	Name        string       `json:"name" gorm:"type:varchar(25)"`
	Detail      string       `json:"detail" gorm:"type:varchar(25)"`
	Price       float32      `json:"price"`
	AccountID   uint         `form:"account_id" json:"account_id" gorm:"not null;"`
	Account     Account      `gorm:"ForeignKey:id"`
	Image       string       `json:"image" gorm:"type:varchar(100)"`
	SubProducts []SubProduct `gorm:"foreignkey:ProductID;" json:"sub_products"`
}

// SubProduct product set.
type SubProduct struct {
	gorm.Model
	From      string `json:"from"`
	To        string `json:"to"`
	Day       string `json:"day"`
	Amount    int    `json:"amount"`
	ProductID uint   `json:"product_id"`
}

// SaveProduct is function create Product.
func (product *Product) SaveProduct() *Product {
	if err := DB().Create(&product).Error; err != nil {
		fmt.Println(err)
		return nil
	}
	return product
}
