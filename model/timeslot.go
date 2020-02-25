package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
)

type TimeSlot struct {
	orm.ModelBase

	TimeStart  time.Time  `form:"time_start" json:"time_start"`
	TimeEnd    time.Time  `form:"time_end" json:"time_end"`
	TimeDay    int        `json:"time_day"`
	TimeActive bool       `json:"time_active" sql:"default:true" gorm:"default:true"`
	EmployeeID uint       `json:"employee_id"`
	Employee   Employee   `json:"employee" gorm:"ForeignKey:EmployeeID"`
	AccountID  uint       `json:"account_id"`
	Bookings   []*Booking `json:"bookings"`
	Account    Account    `json:"account" gorm:"ForeignKey:AccountID"`
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
