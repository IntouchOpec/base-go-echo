package model

import "github.com/hb-go/gorm"

// Customer follow Line OA.
type Customer struct {
	gorm.Model
	// ID          uint    `json:"id,omitempty"`
	LineID      string  `json:"line_id" gorm:"type:varchar(25)"`
	Email       string  `json:"email" gorm:"type:varchar(25)"`
	PhoneNumber string  `json:"phone_number" gorm:"type:varchar(25)"`
	AccountID   uint    `form:"account_id" json:"account_id" gorm:"not null;"`
	Account     Account `gorm:"ForeignKey:id"`
}

// SaveCustomer is function create Customer.
func (customer *Customer) SaveCustomer() *Customer {
	if err := DB().Create(&customer).Error; err != nil {
		return nil
	}
	return customer
}
