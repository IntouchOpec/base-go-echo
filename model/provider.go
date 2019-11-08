package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type Provider struct {
	orm.ModelBase

	ProvNmae      string  `json:"prov_name" gorm:"type:varchar(25)"`
	ProvDetail    string  `json:"prov_detail"`
	ProvLineID    string  `json:"prov_line_id" gorm:"type:varchar(50)"`
	ProvAccountID uint    `json:"prov_account_id"`
	Account       Account `json:"account" gorm:"ForeignKey:ProvAccountID"`
}
