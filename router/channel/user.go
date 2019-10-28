package channel

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// UserListHandler
func UserListHandler(c *Context) error {
	lineID := c.Param("lineID")
	chatChannel := model.ChatChannel{}
	if err := model.DB().Where("line_ID = ?", lineID).Find(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	customers := []*model.Customer{}

	if err := model.DB().Where("chat_channel_id = ?", chatChannel.ID).Find(&customers).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "customer-list", echo.Map{
		"list":  customers,
		"title": "customer",
	})
	return err
}
