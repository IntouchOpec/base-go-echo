package api

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// CreateProduct is route handle for create product
func CreateProduct(c echo.Context) error {
	product := model.Product{}
	if err := c.Bind(&product); err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	product.SaveProduct()
	c.JSON(200, product)

	return nil
}
