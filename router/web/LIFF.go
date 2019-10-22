package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

func LIIFLiatHandler(c *Context) error {
	a := auth.Default(c)
	var chatChannel model.ChatChannel
	model.DB().Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel)
	bot, _ := linebot.New(chatChannel.ChannelID, chatChannel.ChannelSecret)

	res, _ := bot.GetLIFF().Do()
	return c.Render(http.StatusOK, "LIFF-list", echo.Map{
		"list":  res,
		"title": "LIFF",
	})
}

func LIFFCreateHandler(c *Context) error {
	a := auth.Default(c)
	var chatChannel model.ChatChannel
	model.DB().Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel)
	LIFFViewTypes := []linebot.LIFFViewType{linebot.LIFFViewTypeCompact, linebot.LIFFViewTypeTall, linebot.LIFFViewTypeFull}
	return c.Render(http.StatusOK, "LIFF-form", echo.Map{
		"title":         "LIFF",
		"LIFFViewTypes": LIFFViewTypes,
	})
}

type LIFFForm struct {
	Type linebot.LIFFViewType `form:"type" json:"type"`
	URL  string               `form:"url" json:"url"`
}

func LIFFPostHanlder(c *Context) error {
	a := auth.Default(c)

	var LIFFview LIFFForm
	var chatChannel model.ChatChannel

	if err := c.Bind(&LIFFview); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	LIFFModel := linebot.View{Type: LIFFview.Type, URL: LIFFview.URL}

	model.DB().Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel)

	bot, _ := linebot.New(chatChannel.ChannelID, chatChannel.ChannelSecret)
	res, err := bot.AddLIFF(LIFFModel).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, res)

}
