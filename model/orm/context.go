package orm

import (
	"time"

	"github.com/opentracing/opentracing-go"
)

type context interface {
	OpenTracingSpan() opentracing.Span
}

type ModelBase struct {
	ID        uint       `gorm:"primary_key"; "auto_increment" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
	Deleted   bool       `json:"deleted"`
}
