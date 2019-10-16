package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
)

// ChatChannelListHandler
func ChatChannelListHandler(c *Context) error {
	chatChannels := []*model.ChatChannel{}
	a := auth.Default(c)
	model.DB().Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannels)
	return c.Render(http.StatusOK, "chat-channel-list", echo.Map{
		"title": "chat_channel",
		"list":  chatChannels,
	})
}

func ChatChannelGetChannelAccessTokenHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	db.Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id)
	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	res, err := bot.IssueAccessToken(chatChannel.ChannelID, chatChannel.ChannelSecret).Do()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	chatChannel.ChannelAccessToken = res.AccessToken

	if err := db.Save(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, res)
}

// ChatChannelDetailHandler
func ChatChannelDetailHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	model.DB().Preload("Products").Preload("Customers").Preload("ActionLogs").Preload("EventLogs").Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id)
	customerSum := len(chatChannel.Customers)
	productSum := len(chatChannel.Products)
	return c.Render(http.StatusOK, "chat-channel-detail", echo.Map{
		"title":       "chat_channel",
		"detail":      chatChannel,
		"customerSum": customerSum,
		"productSum":  productSum,
	})
}

// ChatChannelCreateHandler
func ChatChannelCreateViewHandler(c *Context) error {
	typeChatChannels := []string{"Facebook", "Line"}
	csrfValue := c.Get("_csrf")
	return c.Render(http.StatusOK, "chat-channel-form", echo.Map{
		"title":            "chat_channel",
		"typeChatChannels": typeChatChannels,
		"mode":             "Create",
		"_csrf":            csrfValue,
	})
}

// ChatChannelCreatePostHandler
func ChatChannelCreatePostHandler(c *Context) error {
	chatChannel := model.ChatChannel{}
	if err := c.Bind(&chatChannel).Error; err != nil {
		fmt.Println("errr", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, chatChannel)
}

// ChatChannelEditHandler
func ChatChannelEditHandler(c *Context) error {
	id := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	db.Preload("Products").Preload("Customers").Preload("ActionLogs").Preload("EventLogs").Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel, id)
	customerSum := len(chatChannel.Customers)
	productSum := len(chatChannel.Products)
	typeChatChannels := []string{"Facebook", "Line"}
	return c.Render(http.StatusOK, "chat-channel-form", echo.Map{
		"title":            "chat_channel",
		"detail":           chatChannel,
		"customerSum":      customerSum,
		"productSum":       productSum,
		"typeChatChannels": typeChatChannels,
		"mode":             "Edit",
	})
}
