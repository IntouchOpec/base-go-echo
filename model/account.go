package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type AccBookingType int8
type AccTransactionType string
type AccTypePayment string

const (
	AccBookingByTimeSlot AccBookingType = 0
	AccBookingByItem     AccBookingType = 1
)

const (
	AccTransactionMan  AccTransactionType = "man"
	AccTransactionAuto AccTransactionType = "auto"
)

const (
	AccTypePaymentBooking AccTypePayment = "booking"
	AccTypePaymentNow     AccTypePayment = "now"
)

// Account struct.
type Account struct {
	orm.ModelBase

	AccProjectID        string             `json:"acc_project_id" grom:"type:varchar(100)"`
	AccAuthJSONFilePath string             `json:"acc_auth_json_file_path" grom:"type:varchar(100)"`
	AccLang             string             `json:"acc_lang" grom:"type:varchar(100)"`
	AccTimeZone         string             `json:"acc_time_zone" grom:"type:varchar(100)"`
	AccName             string             `json:"acc_name" gorm:"type:varchar(25)"`
	AccTransactionType  AccTransactionType `json:"acc_transaction_type" gorm:"type:varchar(25)"`
	AccBookingType      AccBookingType     `json:"acc_booking_type" gorm:"type:varchar(10)"`
	AccTypePayment      AccTypePayment     `json:"acc_type_payment" gorm:"type:varchar(10)"`
	Settings            []*Setting         `json:"settings" gorm:"many2many:account_setting"`
	ChatChannels        []*ChatChannel     `json:"chat_channels"`
}

// func (account *Account) BeforeCreate(scope *gorm.Scope) error {
// 	scope.SetColumn("ID", uuid.New())
// 	return nil
// }

func GetAccountByName(name string) bool {
	account := Account{}

	if err := DB().Where("name = ?", name).Find(&account).Error; err != nil {
		return false
	}

	return true
}

// GetAccount query account list.
func GetAccount() []*Account {
	accounts := []*Account{}

	DB().Find(&accounts)

	return accounts
}

// GetAccountByID find account by id.
func (acc *Account) GetAccountByID(id string) *Account {
	account := Account{}

	if err := DB().Find(&account, id).Error; err != nil {
		return nil
	}

	return &account
}

// CreateAccount is function create accout.
func (acc *Account) CreateAccount() *Account {
	newDb, err := newDB()

	if err != nil {
		return nil
	}

	if err := newDb.Create(&acc).Error; err != nil {
		return nil
	}
	return acc
}

// UpdateAccount edit account soucre.
func (acc *Account) UpdateAccount() *Account {

	if err := DB().Save(&acc).Error; err != nil {
		return nil
	}

	return acc
}

func (acc *Account) RemoveAccount(id string) *Account {
	if err := DB().Find(&acc, id).Error; err != nil {
		return nil
	}
	if err := DB().Delete(&acc, id).Error; err != nil {
		return nil
	}
	return acc
}
