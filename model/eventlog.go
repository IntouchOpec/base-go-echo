package model

import "github.com/hb-go/gorm"

// EventLog log event line OA.
type EventLog struct {
	gorm.Model
	ReplyToken    string       `json:"reply_token" gorm:"type:varchar(255)"`
	Type          string       `json:"type" gorm:"type:varchar(10)"`
	LineID        string       `json:"line_id" gorm:"type:varchar(255)"`
	ChatChannelID uint         `form:"chat_channel_id" json:"chat_channel_id" gorm:"not null;"`
	ChatChannel   *ChatChannel `gorm:"ForeignKey:ChatChannelID"`
	CustomerID    uint         `json:"customer_id"`
	Customer      *Customer    `json:"customer" gorm:"ForeignKey:CustomerID"`
	Text          string       `json:"text"`
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
