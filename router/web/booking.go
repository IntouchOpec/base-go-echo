package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/hb-go/gorm"
	"github.com/labstack/echo"
)

// BookingListHandler
func BookingListHandler(c *Context) error {
	account := c.Param("account")
	chatChannel := model.ChatChannel{}
	bookings := []*model.Booking{}

	if err := model.DB().Preload("Account", func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", account)
	}).Find(&chatChannel).Error; err != nil {
		return c.Render(http.StatusOK, "booking-list", echo.Map{
			"list":  bookings,
			"title": "Book",
		})
	}
	if err := model.DB().Preload("SubProduct", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Product")
	}).Preload("Customer").Where("chat_channel_id = ?", chatChannel.ID).Find(&bookings).Error; err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "booking-list", echo.Map{
		"list":  bookings,
		"title": "Book",
	})
	return err
}
