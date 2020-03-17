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

	ActName       string        `json:"act_name" gorm:"type:varchar(25)"`
	ActStatus     ActionStatus  `json:"act_status" gorm:"type:varchar(10)"`
	ActChannel    ActionChannel `json:"act_channel" gorm:"type:varchar(10)"`
	ActUserID     string        `json:"act_user_id" gorm:"type:varchar(55)"`
	ChatChannelID uint          `json:"chat_channel_id"`
	CustomerID    uint          `json:"act_customer_id"`
	ChatChannel   ChatChannel   `json:"chat_channel" gorm:"ForeignKey:ChatChannelID;"`
	Customer      *Customer     `json:"customer" gorm:"ForeignKey:CustomerID"`
}

// CreateAction create action record
func (act *ActionLog) CreateAction() (*ActionLog, error) {
	if err := db.Create(&act).Error; err != nil {
		return nil, err
	}
	return act, nil
}

// func GetActionList(accID string) ([]*ActionLog, error) {
// 	acts := []*ActionLog{}
// 	db := DB().Where("chat_channel_id = ?", accID)
// 	if err := Cache(db).Preload("Customer").Find(&acts).Error; err != nil {
// 		return nil, err
// 	}
// 	return acts, nil
// }
