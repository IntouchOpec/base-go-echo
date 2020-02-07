package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/lib"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

type AccBooking struct {
	ID      model.AccBookingType
	Text    string
	Checked bool
}

func SettingHandler(c *Context) error {
	var acc model.Account
	accID := auth.Default(c).GetAccountID()
	db := model.DB()

	if err := db.Preload("Settings").Find(&acc, accID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	accTypePayments := []model.AccTypePayment{
		model.AccTypePaymentBooking,
		model.AccTypePaymentNow,
	}
	accBookingConfirmTpyes := []model.AccTransactionConfirmType{model.AccTransactionMan, model.AccTransactionAuto}
	bookingTypes := []AccBooking{
		AccBooking{ID: model.AccBookingByTimeSlot, Text: "TimeSlot"},
		AccBooking{ID: model.AccBookingByItem, Text: "Item"},
		AccBooking{ID: model.AccBookingByNow, Text: "Now"},
	}
	var ints []int
	err := json.Unmarshal([]byte(fmt.Sprintf("[%s]", acc.AccBookingType)), &ints)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", ints)
	for _, char := range ints {
		bookingTypes[char].Checked = true
	}

	return c.Render(http.StatusOK, "setting", echo.Map{
		"detail":                 acc,
		"accTypePayments":        accTypePayments,
		"title":                  "setting",
		"method":                 "PUT",
		"accBookingConfirmTpyes": accBookingConfirmTpyes,
		"bookingTypes":           bookingTypes,
	})
}

func SettingAccountPutHandler(c *Context) error {
	accID := auth.Default(c).GetAccountID()
	var acc model.Account
	db := model.DB()
	if err := db.Find(&acc, accID).Error; err != nil {
		fmt.Println("err")
		return c.JSON(http.StatusBadRequest, err)
	}
	jsonPath := c.FormValue("file")
	if jsonPath != "" {
		ctx := context.Background()
		jsonPath, err := lib.UploadGoolgeStorage(ctx, jsonPath, "json/AuthJSONGoogle/")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		acc.AccAuthJSONFilePath = jsonPath
	}
	if err := c.Bind(&acc); err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println("====1")
	if err := db.Save(&acc).Error; err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, acc)
}

func RemoaveAuthJSONFile(c *Context) error {
	accID := auth.Default(c).GetAccountID()
	var acc model.Account
	db := model.DB()
	if err := db.Find(&acc, accID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	if _, err := lib.RemoveFileGoolgeStorage(ctx, "triple-t", acc.AccAuthJSONFilePath); err != nil {
		// err := lib.DeleteFile(acc.AccAuthJSONFilePath)
		// if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	acc.AccAuthJSONFilePath = ""
	if err := db.Save(&acc).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, acc)
}

func SettingCreateViewHandler(c *Context) error {
	var setting model.Setting
	return c.Render(http.StatusOK, "setting-form", echo.Map{
		"detail": setting,
		"method": "POST",
	})
}

func SettingPostViewHandler(c *Context) error {
	var setting model.Setting
	if err := c.Bind(&setting); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := model.DB().Create(&setting).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"detail":   setting,
		"redirect": c.QueryParam("redirect"),
	})
}

func SettingEditViewHandler(c *Context) error {
	var setting model.Setting
	id := c.Param("id")
	db := model.DB()
	if err := db.Find(&setting, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.Render(http.StatusOK, "setting-form", echo.Map{
		"detail": setting,
		"method": "PUT",
	})
}

func SettingPutHandler(c *Context) error {
	var setting model.Setting
	id := c.Param("id")
	db := model.DB()

	if err := db.Find(&setting, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := c.Bind(&setting); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Save(&setting).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data":     setting,
		"redirect": "/admin/",
	})
}
