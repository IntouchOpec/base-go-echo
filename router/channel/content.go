package channel

import (
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

func ContentListHandler(c echo.Context) (linebot.SendingMessage, error) {
	flexContainerStr := ""
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexContainerStr))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("service", flexContainer), err
}
