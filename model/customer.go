package model

import "github.com/hb-go/gorm"

// Customer follow Line OA.
type Customer struct {
	gorm.Model

	FullName      string       `json:"full_name" gorm:"type:varchar(50);"`
	PictureURL    string       `json:"picture_url"`
	DisplayName   string       `json:"display_name"`
	LineID        string       `json:"line_id" gorm:"type:varchar(255)"`
	Email         string       `json:"email" gorm:"type:varchar(25)"`
	PhoneNumber   string       `json:"phone_number" gorm:"type:varchar(25)"`
	ChatChannelID uint         `json:"chat_channel_id" gorm:"not null;"`
	ChatChannel   ChatChannel  `json:"chat_channel" gorm:"ForeignKey:ChatChannelID"`
	Promotions    []*Promotion `gorm:"many2many:promotion_customer" json:"promotions"`
	EventLogs     []*EventLog  `json:"even_logs"`
	ActionLogs    []*ActionLog `json:"action_logs"`
	Bookings      []*Booking   `json:"bookings"`
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
func (customer *Customer) UpdateCustomerByAtt(pictureURL string, displayName string, email string, phoneNumber string, FillName string) *Customer {
	if err := DB().Preload("Promotions").Where("line_id = ? and chat_channel_id = ?", customer.LineID, customer.ChatChannelID).Find(&customer).Error; err != nil {
		return nil
	}
	customer.PictureURL = pictureURL
	customer.DisplayName = displayName
	customer.Email = email
	customer.PhoneNumber = phoneNumber

	if err := DB().Save(&customer).Error; err != nil {
		return nil
	}

	return customer
}

func GetCustomerList(page, size, chatChannelID int) *[]Customer {
	customer := []Customer{}
	if err := DB().Where("chat_channel_id = ?", chatChannelID).Offset((page - 1) * size).Limit(size).Preload("Bookings").Preload("EventLogs").Preload("ActionLogs").Find(&customer).Error; err != nil {
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
