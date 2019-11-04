package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type Place struct {
	orm.ModelBase

	Name   string `json:"name" gorm:"type:varchar(50)"`
	Detail string `json:"detail"`
}

func (pla *Place) CreatePlace() (*Place, error) {
	if err := DB().Create(&pla).Error; err != nil {
		return nil, err
	}
	return pla, nil
}

func (pla *Place) Update(id string) (*Place, error) {
	if err := DB().Find(&pla).Error; err != nil {
		return nil, err
	}
	if err := DB().Save(&pla).Error; err != nil {
		return nil, err
	}
	return pla, nil
}

func GetPlaceList(AccountID uint) ([]*Place, error) {
	places := []*Place{}
	if err := DB().Where("AccountID = ?", AccountID).Find(&places).Error; err != nil {
		return nil, err
	}
	return places, nil
}

func GetPlaceDetail(id string) (*Place, error) {
	place := Place{}
	if err := DB().Find(&place, id).Error; err != nil {
		return nil, err
	}
	return &place, nil
}
