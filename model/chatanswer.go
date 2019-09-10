package model

import (
	"github.com/hb-go/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ChatAnswer when customer send message use text find in ChatAnswer.
type ChatAnswer struct {
	gorm.Model
	// ID                uint               `gorm:"primary_key" json:"id"`
	Input     string              `json:"input"`
	TypeInput string              `json:"type_input"`
	Reply     string              `json:"reply"`
	TypeReply linebot.MessageType `json:"type_reply"`
	Active    bool                `json:"active"`
	Account   Account             `gorm:"ForeignKey:id"`
	Source    string              `json:"source"`
	AccountID uint                `form:"account_id" json:"account_id" gorm:"not null;"`
}

// SaveChatAnswer is function create chat answer.
func (cha *ChatAnswer) SaveChatAnswer() *ChatAnswer {
	if err := DB().Create(&cha).Error; err != nil {
		return nil
	}
	return cha
}
