package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

// ActionStatus  action status
type ActionStatus string

const (
	// StatusSuccess when status success
	StatusSuccess ActionStatus = "success"
	// StatusFail when status error
	StatusFail ActionStatus = "fail"
)

type ActionChannel string

const (
	ActionChannelFacebook ActionChannel = "facebook"
	ActionChannelLine     ActionChannel = "line"
	ActionChannelWeb      ActionChannel = "web"
	ActionChannelAPI      ActionChannel = "api"
)

// Action source action all
type ActionLog struct {
	orm.ModelBase

	ActName          string        `json:name" gorm:"type:varchar(25)"`
	ActStatus        ActionStatus  `json:"act_status" gorm:"type:varchar(10)"`
	ActChannel       ActionChannel `json:"act_channel" gorm:"type:varchar(10)"`
	ActUserID        string        `json:"act_user_id" gorm:"type:varchar(55)"`
	ActChatChannelID uint          `json:"act_chat_channel_id"`
	ChatChannel      ChatChannel   `json:"chat_channel" gorm:"ForeignKey:ServiceID;"`
	CustomerID       uint          `json:"customer_id"`
	Customer         *Customer     `json:"customer" gorm:"ForeignKey:CustomerID"`
}

// CreateAction create action record
func (act *ActionLog) CreateAction() (*ActionLog, error) {
	// if err := db.Create(&act).Error; err != nil {
	// 	return nil, err
	// }
	return act, nil
}
