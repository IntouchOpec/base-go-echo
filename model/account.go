package model

import (
	"database/sql"

	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AccBookingType int8
type AccTransactionConfirmType string
type AccTypePayment string

const (
	AccBookingByTimeSlot AccBookingType = 0
	AccBookingByItem     AccBookingType = 1
	AccBookingByNow      AccBookingType = 2
)

const (
	AccTransactionMan  AccTransactionConfirmType = "man"
	AccTransactionAuto AccTransactionConfirmType = "auto"
)

const (
	AccTypePaymentBooking AccTypePayment = "booking"
	AccTypePaymentNow     AccTypePayment = "now"
)

// Account struct.
type Account struct {
	orm.ModelBase

	AccProjectID              string                    `form:"acc_project_id" json:"acc_project_id" grom:"type:varchar(100)"`
	AccAuthJSONFilePath       string                    `form:"acc_auth_json_file_path" json:"acc_auth_json_file_path" grom:"type:varchar(100)"`
	AccProjectIDDialogflow    string                    `json:"acc_project_id_dialogflow"`
	AccLang                   string                    `form:"acc_lang" json:"acc_lang" grom:"type:varchar(100)"`
	AccTimeZone               string                    `form:"acc_time_zone" json:"acc_time_zone" grom:"type:varchar(100)"`
	AccName                   string                    `form:"acc_name" json:"acc_name" gorm:"type:varchar(25)"`
	AccAmountPayment          int                       `form:"acc_amount_payment" json:"acc_amount_payment"`
	AccTransactionConfirmType AccTransactionConfirmType `form:"acc_transaction_confirm_type" json:"acc_transaction_confirm_type" gorm:"type:varchar(25)"`
	AccBookingType            string                    `form:"acc_booking_type" json:"acc_booking_type" gorm:"type:varchar(32)"`
	AccTypePayment            AccTypePayment            `form:"acc_type_payment" json:"acc_type_payment" gorm:"type:varchar(10)"`
	Settings                  []*Setting                `json:"settings" gorm:"many2many:account_setting"`
	ChatChannels              []*ChatChannel            `json:"chat_channels"`
}

type AccountLine struct {
	ID                        uint
	AccProjectID              string                    `json:"acc_project_id"`
	AccAuthJSONFilePath       string                    `json:"acc_auth_json_file_path"`
	AccProjectIDDialogflow    string                    `json:"acc_project_id_dialogflow"`
	AccAmountPayment          int                       `json:"acc_amount_payment"`
	AccLang                   string                    `json:"acc_lang"`
	AccTimeZone               string                    `json:"acc_time_zone"`
	AccName                   string                    `json:"acc_name"`
	AccTransactionConfirmType AccTransactionConfirmType `json:"acc_transaction_confirm_type"`
	AccBookingType            string                    `json:"acc_booking_type"`
	AccTypePayment            AccTypePayment            `json:"acc_type_payment"`
	ChaChannelSecret          string                    `json:"cha_channel_secret"`
	ChaChannelAccessToken     string                    `json:"cha_channel_access_token"`
	ChaName                   string                    `json:"cha_name"`
	ChaAddress                string                    `json:"cha_address"`
	Settings                  map[string]string
	// Latitude                  string                    `json:"latitude"`
	// Longitude                 string                    `json:"longitude"`
}

func AccountLineGet(name, channelID string, db *sql.DB) *AccountLine {
	var al AccountLine
	row := db.QueryRow(`
	SELECT 
		ac.id,
		acc_project_id,
		acc_auth_json_file_path ,
		acc_project_id_dialogflow ,
		acc_amount_payment ,
		acc_lang ,
		acc_time_zone ,
		acc_name ,
		acc_transaction_confirm_type ,
		acc_booking_type ,
		acc_type_payment ,
		cha_channel_secret ,
		cha_channel_access_token,
		cha_name,
		cha_address
	FROM accounts AS ac 
	INNER JOIN chat_channels AS cc ON cc.account_id = ac.id AND cc.deleted_at IS NULL AND cha_channel_id = $1
	WHERE ac.deleted_at IS NULL AND acc_name = $2 `,
		channelID, name)
	err := row.Scan(&al.ID,
		&al.AccProjectID,
		&al.AccAuthJSONFilePath,
		&al.AccProjectIDDialogflow,
		&al.AccAmountPayment,
		&al.AccLang,
		&al.AccTimeZone,
		&al.AccName,
		&al.AccTransactionConfirmType,
		&al.AccBookingType,
		&al.AccTypePayment,
		&al.ChaChannelSecret,
		&al.ChaChannelAccessToken,
		&al.ChaName,
		&al.ChaAddress)
	if err != nil {
		return nil
	}
	rows, err := db.Query(`SELECT st.name,st.value FROM settings AS st
	INNER JOIN setting_chat_channel AS scc ON st.id = scc.setting_id
	INNER JOIN chat_channels AS cc ON scc.chat_channel_id = cc.id AND cc.id = $1 AND cc.deleted_at IS NULL
	WHERE st.deleted_at IS NULL`, al.ID)
	if err != nil {
		return nil
	}
	al.Settings = make(map[string]string)
	for rows.Next() {
		var tempSetting Setting
		_ = rows.Scan(&tempSetting.Name, &tempSetting.Value)
		al.Settings[tempSetting.Name] = tempSetting.Value
	}
	return &al
}

func (account *Account) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New())
	return nil
}

func GetAccountByName(name string) bool {
	account := Account{}

	if err := DB().Where("name = ?", name).Find(&account).Error; err != nil {
		return false
	}

	return true
}

// GetAccount query account list.
func GetAccount() []*Account {
	accounts := []*Account{}

	DB().Find(&accounts)

	return accounts
}

// GetAccountByID find account by id.
func (acc *Account) GetAccountByID(id string) *Account {
	account := Account{}

	if err := DB().Find(&account, id).Error; err != nil {
		return nil
	}

	return &account
}

// CreateAccount is function create accout.
func (acc *Account) CreateAccount() *Account {
	newDb, err := newDB()

	if err != nil {
		return nil
	}

	if err := newDb.Create(&acc).Error; err != nil {
		return nil
	}
	return acc
}

// UpdateAccount edit account soucre.
func (acc *Account) UpdateAccount() *Account {

	if err := DB().Save(&acc).Error; err != nil {
		return nil
	}

	return acc
}

func (acc *Account) RemoveAccount(id string) *Account {
	if err := DB().Find(&acc, id).Error; err != nil {
		return nil
	}
	if err := DB().Delete(&acc, id).Error; err != nil {
		return nil
	}
	return acc
}
