package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

// Account struct.
type Account struct {
	orm.ModelBase

	// ID   uint64 `json:"id,omitempty"`
	Name         string         `json:"name" gorm:"type:varchar(25)"`
	Settings     []*Setting     `json:"settings" gorm:"many2many:account_setting"`
	ChatChannels []*ChatChannel `json:"chat_channels"`
}

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
func (acc *Account) UpdateAccount(id string) *Account {

	if err := DB().Find(&acc, id).Error; err != nil {
		return nil
	}

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
