package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

// serviceListHandler
func CustomerTypeListHandler(c *Context) error {
	customerTypes := []*model.CustomerType{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()

	filterCustomerType := db.Model(&customerTypes).Where("account_id = ?", a.User.GetAccountID()).Count(&total)

	pagination := MakePagination(total, page, limit)
	filterCustomerType.Limit(pagination.Record).Offset(pagination.Offset).Find(&customerTypes)

	err := c.Render(http.StatusOK, "service-list", echo.Map{
		"list":       customerTypes,
		"title":      "customer_type",
		"pagination": pagination,
	})
	return err
}

func CustomerTypeEditViewHandler(c *Context) error {
	CustomerType := model.CustomerType{}
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()

	model.DB().Where("account_id = ?", accID).Find(&CustomerType, id)
	return c.Render(http.StatusOK, "service-form", echo.Map{
		"detail": CustomerType,
		"title":  "customer_type",
		"method": "PUT",
	})
}

func CustomerTypeEditPutHandler(c *Context) error {
	customerType := serviceForm{}
	image := c.FormValue("image")
	id := c.Param("id")
	var err error
	if image == "" {
		file := c.FormValue("file")
		image, _, err = lib.UploadteImage(file)
	}

	if err := c.Bind(&customerType); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	accID := auth.Default(c).GetAccountID()
	serviceModel := model.CustomerType{}
	db := model.DB()
	if err := db.Where("account_id = ?", accID).Find(&serviceModel, id).Error; err != nil {
		return err
	}
	if err := db.Save(&serviceModel).Error; err != nil {
		return err
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"data":     serviceModel,
		"redirect": fmt.Sprintf("/admin/service/%d", serviceModel.ID),
	})
}

func CustomerTypeCreateHandler(c *Context) error {
	CustomerType := model.CustomerType{}

	err := c.Render(http.StatusOK, "service-form", echo.Map{
		"detail": CustomerType,
		"title":  "customer_type",
		"method": "POST",
	})
	return err
}

func CustomerTypePostHandler(c *Context) error {
	CustomerType := model.CustomerType{}
	db := model.DB()
	if err := c.Bind(&CustomerType); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Create(&CustomerType).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": fmt.Sprintf("/admin/customer_type"),
	})
}

func CustomerTypeEditHandler(c *Context) error {
	customerType := model.CustomerType{}
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Preload("Account").Preload("CustomerTypeSlots").Preload("ChatChannels").Where("account_id = ? ", a.User.GetAccountID()).Find(&customerType, id)
	err := c.Render(http.StatusOK, "service-form", echo.Map{
		"method": "PUT",
		"detail": customerType,
		"title":  "customer_type",
	})
	return err
}

func CustomerTypeDeleteHandler(c *Context) error {
	id := c.Param("id")
	pro := model.DeleteserviceByID(id)
	err := c.JSON(http.StatusOK, pro)
	return err
}
