package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/module/auth"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

func EmployeeListHandler(c *Context) error {
	employee := []*model.Employee{}
	a := auth.Default(c)
	db := model.DB()
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	filterEmployee := db.Model(&employee).Where("account_id = ?", a.GetAccountID()).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterEmployee.Limit(pagination.Record).Offset(pagination.Offset).Find(&employee)

	return c.Render(http.StatusOK, "employee-list", echo.Map{
		"title":      "employee",
		"list":       employee,
		"pagination": pagination,
	})
}

func EmployeeDetailHandler(c *Context) error {
	id := c.Param("id")
	a := auth.Default(c)

	employee, err := model.GetEmployeeDetail(id, a.GetAccountID())
	bookings := []model.Booking{}
	db := model.DB()
	db.Find(&bookings)
	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "employee-detail", echo.Map{
		"title":  "employee",
		"detail": employee,
	})
}

func EmployeeCreateHandler(c *Context) error {
	employee := model.Employee{}
	chatChannels := []model.ChatChannel{}
	a := auth.Default(c).GetAccountID()

	db := model.DB()

	db.Where("account_id = ?", a).Find(&chatChannels)

	return c.Render(http.StatusOK, "employee-form", echo.Map{
		"title":        "employee",
		"chatChannels": chatChannels,
		"detail":       employee,
		"method":       "POST",
	})
}

func EmployeePUTHandler(c *Context) error {
	employee := model.Employee{}

	if err := c.Bind(&employee); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := employee.UpdateEmployee(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"detail":   employee,
		"redirect": fmt.Sprintf("/admin/employee/%d", employee.ID),
	})
}

func EmployeePostHandler(c *Context) error {
	employee := model.Employee{}
	a := auth.Default(c)
	file := c.FormValue("file")
	ctx := context.Background()

	imagePath, err := lib.UploadGoolgeStorage(ctx, file, "images/employee/")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	employee.ProvImage = imagePath
	employee.AccountID = a.GetAccountID()

	if err := c.Bind(&employee); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = employee.CreateEmployee()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/employee/%d", employee.ID)

	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
	})
}

func EmployeePutHandler(c *Context) error {
	employee := model.Employee{}
	a := auth.Default(c)
	file := c.FormValue("file")
	ctx := context.Background()
	imagePath, err := lib.UploadGoolgeStorage(ctx, file, "images/employee/")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	employee.ProvImage = imagePath
	employee.AccountID = a.GetAccountID()

	if err := c.Bind(&employee); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = employee.CreateEmployee()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/employee/%d", employee.ID)

	return c.JSON(http.StatusCreated, redirect)
}

func EmployeeEditHandler(c *Context) error {
	id := c.Param("id")
	employee := model.Employee{}
	a := auth.Default(c)

	chatChannels := []model.ChatChannel{}

	db := model.DB()

	db.Where("account_id = ?", a).Find(&chatChannels)
	if err := db.Where("account_id = ?", a.GetAccountID()).Preload("ChatChannel").Find(&employee, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "employee-form", echo.Map{
		"title":        "employee",
		"chatChannels": chatChannels,
		"method":       "PUT",
		"detail":       employee,
	})
}

func EmployeeDeleteHandler(c *Context) error {
	id := c.Param("id")

	chatChannel, err := model.RemoveEmployee(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, chatChannel)
}

func EmployeeAddServiceHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	employee := model.Employee{}

	if err := model.DB().Where("account_id = ?", accID).Find(&employee, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	services, err := model.GetServiceList(accID)

	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	return c.Render(http.StatusOK, "employee-service-form", echo.Map{
		"method":   "POST",
		"title":    "employee",
		"services": services,
		"employee": employee,
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

func EmployeeSerciveListHandler(c *Context) error {
	id := c.Param("prov_id")
	accID := auth.Default(c).GetAccountID()
	employee, err := model.GetEmployeeServiceTimeSlotList(id, accID)

	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "employee-service-detail", echo.Map{
		"title":  "employee",
		"detail": employee,
	})
}

func EmployeeAddServicePostHandler(c *Context) error {
	var provService model.EmployeeService
	db := model.DB()
	accID := auth.Default(c).GetAccountID()
	price, err := strconv.ParseFloat(c.FormValue("price"), 10)
	serviceID, err := strconv.ParseUint(c.FormValue("service_id"), 10, 32)
	employeeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	provService.PSPrice = price
	provService.ServiceID = uint(serviceID)
	provService.ID = 0
	provService.AccountID = accID
	provService.EmployeeID = uint(employeeID)
	if err := db.Create(&provService).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/employee/%d", employeeID)
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
		"provs":    provService,
	})
}

func EmployeeAddBookingHandler(c *Context) error {
	id := c.Param("Prov_id")
	a := auth.Default(c)
	employee := model.Employee{}
	db := model.DB()
	var chatChannels []model.ChatChannel
	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&chatChannels).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	if err := db.Where("account_id = ?", a.GetAccountID()).Find(&employee, id).Error; err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	return c.Render(http.StatusOK, "employee-list", echo.Map{
		"title":        "employee",
		"chatChannels": chatChannels,
	})
}

func EmployeeDeleteImageHandler(c *Context) error {
	id := c.Param("Prov_id")
	a := auth.Default(c)
	employee := model.Employee{}

	if err := model.DB().Where("account_id = ?", a.GetAccountID()).Find(&employee, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	ctx := context.Background()
	if _, err := lib.RemoveFileGoolgeStorage(ctx, "triple-t", employee.ProvImage); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := employee.RemoveImage(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": employee,
	})
}
