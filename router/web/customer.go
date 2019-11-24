package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

// CustomerListHandler
func CustomerListHandler(c *Context) error {
	customers := []*model.Customer{}
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	a := auth.Default(c)
	filterCustomer := db.Where("cus_account_id = ?", a.GetAccountID()).Find(&customers).Count(&total)
	filterCustomer.Limit(limit).Offset(page).Find(&customers)
	pagination := MakePagination(total, page, limit)

	err := c.Render(http.StatusOK, "customer-list", echo.Map{
		"list":       customers,
		"title":      "customer",
		"pagination": pagination,
	})
	return err
}

func CustomerDetailHandler(c *Context) error {
	id := c.Param("id")
	customer := model.Customer{}
	if err := model.DB().Preload("ActionLogs").Find(&customer, id).Error; err != nil {
		return c.Render(http.StatusOK, "customer-detail", echo.Map{
			"customer": customer,
			"title":    "customer",
		})
	}
	err := c.Render(http.StatusOK, "customer-detail", echo.Map{
		"customer": customer,
		"title":    "customer",
	})
	return err
}

func CustomerDeleteHandler(c *Context) error {
	id := c.Param("id")

	customer := model.DeleteCustomer(id)
	return c.JSON(http.StatusOK, customer)
}
