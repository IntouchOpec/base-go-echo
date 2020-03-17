package model

import (
	"database/sql"
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

	PayStatus PayStatus `json:"pay_status"`
	// PayChannel    PaymentChannelType `json:"pay_channel" gorm:"type:varchar(25)"`
	PayAmount     float64     `json:"pay_amount"`
	PayImage      string      `json:"pay_image" gorm:"type:varchar(50)"`
	PayType       PayType     `json:"pay_type" gorm:"type:varchar(25)"`
	PayAt         time.Time   `json:"pay_at"`
	PayTracking   string      `json:"pay_tracking" gorm:"type:varchar(50)"`
	TransactionID uint        `json:"transaction_id" gorm:"ForeignKey:ChatChannelID"`
	Transaction   Transaction `json:"transaction" gorm:"ForeignKey:TransactionID"`
}

func (p *Payment) Create(tx *sql.Tx) error {
	stmt := "INSERT INTO payments (pay_status ,pay_amount ,pay_image ,pay_type ,pay_at ,pay_tracking,transaction_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	_, err := tx.Exec(stmt, p.PayStatus, p.PayAmount, p.PayImage, p.PayType, p.PayAt, p.PayTracking, p.TransactionID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
