package channel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/omise/omise-go"
)

const (
	// Read these from environment variables or configuration files!
	OmisePublicKey = "pkey_test_521w1g1t7w4x4rd22z0"
	OmiseSecretKey = "skey_test_521w1g1t6yh7sx4pu8n"
)

type Event struct {
	Base

	Key  string      `json:"key"`
	Data interface{} `json:"data"`
}

type Base struct {
	Object   string    `json:"object"`
	ID       string    `json:"id"`
	Live     bool      `json:"livemode"`
	Location *string   `json:"location"`
	Created  time.Time `json:"created"`
}

type eventShim struct {
	Base
	Key  string    `json:"key"`
	Data *Deletion `json:"data"`
}

type Deletion struct {
	Base
	Deleted bool `json:"deleted"`
}

type omiseContext struct {
	Context *echo.Context
	Event   *eventShim
	DB      *gorm.DB
}

func HandlerOmiseWebHook(c echo.Context) error {
	fmt.Println(c.Request())
	context := omiseContext{}
	context.Context = &c
	db := model.DB()
	context.DB = db
	event := &eventShim{}
	if err := json.NewDecoder(c.Request().Body).Decode(event); err != nil {
		// return err
	}
	ev, err := json.Marshal(event)
	omise := model.OmiseLog{Json: ev}
	if err := db.Create(&omise).Error; err != nil {
		fmt.Println("err", err)
		// return err
	}
	if err != nil {
		fmt.Println("err,", err)
		// return err
	}
	if err := context.dataInstanceFromType(event.Data.Object, ev); err != nil {
		// return err
	}
	return c.JSON(http.StatusOK, echo.Map{})
}

func (c *omiseContext) dataInstanceFromType(typ string, body []byte) error {
	switch typ {
	case "charge":
		charge := eventUnmarshal(body, &omise.Charge{})
		fmt.Println(charge)
	case "customer":
		customer := eventUnmarshal(body, &omise.Customer{})
		fmt.Println("contomer", customer)
	case "card":
		card := eventUnmarshal(body, &omise.Card{})
		fmt.Println("card", card)
	case "dispute":
		dispute := eventUnmarshal(body, &omise.Dispute{})
		fmt.Println("dispute", dispute)
	case "recipient":
		recipient := eventUnmarshal(body, &omise.Recipient{})
		fmt.Println("recipient", recipient)
	case "refund":
		refund := eventUnmarshal(body, &omise.Refund{})
		fmt.Println("refund", refund)
	case "transfer":
		transfer := eventUnmarshal(body, &omise.Transfer{})
		fmt.Println("transfer", transfer)
	}
	return nil
}

func eventUnmarshal(data []byte, eventStruct interface{}) interface{} {
	if err := json.Unmarshal(data, eventStruct); err != nil {
		fmt.Println(err)
	}
	return eventStruct
}
