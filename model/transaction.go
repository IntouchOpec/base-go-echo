package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
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

type TranState string

const (
	TranStateNow         TranState = "now"
	TranStateAppointment TranState = "appointment"
)

type Transaction struct {
	orm.ModelBase

	TranState        TranState      `json:"tran_state" gorm:"type:varchar(15)"`
	TranDocumentCode string         `json:"tran_document_code"`
	TranStatus       TranStatusType `json:"tran_status" gorm:"type:varchar(50)"`
	TranRemark       string         `json:"tran_remark"`
	TranTotal        float64        `json:"tran_total"`
	AccountID        uint           `json:"account_id"`
	ChatChannelID    uint           `json:"channel_id"`
	CustomerID       uint           `json:"customer_id"`
	TranLineID       string         `json:"tran_line_id" gorm:"type:varchar(50)"`
	Customer         Customer       `json:"customer" gorm:"ForeignKey:CustomerID"`
	Account          Account        `json:"account" gorm:"ForeignKey:AccountID"`
	Bookings         []*Booking     `json:"bookings"`
	Payments         []*Payment     `json:"payments"`
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

func (tran *Transaction) GetLastTransactionSql(db *sql.Tx) error {
	row := db.QueryRow(`
	SELECT 
		tran_document_code
	FROM transactions 
	WHERE deleted_at is null AND account_id = $1
	ORDER BY tran_document_code DESC LIMIT 1 `, tran.AccountID)

	row.Scan(&tran.TranDocumentCode)
	// if err != nil {
	// 	fmt.Println("==-=-=-=")
	// 	return err
	// }
	tran.TranDocumentCode, _ = genTranDocumentCode(tran.TranDocumentCode)
	return nil
}

func genTranDocumentCode(lastTranDocumentCode string) (string, error) {
	fmt.Println(lastTranDocumentCode, "lpadlpskdp")
	now := time.Now()
	day := fmt.Sprintf("%d", now.Day())
	month := fmt.Sprintf("%d", now.Month())
	var code string
	if len(day) == 1 {
		day = fmt.Sprintf("0%s", day)
	}
	if len(month) == 1 {
		month = fmt.Sprintf("0%s", month)
	}
	dayCode := fmt.Sprintf(dayFormatDoc, now.Year(), month, day)
	if lastTranDocumentCode == "" {
		return fmt.Sprintf(formatDocumentCode, dayCode+fmt.Sprintf(codeFormatDoc, 1)), errors.New("")
	}
	if lastTranDocumentCode[1:7] == dayCode[0:6] {
		i, err := strconv.Atoi(lastTranDocumentCode[9:])
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
	return nil
}

func (tran *Transaction) LineBooking(db *gorm.DB) error {
	if err := db.Create(&tran).Error; err != nil {
		return err
	}
	return nil
}

func (tran *Transaction) LineBookingServiceNow(db *gorm.DB, book Booking) error {
	if err := db.Model(&tran).Association("Bookings").Append(&book).Error; err != nil {
		return err
	}
	return nil
}

func (tran *Transaction) LineBookingServiceAppointment(db *gorm.DB, book Booking) error {
	if err := db.Model(&tran).Association("Bookings").Append(&book).Error; err != nil {
		return err
	}
	return nil
}

func (serI *ServiceItem) GetServiceItemPreEmploAndPla(serID string, db *gorm.DB) error {
	db.Preload("Service", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Employees").Preload("Places")
	}).Find(&serI, serID)
	return nil
}

func (tran *Transaction) CreateSql(db *sql.Tx) error {
	if err := tran.GetLastTransactionSql(db); err != nil {
		fmt.Println(err)
		return err
	}
	stmt, err := db.Prepare(`INSERT INTO transactions 
			(tran_state, tran_document_code, tran_status, tran_remark, tran_total, account_id, chat_channel_id, customer_id, tran_line_id, created_at) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9, now()) 
			RETURNING id`)
	if err != nil {
		fmt.Println(err, "12313")
		return err
	}
	defer stmt.Close()
	result := stmt.QueryRow(tran.TranState, tran.TranDocumentCode, tran.TranStatus, tran.TranRemark, tran.TranTotal, tran.AccountID, tran.ChatChannelID, tran.CustomerID, tran.TranLineID)
	if err := result.Scan(&tran.ID); err != nil {
		fmt.Println(err, ":e13123")
		return err
	}
	return nil
}

func (tran *Transaction) UpdateStatus(tx *sql.Tx) (int64, error) {
	stmt, err := tx.Prepare("UPDATE transactions SET tran_status = $1, updated_at = now() WHERE id = $2")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(tran.TranStatus, tran.ID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	rEf, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rEf, nil
}

type Receipt struct {
	Name  string
	Price string
}

func (tran *Transaction) GetReceipt(db *sql.DB) ([]Receipt, error) {
	var reList []Receipt
	rows, err := db.Query(`
	SELECT 
		tr.id AS tran_id, tran_total,tran_document_code,
		ser_name, ss_time,ss_price,
		pac_name, pac_time, pac_price
	FROM transactions AS tr 
	INNER JOIN bookings AS bo ON tr.id = bo.transaction_id AND bo.deleted_at IS NULL 
	LEFT JOIN booking_service_items AS bsi ON bo.id = bsi.booking_id 
	LEFT JOIN service_items AS si ON  si.id = bsi.service_item_id AND si.deleted_at IS NULL 
	LEFT JOIN services AS s ON s.id = si.service_id AND s.deleted_at IS NULL 
	LEFT JOIN booking_packages AS bp ON bo.id = bp.booking_id  
	LEFT JOIN packages AS pa ON pa.id = bp.package_id AND pa.deleted_at IS NULL 
	WHERE  tr.deleted_at IS NULL AND tr.id = $1 `, tran.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var re Receipt
		var SSName string
		var SSTime time.Duration
		var SSPrice float64
		var PName string
		var PTime time.Duration
		var PPrice float64
		rows.Scan(&tran.ID, &tran.TranTotal, &tran.TranDocumentCode, &SSName, &SSTime, &SSPrice, &PName, &PTime, &PPrice)
		// if PPrice != 0 {
		// 	re.Name = fmt.Sprintf("%s-%s", PName, PTime.String())
		// 	re.Price = fmt.Sprintf("%.0f", PPrice)
		// } else {
		re.Name = fmt.Sprintf("%s-%s", SSName, SSTime.String())
		re.Price = fmt.Sprintf("%.0f", SSPrice)

		// }
		reList = append(reList, re)
	}
	return reList, nil
}
