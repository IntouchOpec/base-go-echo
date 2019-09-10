package api

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// CreateChatAnswer create ChatAnswer.
func CreateChatAnswer(c echo.Context) error {
	chatAnswer := model.ChatAnswer{}
	if err := c.Bind(&chatAnswer); err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	chatAnswer.SaveChatAnswer()
	return c.JSON(http.StatusOK, chatAnswer)
}
