package model

import (
	"github.com/hb-go/gorm"
)

// Account struct.
type Account struct {
	// ID   uint64 `json:"id,omitempty"`
	Name string `json:"name" gorm:"not null; type:varchar(25)"`
	gorm.Model
}

// GetAccount query account list.
func (acc *Account) GetAccount() *[]Account {
	accounts := []Account{}

	DB().Find(&accounts)

	return &accounts
}

// GetAccountByID find account by id.
func (acc *Account) GetAccountByID(id uint) *Account {
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
func (acc *Account) UpdateAccount(id uint) *Account {

	if err := DB().Find(&acc, id).Error; err != nil {
		return nil
	}

	if err := DB().Save(&acc).Error; err != nil {
		return nil
	}

	return acc
}
