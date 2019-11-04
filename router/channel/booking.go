package channel

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/hb-go/gorm"
	"github.com/labstack/echo"
)

// BookingListHandler
func BookingListHandler(c *Context) error {
	lineID := c.Param("lineID")
	chatChannel := model.ChatChannel{}
	if err := model.DB().Where("line_ID = ?", lineID).Find(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	bookings := []*model.Booking{}

	if err := model.DB().Preload("ServiceSlot", func(db *gorm.DB) *gorm.DB {
		return db.Preload("service")
	}).Preload("Customer").Where("chat_channel_id = ?", chatChannel.ID).Find(&bookings).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "booking-list", echo.Map{
		"list":  bookings,
		"title": "Book",
	})
	return err
}
