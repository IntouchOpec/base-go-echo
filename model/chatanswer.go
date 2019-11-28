package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ChatAnswer when customer send message use text find in ChatAnswer.
type ChatAnswer struct {
	orm.ModelBase
	// ID                uint               `gorm:"primary_key" json:"id"`
	AnsInput     string              `form:"input" json:"ans_input"`
	AnsTypeInput string              `form:"type_input" json:"ans_type_input"`
	AnsReply     string              `form:"reply" json:"ans_reply"`
	AnsTypeReply linebot.MessageType `form:"type_reply" json:"ans_type_reply"`
	AnsActive    bool                `form:"active" json:"ans_active"`
	AnsSource    string              `form:"source" json:"ans_source"`
	AccountID    uint                `form:"account_id" json:"account_id" gorm:"not null;"`
	ChatChannels []*ChatChannel      `json:"chat_channels" gorm:"many2many:chat_answer_chat_channel"`
	Account      Account             `json:"account" gorm:"ForeignKey:AccountID"`
}

// SaveChatAnswer is function create chat answer.
func (cha *ChatAnswer) SaveChatAnswer() *ChatAnswer {
	if err := DB().Create(&cha).Error; err != nil {
		return nil
	}
	return cha
}

// GetChatAnswerList is get list ChatAnswer where chat_channel_id.
func GetChatAnswerList(ChatChannelID string) *[]ChatAnswer {
	cha := []ChatAnswer{}
	if err := DB().Where("chat_channel_id = ?", ChatChannelID).Find(&cha).Error; err != nil {
		return nil
	}
	return &cha
}

func GetChatAnswer(id string) *ChatAnswer {
	cha := ChatAnswer{}
	if err := DB().Find(&cha).Error; err != nil {
		return nil
	}
	return &cha
}

func (cha *ChatAnswer) UpdateChatAnswer(id string) *ChatAnswer {
	if err := DB().Find(&cha, id).Error; err != nil {
		return nil
	}

	if err := DB().Save(&cha).Error; err != nil {
		return nil
	}

	return cha
}

func DeleteChatAnswer(id string) *ChatAnswer {
	cha := ChatAnswer{}
	if err := DB().Find(&cha, id).Error; err != nil {
		return nil
	}

	if err := DB().Delete(&cha).Error; err != nil {
		return nil
	}

	return &cha
}
