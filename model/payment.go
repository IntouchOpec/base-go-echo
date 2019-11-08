package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type PaymentChannelType string

const ()

type Payment struct {
	orm.ModelBase

	PaymentChannel PaymentChannelType `json:"payment_channel" gorm:"type:varchar(25)"`
	PayAmount      float64            `json:"pay_amount"`
	PayImage       string             `json:"pay_image" gorm:"type:varchar(50)"`
	PayType        string             `json:"pay_type" gorm:"type:varchar(25)"`
	PayAt          time.Time          `json:"pay_at"`
	PayTracking    string             `json:"pay_tracking" gorm:"type:varchar(50)"`
}
