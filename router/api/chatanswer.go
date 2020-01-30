package api

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// CreateChatAnswer create ChatAnswer.
func CreateChatAnswer(c echo.Context) error {
	chatAnswer := model.ChatAnswer{}
	if err := c.Bind(&chatAnswer); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	chatAnswer.SaveChatAnswer()
	return c.JSON(http.StatusOK, chatAnswer)
}

func GetChatAnswerList(c echo.Context) error {
	chatChannelID := c.Param("chatChannelID")
	chatAnswer := model.GetChatAnswerList(chatChannelID)
	return c.JSON(http.StatusOK, chatAnswer)
}

func GetChatAnswerDetail(c echo.Context) error {
	id := c.Param("id")
	chatAnswer := model.GetChatAnswer(id)
	return c.JSON(http.StatusOK, chatAnswer)
}

func UpdateChatAnswers(c echo.Context) error {
	id := c.Param("id")
	chatAnswer := model.ChatAnswer{}
	if err := c.Bind(&chatAnswer); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	chatAnswer.UpdateChatAnswer(id)
	return c.JSON(http.StatusOK, chatAnswer)
}

func DeleteChatAnswer(c echo.Context) error {
	id := c.Param("id")
	cha := model.DeleteChatAnswer(id)
	return c.JSON(http.StatusOK, cha)
}
