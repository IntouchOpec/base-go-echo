package channel

import (
	"github.com/labstack/echo"
)

// Routers channel.
func Routers() *echo.Echo {
	e := echo.New()

	e.POST("/callback/:account/:ChannelID", HandleWebHookLineAPI)

	return e
}
