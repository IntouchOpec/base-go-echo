package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

func PackageListHandler(c *Context) error {
	packageModels := []*model.Package{}
	a := auth.Default(c)
	db := model.DB()
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	filterPackage := db.Model(&packageModels).Where("account_id = ?", a.GetAccountID()).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterPackage.Limit(pagination.Record).Offset(pagination.Offset).Find(&packageModels)

	return c.Render(http.StatusOK, "package-list", echo.Map{
		"title":      "package",
		"list":       packageModels,
		"pagination": pagination,
	})
}

func PackageIsactiveHandler(c *Context) error {
	packageModel := model.Package{}
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	b, err := strconv.ParseBool(c.FormValue("pac_is_acive"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Where("account_id = ?", accID).Find(&packageModel, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	packageModel.PacIsActive = b
	if err := db.Save(&packageModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": packageModel,
	})
}

func PackageDetailHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	pack := model.Package{}
	db := model.DB()
	err := db.Preload("ServiceItems").Where("account_id = ?", accID).Find(&pack, id).Error
	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	list := []model.ServiceItem{}
	for _, serI := range pack.ServiceItems {
		item := model.ServiceItem{}
		db.Preload("Service").Find(&item, serI.ID)
		list = append(list, item)
	}
	return c.Render(http.StatusOK, "package-detail", echo.Map{
		"title":  "package",
		"detail": pack,
		"list":   list,
	})
}

func PackageCreateHandler(c *Context) error {
	packageModel := model.Package{}
	return c.Render(http.StatusOK, "package-form", echo.Map{
		"title":  "package",
		"detail": packageModel,
		"method": "POST",
	})
}

func PackageEditViewHandler(c *Context) error {
	packageModel := model.Package{}
	id := c.Param("id")
	a := auth.Default(c).GetAccountID()

	db := model.DB()
	if err := db.Where("account_id = ?", a).Find(&packageModel, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	return c.Render(http.StatusOK, "package-form", echo.Map{
		"title":  "package",
		"detail": packageModel,
		"method": "PUT",
	})
}

func PackagePutHandler(c *Context) error {
	packageModel := model.Package{}
	packageID := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	if err := db.Where("account_id = ?", accID).Find(&packageModel, packageID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Bind(&packageModel); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Save(&packageModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"detail":   packageModel,
		"redirect": fmt.Sprintf("/admin/package/%d", packageModel.ID),
	})
}

func PackagePostHandler(c *Context) error {
	packageModel := model.Package{}
	accID := auth.Default(c).GetAccountID()
	fmt.Println(accID)
	file := c.FormValue("file")
	ctx := context.Background()
	db := model.DB()
	imagePath, err := lib.UploadGoolgeStorage(ctx, file, "images/package/")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Bind(&packageModel); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	packageModel.PacImage = imagePath
	packageModel.AccountID = accID
	if err := db.Create(&packageModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/package/%d", packageModel.ID)
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
		"data":     packageModel,
	})
}

func PackageEditHandler(c *Context) error {
	id := c.Param("id")
	packageModel := model.Package{}
	a := auth.Default(c)
	db := model.DB()
	if err := db.Where("account_id = ?", a.GetAccountID()).Preload("ChatChannel").Find(&packageModel, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "package-form", echo.Map{
		"title":  "package",
		"method": "PUT",
		"detail": packageModel,
	})
}

func PackageDeleteHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	packageModel := model.Package{}
	db := model.DB()
	if err := db.Where("account_id = ?", accID).Find(&packageModel, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Delete(&packageModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"detail": packageModel,
	})
}

func PackageDeleteImageHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	packageModel := model.Package{}
	db := model.DB()
	if err := db.Where("account_id = ?", accID).Find(&packageModel, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	ctx := context.Background()
	if _, err := lib.RemoveFileGoolgeStorage(ctx, "triple-t", packageModel.PacImage); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	packageModel.PacImage = ""
	if err := db.Save(&packageModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"detail": packageModel,
	})
}

func PackageServiceCreateHandler(c *Context) error {
	accID := auth.Default(c).GetAccountID()
	var serviceItems []*model.ServiceItem
	var packageModel model.Package
	id := c.Param("id")
	db := model.DB()
	var serviceItemIDs []uint
	if err := db.Preload("ServiceItems").Where("account_id = ?", accID).Find(&packageModel, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	serviceItemIDs = append(serviceItemIDs, 0)
	for _, ServiceItem := range packageModel.ServiceItems {
		serviceItemIDs = append(serviceItemIDs, ServiceItem.ID)
	}

	if err := db.Preload("Service").Where("account_id = ? and id not in (?)", accID, serviceItemIDs).Find(&serviceItems).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "package-service-form", echo.Map{
		"method":       "POST",
		"data":         packageModel,
		"serviceItems": serviceItems,
		"title":        "package",
	})
}
func PackageServiceCreatePostHandler(c *Context) error {
	serviceID := c.FormValue("service_id")
	packageID := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	var serviceItem model.ServiceItem
	var packageModel model.Package
	db := model.DB()
	if err := db.Where("account_id = ?", accID).Find(&serviceItem, serviceID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Where("account_id = ?", accID).Find(&packageModel, packageID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Model(&packageModel).Association("ServiceItems").Append(&serviceItem).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Preload("ServiceItems").Where("account_id = ?", accID).Find(&packageModel, packageID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	var sumTimeUsed time.Duration
	if packageModel.ServiceItems != nil {
		for _, service := range packageModel.ServiceItems {
			sumTimeUsed += service.SSTime
		}
	}

	packageModel.PacTime = sumTimeUsed

	if err := db.Save(&packageModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/package/%d", packageModel.ID)

	return c.JSON(http.StatusCreated, echo.Map{
		"data":     packageModel,
		"redirect": redirect,
	})
}
func PackageServiceDeleteHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	serviceID := c.Param("service_id")
	var serviceItem model.ServiceItem
	var packageModel model.Package
	db := model.DB()
	if err := db.Where("account_id = ?", accID).Find(&serviceItem, serviceID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Where("account_id = ?", accID).Find(&packageModel, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Model(&packageModel).Association("ServiceItems").Delete(&serviceItem).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": packageModel,
	})
}
