package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// CustomerListHandler
func CustomerListHandler(c *Context) error {
	lineID := c.Param("lineID")
	chatChannel := model.ChatChannel{}
	if err := model.DB().Where("line_ID = ?", lineID).Find(&chatChannel).Error; err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	customers := []*model.Customer{}

	if err := model.DB().Where("chat_channel_id = ?", chatChannel.ID).Find(&customers).Error; err != nil {
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
