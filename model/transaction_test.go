package model

import (
	"fmt"
	"testing"

	"github.com/IntouchOpec/base-go-echo/model"
)

func TestGetReceipt(t *testing.T) {
	var tr model.Transaction
	db := model.SqlDB()
	fmt.Println(db)
	tr.GetReceipt(db)
}
