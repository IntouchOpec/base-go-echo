package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func SettingHandler(c *Context) error {
	promotions := []*model.Promotion{}

	if err := model.DB().Find(&promotions).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "setting", echo.Map{
		"title": "setting",
	})
	return err
}
