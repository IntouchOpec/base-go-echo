package channel

import (
	"fmt"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/lib/lineapi"
	"github.com/IntouchOpec/base-go-echo/model"
)

func ChooseService(c *Context) {
	var packageModels []*model.Package
	var services []*model.Service
	var m string
	var total int
	now := time.Now()
	format := "2006-01-02"
	initial := now.Format(format)
	max := now.AddDate(0, 3, 0).Format(format)
	min := now.Format(format)
	packsFilter := c.DB.Model(&packageModels).Where("account_id = ? and pac_is_active = ?", c.Account.ID, true).Count(&total)
	packsFilter.Limit(9).Offset(1).Order("pac_order").Find(&packageModels)
	for _, pack := range packageModels {
		m += fmt.Sprintf(serviceTemplate,
			fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, pack.PacImage),
			pack.PacName,
			pack.PacPrice,
			fmt.Sprintf("action=booking&type=now&package_id=%d", pack.ID),
			fmt.Sprintf(buttonServiceTemplate,
				fmt.Sprintf("action=booking&type=appointment&package_id=%d", pack.ID),
				initial,
				max,
				min,
				"booking")) + ","
	}
	if len(packageModels) < 9 {
		serviceFilter := c.DB.Model(&services).Where("account_id = ? and ser_active = ?", c.Account.ID, true).Count(&total)
		serviceFilter.Limit(9).Preload("ServiceItems", "ss_is_active = ?", true).Find(&services)
		for index, ser := range services {
			var button string
			var id uint
			for _, serI := range ser.ServiceItems {
				id = serI.ID
				button += fmt.Sprintf(buttonServiceTemplate,
					fmt.Sprintf("action=booking&type=appointment&service_item_id=%d", serI.ID),
					initial,
					max,
					min,
					serI.SSTime.String()) + ","
			}
			if len(button) > 0 {
				button = button[:len(button)-1]
				// fmt.Println(button)
			}
			fmt.Println(button, index)
			m += fmt.Sprintf(serviceTemplate,
				fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, ser.SerImage),
				ser.SerName,
				ser.SerPrice,
				fmt.Sprintf("action=booking&type=now&service_item_id=%d", id), button) + ","
		}
	}
	m = fmt.Sprintf(`{ "replyToken": "%s", "messages":[ { "type": "flex",  "altText":  "รายการบริการ",  "contents": { "type": "carousel", "contents": [ %s ] } }]}`, c.Event.ReplyToken, m[:len(m)-1])

	err := lineapi.SendMessageCustom("reply", c.ChatChannel.ChaChannelAccessToken, m)
	if err != nil {
		fmt.Println(err, "====")
	}
	fmt.Println(err)
}
