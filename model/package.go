package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type Package struct {
	orm.ModelBase

	PacName     string     `form:"pac_name" json:"pac_name" gorm:"type:varchar(50)"`
	PacDetail   string     `form:"pac_detail" json:"pac_detail" `
	PacOrder    int        `form:"pac_order" json:"pac_order"`
	PacType     string     `form:"pac_type" json:"pac_type" gorm:"type:varchar(50)"`
	PacImage    string     `form:"pac_image" json:"pac_image" gorm:"type:varchar(255)"`
	PacIsActive bool       `form:"pac_is_active" json:"pac_is_active" sql:"default:true" gorm:"default:true"`
	Service     []*Service `json:"service" gorm:"many2many:package_service"`
	AccountID   uint       `json:"account_id"`
	Account     Account    `json:"account" gorm:"ForeignKey:AccountID"`
}
