package channel

import (
	"database/sql"
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
			pa.id, pac_name, si.service_id , pa.pac_price, pa.pac_time, pa.pac_image
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
			rows.Scan(&p.ID, &p.Name, &serID, &p.Price, &p.TimeUse, &p.Image)
			serIDs = append(serIDs, serID)
		}
		switch c.PostbackAction.Type {
		case "now":
			template, err = PackNow(c, p, serIDs, setting)
			if err != nil {
				return err
			}
		case "appointment":
			template, err = PackAppointment(c, p, serIDs)
			if err != nil {
				return err
			}
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
			return err
		}
		serI.Service = &ser

		switch c.PostbackAction.Type {
		case "now":
			template, err = ServiceItemNow(c, serI, ser.SerName, setting)
			if err != nil {
				fmt.Println(err, "231312")
				return err
			}
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

func ServiceItemNow(c *Context, serI model.ServiceItem, serName string, setting map[string]string) (string, error) {
	b, err := bindBookingServiceItemNow(c, model.BookingTypeServiceItem)
	if err != nil {
		return "", err
	}
	plaMDs, err := b.ServiceItemNow(c.sqlDb, serI)
	if err != nil {
		return "", err
	}
	tran, err := bindTransaction(c, serI.SSPrice)
	if err != nil {
		return "", err
	}

	tx, err := c.sqlDb.Begin()
	if err != nil {
		return "", err
	}
	if err := tran.CreateSql(tx); err != nil {
		tx.Rollback()
		return "", err
	}
	b.TransactionID = tran.ID
	if err := b.CreateSql(tx); err != nil {
		tx.Rollback()
		return "", err
	}
	ms, vStr, err := b.MasterBookingSer(plaMDs)
	if err != nil {
		return "", err
	}
	if err := model.CreateMasterBooking(vStr, tx, ms); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return "", err
	}

	return fmt.Sprintf(checkoutTemplate,
		c.ChatChannel.ChaImage,
		serName,
		c.ChatChannel.ChaAddress,
		b.BookedDay.Format("2006-01-02"),
		b.BookedStart.Add(7*time.Hour).Format("15:04"),
		b.BookedEnd.Add(7*time.Hour).Format("15:04"),
		setting[model.NameLIFFIDPayment],
		c.Account.AccName,
		tran.TranDocumentCode,
		setting[model.NameLIFFIDPayment]), nil
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
		return "", err
	}
	if err := b.PackNow(c.sqlDb, p, serIDs); err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}
	t, err := bindTransaction(c, p.Price)
	if err != nil {
		return "", err
	}
	tx, err := c.sqlDb.Begin()
	if err := t.CreateSql(tx); err != nil {
		tx.Rollback()
		return "", err
	}
	b.TransactionID = t.ID
	if err := b.CreateSql(tx); err != nil {
		tx.Rollback()
		return "", err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return "", err
	}
	return fmt.Sprintf(checkoutTemplate,
		c.ChatChannel.ChaImage,
		p.Name,
		c.ChatChannel.ChaAddress,
		b.BookedDay.Format("2006-01-02"),
		b.BookedStart.Format("15:04"),
		b.BookedEnd.Format("15:04"),
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
func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

type emploTime struct {
	TsS     time.Time
	TsE     time.Time
	Amount  int
	TsID    uint
	SerID   uint
	EmploID uint
}

type masterPla struct {
	Amount  uint
	Date    time.Time
	PlaceID uint
}

type masterEmplo struct {
	EmployeeID    uint
	Date          time.Time
	Amount        int
	SupportAmount int
}

func PackAppointment(c *Context, p model.Pack, serIDs []string) (string, error) {
	var rows *sql.Rows
	var start time.Time
	var end time.Time
	d, err := time.Parse("2006-01-02", c.Event.Postback.Params.Date)
	if err != nil {
		fmt.Println("date err")
		return "", err
	}
	if err := chackDateBooking(d); err != nil {
		return "", err
	}
	now := time.Now()
	isToDay := DateEqual(now, d)
	if isToDay {
		start, err = model.MakeHour(now)
		if err != nil {
			return "", nil
		}
	}
	var serIDsStr string
	for i := 0; i < len(serIDs); i++ {
		serIDsStr += serIDs[i] + ","
	}
	serIDsStr = serIDsStr[:len(serIDsStr)-1]
	qe := `
		SELECT 
			es.employee_id ,
			ts.time_start, ts.time_end, ts.id,
			es.service_id
		FROM employee_service AS es
		INNER JOIN employees AS e ON e.id = es.employee_id AND e.deleted_at IS NULL
		INNER JOIN time_slots AS ts ON ts.employee_id = es.employee_id AND ts.deleted_at IS NULL
			AND ts.deleted_at IS NULL 
			AND ts.time_active = true
			%s
		WHERE es.service_id %s
		ORDER BY es.employee_id`

	var wherIn string
	if isToDay {
		wherIn = fmt.Sprintf(`IN (%s)`, serIDsStr)
		qe = fmt.Sprintf(qe, `AND time_day = $1 AND ts.time_start < $2`, wherIn)
		rows, err = c.sqlDb.Query(qe, int(d.Weekday()), start)
	} else {
		wherIn = fmt.Sprintf(`IN (%s)`, serIDsStr)
		qe = fmt.Sprintf(qe, `AND time_day = $1`, wherIn)
		rows, err = c.sqlDb.Query(qe, int(d.Weekday()))
	}
	if err != nil {
		return "", err
	}

	var epts []emploTime
	var emploIDs string
	for rows.Next() {
		var ept emploTime
		rows.Scan(&ept.EmploID, &ept.TsS, &ept.TsE, &ept.TsID, &ept.SerID)
		if end.Before(ept.TsE) {
			end = ept.TsE
		}
		epts = append(epts, ept)
		emploIDs += fmt.Sprintf("%d,", ept.EmploID)
	}
	if len(epts) == 0 {
		return "", errors.New("not found employee")
	}
	qe = `
		SELECT 
			COUNT(*), 
			date_trunc('hour', mb_from) +
			(((date_part('minute', mb_from)::integer / 60::integer) * 60::integer)
			|| ' minutes')::interval AS hour_time, 
			employee_id
		FROM master_bookings
		WHERE mb_day = $1 AND employee_id IN (%s) %s
		GROUP BY hour_time, employee_id
		ORDER BY employee_id, hour_time
	`
	qe = fmt.Sprintf(qe, emploIDs[:len(emploIDs)-1], "")
	rows, err = c.sqlDb.Query(qe, d)

	if err != nil {
		return "", err
	}

	var mes []masterEmplo
	for rows.Next() {
		var me masterEmplo
		rows.Scan(&me.Amount, &me.Date, &me.EmployeeID)
		mes = append(mes, me)
	}
	qe = `
		SELECT 
			pl.id, pl.plac_amount
		FROM places AS pl
		INNER JOIN place_service AS ps ON ps.place_id = pl.id AND service_id %s
		WHERE pl.deleted_at IS NULL 
		GROUP BY pl.id, pl.plac_amount
		ORDER BY pl.id
	`
	wherIn = fmt.Sprintf(`IN (%s)`, serIDsStr)
	qe = fmt.Sprintf(qe, wherIn)

	rows, err = c.sqlDb.Query(qe)
	if err != nil {
		return "", err
	}
	var ps []model.Pla
	var placIDs string
	for rows.Next() {
		var p model.Pla
		rows.Scan(&p.ID)
		placIDs += fmt.Sprintf("%d,", p.ID)
		ps = append(ps, p)
	}

	qe = `
		SELECT 
			COUNT(*), 
			date_trunc('hour', mb_from) +
			(((date_part('minute', mb_from)::integer / 60::integer) * 60::integer)
			|| ' minutes')::interval AS hour_time, 
			place_id
		FROM master_bookings
		WHERE mb_day = $1 AND place_id IN ($2) %s
		GROUP BY hour_time, place_id
		ORDER BY place_id, hour_time
	`
	var mps []masterPla

	for rows.Next() {
		var mp masterPla
		rows.Scan(&mp.Amount, &mp.Date, &mp.PlaceID)
		mps = append(mps, mp)
	}

	var button string
	var cont string
	end = end.Add(7 * time.Hour)
	slot := end.Sub(start) / (60 * time.Minute)

	for i := 0; i < int(slot); i++ {
		tim := start.Add((time.Duration(60*i) * time.Minute))
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
		action := fmt.Sprintf("action=%s&package_id=%d&date=%s&time=%s",
			"checkout", p.ID, c.Event.Postback.Params.Date, bookingTime)
		button += fmt.Sprintf(buttonTimeSlotTemplate,
			bookingTime, action) + ","
		if i%2 != 0 {
			cont += fmt.Sprintf(layoutTimeSlotTemplate, button[:len(button)-1]) + ","
			button = ""
		} else if int(slot)-1 == i {
			cont += fmt.Sprintf(layoutTimeSlotTemplate, button[:len(button)-1]) + ","
		}
	}

	return fmt.Sprintf(timeSlotTemplate, fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, p.Image), cont[:len(cont)-1]), nil
}

func ServiceItemAppointment(c *Context, serI model.ServiceItem) string {
	d, err := time.Parse("2006-01-02", c.Event.Postback.Params.Date)
	var start time.Time

	if err != nil {
		fmt.Println("date err")
		return ""
	}
	var emplos []emplo
	var template string
	if err := chackDateBooking(d); err != nil {
		return ""
	}
	now := time.Now()
	isToDay := DateEqual(now, d)
	if isToDay {
		start, err = model.MakeHour(now)
		start = start.Add(-7 * time.Hour)
		if err != nil {
			return ""
		}
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
		slot := emp.timeEnd.Sub(start) / (60 * time.Minute)
		var button string
		var cont string
		fmt.Println(slot, start, emp.timeEnd, emp.timeStart)
		for i := 0; i < int(slot); i++ {
			tim := start.Add((time.Duration(60*i) * time.Minute) + (7 * time.Hour))
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
		fmt.Println(template, cont)
		if cont != "" {
			template += fmt.Sprintf(timeSlotTemplate, fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, emp.image), cont[:len(cont)-1]) + ","
		}
	}
	return fmt.Sprintf(carouselTemplate, template[:len(template)-1])
}

func ChackOutHandler(c *Context) (linebot.SendingMessage, error) {
	var flexContainerStr string
	// var slot model.TimeSlot
	setting := c.ChatChannel.GetSetting([]string{model.NameLIFFIDPayment})
	if c.PostbackAction.ServiceItemID != "" {
		var serI model.ServiceItem
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
		if err != nil {

		}
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
	} else {
		p, serIDs, err := model.GetPack(c.sqlDb, c.PostbackAction.PackageID)
		if err != nil {
			return nil, err
		}
		b, err := bindBookingPackAppoint(c, model.BookingTypePackage, p)
		if err != nil {

		}
		if err := b.MakeMBPacks(c.sqlDb, p, serIDs); err != nil {
			return nil, err
		}

		tr, err := bindTransaction(c, p.Price)
		if err != nil {
			return nil, err
		}

		tx, err := c.sqlDb.Begin()
		if err := tr.CreateSql(tx); err != nil {
			return nil, err
		}
		b.TransactionID = tr.ID
		if err := b.CreateSql(tx); err != nil {
			fmt.Println(err, ":err6")
			tx.Rollback()
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
		}
		flexContainerStr = fmt.Sprintf(checkoutTemplate,
			c.ChatChannel.ChaImage,
			p.Name,
			c.ChatChannel.ChaAddress,
			b.BookedDay.Format("2006-01-02"),
			fmt.Sprintf("%s", b.BookedStart.Format("15:04")),
			fmt.Sprintf("%s", b.BookedEnd.Format("15:04")),
			setting[model.NameLIFFIDPayment],
			c.Account.AccName,
			tr.TranDocumentCode,
			setting[model.NameLIFFIDPayment])
	}

	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexContainerStr))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("buy page", flexContainer), nil
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

func bindBookingPackAppoint(c *Context, bookingType model.BookingType, p *model.Pack) (model.Booking, error) {
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
	start, err := time.Parse(time.RFC3339, "2012-11-01T"+c.PostbackAction.TimeStr+":41+00:00")
	if err != nil {
		return b, err
	}
	b.BookedDay = day
	b.BookedStart = start
	b.BookedEnd = start.Add(p.TimeUse)
	var bp model.BookingPackage
	b.BookingPackage.PackageID = p.ID
	b.BookingPackage = &bp
	return b, nil
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
	u64, err := strconv.ParseUint(c.PostbackAction.PackageID, 10, 32)
	if err != nil {
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
