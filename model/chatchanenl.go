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
	Image              string  `json:"image" gorm:"type:varchar(255)"`
	WebSite            string  `json:"website" gorm:"type:varchar(255)"`
	WelcomeMessage     string  `json:"welcome_message" gorm:"type:varchar(100)"`
	Address            string  `json:"address" gorm:"type:varchar(100)"`
}

// SaveChatChannel router create chatchannel.
func (cha *ChatChannel) SaveChatChannel() *ChatChannel {
	if err := DB().Create(&cha).Error; err != nil {
		return nil
	}
	return cha
}

// GetChatChannel query account list.
func GetChatChannel(chatChannelID, size, page int) *[]ChatChannel {
	chatChannels := []ChatChannel{}
	if err := DB().Where("AccountID = ?", chatChannelID).Offset((page - 1) * size).Limit(size).Find(&chatChannels).Error; err != nil {
		return nil
	}
	return &chatChannels
}

func GetChatChannelByID(id string) *ChatChannel {
	chatChannel := ChatChannel{}
	if err := DB().Find(&chatChannel, id).Error; err != nil {
		return nil
	}
	return &chatChannel
}

// EditChatChannel update ChatChannel.
func (cha *ChatChannel) EditChatChannel(id int) *ChatChannel {

	if err := DB().Find(&cha, id).Error; err != nil {
		return nil
	}

	if err := DB().Save(&cha).Error; err != nil {
		return nil
	}

	return cha
}

// DeleteChatChannel delete ChatChannels
func DeleteChatChannel(id int) *ChatChannel {
	cha := ChatChannel{}
	if err := DB().Find(&cha, id).Error; err != nil {
		return nil
	}

	if err := DB().Delete(&cha).Error; err != nil {
		return nil
	}

	return &cha
}
