package model

import (
	"fmt"

	"github.com/hb-go/gorm"
)

// Product souce product and service.
type Product struct {
	gorm.Model
	Name          string       `json:"name" gorm:"type:varchar(25)"`
	Detail        string       `json:"detail" gorm:"type:varchar(25)"`
	Price         float32      `json:"price"`
	AccountID     uint         `form:"account_id" json:"account_id" gorm:"not null;"`
	Account       Account      `gorm:"ForeignKey:id"`
	Image         string       `json:"image" gorm:"type:varchar(255)"`
	ChatChannelID uint         `json:"chat_channel_id"`
	Chatchannel   ChatChannel  `gorm:"foreignkey:chatchannelID;" json:"chat_channels"`
	SubProducts   []SubProduct `gorm:"foreignkey:ProductID;" json:"sub_products"`
}

// SubProduct product set.
type SubProduct struct {
	gorm.Model
	Start     string  `json:"start"`
	End       string  `json:"end"`
	Day       int     `json:"day"`
	Amount    int     `json:"amount"`
	Product   Product `json:"product"`
	ProductID uint    `json:"product_id"`
}

// SaveProduct is function create Product.
func (product *Product) SaveProduct() *Product {
	if err := DB().Create(&product).Error; err != nil {
		fmt.Println(err)
		return nil
	}
	return product
}

func (product *Product) UpdateProduct(id int) *Product {

	if err := DB().Find(&product, id).Error; err != nil {
		return nil
	}

	if err := DB().Save(&product).Error; err != nil {
		return nil
	}

	return product
}

func (subProduct *SubProduct) CreateSubProduct() *SubProduct {
	if err := DB().Create(&subProduct).Error; err != nil {
		return nil
	}
	return subProduct
}

func (subProduct *SubProduct) UpdateSubProduct(id int) *SubProduct {
	if err := DB().Find(&subProduct, id).Error; err != nil {
		return nil
	}

	if err := DB().Save(&subProduct).Error; err != nil {
		return nil
	}

	return subProduct
}

func GetProduct(chatchannelID int) *[]Product {
	products := []Product{}
	if err := DB().Where("").Find(&products).Error; err != nil {
		return nil
	}
	return &products
}

func GetProductByID(chatchannelID, id int) *Product {
	product := Product{}
	if err := DB().Where("ChatChannelID = ?", chatchannelID).Find(&product, id).Error; err != nil {
		return nil
	}
	return &product
}
