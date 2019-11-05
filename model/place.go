package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type Place struct {
	orm.ModelBase

	Name   string `json:"name" form:"name" gorm:"type:varchar(50)"`
	Detail string `json:"detail" form:"detail"`
	Active bool   `json:"active" form:"active"`
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
