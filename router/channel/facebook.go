package channel

import (
	"net/http"

	"github.com/labstack/echo"
)

func HandleWebHookFacebookAPI(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
