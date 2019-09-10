package model

import "github.com/hb-go/gorm"

// EventLog log event line OA.
type EventLog struct {
	gorm.Model
	ReplyToken string  `json:"reply_token" gorm:"type:varchar(255)"`
	Type       string  `json:"type" gorm:"type:varchar(10)"`
	LineID     string  `json:"line_id" gorm:"type:varchar(25)"`
	AccountID  uint    `form:"account_id" json:"account_id" gorm:"not null;"`
	Account    Account `gorm:"ForeignKey:id"`
}

// SaveEventLog is function create EventLog.
func (eventlog *EventLog) SaveEventLog() *EventLog {
	if err := DB().Create(&eventlog).Error; err != nil {
		return nil
	}
	return eventlog
}
