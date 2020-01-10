package web

import (
	"context"
	"net/http"

	. "github.com/IntouchOpec/base-go-echo/conf"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

func FileListHandler(c *Context) error {
	fils := []model.File{}
	a := auth.Default(c)
	var total int
	var pagination Pagination
	queryPar := c.QueryParams()
	db := model.DB()
	page, limit := SetPagination(queryPar)

	filterFil := db.Model(&fils).Where("account_id = ?", a.GetAccountID()).Count(&total)
	pagination = MakePagination(total, page, limit)

	filterFil.Limit(pagination.Record).Offset(pagination.Offset).Find(&fils)
	return c.Render(http.StatusOK, "file-list", echo.Map{
		"pagination": pagination,
		"list":       fils,
		"title":      "upload_file",
		"host":       "web." + Conf.Server.Domain,
	})
}

func FileCreateHandler(c *Context) error {
	db := model.DB()
	accountID := auth.Default(c).GetAccountID()
	file := c.FormValue("file")
	fileUrl, _, err := lib.UploadteImage(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	filModel := model.File{
		AccountID: accountID,
		Path:      fileUrl,
	}
	if err := db.Create(&filModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"data": filModel,
		"host": "web." + Conf.Server.Domain,
	})

}

func FileRemoveHandler(c *Context) error {
	db := model.DB()
	accountID := auth.Default(c).GetAccountID()
	fileID := c.FormValue("file_id")
	var file model.File
	if err := db.Where("account_id = ?", accountID).Find(&file, fileID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Delete(&file).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, file)
}

func GetFileGoogleStorage(c *Context) error {
	imagePath := c.QueryParam("imagePath")
	ctx := context.Background()
	image, size, err := lib.GetGoolgeStorage(ctx, "triple-t", imagePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set(echo.HeaderContentLength, size)
	c.Response().WriteHeader(http.StatusOK)
	c.Response().Write(image)
	c.Response().Flush()
	return nil
}

func UploadFileGoogleStorage(c *Context) error {
	ctx := context.Background()
	code := c.FormValue("file")
	if err := lib.UploadGoolgeStorage(ctx, code, "/images/"); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{})
}
