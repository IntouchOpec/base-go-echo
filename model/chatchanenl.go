package model

import (
	"github.com/hb-go/gorm"
)

// ChatChannels list ChatChannel
type ChatChannels []ChatChannel

// ChatChannel is soucre for chatbot.
type ChatChannel struct {
	gorm.Model
	// ID                 uint    `gorm:"primary_key" json:"id"`
	ChannelID          string  `json:"channel_id" gorm:"type:varchar(25);unique_index"`
	Name               string  `json:"name" gorm:"type:varchar(25)"`
	LineID             string  `json:"line_id" gorm:"type:varchar(25);unique_index"`
	ChannelSecret      string  `json:"channel_secret" gorm:"type:varchar(255)"`
	ChannelAccessToken string  `json:"channel_access_token" gorm:"type:varchar(255)"`
	Tpye               string  `json:"type" gorm:"type:varchar(10)"`
	AccountID          uint    `form:"account_id" json:"account_id" gorm:"not null;"`
	Account            Account `gorm:"ForeignKey:id"`
}

// SaveChatChannel router create chatchannel.
func (cha *ChatChannel) SaveChatChannel() *ChatChannel {
	if err := DB().Create(&cha).Error; err != nil {
		return nil
	}
	return cha
}

// GetChatChannel query account list.
func (cha *ChatChannels) GetChatChannel() *ChatChannels {
	DB().Find(&cha)
	return cha
}
