package channel

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/lib/lineapi"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Booking struct {
	WeekDay int
	deta    *time.Time
}

// if err := c.DB.Preload("Service", func(db *gorm.DB) *gorm.DB {
// 	return db.Preload("Employees", func(db *gorm.DB) *gorm.DB {
// 		return db.Order("id")
// 	}).Preload("Places", func(db *gorm.DB) *gorm.DB {
// 		return db.Order("id")
// 	})
// }).Find(&serI, c.PostbackAction.ServiceItemID).Error; err != nil {
// 	return err
// }
type se struct {
	id       uint
	emploIDs []uint
	placeIDs []uint
}

func BookingHandler(c *Context) error {
	fmt.Println(c.PostbackAction.Type, "======", c.PostbackAction.ServiceItemID)
	setting := c.ChatChannel.GetSetting([]string{model.NameLIFFIDPayment})
	var template string
	if c.PostbackAction.PackageID != "" {
		var p model.Pack
		var serIDs []string
		rows, err := c.sqlDb.Query(`
		SELECT 
			pa.id, pac_name, si.service_id , pa.pac_price, pa.pac_time
		FROM packages AS pa
		INNER JOIN package_service_item AS psi ON psi.package_id = pa.id
		INNER JOIN service_items AS si ON si.id = psi.service_item_id AND si.deleted_at IS NULL
		INNER JOIN services AS s ON s.id = si.service_id AND s.deleted_at IS NULL
		WHERE pa.deleted_at IS NULL AND pa.id = $1 AND pac_is_active = true`, c.PostbackAction.PackageID)
		if err != nil {
			return err
		}
		for rows.Next() {
			var serID string
			rows.Scan(&p.ID, &p.Name, &serID, &p.Price, &p.TimeUse)
			serIDs = append(serIDs, serID)
		}
		switch c.PostbackAction.Type {
		case "now":
			template, err = PackNow(c, p, serIDs, setting)
			if err != nil {
				return err
			}
		case "appointment":
			template = PackAppointment(c, p)
		}
	} else if c.PostbackAction.ServiceItemID != "" {
		var serI model.ServiceItem
		row := c.sqlDb.QueryRow(`
			SELECT
				s.id AS service_id, ser_name, ser_image,
				si.id AS service_item_id, ss_price, ss_time,
				s.account_id, si.account_id
			FROM service_items AS si
			INNER JOIN services AS s ON s.id = si.service_id AND s.deleted_at IS NULL 
			WHERE si.deleted_at IS NULL AND si.id = $1`, c.PostbackAction.ServiceItemID)

		var ser model.Service
		err := row.Scan(&ser.ID, &ser.SerName, &ser.SerImage, &serI.ID, &serI.SSPrice, &serI.SSTime, &ser.AccountID, &serI.AccountID)
		if err != nil {
			fmt.Println("err", err)
			return err
		}
		serI.Service = &ser

		switch c.PostbackAction.Type {
		case "now":
			b, err := bindBookingServiceItemNow(c, model.BookingTypeServiceItem)
			if err != nil {
				fmt.Println("err1", err)
				return err
			}
			if err := b.Now(c.sqlDb, serI); err != nil {
				fmt.Println("err2", err)

				return err
			}
			tran, err := bindTransaction(c, serI.SSPrice)
			if err != nil {
				fmt.Println("err3", err)
				return err
			}
			tx := c.DB.Begin()
			if err := tran.LineBooking(tx); err != nil {
				fmt.Println("err4", err)
				tx.Rollback()
				return err
			}
			if err := tran.LineBookingServiceNow(tx, b); err != nil {
				fmt.Println(err)
				tx.Rollback()
				return err
			}
			if err := tx.Commit().Error; err != nil {
				return err
			}
			template = fmt.Sprintf(checkoutTemplate,
				c.ChatChannel.ChaImage,
				ser.SerName,
				c.ChatChannel.ChaAddress,
				b.BookedDay.Format("2006-01-02"),
				b.BookedStart.Add(7*time.Hour).Format("15:04"),
				b.BookedEnd.Add(7*time.Hour).Format("15:04"),
				setting[model.NameLIFFIDPayment],
				c.Account.AccName,
				tran.TranDocumentCode,
				setting[model.NameLIFFIDPayment])
		case "appointment":
			template = ServiceItemAppointment(c, serI)
		}
	}

	template = fmt.Sprintf(`{ "replyToken": "%s", "messages":[ { "type": "flex",  "altText":  "รายการบริการ",  "contents": %s }]}`, c.Event.ReplyToken, template)
	err := lineapi.SendMessageCustom("reply", c.ChatChannel.ChaChannelAccessToken, template)
	if err != nil {
		return nil
	}
	return nil
}

type emplo struct {
	id        string
	timeEnd   time.Time
	timeStart time.Time
	tsID      uint
	image     string
}

func PackNow(c *Context, p model.Pack, serIDs []string, setting map[string]string) (string, error) {
	b, err := bindBookingPackageNow(c, model.BookingTypePackage)
	if err != nil {
		fmt.Println(err, ":err1")
		return "", err
	}
	if err := b.PackNow(c.sqlDb, p, serIDs); err != nil {
		fmt.Println(err, ":err2")
		return "", err
	}
	fmt.Println(err)
	if err != nil {
		fmt.Println(err, ":err3")
	}
	t, err := bindTransaction(c, p.Price)
	if err != nil {
		fmt.Println(err, ":err4")
	}
	tx, err := c.sqlDb.Begin()
	if err := t.CreateSql(tx); err != nil {
		fmt.Println(err, ":err5")
		tx.Rollback()
		return "", err
	}
	b.TransactionID = t.ID
	if err := b.CreateSql(tx); err != nil {
		fmt.Println(err, ":err6")
		tx.Rollback()
		return "", err
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err, ":err7")
		return "", err
	}
	return fmt.Sprintf(checkoutTemplate,
		c.ChatChannel.ChaImage,
		p.Name,
		c.ChatChannel.ChaAddress,
		b.BookedDay.Format("2006-01-02"),
		b.BookedStart.Add(7*time.Hour).Format("15:04"),
		b.BookedEnd.Add(7*time.Hour).Format("15:04"),
		setting[model.NameLIFFIDPayment],
		c.Account.AccName,
		t.TranDocumentCode,
		setting[model.NameLIFFIDPayment]), nil
}

func chackDateBooking(d time.Time) error {
	now := time.Now()
	monthBoo := d.Month()
	monthNow := now.Month()
	if monthBoo < monthNow {
		return errors.New("")
	} else if monthBoo == monthNow {
		if d.Day() < now.Day() {
			return errors.New("")
		}
	}
	return nil
}

func PackAppointment(c *Context, p model.Pack) string {
	d, err := time.Parse("2006-01-02", c.Event.Postback.Params.Date)
	if err != nil {
		fmt.Println("date err")
		// return ""
	}
	if err := chackDateBooking(d); err != nil {
		return ""
	}
	return ""
}
func ServiceItemAppointment(c *Context, serI model.ServiceItem) string {
	d, err := time.Parse("2006-01-02", c.Event.Postback.Params.Date)
	if err != nil {
		fmt.Println("date err")
		// return ""
	}
	var emplos []emplo
	var template string
	if err := chackDateBooking(d); err != nil {
		return ""
	}

	rows, err := c.sqlDb.Query(`	
			SELECT 
				es.employee_id ,
				ts.time_end, ts.time_start, ts.id,
				e.empo_image
			FROM employee_service AS es
			INNER JOIN employees AS e ON e.id = es.employee_id AND e.deleted_at IS NULL
			INNER JOIN time_slots AS ts ON ts.employee_id = es.employee_id AND ts.deleted_at IS NULL
				AND ts.deleted_at IS NULL 
				AND time_day = $1 
				AND ts.time_active = true
			WHERE es.service_id = $2
			ORDER BY es.employee_id`, int(d.Weekday()), serI.Service.ID)
	if err != nil {
		fmt.Println("sql find time slot err")
		// now not have employee work
		return ""
	}
	if err == nil {
		for rows.Next() {
			var empl emplo
			rows.Scan(&empl.id, &empl.timeEnd, &empl.timeStart, &empl.tsID, &empl.image)
			emplos = append(emplos, empl)
		}
	}

	rows, err = c.sqlDb.Query(`
			SELECT 
				pl.id, pl.plac_amount
			FROM places AS pl
			INNER JOIN place_service AS ps ON ps.place_id = pl.id AND service_id = $1
			WHERE pl.deleted_at IS NULL 
			GROUP BY pl.id, pl.plac_amount
			ORDER BY pl.id`, serI.Service.ID)
	if err != nil {
		fmt.Println(err)
		// now not have place
		return ""
	}
	for _, emp := range emplos {
		slot := emp.timeEnd.Sub(emp.timeStart) / (60 * time.Minute)
		var button string
		var cont string
		for i := 0; i < int(slot); i++ {
			tim := emp.timeStart.Add((time.Duration(60*i) * time.Minute) + (7 * time.Hour))
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
			// fmt.Println(bookingTime, tim.Minute())
			action := fmt.Sprintf("action=%s&service_item_id=%d&employee_id=%s&date=%s&time=%s&time_slot_id=%d",
				"checkout", serI.ID, emp.id, c.Event.Postback.Params.Date, bookingTime, emp.tsID)
			button += fmt.Sprintf(buttonTimeSlotTemplate,
				bookingTime, action) + ","
			if i%2 != 0 {
				cont += fmt.Sprintf(layoutTimeSlotTemplate, button[:len(button)-1]) + ","
				button = ""
			} else if int(slot)-1 == i {
				cont += fmt.Sprintf(layoutTimeSlotTemplate, button[:len(button)-1]) + ","
			}
		}
		template += fmt.Sprintf(timeSlotTemplate, fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, emp.image), cont[:len(cont)-1]) + ","
	}
	return fmt.Sprintf(carouselTemplate, template[:len(template)-1])
}

func ChackOutHandler(c *Context) (linebot.SendingMessage, error) {
	var flexContainerStr string
	var serI model.ServiceItem
	// var slot model.TimeSlot
	fmt.Println(c.PostbackAction.ServiceItemID)
	// c.sqlDb.Query("SELECT  FROM service_items AS si INNER JOIN services AS s ON s.id = si.service_item_id AND s.deleted_at IS NULL WHERE si.deleted_at IS NULL")
	if err := c.DB.Preload("Service").Find(&serI, c.PostbackAction.ServiceItemID).Error; err != nil {
		fmt.Println("err find service")
		return nil, errors.New("error")
	}
	// if err := c.DB.Find(&slot, c.PostbackAction.TimeSlotID).Error; err != nil {
	// 	return nil, errors.New("error")
	// }
	tx := c.DB.Begin()
	b, err := bindBookingServiceItemAppointment(c, model.BookingTypeServiceItem, serI)
	tran, err := bindTransaction(c, serI.SSPrice)
	if err := tran.LineBooking(tx); err != nil {
		fmt.Println(err)
		tx.Rollback()
		return nil, err
	}
	if err := tran.LineBookingServiceAppointment(tx, b); err != nil {
		fmt.Println(err)
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	setting := c.ChatChannel.GetSetting([]string{model.NameLIFFIDPayment})
	fmt.Println(b.BookedStart.Format("15:04"), b.BookedEnd.Format("15:04"))
	flexContainerStr = fmt.Sprintf(checkoutTemplate,
		c.ChatChannel.ChaImage,
		serI.Service.SerName,
		c.ChatChannel.ChaAddress,
		b.BookedDay.Format("2006-01-02"),
		fmt.Sprintf("%s", b.BookedStart.Format("15:04")),
		fmt.Sprintf("%s", b.BookedEnd.Format("15:04")),
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

func bindTransaction(c *Context, Total float64) (*model.Transaction, error) {
	var tran model.Transaction
	tran.ChatChannelID = c.ChatChannel.ID
	tran.TranTotal = Total
	tran.AccountID = c.ChatChannel.AccountID
	tran.CustomerID = c.Customer.ID
	tran.TranLineID = c.Event.Source.UserID
	tran.TranStatus = model.TranStatusPanding
	return &tran, nil
}

func bindMSPlace(c *Context, placeID uint, placAmount int) (*model.MasterPlace, error) {
	var MSPlace model.MasterPlace
	// day, err := time.Parse("2006-01-02", c.PostbackAction.DateStr)
	// if err != nil {
	// 	return nil, err
	// }
	MSPlace.PlaceID = placeID
	// MSPlace.MPlaDay = day
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

func bindBookingServiceItemAppointment(c *Context, bookingType model.BookingType, serI model.ServiceItem) (model.Booking, error) {
	var b model.Booking

	b.BookingType = bookingType
	b.ChatChannelID = c.ChatChannel.ID
	b.BooLineID = c.Event.Source.UserID
	b.CustomerID = c.Customer.ID
	b.BooStatus = model.BookingStatusPandding
	b.AccountID = c.Account.ID
	day, err := time.Parse(time.RFC3339, c.PostbackAction.DateStr+"T00:00:00+00:00")
	if err != nil {
		return b, err
	}
	fmt.Println(time.Since(day))
	// if time.Since(day) {

	// }

	start, err := time.Parse(time.RFC3339, "2012-11-01T"+c.PostbackAction.TimeStr+":41+00:00")
	if err != nil {
		return b, err
	}
	b.BookedDay = day
	b.BookedStart = start
	b.BookedEnd = start.Add(serI.SSTime)
	emID, err := strconv.ParseUint(c.PostbackAction.EmployeeID, 10, 32)
	if err != nil {
		return b, err
	}
	serIID, err := strconv.ParseUint(c.PostbackAction.ServiceItemID, 10, 32)
	if err != nil {
		return b, err
	}
	tsID, err := strconv.ParseUint(c.PostbackAction.TimeSlotID, 10, 32)
	if err != nil {
		return b, err
	}
	b.BookingServiceItem.TimeSlotID = uint(tsID)
	b.BookingServiceItem.ServiceItemID = uint(serIID)
	b.BookingServiceItem.EmployeeID = uint(emID)
	return b, nil
}

func bindBookingServiceItemNow(c *Context, bookingType model.BookingType) (model.Booking, error) {
	var b model.Booking
	b.BookingType = bookingType
	b.ChatChannelID = c.ChatChannel.ID
	b.BooLineID = c.Event.Source.UserID
	b.CustomerID = c.Customer.ID
	b.BooStatus = model.BookingStatusPandding
	b.AccountID = c.Account.ID
	u64, err := strconv.ParseUint(c.PostbackAction.ServiceItemID, 10, 32)
	if err != nil {
		fmt.Println(err)
		return b, err
	}
	b.BookingServiceItem.ServiceItemID = uint(u64)
	return b, nil
}

func bindBookingPackageNow(c *Context, bookingType model.BookingType) (model.Booking, error) {
	var b model.Booking
	b.BookingType = bookingType
	b.ChatChannelID = c.ChatChannel.ID
	b.BooLineID = c.Event.Source.UserID
	b.CustomerID = c.Customer.ID
	b.BooStatus = model.BookingStatusPandding
	b.AccountID = c.Account.ID
	fmt.Println(c.PostbackAction.PackageID)
	u64, err := strconv.ParseUint(c.PostbackAction.PackageID, 10, 32)
	if err != nil {
		fmt.Println(err)
		return b, err
	}
	var bp model.BookingPackage
	bp.PackageID = uint(u64)
	b.BookingPackage = &bp
	return b, nil
}

func bindBookingPackageAppointmant(c *Context, bookingType model.BookingType) (model.Booking, error) {
	var b model.Booking
	u64, err := strconv.ParseUint(c.PostbackAction.PackageID, 10, 32)
	if err != nil {
		fmt.Println(err)
		return b, err
	}
	b.BookingPackage.PackageID = uint(u64)
	u64, err = strconv.ParseUint(c.PostbackAction.TimeSlotID, 10, 32)
	if err != nil {
		fmt.Println(err)
		return b, err
	}
	b.BookingPackage.TimeSlotID = uint(u64)
	b.BookingType = bookingType
	b.ChatChannelID = c.ChatChannel.ID
	b.BooLineID = c.Event.Source.UserID
	b.CustomerID = c.Customer.ID
	b.BooStatus = model.BookingStatusPandding
	b.AccountID = c.Account.ID
	return b, nil
}
