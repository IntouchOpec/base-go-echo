package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

func LIIFListHandler(c *Context) error {
	a := auth.Default(c)
	chatChannel := model.ChatChannel{}
	chatChannels := []model.ChatChannel{}
	db := model.DB()
	ChatChannelID := c.QueryParam("chat_channel_id")
	filterChatChannel := db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannels)
	filterChatChannel.Find(&chatChannel, ChatChannelID)

	bot, err := linebot.New(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		return c.Render(http.StatusBadRequest, "LIFF-list", echo.Map{
			"title":        "LIFF",
			"detail":       chatChannel,
			"chatChannels": chatChannels,
		})
	}

	res, err := bot.GetLIFF().Do()
	if err != nil {
		return c.Render(http.StatusBadRequest, "LIFF-list", echo.Map{
			"err":          err,
			"title":        "LIFF",
			"detail":       chatChannel,
			"chatChannels": chatChannels,
		})
	}
	return c.Render(http.StatusOK, "LIFF-list", echo.Map{
		"list":         &res,
		"detail":       chatChannel,
		"title":        "LIFF",
		"chatChannels": chatChannels,
	})
}

func LIFFCreateHandler(c *Context) error {
	a := auth.Default(c)
	chatChannels := []model.ChatChannel{}

	model.DB().Where("account_id = ?", a.GetAccountID()).Find(&chatChannels)
	LIFFViewTypes := []linebot.LIFFViewType{linebot.LIFFViewTypeCompact, linebot.LIFFViewTypeTall, linebot.LIFFViewTypeFull}
	return c.Render(http.StatusOK, "LIFF-form", echo.Map{"method": "PUT",
		"title":         "LIFF",
		"LIFFViewTypes": LIFFViewTypes,
		"chatChannels":  chatChannels,
	})
}

type LIFFForm struct {
	Type linebot.LIFFViewType `form:"type" json:"type"`
	URL  string               `form:"url" json:"url"`
}

func LIFFPostHandler(c *Context) error {
	a := auth.Default(c)

	var LIFFview LIFFForm
	var chatChannel model.ChatChannel
	chatChannelID := c.QueryParam("chat_channel_id")

	if err := c.Bind(&LIFFview); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	LIFFModel := linebot.View{Type: LIFFview.Type, URL: LIFFview.URL}

	if err := model.DB().Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	bot, _ := linebot.New(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	res, err := bot.AddLIFF(LIFFModel).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/LIFF")
	return c.JSON(http.StatusCreated, echo.Map{
		"detail":   res,
		"redirect": redirect,
	})

}

func LIFFRemoveHandler(c *Context) error {
	a := auth.Default(c)
	var chatChannel model.ChatChannel
	chatChannelID := c.QueryParam("chat_channel_id")
	db := model.DB()
	id := c.Param("id")
	db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID)
	bot, err := linebot.New(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	res, err := bot.DeleteRichMenu(id).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, res)
}
