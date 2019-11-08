package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type PlaceType string

const (
	PlaceRoom PlaceType = "room"
	// PlaceRoom PlaceType = "room"
)

type Place struct {
	orm.ModelBase

	PlacName      string  `json:"plac_name" form:"name" gorm:"type:varchar(50)"`
	PlacDetail    string  `json:"plac_detail" form:"detail"`
	PlacActive    bool    `json:"plac_active" form:"active"`
	PlacType      string  `json:"plac_type" gorm:"type:varchar(50)"`
	PlacAmount    int     `json:"plac_amount"`
	PlacAccountID uint    `json:"plac_account_id"`
	Account       Account `json:"account" gorm:"ForeignKey:PlacAccountID"`
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
