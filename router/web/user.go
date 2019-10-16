package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

func UserListHandler(c *Context) error {
	users := []*model.User{}

	if err := model.DB().Find(&users).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "user-list", echo.Map{
		"list":  users,
		"title": "user",
	})
	return err
}

func UserDetailHandler(c *Context) error {
	id := c.Param("id")
	var user model.User
	a := auth.Default(c)
	if err := model.DB().Where("account_id = ?", a.GetAccountID()).Find(&user, id).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "user-detail", echo.Map{
		"detail": user,
		"title":  "user",
	})
	return err
}
