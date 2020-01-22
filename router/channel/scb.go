package channel

import (
	"net/http"

	"github.com/labstack/echo"
)

func HandlerSBCWebHook(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{})
}
