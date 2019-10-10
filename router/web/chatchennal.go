package web

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/jinzhu/gorm"

	"github.com/IntouchOpec/base-go-echo/model"
)

func ChatChannelListHandler(c *Context) error {
	account := c.Param("account")
	chatChannels := []*model.ChatChannel{}
	model.DB().Preload("Account", func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", account)
	}).Find(&chatChannels)
	return c.Render(http.StatusOK, "chatchannel-list", echo.Map{
		"title": "chatchannel",
		"list":  chatChannels,
	})
}
