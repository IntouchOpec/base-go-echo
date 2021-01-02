package channel

import (
	"fmt"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

func LocationHandler(c *Context) (linebot.SendingMessage, error) {
	Latitude, err := strconv.ParseFloat(c.AccountLine.Settings["Latitude"], 64)
	fmt.Println(Latitude)
	if err != nil {
		fmt.Println(err, "err location")
		return nil, err
	}
	Longitude, err := strconv.ParseFloat(c.AccountLine.Settings["Latitude"], 64)
	if err != nil {
		fmt.Println(err, "Test")
		return nil, err
	}

	return linebot.NewLocationMessage(c.AccountLine.ChaName, c.AccountLine.ChaAddress, Latitude, Longitude), nil
}
