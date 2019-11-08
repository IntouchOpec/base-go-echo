package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type TimeSlot struct {
	orm.ModelBase

	TimeStart             string          `json:"time_start" gorm:"type:varchar(10)"`
	TimeEnd               string          `json:"time_end" gorm:"type:varchar(10)"`
	TimeDay               int             `json:"time_day"`
	TimeAmount            int             `json:"time_amount"`
	TimeActive            bool            `json:"time_active"`
	TimeProviderServiceID uint            `json:"time_provider_service_id"`
	TimeAccountID         uint            `json:"time_account_id"`
	ProviderService       ProviderService `json:"provider_service" `
	Bookings              []*Booking      `json:"bookings"`
	Account               Account         `json:"account"`
}
