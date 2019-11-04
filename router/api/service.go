package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// Createserviceis route handle for create service
func Createservice(c echo.Context) error {
	service := model.Service{}
	if err := c.Bind(&service); err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	service.Saveservice()
	c.JSON(200, service)

	return nil
}

func GetserviceList(c echo.Context) error {
	chatchannelID := c.Param("chatChannelID")
	chatChannelIDInt, _ := strconv.Atoi(chatchannelID)

	services := model.Getservice(chatChannelIDInt)
	return c.JSON(http.StatusOK, services)
}

func Getservice(c echo.Context) error {
	id := c.Param("id")
	chatchannelID := c.Param("chatChannelID")
	chatChannelIDInt, _ := strconv.Atoi(chatchannelID)
	idInt, _ := strconv.Atoi(id)
	service := model.GetserviceByID(chatChannelIDInt, idInt)
	return c.JSON(http.StatusOK, service)
}

func Updateservice(c echo.Context) error {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	service := model.Service{}
	if err := c.Bind(&service).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	service.Updateservice(idInt)
	return c.JSON(http.StatusOK, service)
}

func UpdateserviceSlot(c echo.Context) error {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	subCustomer := model.ServiceSlot{}
	if err := c.Bind(&subCustomer).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	subCustomer.UpdateServiceSlot(idInt)

	return c.JSON(http.StatusOK, subCustomer)
}
