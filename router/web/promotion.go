package web

import (
	"net/http"
	"time"

	"github.com/IntouchOpec/base-go-echo/lib"
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

	model.DB().Preload("Account").Preload("services").Preload("Customers").Preload("ChatChannels").Where("account_id = ?", a.User.GetAccountID()).Find(&promotion, id)
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

type PromotionForm struct {
	Title         string    `form:"title"`
	TypePromotion string    `form:"type_promotion"`
	Discount      int       `form:"discount"`
	Amount        int       `form:"amount"`
	Code          string    `form:"code"`
	Name          string    `form:"name"`
	StartDate     time.Time `form:"start_time"`
	EndDate       time.Time `form:"end_time"`
	Condition     string    `form:"condition"`
}

func PromotionPostHandler(c *Context) error {
	file := c.FormValue("file")
	a := auth.Default(c)
	promotion := PromotionForm{}
	if err := c.Bind(&promotion); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	file, err := lib.UploadteImage(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	promotionModel := model.Promotion{
		Title:         promotion.Title,
		TypePromotion: promotion.TypePromotion,
		Discount:      promotion.Discount,
		Amount:        promotion.Amount,
		Code:          promotion.Code,
		Name:          promotion.Name,
		StartDate:     promotion.StartDate,
		EndDate:       promotion.EndDate,
		Condition:     promotion.Condition,
		Image:         file,
		AccountID:     a.User.GetAccountID(),
	}

	promotionModel.SavePromotion()
	return c.JSON(http.StatusCreated, promotionModel)
}

func PromotionEditHandler(c *Context) error {
	promotion := model.Promotion{}
	id := c.Param("id")
	a := auth.Default(c)

	model.DB().Preload("Account").Preload("services").Preload("Customers").Preload("ChatChannels").Where("account_id = ?", a.User.GetAccountID()).Find(&promotion, id)
	return c.Render(http.StatusOK, "promotion-form", echo.Map{
		"detail": promotion,
		"title":  "promotion",
	})
}
