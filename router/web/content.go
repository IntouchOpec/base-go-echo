package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

func GetContentHandler(c *Context) error {
	contentID := c.QueryParam("contentID")
	var content model.Content
	db := model.DB()
	if err := db.Find(&content, contentID).Error; err != nil {
		return c.Render(http.StatusOK, "content-liff", echo.Map{})
	}
	return c.Render(http.StatusOK, "content-liff", echo.Map{})
}

func ContentListHandler(c *Context) error {
	contents := []*model.Content{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()
	filterChatAns := db.Where("account_id = ?", a.GetAccountID()).Find(&contents).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterChatAns.Limit(pagination.Record).Offset(pagination.Offset).Find(&contents)
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

	content.AccountID = a.GetAccountID()
	err := content.SaveContent()
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
	a := auth.Default(c)
	db := model.DB()
	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&content, id).Error; err != nil {
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

func ContentDeleteHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	if err := model.DeleteContent(id, accID); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{})
}
