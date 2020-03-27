package channel

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/IntouchOpec/base-go-echo/lib/lineapi"
	"github.com/IntouchOpec/base-go-echo/middleware/session"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
	DB             *gorm.DB
	sqlDb          *sql.DB
	AccountLine    *model.AccountLine
	Source         *linebot.EventSource
	Massage        *linebot.TextMessage
	Account        model.Account
	ChatChannel    model.ChatChannel
	Customer       model.Customer
	ClientLine     *lineapi.ClientLine
	Event          *linebot.Event
	PostbackAction *PostbackAction
}

type PostbackAction struct {
	Action        string `json:"action"`
	Type          string `json:"type"`
	DateStr       string `json:"date"`
	TimeStr       string `json:"time"`
	ServiceID     string `json:"service_id"`
	ServiceItemID string `json:"service_item_id"`
	PackageID     string `json:"package_id"`
	TimeSlotID    string `json:"time_slot_id"`
	EmployeeID    string `json:"employee_id"`
}

type Pagination struct {
	Page      int  `json:"page"`
	Previous  bool `json:"previous"`
	Next      bool `json:"next"`
	StartPage int  `json:"start_page"`
	Record    int  `json:"record"`
	Offset    int  `json:"offset"`
}

type (
	HandlerFunc func(*Context) string
)

func (ctx *Context) Session() session.Session {
	return session.Default(ctx)
}

func (pagi Pagination) SetPagination() {
	if pagi.Offset == 0 {
		pagi.Offset = 9
	}
	if pagi.Page == 0 {
		pagi.Page = 1
	}
}

func (pagi *Pagination) MakePagination(total, limit int) {
	pagi.Previous = false
	pagi.Next = false
	countPage := total / limit

	pagi.StartPage = 1
	if pagi.Page > 1 {
		pagi.Previous = true
	}
	if pagi.Page > countPage {
		pagi.Next = true
	}
	pagi.StartPage = pagi.Page + 1
	if pagi.StartPage <= 2 {
		pagi.StartPage = 1
	}
	pagi.Record = total - limit*pagi.Page
	pagi.Offset = limit * pagi.Page
	pagi.Record = 9
}

func (pagi *Pagination) ParseQueryUnmarshal(param string) error {
	str, err := convertParamToJsonString(param)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(str), &pagi); err != nil {
		return err
	}
	return nil
}

func (post *PostbackAction) ParseQueryUnmarshal(param string) error {
	str, err := convertParamToJsonString(param)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(str), &post); err != nil {
		return err
	}
	return nil
}

func convertParamToJsonString(param string) (string, error) {
	u, err := url.Parse(fmt.Sprintf("?%s", param))
	if err != nil {
		return "", err
	}
	var str string
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	for key, value := range q {
		str += fmt.Sprintf(`"%s": "%s",`, key, value[0])
	}
	str = fmt.Sprintf(fmt.Sprintf("{%s}", str[:len(str)-1]))
	return str, nil
}

func (pagi *Pagination) MakePaginationTemplate(action string) string {
	var button string
	if pagi.Next == true {
		button = fmt.Sprintf(buttonTamplate, "ถัดไป", fmt.Sprintf("action=%s&page=%d", action, pagi.Page+1))
	}
	if pagi.Previous == true {
		if pagi.Next == true {
			button += ","
		}
		button += fmt.Sprintf(buttonTamplate, "ย้อนกลับ", fmt.Sprintf("action=%s&page=%d", action, pagi.Page-1))
	}

	return fmt.Sprintf(paginationTemplate, button)
}

var buttonTamplate string = `
	{ "type": "button", "margin": "xs", "style": "primary", "action": 
		{ "type": "postback", "label": "%s", "data": "%s" } }`

var paginationTemplate string = `{ "type": "bubble", "direction": "ltr", "body": { "type": "box", "layout": "vertical", "contents": [ %s ] } }`
