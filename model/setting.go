package model

import "github.com/hb-go/gorm"

// Setting souce value.
type Setting struct {
	gorm.Model

	Value        string         `json:"value" gorm:"type:varchar(25)"`
	Name         string         `json:"name" gorm:"type:varchar(25)"`
	ChatChannels []*ChatChannel `gorm:"many2many:setting_chat_channel;" json:"chat_channels"`
}

// SaveSetting is function create Setting.
func (setting *Setting) SaveSetting() *Setting {
	if err := DB().Create(&setting).Error; err != nil {
		return nil
	}
	return setting
}

func (setting *Setting) UpdateSetting() *Setting {
	if err := DB().Find(&setting).Error; err != nil {
		return nil
	}
	if err := DB().Save(&setting).Error; err != nil {
		return nil
	}
	return setting
}
