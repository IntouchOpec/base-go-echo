package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/lib/lineapi"
	"github.com/IntouchOpec/base-go-echo/lib/uploadgoolgestorage"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

func RichMenuListHandler(c *Context) error {
	chatChannels := []model.ChatChannel{}
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	filterChatChannelDB := db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannels)
	chatChannelID := c.QueryParam("chat_channel_id")

	filterChatChannelDB.First(&chatChannel, chatChannelID)

	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
			"chatChannels": chatChannels,
			"detail":       chatChannel,
			"err":          err,
			"title":        "rich-menu",
		})
	}

	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
			"detail":       chatChannel,
			"chatChannels": chatChannels,
			"err":          err,
			"title":        "rich-menu",
		})
	}

	return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
		"detail":       chatChannel,
		"list":         res,
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

	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

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
		"chatChannel":   chatChannel,
	})
}

func RichMenuActiveHandler(c *Context) error {
	richID := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	chatChannelID := c.FormValue("chat_channel_id")

	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	res, err := bot.SetDefaultRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, res)
}

func RichMenuAddCustomerHandler(c *Context) error {
	chatChannel := model.ChatChannel{}
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	id := c.Param("id")
	chatChannelID := c.QueryParam("chat_channel_id")
	lineID := c.FormValue("line_id")
	db.Where("account_id = ?", accID).Find(&chatChannel, chatChannelID)
	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	res, err := bot.BulkLinkRichMenu(id, lineID).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, res)
}

func RichMenuImageHandler(c *Context) error {
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	id := c.Param("id")
	chatChannelID := c.QueryParam("chat_channel_id")
	file := c.FormValue("file")

	_, fileURL, err := uploadgoolgestorage.UploadteImage(file)
	ctx := context.Background()
	image, err := uploadgoolgestorage.UploadGoolgeStorage(ctx, file, "images/richMenu/")

	if err != nil {
		fmt.Println("err1", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID)

	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		fmt.Println("err2", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	_, err = bot.UploadRichMenuImage(id, fileURL).Do()
	if err != nil {
		fmt.Println("err3", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	setting := model.Setting{Name: id, Value: image}
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

	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	res, err := bot.DownloadRichMenuImage(id).Do()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
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
		"method":       "POST",
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

	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

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

// func RichMenuPutHandler(c *Context) error {
// 	richMenu := RichMenuReq{}
// 	chatChannel := model.ChatChannel{}
// 	chatChennalID := c.FormValue("chat_channel_id")
// 	db := model.DB()

// 	if err := c.Bind(&richMenu); err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}

// 	if err := db.Find(&chatChannel, chatChennalID).Error; err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}

// 	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}

// 	richMenuSize := linebot.RichMenuSize{}
// 	err = json.Unmarshal([]byte(richMenu.Size), &richMenuSize)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}

// 	areas := []linebot.AreaDetail{}
// 	err = json.Unmarshal([]byte(richMenu.Areas), &areas)

// 	if err != nil {
// 		fmt.Println(err)
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	richMenuLineBot := linebot.RichMenu{
// 		Size:        richMenuSize,
// 		Name:        richMenu.Name,
// 		ChatBarText: richMenu.ChatBarText,
// 		Areas:       areas,
// 	}

// 	richID := c.Param("id")

// 	res, err := bot.rich
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}

// 	redirect := fmt.Sprintf("/admin/richmenu/%s?chat_channel_id=%d", res.RichMenuID, chatChannel.ID)
// 	return c.JSON(http.StatusCreated, echo.Map{
// 		"data":     res,
// 		"redirect": redirect,
// 	})
// }

func RichMenuEditView(c *Context) error {
	richID := c.Param("id")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	chatChannelID := c.QueryParam("chat_channel_id")

	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	res, err := bot.GetRichMenu(richID).Do()

	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}

	setting := model.Setting{}
	db.Where("name = ?", res.RichMenuID).Find(&setting)
	return c.Render(http.StatusOK, "rich-menu-form", echo.Map{
		"detail":        &res,
		"method":        "PUT",
		"title":         "rich-menu",
		"ImageRichMenu": setting,
		"chatChannel":   chatChannel,
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

	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
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

// 	bot, err := lineapi.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)

// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	bot.DownloadRichMenuImage()
// 	return c.JSON(http.StatusOK, )
// }
