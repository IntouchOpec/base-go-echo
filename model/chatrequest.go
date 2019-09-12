package model

import "github.com/hb-go/gorm"

// ChatRequest log castomer when customer send something to line OA.
type ChatRequest struct {
	gorm.Model
	// ID           uint       `gorm:"primary_key" json:"id"`
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
func (cha *ChatRequest) GetChatRequest() *ChatRequest {
	DB().Find(&cha)
	return cha
}

// EditChatRequest update ChatRequest.
func (cha *ChatRequest) EditChatRequest(id int) *ChatRequest {

	if err := DB().Find(&cha, id).Error; err != nil {
		return nil
	}

	if err := DB().Save(&cha).Error; err != nil {
		return nil
	}

	return cha
}

// DeleteChatRequest delete ChatRequest
func (cha *ChatRequest) DeleteChatRequest() *ChatRequest {
	if err := DB().Find(&cha).Error; err != nil {
		return nil
	}

	if err := DB().Delete(&cha).Error; err != nil {
		return nil
	}

	return cha
}
