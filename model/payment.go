package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type PaymentChannelType string

type PayType string

const (
	PayTypeOmise PayType = "omise"
)

type PayStatus int

const (
	PayStatusReject  PayStatus = -1
	PayStatusPanding PayStatus = 0
	PayStatusSuccess PayStatus = 1
)

type Payment struct {
	orm.ModelBase

	PayStatus      PayStatus          `json:"pay_status"`
	PaymentChannel PaymentChannelType `json:"payment_channel" gorm:"type:varchar(25)"`
	PayAmount      float64            `json:"pay_amount"`
	PayImage       string             `json:"pay_image" gorm:"type:varchar(50)"`
	PayType        PayType            `json:"pay_type" gorm:"type:varchar(25)"`
	PayAt          time.Time          `json:"pay_at"`
	PayTracking    string             `json:"pay_tracking" gorm:"type:varchar(50)"`
	TransactionID  uint               `json:"transaction_id" gorm:"ForeignKey:ChatChannelID"`
	Transaction    Transaction        `json:"transaction" gorm:"ForeignKey:TransactionID"`
}
