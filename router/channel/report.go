package channel

import (
	"fmt"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
)

func ReportListHandler(c *Context) (linebot.SendingMessage, error) {
	var template string
	var transactions []*model.Transaction
	c.DB.Limit(10).Preload("ChatChannel").Preload("Bookings").Where("tran_line_id = ? and chat_channel_id = ?", c.Event.Source.UserID, c.ChatChannel.ID).Find(&transactions)
	var bookingTimeSlot model.BookingTimeSlot
	var bookingServiceItem model.BookingServiceItem
	var bookingPackage model.BookingPackage
	var timeStart string
	var timeEnd string
	var name string

	setting := c.ChatChannel.GetSetting([]string{model.NameLIFFIDReport})

	for _, transaction := range transactions {
		for _, booking := range transaction.Bookings {
			switch booking.BookingType {
			case model.BookingTypeSlotTime:
				c.DB.Preload("TimeSlot", func(db *gorm.DB) *gorm.DB {
					return db.Preload("EmployeeService", func(db *gorm.DB) *gorm.DB {
						return db.Preload("Employee").Preload("Service")
					})
				}).Find(&bookingTimeSlot)
				timeStart = bookingTimeSlot.TimeSlot.TimeStart
				timeEnd = bookingTimeSlot.TimeSlot.TimeEnd
				name = bookingTimeSlot.TimeSlot.EmployeeService.Service.SerName
			case model.BookingTypeServiceItem:
				c.DB.Preload("ServiceItem", func(db *gorm.DB) *gorm.DB {
					return db.Preload("Service")
				}).Find(&bookingServiceItem)
				timeStart = booking.BookedDate.Format("15:01")
				timeEnd = booking.BookedDate.Format("15:01")
				name = bookingServiceItem.ServiceItem.Service.SerName
			case model.BookingTypePackage:
				c.DB.Preload("Package").Find(&bookingPackage)
				timeStart = booking.BookedDate.Format("15:01")
				timeEnd = booking.BookedDate.Format("15:01")
				name = bookingPackage.Package.PacName
			}
		}
		template += fmt.Sprintf(ReviewTemplate,
			transaction.TranDocumentCode,
			transaction.ChatChannel.ChaAddress,
			timeStart,
			timeEnd,
			name,
			fmt.Sprintf("line://app/%s?transactionID=%d", setting[model.NameLIFFIDReport], transaction.ID),
		) + ","
	}
	template = fmt.Sprintf(carouselTemplate, template[:len(template)-1])

	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(template))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("sss", flexContainer), nil
}

var ReviewTemplate string = `{
	"type": "bubble",
	"body": {
	  "type": "box",
	  "layout": "vertical",
	  "contents": [
		{
		  "type": "text",
		  "weight": "bold",
		  "size": "xl",
		  "text": "Review"
		},
		{
		  "type": "box",
		  "layout": "vertical",
		  "margin": "lg",
		  "spacing": "sm",
		  "contents": [
			{
			  "type": "box",
			  "layout": "baseline",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "text",
				  "text": "DocID",
				  "color": "#aaaaaa",
				  "size": "sm",
				  "flex": 2
				},
				{
				  "type": "text",
				  "text": "%s",
				  "wrap": true,
				  "color": "#666666",
				  "size": "sm",
				  "flex": 5
				}
			  ]
			},
			{
			  "type": "box",
			  "layout": "baseline",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "text",
				  "text": "Place",
				  "color": "#aaaaaa",
				  "size": "sm",
				  "flex": 2
				},
				{
				  "type": "text",
				  "text": "%s",
				  "wrap": true,
				  "color": "#666666",
				  "size": "sm",
				  "flex": 5
				}
			  ]
			},
			{
			  "type": "box",
			  "layout": "baseline",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "text",
				  "text": "Time",
				  "color": "#aaaaaa",
				  "size": "sm",
				  "flex": 2
				},
				{
				  "type": "text",
				  "text": "%s - %s",
				  "wrap": true,
				  "color": "#666666",
				  "size": "sm",
				  "flex": 5
				}
			  ]
			},
			{
			  "type": "box",
			  "layout": "baseline",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "text",
				  "color": "#aaaaaa",
				  "size": "sm",
				  "flex": 2,
				  "text": "service"
				},
				{
				  "type": "text",
				  "text": "%s",
				  "wrap": true,
				  "color": "#666666",
				  "size": "sm",
				  "flex": 5
				}
			  ]
			}
		  ]
		}
	  ]
	},
	"footer": {
	  "type": "box",
	  "layout": "vertical",
	  "spacing": "sm",
	  "contents": [
		{
		  "type": "button",
		  "action": {
			"type": "uri",
			"label": "Review",
			"uri": "%s"
		  }
		}
	  ],
	  "flex": 0
	}
  }`
