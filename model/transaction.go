package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

type TranStatusType int

const (
	TranStatusReject         TranStatusType = -1
	TranStatusPanding        TranStatusType = 0
	TranStatusApproveBooking TranStatusType = 1
	TranStatusPaid           TranStatusType = 2
	TranStatus               TranStatusType = 3
)

type Transaction struct {
	orm.ModelBase
	TranStatus    TranStatusType `json:"tran_status" gorm:"type:varchar(50)"`
	TranRemark    string         `json:"tran_remark"`
	TranTotal     float64        `json:"tran_total"`
	AccountID     uint           `json:"account_id"`
	ChatChannelID uint           `json:"channel_id"`
	CustomerID    uint           `json:"customer_id"`
	TranLineID    string         `json:"tran_line_id" gorm:"type:varchar(50)"`
	Bookings      []Booking      `json:"bookings"`
	Customer      Customer       `json:"customer" gorm:"ForeignKey:CustomerID"`
	Account       Account        `json:"account" gorm:"ForeignKey:AccountID"`
	Payments      []Payment      `json:"payments"`
	ChatChannel   ChatChannel    `json:"chat_channel" gorm:"ForeignKey:ChatChannelID"`
}

type Report struct {
	orm.ModelBase
	TransactionID uint        `json:"transaction_id"`
	Transaction   Transaction `json:"transaction" gorm:"ForeingKey:TransactionID"`
	Detail        string      `json:"detail"`
}

func (tran *Transaction) Create() error {
	if err := DB().Create(&tran).Error; err != nil {
		return err
	}
	return nil
}

func GetTransactionList(accID uint) ([]*Transaction, error) {
	trans := []*Transaction{}
	if err := DB().Where("tran_account_id = ?", accID).Find(&trans).Error; err != nil {
		return nil, err
	}
	return trans, nil
}

func GetTransactionDetail(accID uint, id string) (*Transaction, error) {
	tran := Transaction{}
	if err := DB().Where("tran_account_id = ?", accID).Find(&tran, id).Error; err != nil {
		return nil, err
	}
	return &tran, nil
}

func (tran *Transaction) UpdateTransaction() error {
	if err := DB().Save(&tran).Error; err != nil {
		return err
	}
	return nil
}

func RemoveTransaction(id string, accID uint) (*Transaction, error) {
	tran := Transaction{}
	if err := DB().Find(&tran).Error; err != nil {

	}
	if err := DB().Delete(&tran).Error; err != nil {
		return nil, err
	}
	return &tran, nil
}
