package api

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// CreateChatChannel route create chat_channel.
func CreateChatChannel(c echo.Context) error {
	cha := model.ChatChannel{}
	if err := c.Bind(&cha); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	cha.SaveChatChannel()
	return c.JSON(http.StatusOK, cha)
}

// GetChannelAccessToken route get channel access token.
func GetChannelAccessToken(c echo.Context) error {
	id := c.Param("id")

	chatChannel := model.ChatChannel{}

	db := model.DB()

	if err := db.Find(&chatChannel, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	res, err := bot.IssueAccessToken(chatChannel.ChannelID, chatChannel.ChannelSecret).Do()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	chatChannel.ChannelAccessToken = res.AccessToken

	if err := db.Save(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, chatChannel)
}
