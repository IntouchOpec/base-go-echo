package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

// timeSlots: [{"time_day":6,"employee_id":5,"time_start":"0001-01-01T11:00:04+06:45","time_end":"0001-01-01T11:00:04+06:45"}]
type timeSlot struct {
	TimeDay    int    `json:"time_day"`
	EmployeeID uint   `json:"employee_id"`
	TimeStart  string `json:"time_start"`
	TimeEnd    string `json:"time_end"`
}

func TimeSlotPostHandler(c *Context) error {
	timeSlotsForm := c.FormValue("timeSlots")
	var timeSlots []timeSlot
	AccountID := auth.Default(c).GetAccountID()
	err := json.Unmarshal([]byte(timeSlotsForm), &timeSlots)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	tx := model.DB().Begin()

	for _, timeSlot := range timeSlots {
		var timeSlotM model.TimeSlot
		start, _ := time.Parse("2006-01-02 15:04:01", fmt.Sprintf("0001-01-01 %s:01", timeSlot.TimeStart))
		end, _ := time.Parse("2006-01-02 15:04:01", fmt.Sprintf("0001-01-01 %s:01", timeSlot.TimeEnd))
		timeSlotM.TimeDay = timeSlot.TimeDay
		timeSlotM.EmployeeID = timeSlot.EmployeeID
		timeSlotM.TimeStart = start.Add(-(7 * time.Hour))
		timeSlotM.TimeEnd = end.Add(-(7 * time.Hour))
		timeSlotM.AccountID = AccountID
		if err = tx.Create(&timeSlotM).Error; err != nil {
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
	redirect := fmt.Sprintf("/admin/employee_service/%d", timeSlot.EmployeeID)
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
