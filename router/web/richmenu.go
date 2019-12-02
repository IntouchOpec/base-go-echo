package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
)

func RichMenuListHandler(c *Context) error {
	chatChannels := []model.ChatChannel{}
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	var rows string
	filterChatChannelDB := db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannels)
	chatChannelID := c.QueryParam("chat_channel_id")

	filterChatChannelDB.First(&chatChannel, chatChannelID)

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		rows = `<tr><td colspan="8" class="text-center">No content</td></tr>`
		return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
			"chatChannels": chatChannels,
			"detail":       chatChannel,
			"list":         rows,
			"title":        "rich-menu",
		})
	}

	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		rows = `<tr><td colspan="8" class="text-center">No content</td></tr>`
		return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
			"detail":       chatChannel,
			"chatChannels": chatChannels,
			"list":         rows,
			"title":        "rich-menu",
		})
	}

	return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
		"detail":       chatChannel,
		"list":         &res,
		"chatChannels": chatChannels,
		"title":        "rich-menu",
	})
}

func RichMenuDetailHandler(c *Context) error {
	richID := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	chatChannelID := c.QueryParam("chat_channel_id")

	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	res, err := bot.GetRichMenu(richID).Do()

	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}

	setting := model.Setting{}
	db.Where("name = ?", res.RichMenuID).Find(&setting)
	return c.Render(http.StatusOK, "rich-menu-detail", echo.Map{
		"detail":        &res,
		"title":         "rich-menu",
		"ImageRichMenu": setting,
	})
}

func RichMenuActiveHandler(c *Context) error {
	richID := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	chatChannelID := c.FormValue("chat_channel_id")
	fmt.Println(chatChannelID)
	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	res, err := bot.SetDefaultRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, res)
}

func RichMenuImageHandler(c *Context) error {
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	id := c.Param("id")
	chatChannelID := c.QueryParam("chat_channel_id")
	file := c.FormValue("file")

	fileURL, file, err := lib.UploadteImage(file)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID)

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	_, err = bot.UploadRichMenuImage(id, file).Do()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	setting := model.Setting{Name: id, Value: fileURL}
	db.Create(&setting)
	return c.JSON(http.StatusCreated, setting)
}

func RichMenuDonwloadImage(c *Context) error {
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	id := c.Param("id")
	chatChannelID := c.QueryParam("chat_channel_id")
	db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID)

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	res, err := bot.DownloadRichMenuImage(id).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println(res.ContentType)
	fmt.Println(res.ContentLength)
	data, err := ioutil.ReadAll(res.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = ioutil.WriteFile("public/assets/images/"+id+".jpg", data, 0666)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	urlImage := "/images/" + id + ".jpg"
	db.Model(&chatChannel).Association("Settings").Append(model.Setting{Name: id, Value: urlImage})

	return c.JSON(http.StatusOK, echo.Map{
		"urlImage": urlImage,
	})
}

type RichMenuSizeType struct {
	Size linebot.RichMenuSize `json:"size"`
	Name string               `json:"name"`
	ID   string               `json:"id"`
}

func RichMenuCreateViewHandler(c *Context) error {
	chatChannels := []model.ChatChannel{}
	richMenu := []RichMenuSizeType{RichMenuSizeType{
		Size: linebot.RichMenuSize{Width: 2500, Height: 843},
		ID:   "2",
		Name: "2500x843",
	}, RichMenuSizeType{
		Size: linebot.RichMenuSize{Width: 2500, Height: 1686},
		ID:   "1",
		Name: "2500x1686",
	}}
	db := model.DB()
	a := auth.Default(c)
	db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannels)

	return c.Render(http.StatusOK, "rich-menu-form", echo.Map{
		"chatChannels": chatChannels,
		"title":        "rich-menu",
		"sizes":        richMenu,
	})
}

type RichMenuReq struct {
	Size          string `form:"Size"`
	Name          string `form:"Name"`
	ChatBarText   string `form:"ChatBarText"`
	Areas         string `form:"Areas"`
	ChatChannelID string `form:"ChatChannelID"`
}

func RichMenuCreateHandler(c *Context) error {
	richMenu := RichMenuReq{}
	chatChannel := model.ChatChannel{}
	chatChennalID := c.FormValue("chat_channel_id")
	db := model.DB()

	if err := c.Bind(&richMenu); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Find(&chatChannel, chatChennalID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	richMenuSize := linebot.RichMenuSize{}
	err = json.Unmarshal([]byte(richMenu.Size), &richMenuSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	areas := []linebot.AreaDetail{}
	err = json.Unmarshal([]byte(richMenu.Areas), &areas)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	richMenuLineBot := linebot.RichMenu{
		Size:        richMenuSize,
		Name:        richMenu.Name,
		ChatBarText: richMenu.ChatBarText,
		Areas:       areas,
	}

	res, err := bot.CreateRichMenu(richMenuLineBot).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	redirect := fmt.Sprintf("/admin/richmenu/%s?chat_channel_id=%d", res.RichMenuID, chatChannel.ID)
	return c.JSON(http.StatusCreated, echo.Map{
		"data":     res,
		"redirect": redirect,
	})
}

func RichMenuDeleteHandler(c *Context) error {
	richID := c.Param("id")

	chatChannel := model.ChatChannel{}
	accountID := auth.Default(c).GetAccountID()
	db := model.DB()
	id := c.QueryParam("id")
	if err := db.Where("account_id = ?", accountID).Find(&chatChannel, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println(richID)
	res, err := bot.DeleteRichMenu(richID).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, res)
}

// func RichMenuRemoveHandler(c *Context) error{

// }
// func RichMenuRemoveImage(c *Context) error {
// 	a := auth.Default(c)
// 	db := model.DB()
// 	chatChannel := model.ChatChannel{}
// 	id :=
// 	chatChennalID := c.QueryParam("chatChennalID")
// 	if err := db.Find(&chatChannel, chatChennalID).Error; err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}

// 	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	bot.DownloadRichMenuImage()
// 	return c.JSON(http.StatusOK, )
// }
