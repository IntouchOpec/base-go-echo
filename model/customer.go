package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

// Customer follow Line OA.
type Customer struct {
	orm.ModelBase

	CusFullName    string       `json:"cus_full_name" gorm:"type:varchar(50);"`
	CusPictureURL  string       `json:"cus_picture_url"`
	CusDisplayName string       `json:"cus_display_name"`
	CusLineID      string       `json:"cus_line_id" gorm:"type:varchar(255)"`
	CusEmail       string       `json:"cus_email" gorm:"type:varchar(25)"`
	CusPhoneNumber string       `json:"cus_phone_number" gorm:"type:varchar(25)"`
	CusAccountID   uint         `json:"cus_chat_channel_id" gorm:"not null;"`
	CustomerTpyeID uint         `json:"customer_type_id"`
	Accout         Account      `json:"account" gorm:"ForeignKey:CusAccountID"`
	CustomerTpye   CustomerTpye `json:"customer" gorm:"ForeignKey:CustomerTpyeID"`
	Promotions     []*Promotion `gorm:"many2many:promotion_customer" json:"promotions"`
	EventLogs      []*EventLog  `json:"even_logs"`
	ActionLogs     []*ActionLog `json:"action_logs"`
	Bookings       []*Booking   `json:"bookings"`
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

type CustomerTpye struct {
	orm.ModelBase
	Name      string  `json:"name" gorm:"type:vachat(25)"`
	AccoutnID uint    `json:"acount_id" gorm:"not null;"`
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
func (customer *Customer) UpdateCustomerByAtt(pictureURL string, displayName string, email string, phoneNumber string, FillName string) *Customer {
	if err := DB().Preload("Promotions").Where("line_id = ? and chat_channel_id = ?", customer.CusLineID, customer.CusAccountID).Find(&customer).Error; err != nil {
		return nil
	}
	customer.CusPictureURL = pictureURL
	customer.CusDisplayName = displayName
	customer.CusEmail = email
	customer.CusPhoneNumber = phoneNumber

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
