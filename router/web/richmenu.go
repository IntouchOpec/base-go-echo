package web

import (
	"encoding/json"
	"fmt"
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
	filterChatChannelDB := db.Where("cha_account_id = ?", a.GetAccountID()).Find(&chatChannels)
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

	for _, row := range res {
		rows = rows + fmt.Sprintf(`
		<tr>
			<td>%s</td>
			<td>%d</td>
			<td>%d</td>
			<td>%t</td>
			<td>%s</td>
			<td>%s</td>
			<td class="td-actions text-right">
			<a href="/admin/richmenu/%s?chat_channel_id=%d">
			  <button type="button" rel="tooltip" class="btn btn-success">
				<i class="material-icons">settings_applications</i>
			  </button>
			</a>
			<a href="/admin/richmenu/%s/edit">
			  <button type="button" rel="tooltip" class="btn btn-warning">
				  <i class="material-icons">edit</i>
			  </button>
			</a>
			<a href="/admin/richmenu/%s/delete">
			  <button type="button" rel="tooltip" class="btn btn-danger">
				  <i class="material-icons">close</i>
			  </button>
			</a>
		  </td>
		</tr>`,
			row.RichMenuID,
			row.Size.Height,
			row.Size.Width,
			row.Selected,
			row.Name,
			row.ChatBarText,
			row.RichMenuID,
			chatChannel.ID,
			row.RichMenuID,
			row.RichMenuID)
	}
	if len(res) == 0 {
		rows = `<tr><td colspan="8" class="text-center">No content</td></tr>`
	}

	return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
		"detail":       chatChannel,
		"list":         rows,
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

	if err := db.Where("cha_account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID).Error; err != nil {
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
	chatChannelID := c.QueryParam("chat_channel_id")

	if err := db.Where("cha_account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID).Error; err != nil {
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
	chatChannelID := c.QueryParam("chatChannelID")
	file := c.FormValue("file")

	fileURL, file, err := lib.UploadteImage(file)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	db.Where("cha_account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID)

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
	db.Where("cha_account_id = ?", a.GetAccountID()).Find(&chatChannels)

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
	richID := c.Param("richID")

	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	qpa := c.QueryParams()
	db.Preload("Account", "ID = ?", a.User.GetAccountID).Where(qpa).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelID, chatChannel.ChaChannelSecret)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.DeleteRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, res)
}

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
