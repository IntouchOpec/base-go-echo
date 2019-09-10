package api

import (
	"io"
	"net/http"
	"os"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

//GetRichMenu
func GetRichMenu(c echo.Context) error {
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
	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, res)
}

// CreateRichMenu
func CreateRichMenu(c echo.Context) error {
	chatChennalID := c.Param("id")
	richMenu := linebot.RichMenu{}
	chatChannel := model.ChatChannel{}

	if err := c.Bind(&richMenu); err != nil {
		return nil
	}

	db := model.DB()

	if err := db.Find(&chatChannel, chatChennalID).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, res)
}

// UploadImageRichMenu
func UploadImageRichMenu(c echo.Context) error {
	id := c.Param("id")
	richMenuID := c.Param("richMenuID")

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	chatChannel := model.ChatChannel{}

	db := model.DB()

	if err := db.Find(&chatChannel, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res := bot.UploadRichMenuImage(richMenuID, file.Filename)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, res)

}

// AtiveRichMenu
func AtiveRichMenu(c echo.Context) error {
	id := c.Param("id")
	richMenuID := c.Param("richMenuID")
	chatChannel := model.ChatChannel{}

	db := model.DB()

	if err := db.Find(&chatChannel, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res := bot.SetDefaultRichMenu(richMenuID)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, res)
}
