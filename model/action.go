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

// ActionType is
type ActionType string

const (
	// TypeFacebook
	TypeActionFacebook ActionType = "facebook"
	// TypeLine
	TypeActionLine ActionType = "line"
	// TypeActionWeb
	TypeActionWeb ActionType = "web"
	// TypeActionAPI
	TypeActionAPI ActionType = "api"
)

// Action source action all
type ActionLog struct {
	orm.ModelBase

	Name          string       `json:"name" gorm:"type:varchar(25)"`
	Status        ActionStatus `json:"status" gorm:"type:varchar(10)"`
	Type          ActionType   `json:"type" gorm:"type:varchar(10)"`
	UserID        string       `json:"user_id" gorm:"type:varchar(55)"`
	ChatChannelID uint         `json:"chat_channel_id"`
	ChatChannel   ChatChannel  `json:"chat_channel" gorm:"ForeignKey:ProductID;"`
	CustomerID    uint         `json:"customer_id"`
	Customer      *Customer    `json:"customer" gorm:"ForeignKey:CustomerID"`
}

// CreateAction create action record
func (act *ActionLog) CreateAction() (*ActionLog, error) {
	// if err := db.Create(&act).Error; err != nil {
	// 	return nil, err
	// }
	return act, nil
}
