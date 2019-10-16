package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func ChatAnswerListHandler(c *Context) error {
	// account := c.Param("account")
	chatAnswer := []*model.ChatAnswer{}
	model.DB().Preload("Account").Find(&chatAnswer)
	return c.Render(http.StatusOK, "chat-answer-list", echo.Map{
		"title": "chat_answer",
		"list":  chatAnswer,
	})
}
