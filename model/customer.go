package model

import (
	"strconv"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

// Customer follow Line OA.
type Customer struct {
	orm.ModelBase

	CusFullName    string         `json:"cus_full_name" gorm:"type:varchar(50);"`
	CusPictureURL  string         `json:"cus_picture_url"`
	CusDisplayName string         `json:"cus_display_name"`
	CusLineID      string         `json:"cus_line_id" gorm:"type:varchar(255)"`
	CusEmail       string         `json:"cus_email" gorm:"type:varchar(25)"`
	CusPhoneNumber string         `json:"cus_phone_number" gorm:"type:varchar(25)"`
	AccountID      uint           `json:"account_id" gorm:"not null;"`
	CustomerTypeID uint           `json:"customer_type_id"`
	Account        Account        `json:"account" gorm:"ForeignKey:AccountID"`
	CustomerType   CustomerType   `json:"customer" gorm:"ForeignKey:CustomerTypeID"`
	ChatChannels   []*ChatChannel `gorm:"many2many:chat_channel_customer" json:"chat_channels"`
	Promotions     []*Promotion   `gorm:"many2many:promotion_customer" json:"promotions"`
	EventLogs      []*EventLog    `json:"even_logs" gorm:"foreignkey:ID"`
	ActionLogs     []*ActionLog   `json:"action_logs" gorm:"foreignkey:ID"`
	Bookings       []*Booking     `json:"bookings"`
	Transactions   []*Transaction `json:"transactions"`
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

type CustomerType struct {
	orm.ModelBase
	Name      string  `json:"name" gorm:"type:varchar(25)"`
	AccountID uint    `json:"account_id" gorm:"not null;"`
	Accout    Account `json:"chat_channel" gorm:"ForeignKey:AccountID"`
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
func (customer *Customer) UpdateCustomerByAtt(pictureURL, displayName, email, phoneNumber, FullName, Type string) (*Customer, error) {
	if err := DB().Where("cus_line_id = ?", customer.CusLineID).Find(&customer).Error; err != nil {
		return nil, err
	}
	customer.CusPictureURL = pictureURL
	customer.CusDisplayName = displayName
	customer.CusEmail = email
	customer.CusFullName = FullName
	customer.CusPhoneNumber = phoneNumber
	u64, _ := strconv.ParseUint(Type, 10, 32)
	customer.CustomerTypeID = uint(u64)

	if err := DB().Save(&customer).Error; err != nil {
		return nil, err
	}

	return customer, nil
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

func DeleteCustomer(id string) *Customer {
	cus := Customer{}
	if err := DB().Find(&cus, id).Error; err != nil {
		return nil
	}

	if err := DB().Delete(&cus).Error; err != nil {
		return nil
	}

	return &cus
}
