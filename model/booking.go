package model

import (
	"errors"
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

// Booking struct save date time
type Booking struct {
	orm.ModelBase
	// ID            uint          `json:"id,omitempty"`
	Queue         int         `json:"queue" `
	LineID        string      `json:"line_id" gorm:"type:varchar(50)"`
	CustomerID    uint        `json:"customer_id"`
	Customer      Customer    `json:"customer" gorm:"ForeignKey:CustomerID"`
	ServiceSlotID uint        `json:"sub_service_id"`
	ServiceSlot   ServiceSlot `json:"sub_service" gorm:"ForeignKey:ServiceSlotID"`
	ChatChannelID uint        `json:"chat_chaneel_id"`
	ChatChannel   ChatChannel `gorm:"ForeignKey:ChatChannelID"`
	BookStatus    int         `json:"booking_status"`
	BookState     int         `json:"booking_state"`
	BookedDate    time.Time   `gorm:"column:booked_date" json:"booked_date"`
}

// BookingStatus is status of booking.
type BookingStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// BookingStatusPandding is booking status pandding for confirm
var BookingStatusPandding = BookingStatus{ID: 1, Name: "pandding"}

// BookingStatusReject is booking status after pandding user pick It.
var BookingStatusReject = BookingStatus{ID: 2, Name: "reject"}

// BookingStatusApprove is status approve.
var BookingStatusApprove = BookingStatus{ID: 3, Name: "approve"}

// BookingState is state of booking.
type BookingState struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SaveBooking is function create chat answer.
func (booking *Booking) SaveBooking() (*Booking, error) {
	db := DB()
	booked := Booking{}
	db.Preload("ServiceSlot").Where("Booked_Date = ? and Sub_service_ID = ?", booking.BookedDate, booking.ServiceSlotID).Last(&booked)
	if booked.ServiceSlot.Amount == 0 {
		booking.Queue = 1
	} else if booked.ServiceSlot.Amount > booked.Queue {
		booking.Queue = booked.Queue + 1
	} else {
		return nil, errors.New("can't insert booking case queue full")
	}

	if err := db.Create(&booking).Error; err != nil {
		return nil, err
	}
	return booking, nil
}

func (booking *Booking) UpdateBooking(id string) *Booking {
	db := DB()
	if err := db.Find(&booking, id).Error; err != nil {
		return nil
	}

	if err := db.Save(&booking).Error; err != nil {
		return nil
	}
	return booking
}

func GetBookingList(chatChannelID string) *[]Booking {
	bookings := []Booking{}
	db := DB()
	if err := db.Where("chat_channel_id = ?", chatChannelID).Find(&bookings).Error; err != nil {
		return nil
	}
	return &bookings
}

func GetBooking(id string) *Booking {
	db := DB()
	booking := Booking{}
	if err := db.Find(&booking, id).Error; err != nil {
		return nil
	}
	return &booking
}

func (booking *Booking) DeleteBooking(id string) *Booking {
	db := DB()
	if err := db.Find(&booking, id).Error; err != nil {
		return nil
	}
	if err := db.Delete(&booking, id).Error; err != nil {
		return nil
	}
	return booking
}

func (book *Booking) BookingAcjectStatus(status string) (*Booking, error) {

	return book, nil
}
