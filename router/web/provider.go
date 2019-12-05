package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/module/auth"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func ProviderListHandler(c *Context) error {
	provider := []*model.Provider{}
	a := auth.Default(c)
	db := model.DB()
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	filterProvider := db.Model(&provider).Where("account_id = ?", a.GetAccountID()).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterProvider.Limit(pagination.Record).Offset(pagination.Offset).Find(&provider)

	return c.Render(http.StatusOK, "provider-list", echo.Map{
		"title":      "provider",
		"list":       provider,
		"pagination": pagination,
	})
}

func ProviderDetailHandler(c *Context) error {
	id := c.Param("id")
	a := auth.Default(c)

	provider, err := model.GetProviderDetail(id, a.GetAccountID())
	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "provider-detail", echo.Map{
		"title":  "provider",
		"detail": provider,
	})
}

func ProviderCreateHandler(c *Context) error {
	provider := model.Provider{}

	return c.Render(http.StatusOK, "provider-form", echo.Map{
		"title":  "provider",
		"detail": provider,
		"method": "POST",
	})
}

func ProviderPUTHandler(c *Context) error {
	provider := model.Provider{}

	return c.Render(http.StatusOK, "provider-form", echo.Map{
		"title":  "provider",
		"detail": provider,
		"method": "POST",
	})
}

func ProviderPostHandler(c *Context) error {
	provider := model.Provider{}
	a := auth.Default(c)
	file := c.FormValue("file")

	fileUrl, _, err := lib.UploadteImage(file)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	provider.ProvImage = fileUrl
	provider.AccountID = a.GetAccountID()

	if err := c.Bind(&provider); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	fmt.Println(provider.ProvName, "===")

	err = provider.CreateProvider()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/provider/%d", provider.ID)

	return c.JSON(http.StatusCreated, redirect)
}

func ProviderPutHandler(c *Context) error {
	provider := model.Provider{}
	a := auth.Default(c)
	file := c.FormValue("file")

	file, _, err := lib.UploadteImage(file)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	provider.ProvImage = file
	provider.AccountID = a.GetAccountID()

	if err := c.Bind(&provider); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	fmt.Println(provider.ProvName, "===")

	err = provider.CreateProvider()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/provider/%d", provider.ID)

	return c.JSON(http.StatusCreated, redirect)
}

func ProviderEditHandler(c *Context) error {
	id := c.Param("id")
	provider := model.Provider{}
	a := auth.Default(c)

	if err := model.DB().Where("account_id = ?", a.GetAccountID()).Find(&provider, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "provider-form", echo.Map{
		"title":  "provider",
		"method": "PUT",
		"detail": provider,
	})
}

func ProviderDeleteHandler(c *Context) error {
	id := c.Param("id")

	chatChannel, err := model.RemoveProvider(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, chatChannel)
}

func ProviderAddServiceHandler(c *Context) error {
	id := c.Param("id")
	a := auth.Default(c)
	provider := model.Provider{}
	if err := model.DB().Where("account_id = ?", a.GetAccountID()).Find(&provider, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	services, err := model.GetServiceList(a.GetAccountID())

	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	return c.Render(http.StatusOK, "provider-service-form", echo.Map{"method": "PUT",
		"title":    "provider",
		"services": services,
		"provider": provider,
	})
}

func ProviderSerciveListHandler(c *Context) error {
	id := c.Param("prov_id")
	a := auth.Default(c)

	provider, err := model.GetProviderServiceTimeSlotList(id, a.GetAccountID())

	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "provider-service-detail", echo.Map{
		"title":  "provider",
		"detail": provider,
	})
}

func ProviderAddServicePostHandler(c *Context) error {
	a := auth.Default(c)
	id := c.Param("id")
	prov, err := model.GetProviderDetail(id, a.GetAccountID())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	provs := model.ProviderService{}
	if err := c.Bind(&provs); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	provs.ProviderID = prov.ID
	provs.ID = 0
	if err := model.DB().Create(&provs).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	redirect := fmt.Sprintf("/admin/provider/%s", id)

	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
	})

}

func ProviderAddBookingHandler(c *Context) error {
	id := c.Param("Prov_id")
	a := auth.Default(c)
	provider := model.Provider{}

	chatChannels, err := model.GetChatChannelList(a.GetAccountID())
	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	if err := model.DB().Where("account_id = ?", a.GetAccountID()).Find(&provider, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	return c.Render(http.StatusOK, "provider-list", echo.Map{
		"title":        "provider",
		"chatChannels": chatChannels,
	})
}
