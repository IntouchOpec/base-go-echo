package api

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func GetBookingList(c echo.Context) error {
	chatChannelID := c.Param("chatChannelID")
	bookings := model.GetBookingList(chatChannelID)
	return c.JSON(http.StatusOK, bookings)
}

func GetBookingDetail(c echo.Context) error {
	id := c.Param("id")
	booking := model.GetBooking(id)
	return c.JSON(http.StatusOK, booking)
}

func UpdateBooking(c echo.Context) error {
	id := c.Param("id")
	booking := model.Booking{}
	if err := c.Bind(&booking); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	booking.UpdateBooking(id)
	return c.JSON(http.StatusOK, booking)
}

func DeleteBooking(c echo.Context) error {
	id := c.Param("id")
	booking := model.Booking{}

	booking.DeleteBooking(id)
	return c.JSON(http.StatusOK, booking)
}
