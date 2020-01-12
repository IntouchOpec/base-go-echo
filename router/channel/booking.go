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
{
	"type": "carousel",
	"contents": [
	  { "type": "bubble",
		"hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s" },
		"body": { "type": "box", "layout": "vertical", "spacing": "sm",
		  "contents": [
			{ "type": "button", "style": "primary", "action": { "type": "message", "label": "ทำนัดอัตโนมัติ", "text": "service auto" } },	
			{ "type": "button", "style": "primary", "action": { "type": "message", "label": "ทำนัดเอง", "text": "service choose" } }
		  ]
		}
	  }
	]
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
	{ "type": "button", "flex": 1, "gravity": "center", "action": { "type": "uri", "label": "See more", "uri": "https://linecorp.com" } }] }}`
var thankyouTemplate string = `{
	"type": "bubble",
	"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
	"body": {
	  "type": "box",
	  "layout": "vertical",
	  "contents": [
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
		{ "type": "button", "style": "link", "height": "sm", "action": { "type": "uri", "label": "CALL", "uri": "https://linecorp.com" }
		},
		{ "type": "spacer", "size": "sm" }
	],
	"flex": 0
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
				{ "type": "text", "text": "ถัดไป", "size": "sm", "color": "#111111", "align": "end", "action": { "type": "message", "label": " ", "text": "%s"} }]
		}, %s]}}`, HeaderCalendat, actionNextMonth, weekdaysStr+`{"type": "separator"},`+calendar[:len(calendar)-1])
	return m
}

func ChooseService(c *Context) (linebot.SendingMessage, error) {
	m := fmt.Sprintf(serviceMassage, c.ChatChannel.ChaImage)
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

func CalandarHandler(c *Context) (linebot.SendingMessage, error) {
	var m string
	text := c.Massage.Text
	if len(text) > 8 {
		m = CalendarTemplate("", "", text[9:19])
	} else {
		m = CalendarTemplate("", "", "")
	}
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
	count = 0
	dateTime := c.Massage.Text[9:19]

	var employeeServices []model.EmployeeService

	day, err := time.Parse("2006-01-02 15:04", dateTime+" 15:00")
	var start time.Time
	var end time.Time

	if err != nil {
		return nil, err
	}

	if err := c.DB.Preload("Places").Where("ser_name = ? and account_id = ?", c.Massage.Text[20:], c.Account.ID).Find(&service).Error; err != nil {
		return nil, err
	}

	var placeIDs []uint

	for _, place := range service.Places {
		placeIDs = append(placeIDs, place.ID)
	}

	if err := c.DB.Preload("Employee").Preload("TimeSlots", func(db *gorm.DB) *gorm.DB {
		return db.Where("time_day = ?", day.Weekday()).Preload("Bookings", "booked_date = ?", day)
	}).Where("service_id = ? and account_id = ?",
		service.ID, c.Account.ID).Find(&employeeServices).Error; err != nil {
		return nil, err
	}
	c.DB.Where("m_pla_day = ?", day).Find(&MSPlaces)
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
				start, err = time.Parse("15:04", timeSlot.TimeStart)
				end, err = time.Parse("15:04", timeSlot.TimeEnd)
				fmt.Println(start, end)
				if len(MSPlaces) != 0 {

				} else {
					actionMessge = "timeslot" + " " + dateTime + " " + timeSlot.TimeStart + "-" + timeSlot.TimeEnd + " " + fmt.Sprint(timeSlot.ID)
					buttonTime = buttonTime + fmt.Sprintf(buttonTimePrimaryTemplate, fmt.Sprintf("%s-%s", timeSlot.TimeStart, timeSlot.TimeEnd), actionMessge)
				}
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

// ThankyouTemplate
func ThankyouTemplate(c *Context) (linebot.SendingMessage, error) {
	var timeSlot model.TimeSlot
	var tran model.Transaction
	var book model.Booking
	var service model.Service
	var MSPlace model.MasterPlace
	tx := c.DB.Begin()
	dateTime := c.Massage.Text[9:19]
	start, err := time.Parse("2006-01-02 15:04", dateTime+" 15:00")

	if err := tx.Preload("EmployeeService").Find(&timeSlot, c.Massage.Text[32:]).Error; err != nil {
		return nil, err
	}

	tx.Preload("Booking", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Bookings", "booked_date = ?", start).Preload("Place", "m_pla_status = ?", model.MPlaStatusOpen)
	}).Where("account_id = ? and time_day = ?", c.Account.ID, start.Weekday()).Find(&service, timeSlot.EmployeeService.ServiceID)
	if len(service.Places) == 0 {
		return nil, errors.New("Not found place")
	}
	book.PlaceID = service.Places[0].ID
	book.ChatChannelID = c.ChatChannel.ID
	book.CustomerID = c.Customer.ID
	book.BooLineID = c.Massage.ID
	book.TimeSlotID = timeSlot.ID
	book.EmployeeID = timeSlot.EmployeeService.EmployeeID
	layout := "2006-01-02 15:00"
	updatedAt, err := time.Parse(layout, c.Massage.Text[9:19]+" 15:00")
	if err != nil {
		return nil, err
	}
	book.BookedDate = updatedAt
	err = tx.Create(&book).Error
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
	MSPlace.PlaceID = book.PlaceID
	MSPlace.MPlaDay = start
	MSPlace.AccountID = c.Account.ID
	form, err := time.Parse("2006-01-02 15:04", "2006-01-02 "+timeSlot.TimeStart)
	to, err := time.Parse("2006-01-02 15:04", "2006-01-02 "+timeSlot.TimeEnd)
	MSPlace.MPlaFrom = form
	MSPlace.MPlaTo = to
	err = tx.FirstOrCreate(&MSPlace).Error
	MSPlace.MPlaAmount = MSPlace.MPlaAmount + 1
	if MSPlace.MPlaAmount == service.Places[0].PlacAmount {
		MSPlace.MPlaStatus = model.MPlaStatusBusy
	}
	tx.Save(&MSPlace)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	thankyou := fmt.Sprintf(thankyouTemplate, c.ChatChannel.ChaImage, c.ChatChannel.ChaAddress, timeSlot.TimeStart, timeSlot.TimeEnd)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(thankyou))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("ตาราง", flexContainer), nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
