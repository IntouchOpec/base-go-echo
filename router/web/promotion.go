package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func PromotionListHandler(c *Context) error {
	promotions := []*model.Promotion{}

	if err := model.DB().Find(&promotions).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "promotion-list", echo.Map{
		"list":  promotions,
		"title": "promotion",
	})
	return err
}
