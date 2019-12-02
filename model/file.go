package model

import (
	"time"

	"github.com/jinzhu/gorm"

	uuid "github.com/satori/go.uuid"
)

type File struct {
	ID        uuid.UUID  `gorm:"primary_key ;type:varchar(255)" json:"id"`
	Path      string     `json:"path"`
	AccountID uint       `json:"account_id"`
	Account   Account    `json:"account"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
	Deleted   bool       `json:"deleted"`
}

func (f *File) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid)
}

func (f *File) BeforeDelete(scope *gorm.Scope) error {
	return scope.SetColumn("delet_ed", true)
}
