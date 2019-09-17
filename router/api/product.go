package api

import (
	"fmt"
	"net/http"
	"strconv"

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

func GetProductList(c echo.Context) error {
	chatchannelID := c.Param("chatChannelID")
	chatChannelIDInt, _ := strconv.Atoi(chatchannelID)

	products := model.GetProduct(chatChannelIDInt)
	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")
	chatchannelID := c.Param("chatChannelID")
	chatChannelIDInt, _ := strconv.Atoi(chatchannelID)
	idInt, _ := strconv.Atoi(id)
	product := model.GetProductByID(chatChannelIDInt, idInt)
	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	product := model.Product{}
	if err := c.Bind(&product).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	product.UpdateProduct(idInt)
	return c.JSON(http.StatusOK, product)
}

func UpdateSubProduct(c echo.Context) error {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	subCustomer := model.SubProduct{}
	if err := c.Bind(&subCustomer).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	subCustomer.UpdateSubProduct(idInt)

	return c.JSON(http.StatusOK, subCustomer)
}
