package model

import (
	"fmt"

	"github.com/hb-go/gorm"
)

// Customer follow Line OA.
type Customer struct {
	gorm.Model
	// ID          uint    `json:"id,omitempty"`
	PictureURL    string      `json:"picture_url"`
	DisplayName   string      `json:"display_name"`
	LineID        string      `json:"line_id" gorm:"type:varchar(255)"`
	Email         string      `json:"email" gorm:"type:varchar(25)"`
	PhoneNumber   string      `json:"phone_number" gorm:"type:varchar(25)"`
	ChatChannelID uint        `form:"chat_channel_id" json:"chat_channel_id" gorm:"not null;"`
	ChatChannel   ChatChannel `gorm:"ForeignKey:id"`
}

// LoginRespose is instacne respose line json
type LoginRespose struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

// SaveCustomer is function create Customer.
func (customer *Customer) SaveCustomer() *Customer {
	if err := DB().Create(&customer).Error; err != nil {
		return nil
	}
	return customer
}

// SaveLoginRespose
func (loginRespose *LoginRespose) SaveLoginRespose() *LoginRespose {
	if err := DB().Create(&loginRespose).Error; err != nil {
		return nil
	}
	return loginRespose
}

func (customer *Customer) GetCustomer(id int) *Customer {
	if err := DB().Create(&customer).Error; err != nil {
		return nil
	}
	return customer
}

func (customer *Customer) UpdateCustomer(id int) *Customer {
	if err := DB().Find(&customer, id).Error; err != nil {
		return nil
	}
	if err := DB().Save(&customer).Error; err != nil {
		return nil
	}
	return customer
}

// UpdateCustomerByAtt update by atti
func (customer *Customer) UpdateCustomerByAtt(pictureURL string, displayName string, email string, phoneNumber string) *Customer {
	if err := DB().Where("line_id = ? and chat_channel_id = ?", customer.LineID, customer.ChatChannelID).Find(&customer).Error; err != nil {
		fmt.Println(err)
		fmt.Println("=======================", customer)
		return nil
	}
	customer.PictureURL = pictureURL
	customer.DisplayName = displayName
	customer.Email = email
	customer.PhoneNumber = phoneNumber

	if err := DB().Save(&customer).Error; err != nil {
		fmt.Println(err, "=========")
		return nil
	}

	return customer
}

func GetCustomerList(page, size, chatChannelID int) *[]Customer {
	customer := []Customer{}
	if err := DB().Where("chat_channel_id = ?", chatChannelID).Offset((page - 1) * size).Limit(size).Find(&customer).Error; err != nil {
		return nil
	}
	return &customer
}

func GetCustomer(customerID int) *Customer {
	customer := Customer{}
	if err := DB().Find(&customer, customerID).Error; err != nil {
		return nil
	}
	return &customer
}
