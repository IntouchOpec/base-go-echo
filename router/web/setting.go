package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/lib"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

func SettingHandler(c *Context) error {
	promotions := []*model.Promotion{}
	var acc model.Account
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	if err := db.Find(&acc, accID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := model.DB().Find(&promotions).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err := c.Render(http.StatusOK, "setting", echo.Map{
		"detail": acc,
		"title":  "setting",
		"method": "PUT",
	})
	return err
}

func SettingPostHandler(c *Context) error {
	accID := auth.Default(c).GetAccountID()
	var acc model.Account

	db := model.DB()
	if err := db.Find(&acc, accID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Bind(&acc); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	file := c.FormValue("file")

	if file == "" {
		fileURL, _, err := lib.UploadFile(acc.AccName, "json")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		acc.AccAuthJSONFilePath = fileURL
	}

	if err := db.Save(&acc).Error; err != nil {
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

	err := lib.DeleteFile(acc.AccAuthJSONFilePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	acc.AccAuthJSONFilePath = ""
	if err := db.Save(&acc).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, acc)
}
