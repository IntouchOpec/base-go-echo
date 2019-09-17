package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// LIFFloginHandler
func LIFFRegisterHandler(c echo.Context) error {
	lineID := c.Param("lineID")
	chatChannel := model.ChatChannel{}
	model.DB().Where("lineID = ?", lineID).Find(&chatChannel)
	fmt.Println(chatChannel)
	err := c.Render(http.StatusOK, "register", map[string]interface{}{})
	return err
}
