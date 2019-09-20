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
	if err := model.DB().Preload("Settings", "name in (?)", "host_api").Where("line_ID = ?", lineID).Find(&chatChannel).Error; err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	custo := model.Customer{}
	if err := model.DB().FirstOrCreate(&custo).Error; err != nil {
		return err
	}
	APIRegister := fmt.Sprintf("https://%s/register/%s", chatChannel.Settings[0].Value, lineID)
	err := c.Render(http.StatusOK, "register", map[string]interface{}{
		"web": APIRegister,
	})
	return err
}

type LineReqRegister struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

func LIIFRegisterSaveCustomer(c echo.Context) error {
	lineID := c.Param("lineID")
	chatChannel := model.ChatChannel{}
	req := LineReqRegister{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := model.DB().Where("line_ID = ?", lineID).Find(&chatChannel).Error; err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	fmt.Println(req.DisplayName)
	custo := model.Customer{LineID: req.UserID, ChatChannelID: chatChannel.ID}
	// pictureURL string, displayName string, email string, phoneNumber string
	custo.UpdateCustomerByAtt(req.PictureURL, req.DisplayName, req.Email, req.Phone)

	return c.JSON(http.StatusOK, custo)
}
