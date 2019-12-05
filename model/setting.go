package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

// Setting souce value.
type Setting struct {
	orm.ModelBase

	Value        string         `json:"value" gorm:"type:varchar(255)"`
	Name         string         `json:"name" gorm:"type:varchar(255)"`
	ChatChannels []*ChatChannel `gorm:"many2many:setting_chat_channel;" json:"chat_channels"`
	Promotions   []*Promotion   `gorm:"many2many:promotion_setting;" json:"promotions"`
	Accounts     []*Account     `gorm:"many2many:account_setting;" json:"accounts"`
}

// SaveSetting is function create Setting.
func (setting *Setting) SaveSetting() *Setting {
	if err := DB().Create(&setting).Error; err != nil {
		return nil
	}
	return setting
}

// UpdateSetting
func (setting *Setting) UpdateSetting() *Setting {
	if err := DB().Save(&setting).Error; err != nil {
		return nil
	}
	return setting
}
