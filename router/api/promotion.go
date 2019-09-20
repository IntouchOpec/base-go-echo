package api

import (
	"net/http"
	"strconv"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func GetPromotionList(c echo.Context) error {
	chatChannelID := c.Param("chatChannelID")
	chatChannelIDInt, _ := strconv.Atoi(chatChannelID)
	promtions := model.GetPromotionList(chatChannelIDInt)
	return c.JSON(http.StatusOK, promtions)
}

func GetPromotion(c echo.Context) error {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	promotion := model.GetPromotion(idInt)
	return c.JSON(http.StatusOK, promotion)
}

func CreatePromotion(c echo.Context) error {
	chatChannelID := c.Param("chatChannelID")
	chatChannelIDInt, _ := strconv.ParseUint(chatChannelID, 10, 64)

	promotion := model.Promotion{ChatChannelID: uint(chatChannelIDInt)}

	if err := c.Bind(&promotion).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	promotion.SavePromotion()
	return c.JSON(http.StatusOK, promotion)
}

func UpdatePromotion(c echo.Context) error {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	promotion := model.Promotion{}
	if err := c.Bind(&promotion).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	promotion.UpdatePromotion(idInt)
	return c.JSON(http.StatusOK, promotion)
}

func DeletePromotion(c echo.Context) error {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	promotion := model.DeletePromotion(idInt)
	return c.JSON(http.StatusOK, promotion)
}