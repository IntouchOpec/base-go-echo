package api

import (
	"net/http"
	"strconv"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func CreatePromotion(c echo.Context) error {
	lineID := c.Param("lineID")

	chatChanne := model.ChatChannel{}
	if err := model.DB().Preload("Account").Find(&chatChanne, lineID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	promotion := model.Promotion{ChatChannelID: chatChanne.ID, AccountID: chatChanne.AccountID}

	if err := c.Bind(&promotion); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	promotion.SavePromotion()
	return c.JSON(http.StatusOK, promotion)
}

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
