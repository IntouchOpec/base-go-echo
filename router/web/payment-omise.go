package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IntouchOpec/base-go-echo/lib/lineapi"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"

	"github.com/labstack/echo"
)

func PaymentOmiseHandler(c echo.Context) error {
	accountName := c.QueryParam("account_name")
	DocCodeTransaction := c.QueryParam("doc_code_transaction")
	LiffID := c.QueryParam("liff_id")
	var transaction model.Transaction
	var account model.Account
	db := model.DB()
	db.Where("acc_name = ?", accountName).Find(&account)
	db.Where("account_id = ? and tran_document_code = ?", account.ID, DocCodeTransaction).Find(&transaction)
	var chatChannel model.ChatChannel
	db.Find(&chatChannel, transaction.ChatChannelID)

	if transaction.TranStatus == model.TranStatusPaid {
		return c.Render(http.StatusOK, "payment-success", echo.Map{
			"accountName":        accountName,
			"DocCodeTransaction": DocCodeTransaction,
			"detail":             transaction,
			"title":              "ชำระเงินเรียบร้อยแล้ว",
		})
	}

	return c.Render(http.StatusOK, "payment-omise", echo.Map{
		"accountName":        accountName,
		"DocCodeTransaction": DocCodeTransaction,
		"detail":             transaction,
		"LiffID":             LiffID,
	})
}

const (
	OmisePublicKey = "pkey_test_5ip8fflleizk5mzvnut"
	OmiseSecretKey = "skey_test_5ip8nm6pyp7ziztxlh9"
)

func ChargeOmiseHandler(c echo.Context) error {
	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	accountName := c.QueryParam("account_name")
	DocCodeTransaction := c.QueryParam("doc_code_transaction")
	var transaction model.Transaction
	var account model.Account
	var chatChannel model.ChatChannel
	db := model.DB()
	db.Where("acc_name = ?", accountName).Find(&account)
	db.Preload("Bookings", func(db *gorm.DB) *gorm.DB {
		return db.Order("time_start, time_end")
	}).Where("account_id = ? and tran_document_code = ?", account.ID, DocCodeTransaction).Find(&transaction)
	// CreateMester()
	CreateMester(transaction.Bookings, db)
	db.Find(&chatChannel, transaction.ChatChannelID)
	if err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	token := c.FormValue("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{})
	}
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   int64(transaction.TranTotal * 100),
		Currency: "thb",
		Card:     token,
	}
	if err := client.Do(charge, createCharge); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	var omiseLog model.OmiseLog
	ev, err := json.Marshal(charge)
	omiseLog.Json = ev
	omiseLog.AccountID = account.ID

	if charge.Status != omise.ChargeSuccessful {
		fmt.Println("err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Save(&omiseLog).Error; err != nil {
		fmt.Println(err)
	}
	var payment model.Payment
	payment.PayAt = charge.Created
	payment.TransactionID = transaction.ID
	payment.PayStatus = model.PayStatusSuccess
	transaction.TranStatus = model.TranStatusPaid
	payment.PayAmount = transaction.TranTotal
	payment.PayType = model.PayTypeOmise
	if err := db.Model(&transaction).Updates(transaction).Error; err != nil {
		fmt.Print(err)
	}
	if err := db.Model(&transaction).Association("Payments").Append(&payment).Error; err != nil {
		fmt.Print(err)
	}

	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	// var bookingTimeSlot model.BookingTimeSlot
	// var bookingServiceItem model.BookingServiceItem
	var bookingPackage model.BookingPackage
	var list string
	for _, booking := range transaction.Bookings {
		// booking.PlaceID
		// booking.TimeSlot.EmployeeID
		db.Preload("TimeSlot").Find(&booking, booking.ID)
		switch booking.BookingType {
		case model.BookingTypeSlotTime:
			// db.Create()
			// db.Preload("TimeSlot", func(db *gorm.DB) *gorm.DB {
			// 	return db.Preload("EmployeeService", func(db *gorm.DB) *gorm.DB {
			// 		return db.Preload("Employee").Preload("Service")
			// 	})
			// }).Find(&bookingTimeSlot)
			// list += fmt.Sprintf(listTemplate,
			// 	bookingTimeSlot.TimeSlot.EmployeeService.Service.SerName,
			// 	bookingTimeSlot.TimeSlot.EmployeeService.PSPrice)
		case model.BookingTypePackage:
			db.Preload("Package").Find(&bookingPackage)
			list += fmt.Sprintf(listTemplate, bookingPackage.Package.PacName, bookingPackage.Package.PacPrice)
		}

	}
	receiptCard := fmt.Sprintf(receiptTemplate, account.AccName, chatChannel.ChaAddress, list, len(transaction.Bookings), transaction.TranTotal, transaction.TranDocumentCode)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(receiptCard))
	message := linebot.NewFlexMessage("ชำรเงินเสำเร็จ", flexContainer)
	res, err := bot.PushMessage(transaction.TranLineID, message).Do()
	if err != nil {
		fmt.Println(res, err)
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{})
}

func CreateMester(bookings []*model.Booking, db *gorm.DB) error {
	var mtPlas []*model.MasterPlace
	var mtEms []*model.MasterEmployee
	var mtPlasIn []*model.MasterPlace
	var mtEmsIn []*model.MasterEmployee
	var booI model.BookingServiceItem
	var booP model.Package

	for index, boo := range bookings {
		diff := boo.BookedEnd.Sub(boo.BookedStart) / (15 * time.Minute)
		switch boo.BookingType {
		case model.BookingTypeServiceItem:
			db.Where("account_id = ? and m_pla_day = ? m_pla_from BETWEEN ? and ? or m_pla_to BETWEEN ? and ?",
				boo.AccountID, boo.AccountID, boo.BookedStart, boo.BookedEnd, boo.BookedStart, boo.BookedEnd).Find(&mtPlas)
			db.Where("account_id = ? and m_pla_day = ? m_emp_from BETWEEN ? and ? or m_emp_to BETWEEN ? and ?",
				boo.AccountID, boo.AccountID, boo.BookedStart, boo.BookedEnd, boo.BookedStart, boo.BookedEnd).Find(&mtEms)
			db.Preload("ServiceItem", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Service", func(db *gorm.DB) *gorm.DB {
					return db.Preload("Employees", func(db *gorm.DB) *gorm.DB {
						return db.Select("id")
					}).Preload("Places")
				})
			}).Where("booking_id = ?", boo.ID).Find(&booI)
			if booI.EmployeeID != 0 {
				var isFind bool = true
				for _, emp := range booI.ServiceItem.Service.Employees {
					for _, mtEm := range mtEms {
						if mtEm.ID == emp.ID {
							isFind = false
							break
						}
					}
					if isFind {
						bookings[index].BookingServiceItem.EmployeeID = emp.ID
					}
				}
			}
			if booI.TimeSlot. != 0 {

			}
		case model.BookingTypePackage:
			db.Preload("ServiceItems", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Service", func(db *gorm.DB) *gorm.DB {
					return db.Preload("Employees").Preload("Places")
				})
			}).Where("booking_id = ?", boo.ID).Find(&booP)
		}
		// boo.BookingServiceItem.ServiceItemID
		// .Preload("BookingServiceItem").Preload("BookingPackage")
		// db.Find()

		for i := 0; i < int(diff); i++ {

		}
	}

	if err := db.Create(&mtPlasIn).Error; err != nil {
		return err
	}
	if err := db.Create(&mtEmsIn).Error; err != nil {
		return err
	}
	return nil
}

var listTemplate string = `
{ "type": "box", "layout": "horizontal",
	"contents": [
		{ "type": "text", "text": "%s", "size": "sm", "color": "#555555", "flex": 0 },
		{ "type": "text", "text": "฿ %f", "size": "sm", "color": "#111111", "align": "end"}]},`

var receiptTemplate string = `{
	"type": "bubble",
	"body": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "text", "text": "RECEIPT", "weight": "bold", "color": "#1DB446", "size": "sm" },
		{ "type": "text", "text": "%s", "weight": "bold", "size": "xxl", "margin": "md" },
		{ "type": "text", "text": "%s", "size": "xs", "color": "#aaaaaa", "wrap": true },
		{
		  "type": "separator",
		  "margin": "xxl"
		},
		{ "type": "box", "layout": "vertical", "margin": "xxl", "spacing": "sm", "contents": [
			%s
			{ "type": "separator", "margin": "xxl" },
			{ "type": "box", "layout": "horizontal", "margin": "xxl", "contents": [
				{ "type": "text", "text": "ITEMS", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "%d", "size": "sm", "color": "#111111", "align": "end" }
			  ]
			},
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "TOTAL", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "฿ %f", "size": "sm", "color": "#111111", "align": "end" } ]
			}
		  ]
		},
		{ "type": "separator", "margin": "xxl" },
		{ "type": "box", "layout": "horizontal", "margin": "md",
		  "contents": [
			{ "type": "text", "text": "PAYMENT ID", "size": "xs", "color": "#aaaaaa", "flex": 0 },
			{ "type": "text", "text": "#%s", "color": "#aaaaaa", "size": "xs", "align": "end" }
		  ]
		}
	  ]
	},
	"styles": {
	  "footer": { "separator": true }
	}
  }`
