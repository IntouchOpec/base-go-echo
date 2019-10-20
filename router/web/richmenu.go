package web

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
)

func RichMenuListHandler(c *Context) error {
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()

	db.Preload("Account", "ID = ?", a.User.GetAccountID).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	richMenus := []*linebot.RichMenuResponse(res)
	// for _, var := range var {

	// }
	// richMenus = append(richMenus, linebot.RichMenuResponse)
	return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
		"list":  richMenus,
		"title": "rich-menu",
	})
}

func RichMenuDetailHandler(c *Context) error {
	richID := c.Param("richID")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()

	db.Preload("Account", "ID = ?", a.User.GetAccountID).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.GetRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.Render(http.StatusOK, "rich-menu-detail", echo.Map{
		"detail": res,
		"title":  "rich-menu",
	})
}

func RichMenuActiveHandler(c *Context) error {
	richID := c.Param("richID")

	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	qpa := c.QueryParams()
	db.Preload("Account", "ID = ?", a.User.GetAccountID).Where(qpa).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.SetDefaultRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, res)
}

func RichMenuImageHandler(c *Context) error {
	richID := c.Param("richID")

	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	qpa := c.QueryParams()
	db.Preload("Account", "ID = ?", a.User.GetAccountID).Where(qpa).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.SetDefaultRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, res)
}

func RichMenuDeleteHandler(c *Context) error {
	richID := c.Param("richID")

	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	qpa := c.QueryParams()
	db.Preload("Account", "ID = ?", a.User.GetAccountID).Where(qpa).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.DeleteRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, res)
}
