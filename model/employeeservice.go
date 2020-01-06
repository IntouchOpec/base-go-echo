package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type EmployeeService struct {
	orm.ModelBase
	PSPrice    float64     `form:"price" json:"ps_price"`
	EmployeeID uint        `json:"employee_id"`
	ServiceID  uint        `form:"service_id" json:"service_id"`
	Employee   Employee    `json:"employee" gorm:"ForeignKey:EmployeeID"`
	Service    Service     `json:"service" gorm:"ForeignKey:ServiceID"`
	TimeSlots  []*TimeSlot `json:"time_slots"`
	Bookings   []*Booking  `json:"bookings"`
}

func GetEmployeeServiceDetail(id string) (*EmployeeService, error) {
	employeeService := EmployeeService{}
	if err := DB().Find(&employeeService).Error; err != nil {
		return nil, err
	}
	return &employeeService, nil
}

func GetEmployeeServiceDetailByFild(Params ...interface{}) (*EmployeeService, error) {
	employeeService := EmployeeService{}
	if err := DB().Find(&employeeService).Error; err != nil {
		return nil, err
	}
	return &employeeService, nil
}

func (prov *EmployeeService) CreateEmployeeService() error {
	if err := DB().Create(&prov).Error; err != nil {
		return err
	}
	return nil
}

func (prov *EmployeeService) UpdateEmployeeService() error {
	if err := DB().Save(&prov).Error; err != nil {
		return err
	}
	return nil
}

func RemoveEmployeeService(id string) (*EmployeeService, error) {
	prov := EmployeeService{}
	if err := DB().Find(&prov, id).Error; err != nil {
		return nil, err
	}

	if err := DB().Delete(&prov).Error; err != nil {
		return nil, err
	}
	return &prov, nil
}
