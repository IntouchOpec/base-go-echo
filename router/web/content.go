package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

func GetContentHandler(c echo.Context) error {
	contentID := c.QueryParam("contentID")
	chaLineID := c.Param("cha_line_id")
	fmt.Println(contentID, chaLineID)
	var content model.Content
	var chatCha model.ChatChannel
	db := model.DB()
	if err := db.Where("cha_line_id = ?", chaLineID).Find(&chatCha).Error; err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, "content-liff", echo.Map{})
	}
	if err := db.Where("account_id = ?", chatCha.AccountID).Find(&content, contentID).Error; err != nil {
		return c.Render(http.StatusOK, "content-liff", echo.Map{})
	}
	return c.Render(http.StatusOK, "content-liff", echo.Map{
		"detail": content,
	})
}

func ContentListHandler(c *Context) error {
	contents := []*model.Content{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	filterChatAns := db.Model(&contents).Where("account_id = ?", a.GetAccountID()).Count(&total)
	pagination := MakePagination(total, page, limit)
	fmt.Println("pagination.Record", pagination.Record)
	filterChatAns.Limit(pagination.Record).Offset(pagination.Offset).Find(&contents)
	fmt.Println("len(contents)", len(contents))
	for i := 0; i < len(contents); i++ {
		fmt.Println(len(contents[i].ConDetail) > 20)
		if len(contents[i].ConDetail) > 20 {
			fmt.Println(contents[i].ConDetail[:20])
			contents[i].ConDetail = contents[i].ConDetail[:20]
		}

	}
	return c.Render(http.StatusOK, "content-list", echo.Map{
		"title":      "content",
		"list":       contents,
		"pagination": pagination,
	})
}

func ContentDetailHandler(c *Context) error {
	id := c.Param("id")
	content := model.Content{}
	a := auth.Default(c)

	model.DB().Where("account_id = ?", a.GetAccountID()).Find(&content, id)
	return c.Render(http.StatusOK, "content-detail", echo.Map{
		"title":  "content",
		"detail": content,
	})
}

func ContentCreateHandler(c *Context) error {
	content := model.Content{}
	return c.Render(http.StatusOK, "content-form", echo.Map{
		"method": "POST",
		"title":  "content",
		"detail": content,
	})
}

func ContentPostHandler(c *Context) error {
	content := model.Content{}
	if err := c.Bind(&content); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	a := auth.Default(c)
	file := c.FormValue("file")
	ctx := context.Background()
	imagePath, err := lib.UploadGoolgeStorage(ctx, file, "images/content/")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	content.ConImage = imagePath
	content.AccountID = a.GetAccountID()
	err = content.SaveContent()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/content/%d", content.ID)
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
		"data":     content,
	})
}

func ContentEditHandler(c *Context) error {
	id := c.Param("id")
	content := model.Content{}
	a := auth.Default(c)
	model.DB().Where("account_id = ?", a.GetAccountID()).Find(&content, id)
	return c.Render(http.StatusOK, "content-form", echo.Map{
		"method": "PUT",
		"title":  "content",
		"detail": content,
	})
}

func ContentPutHandler(c *Context) error {
	id := c.Param("id")
	content := model.Content{}
	a := auth.Default(c).GetAccountID()
	db := model.DB()
	if err := db.Where("account_id = ?", a).Find(&content, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Bind(&content); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Save(&content).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/content/%d", content.ID)

	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
		"data":     content,
	})
}

func ContentPatchHandler(c *Context) error {
	id := c.Param("id")
	content := model.Content{}
	a := auth.Default(c)
	db := model.DB()
	IsActive := c.FormValue("con_is_active")
	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&content, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	b, err := strconv.ParseBool(IsActive)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	content.ConIsActive = b
	if err := db.Save(&content).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/content/%d", content.ID)

	return c.JSON(http.StatusOK, echo.Map{
		"redirect": redirect,
		"data":     content,
	})
}

func ContentDeleteHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	if err := model.DeleteContent(id, accID); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{})
}
