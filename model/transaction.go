package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

type TranStatusType int

const (
	StatusReject         TranStatusType = -1
	StatusPanding        TranStatusType = 0
	StatusApproveBooking TranStatusType = 1
	StatusPaid           TranStatusType = 2
	Status               TranStatusType = 3
)

type Transaction struct {
	orm.ModelBase
	TranStatus    TranStatusType `json:"tran_status"`
	ChatChannelID uint           `json:"chat_channel_id"`
	ChatChannel   ChatChannel    `json:"chat_channel" gorm:"ForeignKey:ChatChannelID"`
}
