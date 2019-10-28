package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/hb-go/json"
	"github.com/labstack/echo"
)

// RegisterCustomerByLine
func RegisterCustomerByLine(c echo.Context) error {
	chatChannelID := c.Param("chatChannelID")

	code := c.QueryParam("code")
	db := model.DB()
	chatChannel := model.ChatChannel{}
	if err := db.Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	enCodeURLLocal := url.QueryEscape("https%3A%2F%2F15e330d8.ngrok.io%2Fauth")
	url := "https://api.line.me/oauth2/v2.1/token"
	payload := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		code, enCodeURLLocal, chatChannel.ChannelID, chatChannel.ChannelSecret)

	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusOK, "")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	// customer := model.Customer{}
	lineRespose := model.LoginRespose{}

	json.Unmarshal(body, &lineRespose)
	lineRespose.SaveLoginRespose()

	claims := jwt.MapClaims{}
	if lineRespose.IDToken == "" {
		return c.JSON(http.StatusBadRequest, "")
	}
	token, err := jwt.ParseWithClaims(lineRespose.IDToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})

	return c.JSON(http.StatusOK, token)
}

func GetCustomerList(c echo.Context) error {
	page := c.QueryParam("page")
	size := c.QueryParam("size")
	chatChannelID := c.Param("chatChannelID")
	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)
	chatChannelIDInt, _ := strconv.Atoi(chatChannelID)

	customers := model.GetCustomerList(pageInt, sizeInt, chatChannelIDInt)
	return c.JSON(http.StatusOK, customers)
}

func GetCustomerDetail(c echo.Context) error {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	customer := model.Customer{}
	customer.GetCustomer(idInt)
	return c.JSON(http.StatusOK, customer)
}

func UpdateCustomer(c echo.Context) error {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	customer := model.Customer{}
	if err := c.Bind(&customer).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	customer.UpdateCustomer(idInt)
	return c.JSON(http.StatusOK, customer)
}
