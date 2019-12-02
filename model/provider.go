package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
)

type Provider struct {
	orm.ModelBase

	ProvName         string             `form:"prov_name" json:"prov_name" gorm:"type:varchar(25)"`
	ProvDetail       string             `form:"prov_detail" json:"prov_detail"`
	ProvLineID       string             `form:"prov_line_id" json:"prov_line_id" gorm:"type:varchar(50)"`
	ProvImage        string             `form:"image" json:"prov_image" gorm:"type:varchar(255)"`
	AccountID        uint               `form:"account_id" json:"account_id"`
	ProviderServices []*ProviderService `json:"provider_services" `
	Account          Account            `json:"account" gorm:"ForeignKey:AccountID"`
}

func (prov *Provider) CreateProvider() error {
	if err := DB().Create(&prov).Error; err != nil {
		return err
	}

	return nil
}

func (prov *Provider) UpdateProvider() error {
	if err := DB().Save(&prov).Error; err != nil {
		return err
	}

	return nil
}

func GetProviderList(accID uint) ([]*Provider, error) {
	provs := []*Provider{}
	if err := DB().Where("account_id = ?", accID).Find(&provs).Error; err != nil {
		return nil, err
	}

	return provs, nil
}

func GetProviderDetail(id string, accID uint) (*Provider, error) {
	prov := Provider{}
	if err := DB().Preload("ProviderServices", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Service")
	}).Where("account_id = ?", accID).Find(&prov, id).Error; err != nil {
		return nil, err
	}

	return &prov, nil
}

func GetProviderServiceTimeSlotList(id string, accID uint) (*Provider, error) {
	prov := Provider{}
	if err := DB().Preload("ProviderServices", func(db *gorm.DB) *gorm.DB {
		return db.Preload("TimeSlots").Preload("Service")
	}).Where("account_id = ?", accID).Find(&prov, id).Error; err != nil {
		return nil, err
	}

	return &prov, nil
}

func RemoveProvider(id string) (*Provider, error) {
	prov := Provider{}
	db := DB()
	if err := db.Find(&prov, id).Error; err != nil {
		return nil, err
	}

	if err := db.Delete(&prov).Error; err != nil {
		return nil, err
	}
	return &prov, nil
}
