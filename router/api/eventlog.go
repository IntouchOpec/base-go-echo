package api

import (
	"net/http"
	"strconv"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func GetEvenLogList(c echo.Context) error {
	page := c.QueryParam("page")
	size := c.QueryParam("size")
	chatChannelID := c.Param("chatChannelID")
	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)
	chatChannelIDInt, _ := strconv.Atoi(chatChannelID)
	eventLogs := model.GetEventLog(pageInt, sizeInt, chatChannelIDInt)
	return c.JSON(http.StatusOK, eventLogs)
}

func GetAllEvenLog(c echo.Context) error {
	page := c.QueryParam("page")
	size := c.QueryParam("size")
	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)
	eventLogs := model.GetAllEventLog(pageInt, sizeInt)
	return c.JSON(http.StatusOK, eventLogs)
}
