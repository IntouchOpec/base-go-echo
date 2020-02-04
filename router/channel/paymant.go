package channel

import "github.com/line/line-bot-sdk-go/linebot"

func PaymentHandler(c *Context) (linebot.SendingMessage, error) {
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(cardPatmentTemplate))
	if err != nil {
		return nil, err
	}
	return linebot.NewFlexMessage("", flexContainer), nil
}
