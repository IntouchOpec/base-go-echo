package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/module/auth"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func TransactionListHandler(c *Context) error {
	Transactions := []*model.Transaction{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	filterTran := db.Where("account_id = ?", a.GetAccountID()).Find(&Transactions).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterTran.Limit(pagination.Record).Offset(pagination.Offset).Find(&Transactions)
	return c.Render(http.StatusOK, "transaction-list", echo.Map{
		"title":      "transaction",
		"list":       Transactions,
		"pagination": pagination,
	})
}

func TransactionDetailHandler(c *Context) error {
	id := c.Param("id")
	Transaction := model.Transaction{}
	a := auth.Default(c)

	model.DB().Preload("Account", "name = ?", a.User.GetAccount()).Find(&Transaction, id)
	return c.Render(http.StatusOK, "transaction-detail", echo.Map{
		"title":  "transaction",
		"detail": Transaction,
	})
}

func TransactionCreateHandler(c *Context) error {
	Transaction := model.Transaction{}
	return c.Render(http.StatusOK, "transaction-form", echo.Map{
		"method": "POST",
		"title":  "transaction",
		"detail": Transaction,
	})
}

func TransactionPostHandler(c *Context) error {
	Transaction := model.Transaction{}
	if err := c.Bind(&Transaction); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	accID := auth.Default(c).GetAccountID()
	Transaction.AccountID = accID
	err := Transaction.Create()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println(Transaction.ID)
	redirect := fmt.Sprintf("/admin/transaction/%d", Transaction.ID)
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
		"data":     Transaction,
	})
}

func TransactionEditHandler(c *Context) error {
	id := c.Param("id")
	Transaction := model.Transaction{}
	a := auth.Default(c)

	model.DB().Preload("Account", "name = ?", a.User.GetAccount()).Find(&Transaction, id)
	return c.Render(http.StatusOK, "transaction-form", echo.Map{"method": "PUT",
		"title":  "transaction",
		"detail": Transaction,
	})
}

func TransactionDeleteHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()

	chatChannel, err := model.RemoveTransaction(id, accID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, chatChannel)
}
