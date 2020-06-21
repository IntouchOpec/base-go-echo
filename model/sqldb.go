package model

import (
	"database/sql"
	"fmt"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"

	_ "github.com/lib/pq"
)

func newSqlDB() (*sql.DB, error) {
	sqlConnection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Conf.DB.Host, Conf.DB.Port, "admin", "todos", Conf.DB.Pwd)
	sqlDb, err := sql.Open("postgres", sqlConnection)
	if err != nil {
		return nil, err
	}
	return sqlDb, nil
}

func SqlDB() *sql.DB {
	newDb, err := newSqlDB()
	if err != nil {
		panic(err)
	}
	newDb.SetMaxOpenConns(20) // Sane default
	newDb.SetMaxIdleConns(0)
	newDb.SetConnMaxLifetime(time.Nanosecond)
	err = newDb.Ping()
	if err != nil {
		return nil
	}
	return newDb
}
