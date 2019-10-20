package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// CustomerListHandler
func CustomerListHandler(c *Context) error {
	customers := []*model.Customer{}

	if err := model.DB().Find(&customers).Error; err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "customer-list", echo.Map{
		"list":  customers,
		"title": "customer",
	})
	return err
}

func CustomerDetailHandler(c *Context) error {
	id := c.Param("id")
	customer := model.Customer{}
	//
	if err := model.DB().Preload("Bookings").Preload("EventLogs").Preload("ActionLogs").Find(&customer, id).Error; err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
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
