package channel

import (
	"net/http"

	"github.com/labstack/echo"
)

func DashboardHandler(c *Context) error {
	err := c.Render(http.StatusOK, "dashboard", echo.Map{
		"title": "dashboard",
	})
	return err
}
