package api

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func GetChatRequestList(c echo.Context) error {
	chatChannelID := c.Param("chatChannelID")
	chatReqList := model.GetChatRequest(chatChannelID)
	return c.JSON(http.StatusOK, chatReqList)
}
