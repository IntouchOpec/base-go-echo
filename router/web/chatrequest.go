package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func ChatRequestListHandler(c *Context) error {
	ChatRequest := []*model.ChatRequest{}
	a := auth.Default(c)

	if err := model.DB().Where("req_account_id = ?",
		a.GetAccountID()).Find(&ChatRequest).Error; err != nil {
		return c.Render(http.StatusOK, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "chat-request-list", echo.Map{
		"title": "chat_request",
		"list":  ChatRequest,
	})
}

func ChatRequestDetailHandler(c *Context) error {
	id := c.Param("id")
	ChatRequest := model.ChatRequest{}
	a := auth.Default(c)

	if err := model.DB().Where("req_account_id = ?",
		a.GetAccountID()).Find(&ChatRequest, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "chat-request-detail", echo.Map{
		"title":  "chat_request",
		"detail": ChatRequest,
	})
}

func ChatRequestCreateHandler(c *Context) error {
	ChatRequest := model.ChatRequest{}
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

	return c.Render(http.StatusOK, "chat-request-form", echo.Map{
		"title":        "chat_request",
		"detail":       ChatRequest,
		"messageTypes": messageTypes,
	})
}

func ChatRequestPostHandler(c *Context) error {
	ChatRequest := model.ChatRequest{}
	if err := c.Bind(&ChatRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := ChatRequest.SaveChatRequest()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/chat_request/%d", ChatRequest.ID)
	return c.JSON(http.StatusOK, redirect)
}

func ChatRequestEditViewHandler(c *Context) error {
	id := c.Param("id")
	ChatRequest := model.ChatRequest{}
	a := auth.Default(c)

	if err := model.DB().Where("req_account_id = ?",
		a.GetAccountID()).Find(&ChatRequest, id).Error; err != nil {
		return c.Render(http.StatusOK, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "chat-request-form", echo.Map{
		"title":  "chat_request",
		"detail": ChatRequest,
	})
}

func ChatRequestEditHandler(c *Context) error {
	id := c.Param("id")
	ChatRequest := model.ChatRequest{}
	a := auth.Default(c)

	if err := model.DB().Where("req_account_id = ?",
		a.User.GetAccount()).Find(&ChatRequest, id).Error; err != nil {
		return c.Render(http.StatusOK, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "chat-request-form", echo.Map{
		"title":  "chat_request",
		"detail": ChatRequest,
	})
}

func ChatRequestDeleteHandler(c *Context) error {
	id := c.Param("id")
	a := auth.Default(c)

	chatChannel, err := model.RemoveChatRequest(id, a.User.GetAccountID())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, chatChannel)
}
