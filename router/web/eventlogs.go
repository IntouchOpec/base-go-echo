package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func EventLogList(c *Context) error {
	eventLogs := []*model.EventLog{}
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	ChatChannelID := c.QueryParam("chat_channel_id")

	filter := db.Model(&eventLogs).Where("chat_channel_id = ?", ChatChannelID).Count(&total)
	pagination := MakePagination(total, page, limit)
	filter.Preload("Customer").Limit(pagination.Record).Offset(pagination.Offset).Find(&eventLogs)

	return c.JSON(http.StatusOK, echo.Map{
		"data":       eventLogs,
		"pagination": pagination,
	})
}
