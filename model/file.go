package model

import (
	"time"

	"github.com/jinzhu/gorm"

	guuid "github.com/google/uuid"
)

type File struct {
	ID        guuid.UUID `gorm:"primary_key ;type:varchar(255)" json:"id"`
	Path      string     `json:"path"`
	AccountID uint       `json:"account_id"`
	Account   Account    `json:"account"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
	Deleted   bool       `json:"deleted"`
}

func (f *File) BeforeCreate(scope *gorm.Scope) error {
	uuid := guuid.New()
	return scope.SetColumn("ID", uuid)
}

func (f *File) BeforeDelete(scope *gorm.Scope) error {
	return scope.SetColumn("delet_ed", true)
}
