package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func TimeSlotCreateHandler(c *Context) error {
	// a := auth.Default(c)
	id := c.Param("employee_id")
	employeeServices := []*model.EmployeeService{}
	model.DB().Preload("Service").Where("employee_id = ?", id).Find(&employeeServices)
	var DayWeeks [7]string
	DayWeeks[0] = "Sunday"
	DayWeeks[1] = "Monday"
	DayWeeks[2] = "Tuesday"
	DayWeeks[3] = "Wednesday"
	DayWeeks[4] = "Thursday"
	DayWeeks[5] = "Friday"
	DayWeeks[6] = "Saturday"
	return c.Render(http.StatusOK, "time-slot-form", echo.Map{
		"method":           "POST",
		"title":            "employee",
		"employeeID":       id,
		"EmployeeServices": employeeServices,
		"DayWeeks":         DayWeeks,
	})
}

func TimeSlotPostHandler(c *Context) error {
	timeSlotsForm := c.FormValue("timeSlots")
	var timeSlots []*model.TimeSlot
	AccountID := auth.Default(c).GetAccountID()
	err := json.Unmarshal([]byte(timeSlotsForm), &timeSlots)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	tx := model.DB().Begin()

	for _, timeSlot := range timeSlots {
		timeSlot.AccountID = AccountID
		if err = tx.Create(&timeSlot).Error; err != nil {
			tx.Rollback()
		}
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = tx.Commit().Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	id := c.Param("employee_id")

	redirect := fmt.Sprintf("/admin/employee_service/%s", id)

	return c.JSON(http.StatusCreated, redirect)
}

func TimeSlotUpdateViewHandler(c *Context) error {
	id := c.Param("id")
	a := auth.Default(c)

	timeSlot, err := model.GetTimeSlotDetail(id, a.GetAccountID())

	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}
	return c.Render(http.StatusOK, "time-slot-form", echo.Map{
		"method": "PUT",
		"title":  "employee",
		"detail": timeSlot,
	})
}

func TimeSlotUpdateHandler(c *Context) error {
	id := c.Param("id")
	a := auth.Default(c)
	timeSlot, err := model.GetTimeSlotDetail(id, a.GetAccountID())

	if err != nil {
		return c.Render(http.StatusNotFound, "404-page", echo.Map{})
	}

	if err := c.Bind(&timeSlot); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := timeSlot.UpdateTimeSlot(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/employee_service/%d", timeSlot.EmployeeService.EmployeeID)
	return c.JSON(http.StatusOK, redirect)
}

func TimeSlotDeleteHandler(c *Context) error {
	id := c.Param("id")

	provi, err := model.RemoveEmployeeService(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, provi)
}
