package web

import (
	"fmt"
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
	filterUser := db.Model(&users).Where("account_id = ?", a.GetAccountID()).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterUser.Limit(pagination.Record).Offset(pagination.Offset).Find(&users)
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

func UserPutHandler(c *Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user.SetPassword()

	if err := user.UpdateUser(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"detail":   user,
		"redirect": fmt.Sprintf("/admin/users/%d", user.ID),
	})
}

func UserPostHandler(c *Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user.SetPassword()
	accID := auth.Default(c).GetAccountID()

	user.Create(accID)
	return c.JSON(http.StatusCreated, echo.Map{
		"detail":   user,
		"redirect": fmt.Sprintf("/admin/users/%d", user.ID),
	})
}

func UserFormHandler(c *Context) error {
	var user model.User
	csrfValue := c.Get("_csrf")

	err := c.Render(http.StatusOK, "user-form", echo.Map{
		"method": "PUT",
		"detail": user,
		"title":  "user",
		"_csrf":  csrfValue,
	})
	return err
}

func UserEditHandler(c *Context) error {
	var user model.User
	id := c.Param("id")
	user.GetById(id)

	err := c.Render(http.StatusOK, "user-form", echo.Map{
		"detail": user,
		"title":  "user",
		"method": "PUT",
	})
	return err
}

func UserDeleteHandler(c *Context) error {
	id := c.Param("id")
	user := model.DeleteUser(id)

	err := c.JSON(http.StatusOK, user)
	return err
}
