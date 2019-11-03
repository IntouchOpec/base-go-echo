package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
)

func RichMenuListHandler(c *Context) error {
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()

	if err := db.Preload("Account", "ID = ?", a.User.GetAccountID).Find(&chatChannel); err != nil {
		return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
			"list":  "",
			"title": "rich-menu",
		})
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	rows := ""

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
			<a href="/admin/richmenu/%s">
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
		</tr>`, row.RichMenuID, row.Size.Height, row.Size.Width, row.Selected, row.Name, row.ChatBarText, row.RichMenuID, row.RichMenuID, row.RichMenuID)
	}
	return c.Render(http.StatusOK, "rich-menu-list", echo.Map{
		"list":  rows,
		"title": "rich-menu",
	})
}

func RichMenuDetailHandler(c *Context) error {
	richID := c.Param("richID")
	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()

	db.Preload("Account", "ID = ?", a.User.GetAccountID).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.GetRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.Render(http.StatusOK, "rich-menu-detail", echo.Map{
		"detail": res,
		"title":  "rich-menu",
	})
}

func RichMenuActiveHandler(c *Context) error {
	richID := c.Param("richID")

	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	qpa := c.QueryParams()
	db.Preload("Account", "ID = ?", a.User.GetAccountID).Where(qpa).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.SetDefaultRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, res)
}

func RichMenuImageHandler(c *Context) error {
	richID := c.Param("richID")

	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	qpa := c.QueryParams()
	db.Preload("Account", "ID = ?", a.User.GetAccountID).Where(qpa).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.SetDefaultRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, res)
}

func RichMenuDeleteHandler(c *Context) error {
	richID := c.Param("richID")

	chatChannel := model.ChatChannel{}
	a := auth.Default(c)
	db := model.DB()
	qpa := c.QueryParams()
	db.Preload("Account", "ID = ?", a.User.GetAccountID).Where(qpa).Find(&chatChannel)

	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res, err := bot.DeleteRichMenu(richID).Do()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, res)
}
