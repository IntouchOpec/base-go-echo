package model

import (
	"fmt"
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
	ProviderServiceID uint            `json:"provider_service_id"`
	AccountID         uint            `json:"account_id"`
	ProviderService   ProviderService `json:"provider_service" gorm:"ForeignKey:ProviderServiceID"`
	Bookings          []*Booking      `json:"bookings"`
	Account           Account         `json:"account" gorm:"ForeignKey:AccountID"`
}

func (tim *TimeSlot) CreateTimeSlot() error {
	if err := DB().Create(&tim).Error; err != nil {
		return err
	}
	return nil
}

func GetTimeSlotByDate(t time.Time) ([]TimeSlot, error) {
	timeSlots := []TimeSlot{}
	ps := []ProviderService{}
	fmt.Println(ps)
	if err := DB().Preload("ProviderService", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Provider").Preload("Service")
	}).Where("time_day = ?", int(t.Weekday())).Find(&timeSlots).Error; err != nil {
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
