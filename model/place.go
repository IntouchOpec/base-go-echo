package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type PlaceType string

const (
	PlaceRoom PlaceType = "room"
	// PlaceRoom PlaceType = "room"
)

const formatDate string = "2006-01-02 15:04"

type Place struct {
	orm.ModelBase

	PlacName     string         `json:"plac_name" form:"name" gorm:"type:varchar(50)"`
	PlacDetail   string         `json:"plac_detail" form:"detail"`
	PlacActive   bool           `json:"plac_active" form:"active"`
	PlacType     string         `form:"type" json:"plac_type" gorm:"type:varchar(50)"`
	PlacAmount   int            `form:"amount" json:"plac_amount"`
	PlacImage    string         `form:"image" json:"plac_image" gorm:"type:varchar(255)"`
	ChatChannels []*ChatChannel `json:"chat_channels" gorm:"many2many:place_chat_channel"`
	AccountID    uint           `json:"account_id"`
	Account      Account        `json:"account" gorm:"ForeignKey:AccountID"`
	MasterPlaces []*MasterPlace `json:"master_places"`
}

type MasterPlace struct {
	orm.ModelBase

	PlaceID    uint      `json:"place_id"`
	MPlaAmount int       `json:"mpla_amount"`
	AccountID  uint      `json:"account_id"`
	MPlaDay    time.Time `json:"mpla_day"`
	MPlaFrom   time.Time `json:"mpla_from"`
	MPlaTo     time.Time `json:"mpla_to"`
	Place      Place     `json:"place" gorm:"ForeignKey:PlaceID"`
	Account    Account   `json:"account" gorm:"ForeignKey:AccountID"`
}

func (pla *Place) CreatePlace() error {
	if err := DB().Create(&pla).Error; err != nil {
		return err
	}
	return nil
}

func (pla *Place) Update() error {
	if err := DB().Save(&pla).Error; err != nil {
		return err
	}
	return nil
}

func GetPlaceList(AccountID uint) ([]*Place, error) {
	places := []*Place{}
	if err := DB().Where("AccountID = ?", AccountID).Find(&places).Error; err != nil {
		return nil, err
	}
	return places, nil
}

func GetPlaceDetail(id string, AccountID uint) (*Place, error) {
	place := Place{}
	if err := DB().Where("account_id = ?", AccountID).Find(&place, id).Error; err != nil {
		return nil, err
	}
	return &place, nil
}

func DeletePlaceByID(id string) (*Place, error) {
	place := Place{}

	if err := DB().Find(&place, id).Error; err != nil {
		return nil, err
	}

	if err := DB().Delete(&place).Error; err != nil {
		return nil, err
	}

	return &place, nil
}

func (mpla *MasterPlace) CreateMasterPlace() (*MasterPlace, error) {
	if err := DB().Create(&mpla).Error; err != nil {
		return nil, err
	}
	return nil, nil
}

func (mpla *MasterPlace) UpdateMasterPlace() error {
	if err := DB().Save(&mpla).Error; err != nil {
		return err
	}
	return nil
}

func GetMasterPlace(from, to, day time.Time) (*MasterPlace, error) {
	startDate, _ := time.Parse(formatDate, "2019-07-09")
	endDate, _ := time.Parse(formatDate, "2019-07-15")

	mpla := MasterPlace{}
	if err := DB().Where("mpla_day = ? and mpla_to >= ? and mpla_to <= ? and mpla_from >= ? and mpla_from <= ?",
		day, startDate, startDate, endDate, endDate).Find(&mpla).Error; err != nil {
		return nil, err
	}
	return &mpla, nil
}

func (pla *Place) RemoveImage() error {
	pla.PlacImage = ""
	if err := DB().Save(&pla).Error; err != nil {
		return err
	}
	return nil
}
