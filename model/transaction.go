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
	TranStatus        TranStatusType `json:"tran_status" gorm:"type:varchar(50)"`
	TranTotal         int            `json:"tran_total"`
	TranAccountID     uint           `json:"tran_account_id"`
	TranChatChannelID uint           `json:"tranchat_channel_id"`
	TranCustomerID    uint           `json:"tran_customer_id"`
	Customer          Customer       `json:"customer" gorm:"ForeignKey:TranCustomerID"`
	Account           Account        `json:"account" gorm:"ForeignKey:TranAccountID"`
	ChatChannel       ChatChannel    `json:"chat_channel" gorm:"ForeignKey:TranChatChannelID"`
}

func (tran *Transaction) CreateTransaction() error {
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

func RemoveTransaction(id string) (*Transaction, error) {
	tran := Transaction{}
	if err := DB().Find(&tran).Error; err != nil {

	}
	if err := DB().Delete(&tran).Error; err != nil {
		return nil, err
	}
	return &tran, nil
}
