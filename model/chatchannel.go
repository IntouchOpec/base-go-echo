package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

// ChatChannels list ChatChannel
type ChatChannels []ChatChannel

// ChatChannel is soucre for chatbot.
type ChatChannel struct {
	orm.ModelBase

	ChaChannelID          string       `json:"cha_channel_id" form:"cha_channel_id" binding:"required" gorm:"type:varchar(255);unique_index"`
	ChaName               string       `json:"cha_name" form:"cha_name" binding:"required" gorm:"type:varchar(25)"`
	ChaLineID             string       `json:"cha_line_id" form:"cha_line_id" binding:"required" gorm:"type:varchar(255);unique_index"`
	ChaChannelSecret      string       `json:"cha_channel_secret" form:"cha_channel_secret" binding:"required" gorm:"type:varchar(255)"`
	ChaChannelAccessToken string       `json:"cha_channel_access_token" form:"cha_channel_access_token" binding:"required" gorm:"type:varchar(255)"`
	ChaType               string       `form:"cha_type" json:"cha_type"  gorm:"type:varchar(10)"`
	ChaPhoneNumber        string       `json:"cha_phone_number" form:"cha_cha_phone_number" binding:"required" gorm:"type:varchar(10)"`
	ChaImage              string       `json:"cha_image" form:"cha_image" binding:"required" gorm:"type:varchar(255)"`
	ChaWebSite            string       `json:"cha_website" form:"cha_website" binding:"required" gorm:"type:varchar(255)"`
	ChaWelcomeMessage     string       `json:"cha_welcome_message" form:"cha_welcome_message" binding:"required" gorm:"type:varchar(100)"`
	ChaAddress            string       `json:"cha_address" form:"cha_address" binding:"required" gorm:"type:varchar(100)"`
	AccountID             uint         `form:"account_id" json:"account_id" gorm:"not null;"`
	VoucherID             uint         `json:"voucher_id"`
	Voucher               *Voucher     `json:"voucher" gorm:"ForeignKey:VoucherID"`
	Account               *Account     `gorm:"ForeignKey:AccountID"`
	Customers             []*Customer  `gorm:"many2many:chat_channel_customer"`
	Settings              []*Setting   `gorm:"many2many:setting_chat_channel;" json:"settings" form:"settings"`
	EventLogs             []*EventLog  `json:"event_logs"`
	ActionLogs            []*ActionLog `json:"action_logs"`
	Services              []*Service   `json:"services" gorm:"many2many:service_chat_channel"`
	Promotions            []*Promotion `json:"promotions" gorm:"many2many:chat_channel_promotion"`
}

// SaveChatChannel router create chatchannel.
func (cha *ChatChannel) SaveChatChannel() error {
	if err := DB().Create(&cha).Error; err != nil {
		return err
	}
	return nil
}

// GetChatChannel query account list.
func GetChatChannel(AccountID uint) *ChatChannel {
	chatChannels := ChatChannel{}
	if err := DB().Find(&chatChannels, 1).Error; err != nil {
		return nil
	}
	return &chatChannels
}

func GetChatChannelList() (*[]ChatChannel, error) {
	chatChannels := []ChatChannel{}
	if err := DB().Find(&chatChannels).Error; err != nil {
		return nil, err
	}
	return &chatChannels, nil
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
func DeleteChatChannel(id string) *ChatChannel {
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
