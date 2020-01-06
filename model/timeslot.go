package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
)

type TimeSlot struct {
	orm.ModelBase

	TimeStart         string          `json:"time_start" gorm:"type:varchar(10)"`
	TimeEnd           string          `json:"time_end" gorm:"type:varchar(10)"`
	TimeDay           int             `json:"time_day"`
	TimeAmount        int             `json:"time_amount"`
	TimeActive        bool            `json:"time_active" gorm:"default:true"`
	EmployeeServiceID uint            `json:"employee_service_id"`
	AccountID         uint            `json:"account_id"`
	EmployeeService   EmployeeService `json:"employee_service" gorm:"ForeignKey:EmployeeServiceID"`
	Account           Account         `json:"account" gorm:"ForeignKey:AccountID"`
}

func (tim *TimeSlot) CreateTimeSlot() error {
	if err := DB().Create(&tim).Error; err != nil {
		return err
	}
	return nil
}

func GetTimeSlotByDate(t time.Time, serName string) ([]TimeSlot, error) {
	timeSlots := []TimeSlot{}
	var service Service
	if err := DB().Where("ser_name = ?", serName).Find(&service).Error; err != nil {
		return nil, err
	}
	if err := DB().Where("time_day = ?", int(t.Weekday())).Preload("EmployeeService", func(db *gorm.DB) *gorm.DB {
		return db.Where("service_id = ?", service.ID).Preload("Employee").Preload("Service")
	}).Find(&timeSlots).Error; err != nil {
		return nil, err
	}
	return timeSlots, nil
}

type TimeSlots []*TimeSlot

func (tims TimeSlots) CreateTimeSlotMultiple() error {
	if err := DB().Create(&tims).Error; err != nil {
		return err
	}
	return nil
}

func (tim *TimeSlot) UpdateTimeSlot(id string) error {
	if err := DB().Save(&tim).Error; err != nil {
		return err
	}
	return nil
}

func GetTimeSlotList(accID uint) ([]*TimeSlot, error) {
	timeSlots := []*TimeSlot{}
	if err := DB().Where("time_account_id = ?", accID).Find(&timeSlots).Error; err != nil {
		return nil, err
	}
	return timeSlots, nil
}

func GetTimeSlotDetail(id string, accID uint) (*TimeSlot, error) {
	timeSlot := TimeSlot{}
	if err := DB().Where("time_account_id = ?", accID).Find(&timeSlot).Error; err != nil {
		return nil, err
	}
	return &timeSlot, nil
}

func RemoveTimeSlot(id string) (*TimeSlot, error) {
	timeSlot := TimeSlot{}

	if err := DB().Find(&timeSlot, id).Error; err != nil {
		return nil, err
	}

	if err := DB().Delete(&timeSlot).Error; err != nil {
		return nil, err
	}

	return &timeSlot, nil
}
