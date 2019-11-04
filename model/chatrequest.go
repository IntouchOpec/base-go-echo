package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

// ChatRequest log castomer when customer send something to line OA.
type ChatRequest struct {
	orm.ModelBase

	LineID       string     `json:"line_id" gorm:"type:varchar(25)"`
	Message      string     `json:"message" gorm:"type:varchar(50)"`
	ReplyToken   string     `json:"reply_token" gorm:"type:varchar(255)"`
	MessageType  string     `json:"message_type" gorm:"type:varchar(25)"`
	ChatAnswerID uint       `form:"chat_answer_id" json:"chat_answer_id" gorm:"not null;"`
	ChatAnswer   ChatAnswer `gorm:"foreignkey:ID"`
	AccountID    uint       `form:"account_id" json:"account_id" gorm:"not null;"`
	Account      Account    `gorm:"ForeignKey:id"`
}

// SaveChatRequest is function create chat answer.
func (chatReq *ChatRequest) SaveChatRequest() *ChatRequest {
	if err := DB().Create(&chatReq).Error; err != nil {
		return nil
	}
	return chatReq
}

// GetChatRequest get ChatRequest list.
func GetChatRequest(chatChannelID string) *ChatRequest {
	chatReq := ChatRequest{}
	DB().Where("ChatChannelID = ?", chatChannelID).Find(&chatReq)
	return &chatReq
}

// EditChatRequest update ChatRequest.
func (chatReq *ChatRequest) EditChatRequest(id int) *ChatRequest {

	if err := DB().Find(&chatReq, id).Error; err != nil {
		return nil
	}

	if err := DB().Save(&chatReq).Error; err != nil {
		return nil
	}

	return chatReq
}

// DeleteChatRequest delete ChatRequest
func (chatReq *ChatRequest) DeleteChatRequest() *ChatRequest {
	if err := DB().Find(&chatReq).Error; err != nil {
		return nil
	}

	if err := DB().Delete(&chatReq).Error; err != nil {
		return nil
	}

	return chatReq
}
