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
	db.Where("account_id = ? and tran_document_code = ?", account.ID, DocCodeTransaction).Find(&transaction)
	if transaction.TranStatus == model.TranStatusPaid {
		return c.Render(http.StatusOK, "payment-success", echo.Map{
			"accountName":        accountName,
			"DocCodeTransaction": DocCodeTransaction,
			"detail":             transaction,
		})
	}

	return c.Render(http.StatusOK, "payment-omise", echo.Map{
		"accountName":        accountName,
		"DocCodeTransaction": DocCodeTransaction,
		"detail":             transaction,
	})
}

const (
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
	db.Where("account_id = ? and tran_document_code = ?", account.ID, DocCodeTransaction).Find(&transaction)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	token := c.FormValue("token")
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   int64(transaction.TranTotal * 100),
		Currency: "thb",
		Card:     token,
	}
	if err := client.Do(charge, createCharge); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	var omiseLog model.OmiseLog
	ev, err := json.Marshal(charge)
	omiseLog.Json = ev
	omiseLog.AccountID = account.ID
	if charge.Status != omise.ChargeSuccessful {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Save(&omiseLog).Error; err != nil {
		fmt.Println(err)

	}
	var payment model.Payment
	payment.PayAt = charge.Created
	payment.TransactionID = transaction.ID
	payment.PayStatus = model.PayStatusSuccess
	transaction.TranStatus = model.TranStatusPaid
	payment.PayAmount = transaction.TranTotal
	payment.PayType = model.PayTypeOmise
	fmt.Println(transaction)
	if err := db.Model(&transaction).Updates(transaction).Error; err != nil {
		fmt.Print(err)
	}
	if err := db.Model(&transaction).Association("Payments").Append(&payment).Error; err != nil {
		fmt.Print(err)
	}

	return c.JSON(http.StatusOK, echo.Map{})
}
