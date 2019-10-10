package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func PromotionListHandler(c *Context) error {
	lineID := c.Param("lineID")
	chatChannel := model.ChatChannel{}
	if err := model.DB().Where("line_ID = ?", lineID).Find(&chatChannel).Error; err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	promotions := []*model.Promotion{}

	if err := model.DB().Where("chat_channel_id = ?", chatChannel.ID).Find(&promotions).Error; err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "promotion-list", echo.Map{
		"list":  promotions,
		"title": "promotion",
	})
	return err
}
