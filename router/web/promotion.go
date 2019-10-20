package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/module/auth"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// PromotionListHandler
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

// PromotionDetailHandler
func PromotionDetailHandler(c *Context) error {
	promotion := model.Promotion{}
	id := c.Param("id")
	a := auth.Default(c)

	model.DB().Preload("Account").Preload("Products").Preload("Customers").Preload("ChatChannels").Where("account_id = ?", a.User.GetAccountID()).Find(&promotion, id)
	sumCustomer := len(promotion.Customers)
	return c.Render(http.StatusOK, "promotion-detail", echo.Map{
		"detail":      promotion,
		"title":       "promotion",
		"sumCustomer": sumCustomer,
	})
}

func PromotionFormHandler(c *Context) error {
	promotion := model.Promotion{}
	promotionTypes := []model.PromotionType{model.PromotionTypePromotion, model.PromotionTypeCoupon, model.PromotionTypeVoucher}
	return c.Render(http.StatusOK, "promotion-form", echo.Map{
		"detail":         promotion,
		"title":          "promotion",
		"promotionTypes": promotionTypes,
	})
}

func PromotionPostHandler(c *Context) error {
	promotion := model.Promotion{}
	if err := c.Bind(&promotion).Error; err != nil {

	}
	return c.JSON(http.StatusCreated, "")
}

func PromotionEditHandler(c *Context) error {
	promotion := model.Promotion{}
	id := c.Param("id")
	a := auth.Default(c)

	model.DB().Preload("Account").Preload("Products").Preload("Customers").Preload("ChatChannels").Where("account_id = ?", a.User.GetAccountID()).Find(&promotion, id)
	return c.Render(http.StatusOK, "promotion-form", echo.Map{
		"detail": promotion,
		"title":  "promotion",
	})
}
