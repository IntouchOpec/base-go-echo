package api

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// CreateAccount route create account.
func CreateAccount(c echo.Context) error {
	account := model.Account{}
	if err := c.Bind(&account); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	fmt.Println(account)
	account.CreateAccount()
	return c.JSON(http.StatusOK, account)
}
