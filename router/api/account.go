package api

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// CreateAccount route create account.
// func CreateAccount(c echo.Context) error {
// 	account := model.Account{}
// 	if err := c.Bind(&account); err != nil {
// 		return c.NoContent(http.StatusBadRequest)
// 	}
// 	account.CreateAccount()
// 	return c.JSON(http.StatusOK, account)
// }

// func GetAccount(c echo.Context) error {
// 	account := model.Account{}
// 	id := c.Param("id")
// 	if err := c.Bind(&account); err != nil {
// 		return c.NoContent(http.StatusBadRequest)
// 	}
// 	account.GetAccountByID(id)
// 	return c.JSON(http.StatusOK, account)
// }

func GetAccontList(c echo.Context) error {
	// page := c.QueryParam("page")
	// size := c.QueryParam("size")

	// pageInt, _ := strconv.Atoi(page)
	// sizeInt, _ := strconv.Atoi(size)

	accounts, err := model.GetAccount()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, accounts)
}

// func UpdateAccount(c echo.Context) error {
// 	id := c.Param("id")
// 	account := model.Account{}
// 	account.GetAccountByID(id)

// 	if err := c.Bind(&account); err != nil {
// 		return c.NoContent(http.StatusBadRequest)
// 	}

// 	account.UpdateAccount()
// 	return c.JSON(http.StatusOK, account)
// }

// func DeleteAccount(c echo.Context) error {
// 	id := c.Param("id")
// 	account := model.Account{}
// 	account.RemoveAccount(id)
// 	return c.JSON(http.StatusOK, account)
// }
