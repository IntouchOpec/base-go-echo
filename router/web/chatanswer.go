package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func ChatAnswerListHandler(c *Context) error {
	chatAnswers := []*model.ChatAnswer{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	filterChatAns := db.Where("account_id = ?", a.GetAccountID()).Find(&chatAnswers).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterChatAns.Limit(pagination.Record).Offset(pagination.Offset).Find(&chatAnswers)
	return c.Render(http.StatusOK, "chat-answer-list", echo.Map{
		"title":      "chat_answer",
		"list":       chatAnswers,
		"pagination": pagination,
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
		"method":       "POST",
		"title":        "chat_answer",
		"detail":       chatAnswer,
		"messageTypes": messageTypes,
	})
}

func ChatAnswerPostHandler(c *Context) error {
	chatAnswer := model.ChatAnswer{}
	if err := c.Bind(&chatAnswer); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	a := auth.Default(c)

	chatAnswer.AccountID = a.GetAccountID()
	err := chatAnswer.SaveChatAnswer()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/chat_answer/%d", chatAnswer.ID)
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
		"data":     chatAnswer,
	})
}

func ChatAnswerEditHandler(c *Context) error {
	id := c.Param("id")
	chatAnswer := model.ChatAnswer{}
	a := auth.Default(c)

	model.DB().Preload("Account", "name = ?", a.User.GetAccount()).Find(&chatAnswer, id)
	return c.Render(http.StatusOK, "chat-answer-form", echo.Map{"method": "PUT",
		"title":  "chat_answer",
		"detail": chatAnswer,
	})
}

func ChatAnswerDeleteHandler(c *Context) error {
	id := c.Param("id")

	chatChannel := model.DeleteChatAnswer(id)
	return c.JSON(http.StatusOK, chatChannel)
}
