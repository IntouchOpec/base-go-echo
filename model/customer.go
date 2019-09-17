package model

import "github.com/hb-go/gorm"

// Customer follow Line OA.
type Customer struct {
	gorm.Model
	// ID          uint    `json:"id,omitempty"`
	LineID      string  `json:"line_id" gorm:"type:varchar(25)"`
	Email       string  `json:"email" gorm:"type:varchar(25)"`
	PhoneNumber string  `json:"phone_number" gorm:"type:varchar(25)"`
	AccountID   uint    `form:"account_id" json:"account_id" gorm:"not null;"`
	Account     Account `gorm:"ForeignKey:id"`
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
