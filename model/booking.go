package model

import (
	"fmt"
	"time"

	"github.com/hb-go/gorm"
)

// Booking struct save date time
type Booking struct {
	gorm.Model
	// ID            uint          `json:"id,omitempty"`
	Queue         int         `json:"queue" `
	LineID        string      `json:"line_id" gorm:"type:varchar(50)"`
	CustomerID    uint        `json:"customer_id"`
	Customer      Customer    `gorm:"foreignkey:ID"`
	SubProductID  uint        `json:"prodict_id"`
	SubProduct    SubProduct  `gorm:"foreignkey:SubProductID"`
	AccountID     uint        `form:"account_id" json:"account_id" gorm:"not null;"`
	Account       Account     `gorm:"ForeignKey:id"`
	ChatChannelID uint        `json:"chat_chaneel_id"`
	ChatChannel   ChatChannel `gorm:"ForeignKey:id"`
	BookingStatus int         `json:"booking_status"`
	BookingState  int         `json:"booking_state"`
	BookingDate   time.Time   `gorm:"column:booking_date" json:"booking_date"`
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
func (booking *Booking) SaveBooking() *Booking {
	db := DB()
	// subProduct := SubProduct{}
	// db.Where("BookingDate = ? and SubProductID = ?", booking.BookingDate, booking.SubProductID).Find(&booking).Related(&subProduct).Last(booking)
	// if subProduct.Amount >= booking.Queue {
	// 	return nil
	// }
	if err := db.Create(&booking).Error; err != nil {
		fmt.Println(err)
		return nil
	}
	return booking
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
