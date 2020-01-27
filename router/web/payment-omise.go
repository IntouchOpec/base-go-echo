package web

import (
	"fmt"
	"net/http"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"

	"github.com/labstack/echo"
)

func PaymentOmiseHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "payment-omise", echo.Map{})
}

const (
	// Read these from environment variables or configuration files!
	OmisePublicKey = "pkey_test_5ip8fflleizk5mzvnut"
	OmiseSecretKey = "skey_test_5ip8nm6pyp7ziztxlh9"
)

// type TokenOmiseReq struct {

// }

func ChargeOmiseHandler(c echo.Context) error {
	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	token := c.FormValue("token")
	fmt.Println(token)
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   100000,
		Currency: "thb",
		Card:     token,
	}
	if err := client.Do(charge, createCharge); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{})
}
