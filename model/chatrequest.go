package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

// ChatRequest log castomer when customer send something to line OA.
type ChatRequest struct {
	orm.ModelBase

	ReqLineID       string     `json:"req_line_id" gorm:"type:varchar(25)"`
	ReqMessage      string     `json:"req_message" gorm:"type:varchar(50)"`
	ReqReplyToken   string     `json:"req_reply_token" gorm:"type:varchar(255)"`
	ReqMessageType  string     `json:"req_message_type" gorm:"type:varchar(25)"`
	ReqChatAnswerID uint       `form:"req_chat_answer_id" json:"chat_answer_id" gorm:"not null;"`
	ChatAnswer      ChatAnswer `json:"chat_answer" gorm:"ForeignKey:ReqChatAnswerID"`
	ReqAccountID    uint       `form:"req_account_id" json:"account_id" gorm:"not null;"`
	Account         Account    `json:"account" gorm:"ForeignKey:ReqAccountID"`
}

// SaveChatRequest is function create chat answer.
func (chatReq *ChatRequest) SaveChatRequest() error {
	if err := DB().Create(&chatReq).Error; err != nil {
		return err
	}
	return nil
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

func RemoveChatRequest(id string, accID uint) (*ChatRequest, error) {
	chatReq := ChatRequest{}
	if err := DB().Where("req_account_id = ?", accID).Find(&chatReq, id).Error; err != nil {
		return nil, err
	}

	if err := DB().Delete(&chatReq).Error; err != nil {
		return nil, err
	}

	return &chatReq, nil
}
