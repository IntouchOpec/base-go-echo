package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

const (
	formatDocumentCode string = "T%s"
	dayFormatDoc       string = "%d%s%s" // YYYYMMDD
	codeFormatDoc      string = "%06d"
)

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
	TranDocumentCode string         `json:"tran_doccument_code"`
	TranStatus       TranStatusType `json:"tran_status" gorm:"type:varchar(50)"`
	TranRemark       string         `json:"tran_remark"`
	TranTotal        float64        `json:"tran_total"`
	AccountID        uint           `json:"account_id"`
	ChatChannelID    uint           `json:"channel_id"`
	CustomerID       uint           `json:"customer_id"`
	TranLineID       string         `json:"tran_line_id" gorm:"type:varchar(50)"`
	Bookings         []Booking      `json:"bookings"`
	Customer         Customer       `json:"customer" gorm:"ForeignKey:CustomerID"`
	Account          Account        `json:"account" gorm:"ForeignKey:AccountID"`
	Payments         []Payment      `json:"payments"`
	ChatChannel      ChatChannel    `json:"chat_channel" gorm:"ForeignKey:ChatChannelID"`
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

func (tran *Transaction) GetLastTransaction() (string, error) {
	tranLast := Transaction{}
	DB().Where("account_id = ?", tran.AccountID).Last(&tranLast)
	return tranLast.TranDocumentCode, nil
}

func genTranDocumentCode(lastTranDocumentCode string) (string, error) {
	now := time.Now()
	day := fmt.Sprintf("%d", now.Day())
	mouth := fmt.Sprintf("%d", now.Month())
	var code string
	if len(day) == 1 {
		day = fmt.Sprintf("0%s", day)
	}
	if len(mouth) == 1 {
		mouth = fmt.Sprintf("0%s", mouth)
	}
	dayCode := fmt.Sprintf(dayFormatDoc, now.Year(), mouth, day)
	if lastTranDocumentCode == "" {
		return fmt.Sprintf(formatDocumentCode, dayCode+fmt.Sprintf(codeFormatDoc, 1)), errors.New("")
	}
	if lastTranDocumentCode[1:7] == dayCode[0:6] {
		i, err := strconv.Atoi(lastTranDocumentCode[8:])
		if err != nil {
			return "", err
		}
		code = fmt.Sprintf(codeFormatDoc, i+1)
	} else {
		code = fmt.Sprintf(codeFormatDoc, 1)
	}
	docCode := fmt.Sprintf(formatDocumentCode, dayCode+code)
	return docCode, nil
}

func (tran *Transaction) BeforeSave() (err error) {
	lastID, err := tran.GetLastTransaction()
	if err != nil {
		return errors.New("")
	}
	tran.TranDocumentCode, err = genTranDocumentCode(lastID)
	return
}
