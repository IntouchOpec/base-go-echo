package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

func UserListHandler(c *Context) error {
	users := []*model.User{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	filterUser := db.Preload("Account").Where("account_id = ?", a.GetAccountID()).Find(&users).Count(&total)
	filterUser.Limit(limit).Offset(page).Find(&users)
	pagination := MakePagination(total, page, limit)
	if err := model.DB().Find(&users).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "user-list", echo.Map{
		"list":       users,
		"title":      "user",
		"pagination": pagination,
	})
	return err
}

func UserDetailHandler(c *Context) error {
	id := c.Param("id")
	var user model.User

	user.GetById(id)
	err := c.Render(http.StatusOK, "user-detail", echo.Map{
		"detail": user,
		"title":  "user",
	})
	return err
}

func UserFormHamdeler(c *Context) error {
	var user model.User

	err := c.Render(http.StatusOK, "user-form", echo.Map{
		"detail": user,
		"title":  "user",
	})
	return err
}

func UserEditHamdeler(c *Context) error {
	var user model.User
	id := c.Param("id")
	user.GetById(id)

	err := c.Render(http.StatusOK, "user-form", echo.Map{
		"detail": user,
		"title":  "user",
	})
	return err
}

func UserDeleteHandler(c *Context) error {
	id := c.Param("id")
	user := model.DeleteUser(id)

	err := c.JSON(http.StatusOK, user)
	return err
}
