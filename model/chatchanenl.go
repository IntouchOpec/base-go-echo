package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

// ChatChannels list ChatChannel
type ChatChannels []ChatChannel

// ChatChannel is soucre for chatbot.
type ChatChannel struct {
	orm.ModelBase

	ChannelID          string       `json:"channel_id" form:"channel_id" binding:"required" gorm:"type:varchar(25);unique_index"`
	Name               string       `json:"name" form:"name" binding:"required" gorm:"type:varchar(25)"`
	LineID             string       `json:"line_id" form:"line_id" binding:"required" gorm:"type:varchar(25);unique_index"`
	ChannelSecret      string       `json:"channel_secret" form:"channel_secret" binding:"required" gorm:"type:varchar(255)"`
	ChannelAccessToken string       `json:"channel_access_token" form:"channel_access_token" binding:"required" gorm:"type:varchar(255)"`
	Type               string       `form:"type" json:"type"  gorm:"type:varchar(10)"`
	PhoneNumber        string       `json:"phone_number" form:"phone_number" binding:"required" gorm:"type:varchar(10)"`
	AccountID          uint         `form:"account_id" json:"account_id" gorm:"not null;"`
	Image              string       `json:"image" form:"image" binding:"required" gorm:"type:varchar(255)"`
	WebSite            string       `json:"website" form:"website" binding:"required" gorm:"type:varchar(255)"`
	WelcomeMessage     string       `json:"welcome_message" form:"welcome_message" binding:"required" gorm:"type:varchar(100)"`
	Address            string       `json:"address" form:"address" binding:"required" gorm:"type:varchar(100)"`
	Account            Account      `gorm:"ForeignKey:AccountID"`
	Settings           []*Setting   `gorm:"many2many:setting_chat_channel;" json:"settings" form:"settings"`
	Customers          []*Customer  `json:"customer"`
	EventLogs          []*EventLog  `json:"event_logs"`
	ActionLogs         []*ActionLog `json:"action_logs"`
	Services           []*Service   `json:"services" gorm:"many2many:service_chat_channel"`
	Promotions         []*Promotion `json:"promotions" gorm:"many2many:chat_channel_promotion"`
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

//
func (cha *ChatChannel) GetSetting(settingNames []string) map[string]string {
	if err := DB().Preload("Settings", "name in (?)", settingNames).Find(&cha).Error; err != nil {
		return nil
	}

	var m map[string]string
	m = make(map[string]string)
	for key := range cha.Settings {
		m[cha.Settings[key].Name] = cha.Settings[key].Value
	}
	return m
}
