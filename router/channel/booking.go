package channel

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	// . "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"

	// . "github.com/IntouchOpec/base-go-echo/conf"

	"github.com/IntouchOpec/base-go-echo/lib/lineapi"
	"github.com/IntouchOpec/base-go-echo/model"
)

type Booking struct {
	WeekDay int
	deta    *time.Time
}

func BookingHandler(c *Context) error {
	fmt.Println(c.PostbackAction.Type, "======", c.PostbackAction.ServiceItemID)
	var template string
	// var textMessage string
	var d time.Time
	var err error

	if c.PostbackAction.PackageID != "" {
		var pack model.Package
		c.DB.Preload("ServiceItems").Find(&pack, c.PostbackAction.PackageID)
		for _, serI := range pack.ServiceItems {
			fmt.Println(serI.ServiceID, "Employees")
		}
	} else if c.PostbackAction.ServiceItemID != "" {
		var serI model.ServiceItem
		var isEmPla bool
		var isEmEmp bool
		var msPla []*model.MasterEmployee
		var msEmp []*model.MasterPlace
		c.DB.Preload("Service", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Employees").Preload("Places")
		}).Find(&serI, c.PostbackAction.ServiceItemID)
		switch c.PostbackAction.Type {
		case "now":
			d = time.Now()
			// var tran model.Transaction
			var timeSs []model.TimeSlot
			var empIDs []uint
			// var isFindSlot bool = true
			var cout int
			var pla model.Place
			for _, emp := range serI.Service.Employees {
				empIDs = append(empIDs, emp.ID)
			}
			ho := d.Hour()
			mi := d.Minute()
			if mi > 0 && mi <= 15 {
				mi = 45
			} else if mi > 15 && mi <= 30 {
				mi = 00
				ho++
			} else if mi > 30 && mi <= 45 {
				ho++
				mi = 30
			} else {
				ho++
				mi = 45
			}
			start, _ := time.Parse("2006-01-02", "1970-01-01")
			start = start.Add((time.Duration(mi) * time.Minute) + (time.Duration(ho-7) * time.Hour))
			end := start.Add(serI.SSTime)
			c.DB.Model(&timeSs).Where(
				"account_id = ? and time_day = ? and employee_id in (?) and time_start BETWEEN ? and ? or time_end BETWEEN ? and ? ",
				c.Account.ID, int(d.Weekday()), empIDs, start, end, start, end).Count(&cout).Find(&timeSs)
			c.DB.Where("account_id = ? and m_pla_day = ? m_pla_from BETWEEN ? and ? or m_pla_to BETWEEN ? and ?", c.Account.ID, int(d.Weekday()), start, end, start, end).Find(&msPla)
			c.DB.Where("account_id = ? and m_pla_day = ? m_emp_from BETWEEN ? and ? or m_emp_to BETWEEN ? and ?", c.Account.ID, int(d.Weekday()), start, end, start, end).Find(&msEmp)
			// validate place and employee

			isEmPla = len(msPla) == 0
			isEmEmp = len(msEmp) == 0
			fmt.Println(timeSs, len(timeSs), cout)
			if len(timeSs) == 0 {
				return errors.New("")
			}
			tx := c.DB.Begin()
			tran, err := bindTransaction(c, tx, serI.SSPrice)
			if err != nil {
				tx.Rollback()
			}
			booking, err := bindBooking(c, tx, tran.ID, pla.ID, model.BookingTypeServiceItem, d, start, end)
			if err != nil {
				tx.Rollback()
			}
			_, err = bindBookingServiceItem(c, tx, booking.ID, serI.ID, timeSs[0].ID)
			if err != nil {
				tx.Rollback()
			}
			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
			}
			setting := c.ChatChannel.GetSetting([]string{model.NameLIFFIDPayment})
			template = fmt.Sprintf(checkoutTemplate,
				c.ChatChannel.ChaImage,
				serI.Service.SerName,
				c.ChatChannel.ChaAddress,
				d.Format("2006-01-02"),
				start.Add(7*time.Hour).Format("15:02"),
				end.Add(7*time.Hour).Format("15:02"),
				setting[model.NameLIFFIDPayment],
				c.Account.AccName,
				tran.TranDocumentCode,
				setting[model.NameLIFFIDPayment])
		case "appointment":
			d, err = time.Parse("2006-01-02", c.Event.Postback.Params.Date)
			if err != nil {
				fmt.Println(err)
			}
			c.DB.Where("account_id = ? and m_pla_day = ?", c.Account.ID, d.Weekday()).Find(&msPla)
			c.DB.Where("account_id = ? and m_pla_day = ?", c.Account.ID, d.Weekday()).Find(&msEmp)
			isEmPla = len(msPla) == 0
			isEmEmp = len(msEmp) == 0
			if isEmPla && isEmEmp {

			}
			for _, emp := range serI.Service.Employees {
				var timeS model.TimeSlot
				c.DB.Where("time_day = ? and employee_id = ?", int(d.Weekday()), emp.ID).Find(&timeS)
				slot := timeS.TimeEnd.Sub(timeS.TimeStart) / (60 * time.Minute)
				var button string
				var cont string
				for i := 0; i < int(slot); i++ {
					tim := timeS.TimeStart.Add((time.Duration(60*i) * time.Minute) + (7 * time.Hour))
					ho := tim.Hour()
					mi := tim.Minute()
					var hostr string
					var mistr string
					if ho < 10 {
						hostr = fmt.Sprintf("0%d", tim.Hour())
					} else {
						hostr = fmt.Sprintf("%d", tim.Hour())
					}
					if mi < 10 {
						mistr = fmt.Sprintf("0%d", tim.Minute())
					} else {
						mistr = fmt.Sprintf("%d", tim.Minute())
					}
					bookingTime := fmt.Sprintf("%s:%s", hostr, mistr)
					button += fmt.Sprintf(buttonTimeSlotTemplate,
						bookingTime,
						fmt.Sprintf("action=%s&service_item_id=%d&employee_id=%d&date=%s&time=%s&time_slot_id=%d",
							"checkout", serI.ID, emp.ID, c.Event.Postback.Params.Date, bookingTime, timeS.ID)) + ","
					if i%2 != 0 {
						cont += fmt.Sprintf(layoutTimeSlotTemplate, button[:len(button)-1]) + ","
						button = ""
					}
				}
				template += fmt.Sprintf(timeSlotTemplate, cont[:len(cont)-1]) + ","
			}
			template = fmt.Sprintf(carouselTemplate, template[:len(template)-1])
		}
	} else {

	}
	template = fmt.Sprintf(`{ "replyToken": "%s", "messages":[ { "type": "flex",  "altText":  "รายการบริการ",  "contents": %s }]}`, c.Event.ReplyToken, template)
	lineapi.SendMessageCustom("reply", c.ChatChannel.ChaChannelAccessToken, template)
	return nil
}

func bookingNow(c *Context, serI model.ServiceItem) (string, error) {
	d := time.Now()
	var msPla []*model.MasterEmployee
	var msEmp []*model.MasterPlace
	// var tran model.Transaction
	var timeSs []model.TimeSlot
	var empIDs []uint
	// var isFindSlot bool = true
	var cout int
	var pla model.Place
	for _, emp := range serI.Service.Employees {
		empIDs = append(empIDs, emp.ID)
	}
	ho := d.Hour()
	mi := d.Minute()
	if mi > 0 && mi <= 15 {
		mi = 45
	} else if mi > 15 && mi <= 30 {
		mi = 00
		ho++
	} else if mi > 30 && mi <= 45 {
		ho++
		mi = 30
	} else {
		ho++
		mi = 45
	}
	start, _ := time.Parse("2006-01-02", "1970-01-01")
	start = start.Add((time.Duration(mi) * time.Minute) + (time.Duration(ho-7) * time.Hour))
	end := start.Add(serI.SSTime)
	c.DB.Model(&timeSs).Where(
		"account_id = ? and time_day = ? and employee_id in (?) and time_start BETWEEN ? and ? or time_end BETWEEN ? and ? ",
		c.Account.ID, int(d.Weekday()), empIDs, start, end, start, end).Count(&cout).Find(&timeSs)
	c.DB.Where("account_id = ? and m_pla_day = ? m_pla_from BETWEEN ? and ? or m_pla_to BETWEEN ? and ?", c.Account.ID, int(d.Weekday()), start, end, start, end).Find(&msPla)
	c.DB.Where("account_id = ? and m_emp_day = ? m_emp_from BETWEEN ? and ? or m_emp_to BETWEEN ? and ?", c.Account.ID, int(d.Weekday()), start, end, start, end).Find(&msEmp)
	// validate place and employee

	// isEmPla = len(msPla) == 0
	// isEmEmp = len(msEmp) == 0
	fmt.Println(timeSs, len(timeSs), cout)
	if len(timeSs) == 0 {
		return "", errors.New("")
	}
	tx := c.DB.Begin()
	tran, err := bindTransaction(c, tx, serI.SSPrice)
	if err != nil {
		tx.Rollback()
	}
	booking, err := bindBooking(c, tx, tran.ID, pla.ID, model.BookingTypeServiceItem, d, start, end)
	if err != nil {
		tx.Rollback()
	}
	_, err = bindBookingServiceItem(c, tx, booking.ID, serI.ID, timeSs[0].ID)
	if err != nil {
		tx.Rollback()
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	setting := c.ChatChannel.GetSetting([]string{model.NameLIFFIDPayment})
	return fmt.Sprintf(checkoutTemplate,
		c.ChatChannel.ChaImage,
		serI.Service.SerName,
		c.ChatChannel.ChaAddress,
		d.Format("2006-01-02"),
		start.Add(7*time.Hour).Format("15:02"),
		end.Add(7*time.Hour).Format("15:02"),
		setting[model.NameLIFFIDPayment],
		c.Account.AccName,
		tran.TranDocumentCode,
		setting[model.NameLIFFIDPayment]), nil
}

func ChackOutHandler(c *Context) (linebot.SendingMessage, error) {
	var flexContainerStr string
	var serI model.ServiceItem
	var pla model.Place
	var slot model.TimeSlot
	d := time.Now()
	if err := c.DB.Preload("Service").Find(&serI, c.PostbackAction.ServiceItemID).Error; err != nil {
		return nil, errors.New("error")
	}
	if err := c.DB.Find(&slot, c.PostbackAction.TimeSlotID).Error; err != nil {
		return nil, errors.New("error")
	}
	day, err := time.Parse("2006-01-02", c.PostbackAction.DateStr)
	fmt.Println("err", err)
	start, err := time.Parse("2006-01-02 15:02:05", "2020-01-02 "+c.PostbackAction.TimeStr+":05")
	fmt.Println("err2", err)
	end := start.Add(serI.SSTime)
	tx := c.DB.Begin()
	tran, err := bindTransaction(c, tx, serI.SSPrice)
	if err != nil {
		return nil, err
	}
	b, err := bindBooking(c, tx, tran.ID, pla.ID, model.BookingTypeServiceItem, day, start, end)
	if err != nil {
		return nil, err
	}
	_, err = bindBookingServiceItem(c, tx, b.ID, serI.ID, slot.ID)
	tx.Commit()
	setting := c.ChatChannel.GetSetting([]string{model.NameLIFFIDPayment})
	fmt.Println(start.Hour(), c.PostbackAction.TimeStr)
	flexContainerStr = fmt.Sprintf(checkoutTemplate,
		c.ChatChannel.ChaImage,
		serI.Service.SerName,
		c.ChatChannel.ChaAddress,
		d.Format("2006-01-02"),
		fmt.Sprintf("%d", start.Hour()+7),
		fmt.Sprintf("%d", end.Hour()+7),
		setting[model.NameLIFFIDPayment],
		c.Account.AccName,
		tran.TranDocumentCode,
		setting[model.NameLIFFIDPayment])
	// flexContainerStr = fmt.Sprintf(carouselTemplate, flexContainerStr)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexContainerStr))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), nil
}

func bindTransaction(c *Context, tx *gorm.DB, Total float64) (*model.Transaction, error) {
	var tran model.Transaction
	fmt.Println(c.Event.Source.UserID)
	tran.ChatChannelID = c.ChatChannel.ID
	tran.TranTotal = Total
	tran.AccountID = c.ChatChannel.AccountID
	tran.CustomerID = c.Customer.ID
	tran.TranLineID = c.Event.Source.UserID
	tran.TranStatus = model.TranStatusPanding
	if err := tx.Create(&tran).Error; err != nil {
		return nil, err
	}
	return &tran, nil
}

func bindMSPlace(c *Context, placeID uint, placAmount int) (*model.MasterPlace, error) {
	var MSPlace model.MasterPlace
	day, err := time.Parse("2006-01-02", c.PostbackAction.DateStr)
	if err != nil {
		return nil, err
	}
	MSPlace.PlaceID = placeID
	MSPlace.MPlaDay = day
	MSPlace.AccountID = c.Account.ID
	form, err := time.Parse("15:04", c.PostbackAction.TimeStr)
	if err != nil {
		return nil, err
	}
	to, err := time.Parse("15:04", c.PostbackAction.TimeStr)
	if err != nil {
		return nil, err
	}
	MSPlace.MPlaFrom = form
	MSPlace.MPlaTo = to
	// MSPlace.Mp = MSPlace.MPlaAmount + 1
	// if MSPlace.MPlaAmount == placAmount {
	// 	MSPlace.MPlaStatus = model.MPlaStatusBusy
	// }
	return &MSPlace, nil
}

func bindBooking(c *Context, tx *gorm.DB, transactionID, placeID uint, bookingType model.BookingType, day, start, end time.Time) (*model.Booking, error) {
	var b model.Booking
	b.TransactionID = transactionID
	b.PlaceID = placeID
	b.BookingType = bookingType
	b.ChatChannelID = c.ChatChannel.ID
	b.BooLineID = c.Event.Source.UserID
	b.CustomerID = c.Customer.ID
	b.BooStatus = model.BookingStatusPandding
	b.AccountID = c.Account.ID
	b.BookedStart = start
	b.BookedEnd = end
	b.BookedDate = day
	if err := tx.Create(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func bindBookingServiceItem(c *Context, tx *gorm.DB, bookID, serviceItemID, timeSlotID uint) (*model.BookingServiceItem, error) {
	var bookingServiceItem model.BookingServiceItem
	bookingServiceItem.BookingID = bookID
	bookingServiceItem.ServiceItemID = serviceItemID
	bookingServiceItem.TimeSlotID = timeSlotID
	if err := tx.Create(&bookingServiceItem).Error; err != nil {
		return nil, err
	}
	return &bookingServiceItem, nil
}

func bindBookingPackage(c *Context, bookID, packageID uint) (model.BookingPackage, error) {
	var bookingPackage model.BookingPackage
	bookingPackage.BookingID = bookID
	bookingPackage.PackageID = packageID
	u64, err := strconv.ParseUint(c.PostbackAction.TimeSlotID, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	bookingPackage.TimeSlotID = uint(u64)
	return bookingPackage, nil
}
