package model

import (
	"fmt"

	"github.com/hb-go/gorm"
)

// Product souce product and service.
type Product struct {
	gorm.Model
	Name         string         `form:"name" json:"name" gorm:"type:varchar(25)"`
	Detail       string         `form:"detail" json:"detail" gorm:"type:varchar(25)"`
	Price        float32        `form:"price" json:"price"`
	Active       bool           `form:"active" json:"active"`
	AccountID    uint           `form:"account_id" json:"account_id" gorm:"not null;"`
	Account      Account        `gorm:"ForeignKey:id"`
	Image        string         `form:"image" json:"image" gorm:"type:varchar(255)"`
	Chatchannels []*ChatChannel `json:"chat_channels" gorm:"many2many:product_chat_channel"`
	SubProducts  []*SubProduct  `gorm:"ForeignKey:ProductID;" json:"sub_products"`
	Promotions   []*Promotion   `json:"promotions" gorm:"many2many:product_promotion;"`
}

// SubProduct product set.
type SubProduct struct {
	gorm.Model
	Start     string     `json:"start"`
	End       string     `json:"end"`
	Day       int        `json:"day"`
	Amount    int        `json:"amount"`
	Active    bool       `json:"active"`
	ProductID uint       `json:"product_id"`
	Bookings  []*Booking `json:"bookings"`
	Product   Product    `json:"product" gorm:"ForeignKey:ProductID;"`
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

func DeleteProductByID(id string) *Product {
	product := Product{}
	if err := DB().Find(&product, id).Error; err != nil {
		return nil
	}
	if err := DB().Delete(&product).Error; err != nil {
		return nil
	}
	return &product
}

func DeleteSubProduct(id string) *SubProduct {
	subProduct := SubProduct{}
	if err := DB().Find(&subProduct, id).Error; err != nil {
		return nil
	}
	if err := DB().Delete(&subProduct).Error; err != nil {
		return nil
	}
	return &subProduct
}
