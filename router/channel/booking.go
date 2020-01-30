package channel

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

var serviceMassage string = `
{ "type": "bubble",
"hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s" },
"body": { "type": "box", "layout": "vertical", "spacing": "sm",
  "contents": [
	{ "type": "button", "style": "primary", "action": {
		"type": "postback", "label": "จองเดี๋ยวนี้", "data": "action=booking_now" } },
	{ "type": "button", "style": "primary", "action": 
		{ "type": "datetimepicker", "label": "จองล่วงหน้า", "data": "action=booking_appointment", "mode": "datetime", "initial": "%s", "max": "%s", "min": "%s" } },	
	{ "type": "button", "style": "primary", "action": {
		"type": "postback", "label": "ทำนัดเอง", "data": "action=choive_man" } }
	]}
}`

var buttonTimeSecondaryTemplate string = `{"type": "button", "style": "secondary", "margin": "sm", "action": { "type": "message", "label": "%s", "text": "%s" }},`
var buttonTimePrimaryTemplate string = `{"type": "button","style": "primary", "margin": "sm", "action": { "type": "message", "label": "%s", "text": "%s" }},`
var buttonTimePrimaryLastTemplate string = `{"type": "button","style": "primary", "margin": "sm", "action": { "type": "message", "label": "%s", "text": "%s" }},`
var slotTimeTemplate string = `,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`
var serviceListTemplate string = `{"type": "bubble", "hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s"},
"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
	{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl" },
	{ "type": "box", "layout": "baseline", "contents": [
		{ "type": "text", "text": "฿%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
	] }
	%s]
}}`
var nextPageTemplate string = `{ "type": "bubble", "body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
	{ "type": "button", "flex": 1, "gravity": "center", "action": { "type": "uri", "label": "See more", "uri": "line://app/1615136604-5wld9ZdL" } }] }}`

var checkoutTemplate string = `{
	"type": "bubble",
	"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
	"body": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "text", "text": "จองสำเร็จ", "weight": "bold", "size": "xl" },
		{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "Place", "color": "#aaaaaa", "size": "sm", "flex": 1 },
				{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
			] },
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "Time", "color": "#aaaaaa", "size": "sm", "flex": 1 },
				{ "type": "text", "text": "%s - %s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
			] }
		  ] }
		  ]
	},
	"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
		{ "type": "button", "style": "link", "height": "sm", "action": { "type": "uri", "label": "ชำระเงิน", "uri": "line://app/%s?account_name=%s&doc_code_transaction=%s" }
		},
		{ "type": "spacer", "size": "sm" }
	],
	"flex": 0
	}
}`

var StatusOpecCardTemplate string = `
{
	"type": "bubble",
	"header": { "type": "box", "layout": "vertical", "flex": 0, "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "image", "url": "%s", "flex": 8, "align": "end", "size": "xxs", "aspectRatio": "1.91:1" },
				{ "type": "text", "text": "%s", "flex": 2, "margin": "xs", "align": "end", "gravity": "top","size": "xxs", "wrap": true }
			  ]
			}
		  ]
		},
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "image", "url": "https://developers.line.biz/assets/images/services/bot-designer-icon.png", "flex": 2 }
		  ]
		}
	  ]
	},
	"body": { "type": "box", "layout": "vertical", "spacing": "md", "action": { "type": "uri", "label": "Action", "uri": "line://app/1615136604-5wld9ZdL" },
	  "contents": [
		{ "type": "text", "contents": [], "size": "xl", "wrap": true, "text": "%s", "weight": "bold" },
		{ "type": "text", "text": "Date, %s", "color": "#000000", "size": "sm" }
	  ]
	},
	"footer": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "spacer", "size": "xxl" },
		{ "type": "button", "style": "primary", "margin": "sm", "action": { "type": "message", "label": "booking", "text": "service auto" } }
	  ]
	}
}`

func CalendarTemplate(firstKeyWordAction, lastKeyWordAction, date string) string {
	var contents string
	var calendar string
	var color string
	year, month, _ := time.Now().Date()
	t := time.Now()
	if len(date) != 0 {
		time2, _ := time.Parse("2006-01-02", date)
		year, month, _ = time2.Date()
	}

	color = "#000000"

	endOfMonth := time.Date(year, month+1, 1, 0, 0, 0, -1, time.UTC)

	Weekday := int(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday())

	for day := 0; day < Weekday; day++ {
		contents = contents + fmt.Sprintf(`{
			"type":    "text",
			"text":    " ",
			"size":    "sm",
			"color":   "#000000",
			"align":   "center",
			"gravity": "center"},`)
	}
	Weekday = int(t.Weekday())

	for day := 1; day <= endOfMonth.Day(); day++ {
		if day == t.Day() && month == t.Month() {
			color = "#1db446"
		} else {
			color = "#000000"
		}
		dayStr := strconv.FormatInt(int64(day), 10)
		monthStr := strconv.FormatInt(int64(month), 10)
		if len(dayStr) == 1 {
			dayStr = fmt.Sprintf("0%s", dayStr)
		}
		if len(monthStr) == 1 {
			monthStr = fmt.Sprintf("0%s", monthStr)
		}
		contents = contents + fmt.Sprintf(`{ "type": "text", "text": "%s", "size": "sm", "color": "%s", "align": "center", "gravity": "center",
					"action": { "type": "message", "label": "%s", "text": "%s-%s-%s"}},`, dayStr, color, day, fmt.Sprintf("%s %d", firstKeyWordAction, year), monthStr, dayStr+" "+lastKeyWordAction)
		contents = contents + `{"type": "separator"},`
		Weekday = int(time.Date(year, month, day, 0, 0, 0, -1, time.UTC).Weekday())
		if endOfMonth.Day() == day {
			for dw := int(endOfMonth.Weekday()); dw < 6; dw++ {
				contents = contents + fmt.Sprintf(`{ "type": "text", "text": " ", "size": "sm", "color": "#000000", "align": "center", "gravity": "center"},`)
			}
		}

		// 6 == saturday
		if (int(Weekday) == 5) || endOfMonth.Day() == day {

			calendar = calendar + fmt.Sprintf(`{
				"type":     "box",
				"layout":   "horizontal",
				"margin":   "md",
				"contents": [%s]
			},`, contents[:len(contents)-1])
			contents = ""
		}
	}
	weekdays := []string{"อา", "จ", "อ", "พ", "พฤ", "ศ", "ส"}
	var weekdaysStr string
	for weekday := 0; weekday < len(weekdays); weekday++ {

		weekdaysStr = weekdaysStr + fmt.Sprintf(`{ "type": "text", "text": "%s", "size": "sm", "color": "#000000", "align": "center" },`, weekdays[weekday])
	}
	weekdaysStr = fmt.Sprintf(`{"type": "box","layout": "horizontal","margin": "md","contents": [%s]},`, weekdaysStr[:len(weekdaysStr)-1])
	HeaderCalendat := fmt.Sprintf("%s %s", month, strconv.FormatInt(int64(year), 10))
	var nextMonth string = strconv.FormatInt(int64(month+1), 10)
	var nextYear int = year
	if nextMonth == "13" {
		nextMonth = "01"
		nextYear = year + 1
	}

	if len(nextMonth) == 1 {
		nextMonth = "0" + nextMonth
	}
	actionNextMonth := fmt.Sprintf("%d-%s-01", nextYear, nextMonth)
	m := fmt.Sprintf(`{"type": "bubble","styles": {"footer": {"separator": true}},
	"body": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "%s", "size": "sm", "weight": "bold", "color": "#1db446", "flex": 0 },
				{ "type": "text", "text": "ถัดไป", "size": "sm", "color": "#111111", "align": "end", "action": { "type": "postback", "label": " ", "data": "action=choive_man&date=%s"} }]
		}, %s]}}`, HeaderCalendat, actionNextMonth, weekdaysStr+`{"type": "separator"},`+calendar[:len(calendar)-1])
	return m
}

func ChooseService(c *Context) (linebot.SendingMessage, error) {
	now := time.Now()
	format := "2006-01-02T15:04"
	m := fmt.Sprintf(serviceMassage, c.ChatChannel.ChaImage, now.Format(format), now.AddDate(0, 3, 0).Format(format), now.Format(format))
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), nil
}

func ServiceList(c *Context) (linebot.SendingMessage, error) {
	var services []model.Service
	var template string
	var image string
	var button string
	if err := c.DB.Where("account_id = ?", c.ChatChannel.AccountID).Limit(10).Find(&services).Error; err != nil {
		return nil, err
	}
	for _, service := range services {
		button = fmt.Sprintf(buttonTimePrimaryTemplate, service.SerName, "Service "+service.SerName)
		image = "https://" + Conf.Server.Domain + service.SerImage
		template += fmt.Sprintf(serviceListTemplate, image, service.SerName, strconv.FormatInt(int64(service.SerPrice), 10), ","+button[:len(button)-1]) + ","
	}
	template = fmt.Sprintf(`{ "type": "carousel", "contents": [%s]}`, template[:len(template)-1])
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), err
}

func CalandarHandler(c *Context, date string) (linebot.SendingMessage, error) {
	m := CalendarTemplate("", "", date)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), nil
}

func SaveServiceHandler(c *Context) (linebot.SendingMessage, error) {
	var service model.Service
	var m string
	if err := c.DB.Where(&model.Service{SerName: c.Massage.Text[8:]}).Find(&service).Error; err != nil {
		return nil, err
	}
	m = CalendarTemplate("booking ", c.Massage.Text[8:], "")
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(m))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), nil
}

// linebot.SendingMessage
func ServiceListLineHandler(c *Context) (linebot.SendingMessage, error) {
	var slotTime string
	var buttonTime string
	var serviceList string
	var actionMessge string
	var count int
	var service model.Service
	var MSPlaces []*model.MasterPlace
	var placeIDs []uint
	count = 0
	dateTime := c.Massage.Text[9:19]

	var employeeServices []model.EmployeeService

	day, err := time.Parse("2006-01-02 15:04", dateTime+" 15:00")

	if err != nil {
		return nil, err
	}

	if err := c.DB.Where("ser_name = ? and account_id = ?", c.Massage.Text[20:], c.Account.ID).Find(&service).Error; err != nil {
		return nil, err
	}
	if err := c.DB.Preload("Employee").Preload("TimeSlots", func(db *gorm.DB) *gorm.DB {
		return db.Where("time_day = ?", day.Weekday()).Preload("Bookings", "booked_date = ?", day)
	}).Where("service_id = ? and account_id = ?",
		service.ID, c.Account.ID).Find(&employeeServices).Error; err != nil {
		return nil, err
	}

	c.DB.Order("m_pla_status").Where("m_pla_day = ? and place_id in (?)", day, placeIDs).Find(&MSPlaces)

	for index, employeeService := range employeeServices {
		if employeeService.Employee.ChatChannelID != c.ChatChannel.ID {
			continue
		}
		for _, timeSlot := range employeeService.TimeSlots {
			if count == 2 {
				slotTime = slotTime + fmt.Sprintf(slotTimeTemplate, buttonTime[:len(buttonTime)-1])
				buttonTime = ""
				count = 0
			}
			if len(timeSlot.Bookings) > 0 {
				buttonTime = buttonTime + fmt.Sprintf(buttonTimeSecondaryTemplate, fmt.Sprintf("%s-%s", timeSlot.TimeStart, timeSlot.TimeEnd), "เต็มแล้ว")
			} else {
				actionMessge = "timeslot" + " " + dateTime + " " + timeSlot.TimeStart + "-" + timeSlot.TimeEnd + " " + fmt.Sprint(timeSlot.ID)
				buttonTime = buttonTime + fmt.Sprintf(buttonTimePrimaryTemplate, fmt.Sprintf("%s-%s", timeSlot.TimeStart, timeSlot.TimeEnd), actionMessge)
			}
			count = count + 1
		}
		if index == len(employeeServices)-1 {
			slotTime = slotTime + fmt.Sprintf(slotTimeTemplate, buttonTime[:len(buttonTime)-1])
			serviceList += fmt.Sprintf(serviceListTemplate, fmt.Sprintf("https://web.%s/file?path=%s", Conf.Server.Domain, employeeService.Employee.ProvImage), employeeService.Employee.ProvName, strconv.FormatInt(int64(employeeService.PSPrice), 10), slotTime)
			break
		}
		serviceList = serviceList + fmt.Sprintf(serviceListTemplate+",",
			fmt.Sprintf("https://web.%s/file?path=%s", Conf.Server.Domain, employeeService.Employee.ProvImage),
			employeeService.Employee.ProvName, strconv.FormatInt(int64(employeeService.PSPrice), 10), slotTime)
		slotTime = ""
		count = 0
		buttonTime = ""

	}
	nextPage := nextPageTemplate

	serviceTamplate := fmt.Sprintf(`{ "type": "carousel", "contents": [%s, %s]}`, serviceList, nextPage)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(serviceTamplate))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), err
}

func BookingServiceItemHandler(c *Context) (linebot.SendingMessage, error) {
	serviceTamplate := fmt.Sprintf(`{ "type": "carousel", "contents": [%s, %s]}`, "serviceList", "nextPage")
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(serviceTamplate))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), err
}

type placeSum struct {
	Amount  int
	PlaceID uint
}

// checkoutTemplate
func BookingTimeSlotHandler(c *Context) (linebot.SendingMessage, error) {
	var timeSlot model.TimeSlot
	var tran model.Transaction
	var book model.Booking
	var service model.Service
	var MSPlace model.MasterPlace
	var MSPlaces []*model.MasterPlace
	var bookingTimeSlot model.BookingTimeSlot

	dateTime := c.Massage.Text[9:19]
	start, err := time.Parse("2006-01-02 15:04", dateTime+" 15:00")

	if err := c.DB.Preload("EmployeeService").Find(&timeSlot, c.Massage.Text[32:]).Error; err != nil {
		return nil, err
	}
	c.DB.Preload("Places").Where("account_id = ?", c.Account.ID).Find(&service, timeSlot.EmployeeService.ServiceID)
	if len(service.Places) == 0 {
		return nil, errors.New("Not found place")
	}
	var placeIDs []uint

	for _, place := range service.Places {
		placeIDs = append(placeIDs, place.ID)
	}

	c.DB.Order("m_pla_status desc, place_id").Where("account_id =? and m_pla_day = ? and m_pla_to BETWEEN ? and ? or m_pla_from BETWEEN ? and ? and place_id in (?) ",
		c.Account.ID, timeSlot.TimeDay, timeSlot.TimeStart, timeSlot.TimeEnd, timeSlot.TimeStart, timeSlot.TimeEnd, placeIDs).Find(&MSPlaces)
	if len(MSPlaces) > 0 {
		return nil, errors.New("")
	}
	var placeSums []placeSum
	for index, MSPlace := range MSPlaces {
		if MSPlace.MPlaStatus == model.MPlaStatusBusy {
			return nil, errors.New("")
		}
		for i, placeSum := range placeSums {
			if placeSum.PlaceID != MSPlace.PlaceID {
				continue
			} else {
				placeSums[i].Amount += MSPlace.MPlaAmount
				break
			}
		}
		if index == 0 {
			placeSums = append(placeSums, placeSum{Amount: MSPlace.MPlaAmount, PlaceID: MSPlace.PlaceID})
		}
	}
	tx := c.DB.Begin()

	book.PlaceID = service.Places[0].ID
	book.ChatChannelID = c.ChatChannel.ID
	book.CustomerID = c.Customer.ID
	book.BooLineID = c.Massage.ID

	layout := "2006-01-02 15:00"
	updatedAt, err := time.Parse(layout, c.Massage.Text[9:19]+" 15:00")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	book.BookedDate = updatedAt
	err = tx.Create(&book).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	bookingTimeSlot.BookingID = book.ID
	bookingTimeSlot.TimeSlotID = timeSlot.ID
	bookingTimeSlot.EmployeeID = timeSlot.EmployeeService.EmployeeID
	err = tx.Create(&bookingTimeSlot).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tran.ChatChannelID = c.ChatChannel.ID
	tran.TranTotal = timeSlot.EmployeeService.PSPrice
	tran.AccountID = c.ChatChannel.AccountID
	tran.CustomerID = c.Customer.ID
	tran.TranLineID = c.Event.UserID
	tran.TranStatus = model.TranStatusPanding
	err = tx.Create(&tran).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Model(&tran).Association("Bookings").Append(&book).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(MSPlaces) == 0 {
		MSPlace.PlaceID = book.PlaceID
		MSPlace.MPlaDay = start
		MSPlace.AccountID = c.Account.ID
		form, err := time.Parse("15:04", timeSlot.TimeStart)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		to, err := time.Parse("15:04", timeSlot.TimeEnd)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		MSPlace.MPlaFrom = form
		MSPlace.MPlaTo = to
		MSPlace.MPlaAmount = MSPlace.MPlaAmount + 1
		if MSPlace.MPlaAmount == service.Places[0].PlacAmount {
			MSPlace.MPlaStatus = model.MPlaStatusBusy
		}
		if err := tx.Create(&MSPlace).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	setting := model.Setting{Value: "LIFFIDPayment"}
	c.DB.Model(&c.ChatChannel).Association("Settings").Find(&setting)
	thankyou := fmt.Sprintf(checkoutTemplate, c.ChatChannel.ChaImage, c.ChatChannel.ChaAddress, timeSlot.TimeStart, timeSlot.TimeEnd, setting.Value, c.Account.AccName, tran.TranDocumentCode)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(thankyou))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("จองสำเร็จ", flexContainer), nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func CheckStatusOpen(c *Context) (linebot.SendingMessage, error) {
	var MSPlaces []*model.MasterPlace
	var flexMessage string
	var status string = "Ready"
	var text string = "ว่าง"
	var icon string = "https://png.pngitem.com/pimgs/s/234-2341004_glossy-red-icon-button-clip-art-at-clker.png"
	now := time.Now()
	c.DB.Order("m_pla_status desc").Where("account_id =? and m_pla_day = ? and m_pla_to < ? and m_pla_from > ? or m_pla_status = ?",
		c.Account.ID, now, now, now, model.MPlaStatusBusy).Find(&MSPlaces)
	var placeSums []placeSum
	for index, MSPlace := range MSPlaces {
		if MSPlace.MPlaStatus == model.MPlaStatusBusy {
			status = "Busy"
			text = "ขณะนี้ไม่ว่าง"
			icon = "https://developers.line.biz/assets/images/services/bot-designer-icon.png"
			break
		}

		for i, placeSum := range placeSums {
			if placeSum.PlaceID != MSPlace.PlaceID {
				continue
			} else {
				placeSums[i].Amount += MSPlace.MPlaAmount
				break
			}
		}
		if index == 0 {
			placeSums = append(placeSums, placeSum{Amount: MSPlace.MPlaAmount, PlaceID: MSPlace.PlaceID})
		}
	}
	flexMessage = fmt.Sprintf(StatusOpecCardTemplate, icon, status, text, now.Format("Mon Jan _2 15:04:05"))
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexMessage))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage(status, flexContainer), nil
}

func BookingServiceHandler(c *Context) (linebot.SendingMessage, error) {
	var tran model.Transaction
	var book model.Booking
	var serviceItem model.ServiceItem
	var MSPlace model.MasterPlace
	var MSPlaces []*model.MasterPlace

	if c.PostbackAction.PackageID != 0 {
		var packageModel model.Package
		err := c.DB.Preload("ServiceItems", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Service", func(*gorm.DB) *gorm.DB {
				return db.Preload("places")
			})
		}).Where("account_id = ?", c.Account.ID).Find(&packageModel, c.PostbackAction.PackageID).Error

		// for _, ServiceItem := range packageModel.ServiceItems {

		// }
		tx := c.DB.Begin()
		var bookingPackage model.BookingPackage
		book.PlaceID = serviceItem.Service.Places[0].ID
		book.ChatChannelID = c.ChatChannel.ID
		book.CustomerID = c.Customer.ID
		book.BooLineID = c.Massage.ID
		layout := "2006-01-02 15:00"
		book.BookingType = model.BookingTypePackage
		updatedAt, err := time.Parse(layout, c.PostbackAction.Day+" 15:00")
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		book.BookedDate = updatedAt
		err = tx.Create(&book).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		bookingPackage.BookingID = book.ID
		bookingPackage.PackageID = packageModel.ID
		err = tx.Create(&bookingPackage).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tran.ChatChannelID = c.ChatChannel.ID
		tran.TranTotal = serviceItem.SSPrice
		tran.AccountID = c.ChatChannel.AccountID
		tran.CustomerID = c.Customer.ID
		tran.TranLineID = c.Event.UserID
		tran.TranStatus = model.TranStatusPanding
		err = tx.Create(&tran).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = tx.Model(&tran).Association("Bookings").Append(&book).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(MSPlaces) == 0 {
			day, err := time.Parse("2006-01-02", c.PostbackAction.Day)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			MSPlace.PlaceID = book.PlaceID
			MSPlace.MPlaDay = day
			MSPlace.AccountID = c.Account.ID
			form, err := time.Parse("15:04", c.PostbackAction.Start)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			to, err := time.Parse("15:04", c.PostbackAction.End)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			MSPlace.MPlaFrom = form
			MSPlace.MPlaTo = to
			MSPlace.MPlaAmount = MSPlace.MPlaAmount + 1
			if MSPlace.MPlaAmount == serviceItem.Service.Places[0].PlacAmount {
				MSPlace.MPlaStatus = model.MPlaStatusBusy
			}
			if err := tx.Create(&MSPlace).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
	} else {
		err := c.DB.Preload("Service", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Places")
		}).Where("account_id = ?", c.Account.ID).Find(&serviceItem, c.PostbackAction.ServiceItemID).Error

		if len(serviceItem.Service.Places) == 0 {
			return nil, errors.New("Not found place")
		}
		var placeIDs []uint
		for _, place := range serviceItem.Service.Places {
			placeIDs = append(placeIDs, place.ID)
		}
		c.DB.Order("m_pla_status desc, place_id").Where("account_id =? and m_pla_day = ? and m_pla_to BETWEEN ? and ? or m_pla_from BETWEEN ? and ? and place_id in (?) ",
			c.Account.ID, c.PostbackAction.Day, c.PostbackAction.Start, c.PostbackAction.End, c.PostbackAction.Start, c.PostbackAction.End, placeIDs).Find(&MSPlaces)
		if len(MSPlaces) > 0 {
			return nil, errors.New("")
		}
		var placeSums []placeSum
		for index, MSPlace := range MSPlaces {
			if MSPlace.MPlaStatus == model.MPlaStatusBusy {
				return nil, errors.New("")
			}
			for i, placeSum := range placeSums {
				if placeSum.PlaceID != MSPlace.PlaceID {
					continue
				} else {
					placeSums[i].Amount += MSPlace.MPlaAmount
					break
				}
			}
			if index == 0 {
				placeSums = append(placeSums, placeSum{Amount: MSPlace.MPlaAmount, PlaceID: MSPlace.PlaceID})
			}
		}
		var bookingServiceItem model.BookingServiceItem
		tx := c.DB.Begin()
		book.BookingType = model.BookingTypeServiceItem
		book.PlaceID = serviceItem.Service.Places[0].ID
		book.ChatChannelID = c.ChatChannel.ID
		book.CustomerID = c.Customer.ID
		book.BooLineID = c.Event.UserID
		// layout := "2006-01-02 15:00"
		// updatedAt, err := time.Parse(layout, c.PostbackAction.Day+" 15:00")
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// book.BookedDate = updatedAt
		err = tx.Create(&book).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		bookingServiceItem.BookingID = book.ID
		bookingServiceItem.ServiceItemID = serviceItem.ID
		err = tx.Create(&bookingServiceItem).Error
		if err != nil {
			tx.Rollback()

			return nil, err
		}
		tran.ChatChannelID = c.ChatChannel.ID
		tran.TranTotal = serviceItem.SSPrice
		tran.AccountID = c.ChatChannel.AccountID
		tran.CustomerID = c.Customer.ID
		tran.TranLineID = c.Event.UserID
		tran.TranStatus = model.TranStatusPanding
		if err := tx.Create(&tran).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		err = tx.Model(&tran).Association("Bookings").Append(&book).Error
		if err != nil {

			tx.Rollback()
			return nil, err
		}

		if len(MSPlaces) == 0 {
			day, err := time.Parse("2006-01-02", c.PostbackAction.Day)
			if err != nil {
				tx.Rollback()

				return nil, err
			}
			MSPlace.PlaceID = book.PlaceID
			MSPlace.MPlaDay = day
			MSPlace.AccountID = c.Account.ID
			form, err := time.Parse("15:04", c.PostbackAction.Start)
			if err != nil {
				tx.Rollback()

				return nil, err
			}
			to, err := time.Parse("15:04", c.PostbackAction.End)
			if err != nil {

				tx.Rollback()
				return nil, err
			}
			MSPlace.MPlaFrom = form
			MSPlace.MPlaTo = to
			MSPlace.MPlaAmount = MSPlace.MPlaAmount + 1
			if MSPlace.MPlaAmount == serviceItem.Service.Places[0].PlacAmount {
				MSPlace.MPlaStatus = model.MPlaStatusBusy
			}
			if err := tx.Create(&MSPlace).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err != nil {

			tx.Rollback()
			return nil, err
		}
		tx.Commit()
	}
	thankyou := fmt.Sprintf(checkoutTemplate, c.ChatChannel.ChaImage, c.ChatChannel.ChaAddress, c.PostbackAction.Start, c.PostbackAction.End, c.Account.AccName, tran.TranDocumentCode)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(thankyou))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("จองสำเร็จ", flexContainer), nil
}

func bindBooking(c *Context, palceID uint) (*model.Booking, error) {
	var book model.Booking
	book.PlaceID = palceID
	book.ChatChannelID = c.ChatChannel.ID
	book.CustomerID = c.Customer.ID
	book.BooLineID = c.Massage.ID
	layout := "2006-01-02 15:00"
	book.BookingType = model.BookingTypePackage
	updatedAt, err := time.Parse(layout, c.PostbackAction.Day+" 15:00")
	if err != nil {
		return nil, err
	}
	book.BookedDate = updatedAt
	return &book, nil
}

func bindTransaction(c *Context, Total float64) (*model.Transaction, error) {
	var tran model.Transaction
	tran.ChatChannelID = c.ChatChannel.ID
	tran.TranTotal = Total
	tran.AccountID = c.ChatChannel.AccountID
	tran.CustomerID = c.Customer.ID
	tran.TranLineID = c.Event.UserID
	tran.TranStatus = model.TranStatusPanding
	return &tran, nil
}

func bindMSPlace(c *Context, placeID uint, placAmount int) (*model.MasterPlace, error) {
	var MSPlace model.MasterPlace
	day, err := time.Parse("2006-01-02", c.PostbackAction.Day)
	if err != nil {
		return nil, err
	}
	MSPlace.PlaceID = placeID
	MSPlace.MPlaDay = day
	MSPlace.AccountID = c.Account.ID
	form, err := time.Parse("15:04", c.PostbackAction.Start)
	if err != nil {
		return nil, err
	}
	to, err := time.Parse("15:04", c.PostbackAction.End)
	if err != nil {
		return nil, err
	}
	MSPlace.MPlaFrom = form
	MSPlace.MPlaTo = to
	MSPlace.MPlaAmount = MSPlace.MPlaAmount + 1
	if MSPlace.MPlaAmount == placAmount {
		MSPlace.MPlaStatus = model.MPlaStatusBusy
	}
	return &MSPlace, nil
}

func bindBookingServiceItem(c *Context, bookID, serviceItemID uint) (*model.BookingServiceItem, error) {
	var bookingServiceItem model.BookingServiceItem
	bookingServiceItem.BookingID = bookID
	bookingServiceItem.ServiceItemID = serviceItemID
	return &bookingServiceItem, nil
}

func bindBookingTimeSlot(c *Context, bookID, timeSlotID, employeeID uint) (*model.BookingTimeSlot, error) {
	var bookingTimeSlot model.BookingTimeSlot
	bookingTimeSlot.BookingID = bookID
	bookingTimeSlot.TimeSlotID = timeSlotID
	bookingTimeSlot.EmployeeID = employeeID
	return &bookingTimeSlot, nil
}

func bindBookingPackage(c *Context, bookID, packageID uint) (*model.BookingPackage, error) {
	var bookingPackage model.BookingPackage
	bookingPackage.BookingID = bookID
	bookingPackage.PackageID = packageID
	return &bookingPackage, nil
}
