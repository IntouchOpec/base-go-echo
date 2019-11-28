package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type Package struct {
	orm.ModelBase

	PacName   string  `json:"pac_name" gorm:"type:varchar(50)"`
	PacDetail string  `json:"pac_detail" `
	PacType   string  `json:"pac_type" gorm:"type:varchar(50)"`
	AccountID string  `json:"account_id"`
	Account   Account `json:"account" gorm:"ForeignKey:AccountID"`
}
