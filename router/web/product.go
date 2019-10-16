package web

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

// ProductListHandler
func ProductListHandler(c *Context) error {
	products := []*model.Product{}
	a := auth.Default(c)
	model.DB().Preload("ChatChannel", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Account", "name = ?", a.User.GetAccountID())
	}).Preload("SubProduct").Find(&products)
	err := c.Render(http.StatusOK, "product-list", echo.Map{
		"list":  products,
		"title": "product",
	})
	return err
}

func ProductDetailHandler(c *Context) error {
	product := model.Product{}
	fmt.Println("====")
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Preload("Account").Preload("SubProducts").Preload("ChatChannels").Where("account_id = ? ", a.User.GetAccountID()).Find(&product, id)
	err := c.Render(http.StatusOK, "product-detail", echo.Map{
		"detail": product,
		"title":  "product",
	})
	return err
}
