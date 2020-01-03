package channel

import (
	"fmt"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

func LocationHandler(c *Context) (linebot.SendingMessage, error) {
	chatChannel := c.ChatChannel
	position := chatChannel.GetSetting([]string{"Latitude", "Longitude"})
	Latitude, err := strconv.ParseFloat(position["Latitude"], 64)
	if err != nil {
		return nil, err
	}
	Longitude, err := strconv.ParseFloat(position["Longitude"], 64)
	if err != nil {
		fmt.Println(err, "Test")
		return nil, err
	}

	return linebot.NewLocationMessage(chatChannel.ChaName, chatChannel.ChaAddress, Latitude, Longitude), nil
}
