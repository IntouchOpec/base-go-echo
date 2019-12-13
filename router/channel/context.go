package channel

import (
	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"

	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
	Massage     string
	Account     model.Account
	ChatChannel model.ChatChannel
	ClientLine  *lib.ClientLine
}

type (
	HandlerFunc func(*Context) string
)
