package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

// EventLog log event line OA.
type EventLog struct {
	orm.ModelBase
	EvenReplyToken    string       `json:"even_reply_token" gorm:"type:varchar(255)"`
	EvenType          string       `json:"even_type" gorm:"type:varchar(10)"`
	EvenLineID        string       `json:"even_line_id" gorm:"type:varchar(255)"`
	EvenChatChannelID uint         `form:"even_chat_channel_id" json:"chat_channel_id" gorm:"not null;"`
	EvenText          string       `json:"even_text"`
	EvenCustomerID    uint         `json:"even_customer_id"`
	ChatChannel       *ChatChannel `json:"chat_channel" gorm:"ForeignKey:EvenChatChannelID"`
	Customer          *Customer    `json:"customer" gorm:"ForeignKey:EvenCustomerID"`
}

// SaveEventLog is function create EventLog.
func (eventlog *EventLog) SaveEventLog() *EventLog {
	if err := DB().Create(&eventlog).Error; err != nil {
		return nil
	}
	return eventlog
}

// GetEventLog
func GetEventLog(page int, size int, chatChannelID int) *[]EventLog {
	eventLogs := []EventLog{}
	DB().Where("chat_channel_id = ? ", chatChannelID).Offset((page - 1) * size).Limit(size).Find(&eventLogs)
	return &eventLogs
}

// GetAllEventLog
func GetAllEventLog(page int, size int) *[]EventLog {
	eventLogs := []EventLog{}
	DB().Offset((page - 1) * size).Limit(size).Find(&eventLogs)
	return &eventLogs
}
