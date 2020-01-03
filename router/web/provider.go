package web

import (
	"fmt"
	"net/http"
	"strconv"

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
	bookings := []model.Booking{}
	db := model.DB()
	db.Find(&bookings)
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

	if err := c.Bind(&provider); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := provider.UpdateProvider(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"detail":   provider,
		"redirect": fmt.Sprintf("/admin/provider/%d", provider.ID),
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

	err = provider.CreateProvider()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/provider/%d", provider.ID)

	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
	})
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
	accID := auth.Default(c).GetAccountID()
	provider := model.Provider{}

	if err := model.DB().Where("account_id = ?", accID).Find(&provider, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	services, err := model.GetServiceList(accID)

	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	return c.Render(http.StatusOK, "provider-service-form", echo.Map{
		"method":   "POST",
		"title":    "provider",
		"services": services,
		"provider": provider,
	})
}
func weeDayString(day int) string {
	var weeDay []string
	weeDay[0] = "Mon"
	weeDay[1] = "Tue"
	weeDay[2] = "Wed"
	weeDay[3] = "Thu"
	weeDay[4] = "Fri"
	weeDay[5] = "Sat"
	weeDay[6] = "Sun"
	return weeDay[day]
}

func ProviderSerciveListHandler(c *Context) error {
	id := c.Param("prov_id")
	accID := auth.Default(c).GetAccountID()
	provider, err := model.GetProviderServiceTimeSlotList(id, accID)

	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "provider-service-detail", echo.Map{
		"title":  "provider",
		"detail": provider,
	})
}

func ProviderAddServicePostHandler(c *Context) error {
	var provService model.ProviderService
	db := model.DB()
	price, err := strconv.ParseFloat(c.FormValue("price"), 10)
	serviceID, err := strconv.ParseUint(c.FormValue("service_id"), 10, 32)
	providerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	provService.PSPrice = price
	provService.ServiceID = uint(serviceID)
	provService.ID = 0
	provService.ProviderID = uint(providerID)
	if err := db.Create(&provService).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/provider/%d", providerID)
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
		"provs":    provService,
	})
}

func ProviderAddBookingHandler(c *Context) error {
	id := c.Param("Prov_id")
	a := auth.Default(c)
	provider := model.Provider{}
	db := model.DB()
	var chatChannels []model.ChatChannel
	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannels).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&provider, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	return c.Render(http.StatusOK, "provider-list", echo.Map{
		"title":        "provider",
		"chatChannels": chatChannels,
	})
}

func ProviderDeleteImageHandler(c *Context) error {
	id := c.Param("Prov_id")
	a := auth.Default(c)
	provider := model.Provider{}

	if err := model.DB().Where("account_id = ?", a.GetAccountID()).Find(&provider, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	if err := lib.DeleteFile(provider.ProvImage); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := provider.RemoveImage(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": provider,
	})
}
