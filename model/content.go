package model

import "github.com/IntouchOpec/base-go-echo/model/orm"

type Content struct {
	orm.ModelBase
	ConTitle  string  `form:"con_title" json:"con_title" gorm:"varchar(50)"`
	ConDetail string  `form:"con_detail" json:"con_detail"`
	ConImage  string  `json:"con_image" gorm:"varchar(50)"`
	UserID    uint    `json:"user_id"`
	User      *User   `json:"user" gorm:"ForeignKey:UserID"`
	AccountID uint    `json:"account_id"`
	Account   Account `json:"account" gorm:"ForeignKey:AccountID"`
}

func DeleteContent(id string, accID uint) error {
	var content Content
	if err := DB().Where("account_id = ?", accID).Find(&content).Error; err != nil {
		return err
	}
	if err := DB().Delete(&content).Error; err != nil {
		return err
	}
	return nil
}

func (con *Content) SaveContent() error {
	if err := DB().Create(&con).Error; err != nil {
		return err
	}
	return nil
}
