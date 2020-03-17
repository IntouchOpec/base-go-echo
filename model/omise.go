package model

import (
	"database/sql"
	"encoding/json"
)

type OmiseLog struct {
	Json      json.RawMessage `sql:"type:json" json:"object,omitempty"`
	AccountID uint            `json:"account_id"`
	Account   Account         `json:"account" gorm:"ForeignKey:AccountID"`
}

func (o *OmiseLog) Create(sql *sql.DB) error {
	_, err := sql.Exec("INSERT INTO omise_logs (json, account_id) VALUES ($1, $2)", o.Json, o.AccountID)
	if err != nil {
		return err
	}
	return nil
}
