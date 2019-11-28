package web

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/IntouchOpec/base-go-echo/lib"
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
	model.DB().Preload("Account").Preload("ChatChannels").Where("plac_account_id = ? ", a.User.GetAccountID()).Find(&place, id)
	err := c.Render(http.StatusOK, "place-detail", echo.Map{
		"detail": place,
		"title":  "place",
	})
	return err
}

func PlaceCreateHandler(c *Context) error {
	Place := model.Place{}
	PlacTypes := []model.PlaceType{model.PlaceRoom}
	err := c.Render(http.StatusOK, "place-form", echo.Map{
		"detail":    Place,
		"title":     "place",
		"PlacTypes": PlacTypes,
	})
	return err
}

func PlacePostHandler(c *Context) error {
	a := auth.Default(c)
	place := model.Place{}
	if err := c.Bind(&place); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	file := c.FormValue("file")
	fileUrl, _, err := lib.UploadteImage(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	place.PlacImage = fileUrl
	place.AccountID = a.GetAccountID()
	if err := place.CreatePlace(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"data":     place,
		"title":    "place",
		"redirect": fmt.Sprintf("/admin/place/%d", place.ID),
	})
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
	model.DB().Where("plac_account_id = ? ", a.User.GetAccountID()).Find(&place, id)
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

func PlaceAddChatChannelViewHanlder(c *Context) error {
	place := model.Place{}
	a := auth.Default(c)
	chatChannels := []model.ChatChannel{}
	db := model.DB()
	db.Where("cha_account_id = ?", a.GetAccountID()).Find(&chatChannels)
	db.Where("plac_account_id = ?", a.GetAccountID()).Find(&place)
	return c.Render(http.StatusOK, "place-chat-channel-form", echo.Map{
		"chatChannels": chatChannels,
		"title":        "place",
	})
}

func PlaceAddChatChannelPostHanlder(c *Context) error {
	place := model.Place{}
	id := c.Param("id")
	a := auth.Default(c)
	chatChannel := model.ChatChannel{}

	chatChannelID := c.FormValue("chat_channel_id")
	fmt.Println(chatChannelID)
	db := model.DB()
	db.Where("cha_account_id = ?", a.GetAccountID()).Find(&chatChannel, chatChannelID)
	if err := db.Where("plac_account_id = ? ", a.GetAccountID()).Find(&place, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println("======", chatChannel)
	if err := db.Model(&place).Association("ChatChannels").Append(&chatChannel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"data":     place,
		"redirect": fmt.Sprintf("/admin/place/%d", place.ID),
	})
}
