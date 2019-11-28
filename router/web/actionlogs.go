package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func ActionLogList(c *Context) error {
	act := []*model.ActionLog{}
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	fmt.Println(page, limit)
	ChatChannelID := c.QueryParam("chat_channel_id")

	filter := db.Model(&act).Where("chat_channel_id = ?", ChatChannelID)
	filter.Count(&total)
	pagination := MakePagination(total, page, limit)

	filter.Preload("Customer").Limit(pagination.Record).Offset(pagination.Offset).Order("id").Find(&act)

	return c.JSON(http.StatusOK, echo.Map{
		"data":       act,
		"pagination": pagination,
	})
}
