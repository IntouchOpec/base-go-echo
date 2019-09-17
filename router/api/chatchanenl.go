package api

import (
	"fmt"
	"net/http"
	"strconv"

	// . "github.com/IntouchOpec/base-go-echo/conf"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

// CreateChatChannel route create chat_channel.
func CreateChatChannel(c echo.Context) error {

	cha := model.ChatChannel{}
	if err := c.Bind(&cha); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	bot, err := lib.ConnectLineBot(cha.ChannelSecret, cha.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	// config web Conf.Server.DomainWeb
	// view := linebot.View{Type: "full", URL: fmt.Sprintf("%s/register/%s", Conf.Server.DomainWeb, cha.LineID)}
	view := linebot.View{Type: "full", URL: "https://d2670202.ngrok.io/register/" + cha.LineID}

	if err := c.Bind(&view); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.AddLIFF(view).Do()

	if err != nil {
		fmt.Println(err)

		return c.NoContent(http.StatusBadRequest)
	}
	cha.SaveChatChannel()
	setting := model.Setting{Name: "LIFFregister", Value: res.LIFFID, ChatChannelID: cha.ID}
	setting.SaveSetting()
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

func GetChatChannelList(c echo.Context) error {
	page := c.QueryParam("page")
	size := c.QueryParam("size")
	chatChannelID := c.Param("chatChannelID")
	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)
	chatChannelIDInt, _ := strconv.Atoi(chatChannelID)

	chatChannels := model.GetChatChannel(chatChannelIDInt, sizeInt, pageInt)

	return c.JSON(http.StatusOK, chatChannels)
}

func GetChatChannelDetail(c echo.Context) error {
	id := c.Param("id")
	chatChannel := model.GetChatChannelByID(id)
	return c.JSON(http.StatusOK, chatChannel)
}

func UpdateChatChannel(c echo.Context) error {
	id := c.Param("id")

	idInt, _ := strconv.Atoi(id)
	chatChannel := model.ChatChannel{}
	if err := c.Bind(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	chatChannel.EditChatChannel(idInt)

	return c.JSON(http.StatusOK, chatChannel)
}

func DeleteChatChannel(c echo.Context) error {
	id := c.Param("id")

	idInt, _ := strconv.Atoi(id)
	chatChannel := model.DeleteChatChannel(idInt)
	return c.JSON(http.StatusOK, chatChannel)
}
