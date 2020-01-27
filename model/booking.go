package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type BookStatus int

const (
	BookingStatusPandding BookStatus = 1
	BookingStatusReject   BookStatus = -1
	BookingStatusApprove  BookStatus = 2
)

type BookingType int

const (
	BookingTypeSlotTime    BookingType = 1
	BookingTypeServiceItem BookingType = 2
	BookingTypePackage     BookingType = 3
)

// Booking struct save date time
type Booking struct {
	orm.ModelBase

	BookingType   BookingType  `json:"booking_type"`
	BooQueue      int          `json:"boo_queue" `
	BooLineID     string       `json:"boo_line_id" gorm:"type:varchar(50)"`
	CustomerID    uint         `json:"customer_id"`
	ChatChannelID uint         `json:"chat_chaneel_id"`
	TransactionID uint         `json:"transaction_id"`
	PlaceID       uint         `json:"place_id"`
	Place         *Place       `json:"place" gorm:"ForeignKey:PlaceID"`
	Transaction   *Transaction `json:"transaction"  gorm:"ForeignKey:TransactionID"`
	Customer      *Customer    `json:"customer" gorm:"ForeignKey:CustomerID"`
	ChatChannel   *ChatChannel `json:"chat_channel" gorm:"ForeignKey:ChatChannelID"`
	BooStatus     *BookStatus  `json:"boo_status"`
	BookedDate    time.Time    `gorm:"column:booked_date" json:"booked_date"`
	AccountID     uint         `json:"account_id"`
	Account       Account      `json:"account" gorm:"ForeignKey:AccountID"`
}

type BookingTimeSlot struct {
	BookingID  uint      `json:"booking_id"`
	Booking    Booking   `json:"booking" gorm:"ForeignKey:BookingID"`
	TimeSlotID uint      `json:"time_slot_id"`
	TimeSlot   *TimeSlot `json:"time_slot" gorm:"ForeignKey:TimeSlotID"`
	EmployeeID uint      `json:"employee_id"`
	Employee   *Employee `json:"employee" gorm:"ForeignKey:EmployeeID"`
}

type BookingServiceItem struct {
	BookingID     uint        `json:"booking_id"`
	Booking       Booking     `json:"booking" gorm:"ForeignKey:BookingID"`
	ServiceItemID uint        `json:"serice_item_id"`
	ServiceItem   ServiceItem `json:"service_item" gorm:"ForeignKey:ServiceItemID"`
}

type BookingPackage struct {
	BookingID uint    `json:"booking_id"`
	Booking   Booking `json:"booking" gorm:"ForeignKey:BookingID"`
	PackageID uint    `json:"package_id"`
	Package   Package `json:"package" gorm:"ForeignKey:PackageID"`
}

// BookingStatus is status of booking.
// type BookingStatus struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

// BookingStatusPandding is booking status pandding for confirm
// var BookingStatusPandding = BookingStatus{ID: 1, Name: "pandding"}

// BookingStatusReject is booking status after pandding user pick It.
// var BookingStatusReject = BookingStatus{ID: 2, Name: "reject"}

// BookingStatusApprove is status approve.
// var BookingStatusApprove = BookingStatus{ID: 3, Name: "approve"}

// BookingState is state of booking.
type BookingState struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SaveBooking is function create chat answer.
// func (booking *Booking) SaveBooking() (*Booking, error) {
// db := DB()
// booked := Booking{}
// db.Preload("ServiceSlot").Where("Booked_Date = ? and Sub_service_ID = ?", booking.BookedDate, booking.TimeSlotID).Last(&booked)
// if booked.TimeSlot.TimeAmount == 0 {
// 	booking.BooQueue = 1
// } else if booked.TimeSlot.TimeAmount > booked.BooQueue {
// 	booking.BooQueue = booked.BooQueue + 1
// } else {
// 	return nil, errors.New("can't insert booking case queue full")
// }

// if err := db.Create(&booking).Error; err != nil {
// 	return nil, err
// }
// return booking, nil
// }

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
