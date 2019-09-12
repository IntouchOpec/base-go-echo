package model

import "github.com/hb-go/gorm"

// Setting souce value.
type Setting struct {
	gorm.Model
	// ID        int     `gorm:"primary_key" json:"id"`
	Value         string  `json:"value" gorm:"unique; type:varchar(25)"`
	Name          string  `json:"name" gorm:"unique; type:varchar(25)"`
	ChatChannelID uint    `form:"chat_channel_id" json:"chat_channel_id" gorm:"not null;"`
	ChatChannel   Account `gorm:"ForeignKey:id"`
}

// SaveSetting is function create Setting.
func (setting *Setting) SaveSetting() *Setting {
	if err := DB().Create(&setting).Error; err != nil {
		return nil
	}
	return setting
}
