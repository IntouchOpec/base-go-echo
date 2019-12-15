package channel

import (
	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
	DB          *gorm.DB
	Massage     *linebot.TextMessage
	Account     model.Account
	ChatChannel model.ChatChannel
	ClientLine  *lib.ClientLine
}

type (
	HandlerFunc func(*Context) string
)
