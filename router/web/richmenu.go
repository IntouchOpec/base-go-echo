package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/module/auth"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
)

func RichmenuListHandler(c *Context) error {
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
	return c.JSON(http.StatusOK, res)
}
