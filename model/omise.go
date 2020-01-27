package model

import "encoding/json"

type OmiseLog struct {
	Json      json.RawMessage `sql:"type:json" json:"object,omitempty"`
	AccountID uint            `json:"account_id"`
	Account   Account         `json:"account" gorm:"ForeignKey:AccountID"`
}
