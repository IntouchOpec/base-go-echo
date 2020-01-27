package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"

	"github.com/labstack/echo"
)

func PaymentOmiseHandler(c echo.Context) error {
	accountName := c.QueryParam("account_name")
	DocCodeTransaction := c.QueryParam("doc_code_transaction")
	var transaction model.Transaction
	var account model.Account
	db := model.DB()
	db.Where("acc_name = ?", accountName).Find(&account)
	db.Where("account_id = ? and tran_doccument_code = ?", account.ID, DocCodeTransaction).Find(&transaction)

	return c.Render(http.StatusOK, "payment-omise", echo.Map{
		"accountName":        accountName,
		"DocCodeTransaction": DocCodeTransaction,
		"detail":             transaction,
	})
}

const (
	// Read these from environment variables or configuration files!
	OmisePublicKey = "pkey_test_5ip8fflleizk5mzvnut"
	OmiseSecretKey = "skey_test_5ip8nm6pyp7ziztxlh9"
)

func ChargeOmiseHandler(c echo.Context) error {
	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	accountName := c.QueryParam("account_name")
	DocCodeTransaction := c.QueryParam("doc_code_transaction")
	var transaction model.Transaction
	var account model.Account
	db := model.DB()
	db.Where("acc_name = ?", accountName).Find(&account)
	db.Where("account_id = ? and tran_doccument_code = ?", account.ID, DocCodeTransaction).Find(&transaction)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	token := c.FormValue("token")
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   100000,
		Currency: "thb",
		Card:     token,
	}
	if err := client.Do(charge, createCharge); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	var omise model.OmiseLog
	ev, err := json.Marshal(charge)
	omise.Json = ev
	omise.AccountID = account.ID
	if err := db.Save(&omise).Error; err != nil {

	}
	transaction.TranStatus = model.TranStatusPaid
	if err := db.Save(&transaction).Error; err != nil {
	}
	fmt.Printf("%+v\n", charge)

	return c.JSON(http.StatusOK, echo.Map{})
}
