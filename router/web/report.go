package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func GetReportViewsHandler(c echo.Context) error {
	var transaction model.Transaction
	transactionID := c.QueryParam("transactionID")
	chatChannelID := c.Param("lineID")
	fmt.Println(transactionID, chatChannelID)
	db := model.DB()
	if err := db.Preload("Bookings").Where("chat_channel_id = ?", chatChannelID).Find(&transaction, transactionID).Error; err != nil {
		return c.Render(http.StatusOK, "report-form", echo.Map{
			"transaction": transaction,
		})
	}
	return c.Render(http.StatusOK, "report-form", echo.Map{
		"transaction": transaction,
	})
}

func CreateReportHandler(c echo.Context) error {
	transactionID := c.QueryParam("transactionID")
	ChatChannelID := c.Param("ChatChannelID")
	var transaction model.Transaction
	var report model.Report
	db := model.DB()
	if err := db.Where("chat_channel_id = ?", ChatChannelID).Find(&transaction, transactionID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	report.TransactionID = transaction.ID
	report.Detail = c.FormValue("detail")
	if err := db.Create(&report).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, report)
}
