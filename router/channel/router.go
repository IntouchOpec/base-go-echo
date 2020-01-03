package channel

import (
	"net/http"

	. "github.com/IntouchOpec/base-go-echo/conf"

	"github.com/labstack/echo"
)

// Routers channel.
func Routers() *echo.Echo {
	e := echo.New()
	e.GET("/", indexHandler)
	e.POST("/callback/:account/:ChannelID", HandleWebHookLineAPI)
	e.POST("/webhook-facebook", HandleWebHookFacebookAPI)
	return e
}

func indexHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "welcome line connect "+Conf.Server.DomainLineChannel)
}
