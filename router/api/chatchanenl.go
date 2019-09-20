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

// SettingType
type SettingType string

const (
	HostAPI         SettingType = "host_api"
	HostWeb         SettingType = "host_web"
	HostLineChaanel SettingType = "host_line_channel"
	LIFFRegister    SettingType = "LIFFregister"
)

// Setting
type Setting struct {
	Value         string            `json:"value" gorm:"unique; type:varchar(25)"`
	Name          SettingType       `json:"name" gorm:"unique; type:varchar(25)"`
	ChatChannelID uint              `form:"chat_channel_id" json:"chat_channel_id" gorm:"not null;"`
	ChatChannel   model.ChatChannel `gorm:"ForeignKey:id"`
}

// CreateChatChannel route create chat_channel.
func CreateChatChannel(c echo.Context) error {

	cha := model.ChatChannel{}
	if err := c.Bind(&cha); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	cha.SaveChatChannel()

	return c.JSON(http.StatusOK, cha)
}

// ActiveRegisterLIFFAPI
func ActiveRegisterLIFFAPI(c echo.Context) error {
	LineID := c.Param("lineID")
	chatChannel := model.ChatChannel{}
	if err := model.DB().Preload("Settings", "name = 'host_web'").Where("line_id = ?", LineID).Find(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	fmt.Println(chatChannel.Settings[0])
	// config web Conf.Server.DomainWeb
	url := fmt.Sprintf("https://%s/register/%s", chatChannel.Settings[0].Value, chatChannel.LineID)

	view := linebot.View{Type: "full", URL: url}
	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.AddLIFF(view).Do()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	setting := model.Setting{Name: "LIFFregister", Value: res.LIFFID}

	if err := model.DB().Model(&chatChannel).Association("Settings").Append(setting).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, chatChannel)
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
