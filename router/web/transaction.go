package web

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/IntouchOpec/base-go-echo/module/auth"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

const rejectTemplate string = `
{
	"type": "bubble",
	"hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "https://cdn2.vectorstock.com/i/1000x1000/71/01/reject-red-grunge-round-vintage-rubber-stamp-vector-9347101.jpg" },
	"body": {
	  "type": "box",
	  "layout": "vertical",
	  "spacing": "sm",
	  "contents": [
		{ "type": "text", "text": "ขออภัยในความไม่สะดวก", "wrap": true, "weight": "bold", "size": "xl" },
		{ "type": "box", "layout": "baseline",
		  "contents": [
			{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
	]}]}
}
`

const approveTemplate string = `
{
	"type": "bubble",
	"hero": { "type": "image", "url": "https://previews.123rf.com/images/outchill/outchill1712/outchill171200534/91258245-done-written-text-on-green-round-rubber-vintage-textured-stamp-.jpg", "size": "full", "aspectRatio": "20:13","aspectMode": "cover" },
	"body": { "type": "box", "layout": "vertical",
	  "contents": [
		{ "type": "text", "text": "ยืนยันการจอง", "weight": "bold", "size": "xl" },
		{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "กรุณาไปก่อนเวลานัด 15 นาที", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
			  ]
			}
		  ]
		}
	  ]
	}
  }
`

func TransactionListHandler(c *Context) error {
	Transactions := []*model.Transaction{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	filterTran := db.Where("account_id = ?", a.GetAccountID()).Preload("Customer").Find(&Transactions).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterTran.Limit(pagination.Record).Offset(pagination.Offset).Find(&Transactions)
	return c.Render(http.StatusOK, "transaction-list", echo.Map{
		"title":      "transaction",
		"list":       Transactions,
		"pagination": pagination,
	})
}

func TransactionPatchHandler(c *Context) error {
	id := c.Param("id")
	transaction := model.Transaction{}
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	if err := db.Where("account_id = ?", accID).Preload("ChatChannel").Find(&transaction, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	status := c.FormValue("tran_status")
	var textMassage string
	var flexContainerText string
	switch status {
	case "1":
		transaction.TranStatus = model.TranStatusApproveBooking
		textMassage = "Approve"
		flexContainerText = approveTemplate
	case "2":

	case "-1":
		remark := c.FormValue("tran_remark")
		transaction.TranStatus = model.TranStatusReject
		transaction.TranRemark = remark
		textMassage = "Reject"
		flexContainerText = fmt.Sprintf(rejectTemplate, remark)
	}
	bot, err := linebot.New(transaction.ChatChannel.ChaChannelSecret, transaction.ChatChannel.ChaChannelAccessToken)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexContainerText))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	message := linebot.NewFlexMessage(textMassage, flexContainer)
	res, err := bot.PushMessage(transaction.TranLineID, message).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Save(&transaction).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": transaction,
		"res":  res,
	})
}

type Booking struct {
	Start string
	End   string
	Name  string
	Price float64
}

func TransactionDetailHandler(c *Context) error {
	id := c.Param("id")
	Transaction := model.Transaction{}
	a := auth.Default(c).GetAccountID()
	db := model.DB()

	var bookingTimeSlot model.BookingTimeSlot
	var bookingServiceItem model.BookingServiceItem
	var bookingPackage model.BookingPackage
	var bookings []Booking
	err := db.Preload("Payments").Preload("ChatChannel").Preload("Bookings").Where("account_id = ?", a).Find(&Transaction, id).Error
	if err != nil {
		fmt.Println(err)
	}
	for _, booking := range Transaction.Bookings {
		switch booking.BookingType {
		case model.BookingTypeSlotTime:
			db.Preload("TimeSlot", func(db *gorm.DB) *gorm.DB {
				return db.Preload("EmployeeService", func(db *gorm.DB) *gorm.DB {
					return db.Preload("Employee").Preload("Service")
				})
			}).Find(&bookingTimeSlot)
			bookings = append(bookings, Booking{Start: bookingTimeSlot.TimeSlot.TimeStart, End: bookingTimeSlot.TimeSlot.TimeEnd, Name: bookingTimeSlot.TimeSlot.EmployeeService.Service.SerName})
		case model.BookingTypeServiceItem:
			db.Preload("ServiceItem", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Service")
			}).Find(&bookingServiceItem)
			bookings = append(bookings, Booking{Start: booking.BookedDate.Format("15:01"), End: booking.BookedDate.Format("15:01"), Name: bookingServiceItem.ServiceItem.Service.SerName})

		case model.BookingTypePackage:
			db.Preload("Package").Find(&bookingPackage)
			bookings = append(bookings, Booking{Start: booking.BookedDate.Format("15:01"), End: booking.BookedDate.Format("15:01"), Name: bookingPackage.Package.PacName})
		}
	}
	return c.Render(http.StatusOK, "transaction-detail", echo.Map{
		"title":    "transaction",
		"detail":   Transaction,
		"bookings": bookings,
	})
}

func TransactionCreateHandler(c *Context) error {
	Transaction := model.Transaction{}
	return c.Render(http.StatusOK, "transaction-form", echo.Map{
		"method": "POST",
		"title":  "transaction",
		"detail": Transaction,
	})
}

func TransactionPostHandler(c *Context) error {
	Transaction := model.Transaction{}
	if err := c.Bind(&Transaction); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	accID := auth.Default(c).GetAccountID()
	Transaction.AccountID = accID
	err := Transaction.Create()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	redirect := fmt.Sprintf("/admin/transaction/%d", Transaction.ID)
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
		"data":     Transaction,
	})
}

func TransactionEditHandler(c *Context) error {
	id := c.Param("id")
	Transaction := model.Transaction{}
	a := auth.Default(c)

	model.DB().Preload("Account", "name = ?", a.User.GetAccount()).Find(&Transaction, id)
	return c.Render(http.StatusOK, "transaction-form", echo.Map{"method": "PUT",
		"title":  "transaction",
		"detail": Transaction,
	})
}

func TransactionDeleteHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()

	chatChannel, err := model.RemoveTransaction(id, accID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, chatChannel)
}
