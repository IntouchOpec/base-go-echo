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

	filterCustomer := db.Model(&customers).Where("account_id = ?", a.GetAccountID()).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterCustomer.Limit(pagination.Record).Offset(pagination.Offset).Find(&customers)

	return c.Render(http.StatusOK, "customer-list", echo.Map{
		"list":       customers,
		"title":      "customer",
		"pagination": pagination,
	})
}

func CustomerDetailHandler(c *Context) error {
	id := c.Param("id")
	customer := model.Customer{}
	db := model.DB()

	var totalEvent int
	var totalAction int
	var paginationEventLogs Pagination
	var paginationActionLogs Pagination
	db.Find(&customer, id)

	eventLogs := []model.EventLog{}
	eventLogsFilter := db.Model(&eventLogs).Where("customer_id = ?", id).Count(&totalEvent)
	paginationEventLogs = MakePagination(totalEvent, 0, 10)
	eventLogsFilter.Preload("Customer").Find(&eventLogs).Limit(10).Offset(0).Order("id")

	actionLogs := []model.ActionLog{}
	filteractionLogs := db.Model(&actionLogs).Where("chat_channel_id = ?", id).Count(&totalAction)
	paginationActionLogs = MakePagination(totalAction, 0, 10)
	filteractionLogs.Preload("Customer").Find(&actionLogs).Limit(10).Offset(0).Order("id")

	return c.Render(http.StatusOK, "customer-detail", echo.Map{
		"customer":             customer,
		"title":                "customer",
		"actionLogs":           actionLogs,
		"eventLogs":            eventLogs,
		"paginationActionLogs": paginationActionLogs,
		"paginationEventLogs":  paginationEventLogs,
	})
}

func CustomerDeleteHandler(c *Context) error {
	id := c.Param("id")

	customer := model.DeleteCustomer(id)
	return c.JSON(http.StatusOK, customer)
}
