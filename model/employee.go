package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
)

type Employee struct {
	orm.ModelBase

	ProvName         string             `form:"prov_name" json:"prov_name" gorm:"type:varchar(25)"`
	ProvDetail       string             `form:"prov_detail" json:"prov_detail"`
	ProvLineID       string             `form:"prov_line_id" json:"prov_line_id" gorm:"type:varchar(50)"`
	ProvImage        string             `form:"image" json:"prov_image" gorm:"type:varchar(255)"`
	AccountID        uint               `form:"account_id" json:"account_id"`
	EmployeeServices []*EmployeeService `json:"employee_services" `
	Account          Account            `json:"account" gorm:"ForeignKey:AccountID"`
}

func (prov *Employee) CreateEmployee() error {
	if err := DB().Create(&prov).Error; err != nil {
		return err
	}

	return nil
}

func (prov *Employee) UpdateEmployee() error {
	if err := DB().Save(&prov).Error; err != nil {
		return err
	}

	return nil
}

func GetEmployeeList(accID uint) ([]*Employee, error) {
	provs := []*Employee{}
	if err := DB().Where("account_id = ?", accID).Find(&provs).Error; err != nil {
		return nil, err
	}

	return provs, nil
}

func GetEmployeeDetail(id string, accID uint) (*Employee, error) {
	prov := Employee{}
	if err := DB().Preload("EmployeeServices", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Service")
	}).Where("account_id = ?", accID).Find(&prov, id).Error; err != nil {
		return nil, err
	}

	return &prov, nil
}

func GetEmployeeServiceTimeSlotList(id string, accID uint) (*Employee, error) {
	prov := Employee{}
	if err := DB().Preload("EmployeeServices", func(db *gorm.DB) *gorm.DB {
		return db.Preload("TimeSlots").Preload("Service")
	}).Where("account_id = ?", accID).Find(&prov, id).Error; err != nil {
		return nil, err
	}

	return &prov, nil
}

func RemoveEmployee(id string) (*Employee, error) {
	prov := Employee{}
	db := DB()
	if err := db.Find(&prov, id).Error; err != nil {
		return nil, err
	}

	if err := db.Delete(&prov).Error; err != nil {
		return nil, err
	}
	return &prov, nil
}

func (pro *Employee) RemoveImage() error {
	pro.ProvImage = ""

	if err := DB().Save(&pro).Error; err != nil {
		return err
	}

	return nil
}
