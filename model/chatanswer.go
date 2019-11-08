package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ChatAnswer when customer send message use text find in ChatAnswer.
type ChatAnswer struct {
	orm.ModelBase
	// ID                uint               `gorm:"primary_key" json:"id"`
	AnsInput         string              `json:"ans_input"`
	AnsTypeInput     string              `json:"ans_type_input"`
	AnsReply         string              `json:"ans_reply"`
	AnsTypeReply     linebot.MessageType `json:"ans_type_reply"`
	AnsActive        bool                `json:"ans_active"`
	AnsSource        string              `json:"ans_source"`
	AnsChatChannelID uint                `form:"ans_chat_channel_id" json:"chat_channel_id" gorm:"not null;"`
	ChatChannel      ChatChannel         `json:"chat_channel" gorm:"ForeignKey:AnsChatChannelID"`
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
