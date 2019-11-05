package web

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

// PlaceListHandler
func PlaceListHandler(c *Context) error {
	places := []*model.Place{}
	a := auth.Default(c)
	model.DB().Preload("ChatChannel", func(db *gorm.DB) *gorm.DB {
		return db.Where("chat_channel_id = ?", a.User.GetAccountID())
	}).Preload("PlaceSlot").Find(&places)
	err := c.Render(http.StatusOK, "place-list", echo.Map{
		"list":  places,
		"title": "place",
	})
	return err
}

// PlaceDetailHandler
func PlaceDetailHandler(c *Context) error {
	place := model.Place{}
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Preload("Account").Where("account_id = ? ", a.User.GetAccountID()).Find(&place, id)
	err := c.Render(http.StatusOK, "place-detail", echo.Map{
		"detail": place,
		"title":  "place",
	})
	return err
}

func PlaceCreateHandler(c *Context) error {
	Place := model.Place{}
	csrfValue := c.Get("_csrf")

	err := c.Render(http.StatusOK, "place-form", echo.Map{
		"detail": Place,
		"title":  "place",
		"_csrf":  csrfValue,
	})
	return err
}

func PlaceEditHandler(c *Context) error {
	// place := model.Place{}
	id := c.Param("id")
	a := auth.Default(c)

	pla, err := model.GetPlaceDetail(id, a.GetAccountID())

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err = c.Bind(&pla); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := pla.Update(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.Render(http.StatusOK, "place-detail", echo.Map{
		"detail": pla,
		"title":  "place",
	})
}

func PlaceEditViewHandler(c *Context) error {
	place := model.Place{}
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Preload("Account").Where("account_id = ? ", a.User.GetAccountID()).Find(&place, id)
	err := c.Render(http.StatusOK, "place-form", echo.Map{
		"detail": place,
		"title":  "place",
	})
	return err
}

func PlaceDeleteHandler(c *Context) error {
	id := c.Param("id")
	pla, err := model.DeletePlaceByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, pla)
}
