package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/line/line-bot-sdk-go/linebot"

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

func ChatAnswerDetailHandler(c *Context) error {
	id := c.Param("id")
	chatAnswer := model.ChatAnswer{}
	a := auth.Default(c)

	model.DB().Preload("Account", "name = ?", a.User.GetAccount()).Find(&chatAnswer, id)
	return c.Render(http.StatusOK, "chat-answer-detail", echo.Map{
		"title":  "chat_answer",
		"detail": chatAnswer,
	})
}

func ChatAnswerCreateHandler(c *Context) error {
	chatAnswer := model.ChatAnswer{}
	messageTypes := []linebot.MessageType{
		linebot.MessageTypeText,
		linebot.MessageTypeImage,
		linebot.MessageTypeVideo,
		linebot.MessageTypeAudio,
		linebot.MessageTypeFile,
		linebot.MessageTypeLocation,
		linebot.MessageTypeSticker,
		linebot.MessageTypeTemplate,
		linebot.MessageTypeImagemap,
		linebot.MessageTypeFlex,
	}

	return c.Render(http.StatusOK, "chat-answer-form", echo.Map{
		"title":        "chat_answer",
		"detail":       chatAnswer,
		"messageTypes": messageTypes,
	})
}

func ChatAnswerPostHandler(c *Context) error {
	chatAnswer := model.ChatAnswer{}
	return c.Render(http.StatusOK, "chat-answer-form", echo.Map{
		"title":  "chat_answer",
		"detail": chatAnswer,
	})
}

func ChatAnswerEditHandler(c *Context) error {
	id := c.Param("id")
	chatAnswer := model.ChatAnswer{}
	a := auth.Default(c)

	model.DB().Preload("Account", "name = ?", a.User.GetAccount()).Find(&chatAnswer, id)
	return c.Render(http.StatusOK, "chat-answer-form", echo.Map{
		"title":  "chat_answer",
		"detail": chatAnswer,
	})
}

func ChatAnswerDeleteHandler(c *Context) error {
	id := c.Param("id")

	chatChannel := model.DeleteChatAnswer(id)
	return c.JSON(http.StatusOK, chatChannel)
}
