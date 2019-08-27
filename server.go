package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

var channel_secret string = "aadbac180706ab87496224bd4f1a0a74"
var channel_accsss_token string = "TWmawWcFmW2nnH9+JfA2l1+m1fHbg/3eDR/TSsHxoOXL2rUB975CKjJWZSHjwGL4CskLD3n33aBXwaKsLuO7a2oqHPlgd0gCOVFEhw0Th43FAvvR305ozcg7GR1NSFknk2F6O8W+o/f2uIenZ/aoWgdB04t89/1O/w1cDnyilFU="

func main() {
	startService()
}

func connectLineBot() *linebot.Client {
	bot, err := linebot.New(
		channel_secret,
		channel_accsss_token,
	)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func startService() {
	e := echo.New()
	bankCoreInfo := NewBankCoreServiceInfo()
	NewServiceHTTPHandler(e, connectLineBot(), bankCoreInfo)
	e.Logger.Fatal(e.Start(":80"))
}

type ServiceInfo struct {
	IPAddress   string
	Port        string
	ServiceName string
}

type ServicesInfo []ServiceInfo

type HTTPCallBackHanlder struct {
	Bot          *linebot.Client
	ServicesInfo *ServicesInfo
	DB           *gorm.DB
}

// NewServiceHTTPHandler provide the inititail set up service path to handle request
func NewServiceHTTPHandler(e *echo.Echo, linebot *linebot.Client, servicesInfo *ServicesInfo) {

	e.GET("/ping", func(c echo.Context) error {

		return c.String(200, "Line boi Service : We are good thank you for asking us.")
	})

	h := CustomerHandler{}

	h.Initialize()

	hanlders := &HTTPCallBackHanlder{Bot: linebot, ServicesInfo: servicesInfo}
	hanlders.Initialize()

	e.POST("/callback", hanlders.Callback)
	e.GET("/customers", h.GetAllCustomer)
	e.POST("/customers", h.SaveCustomer)
	e.GET("/customers/:id", h.GetCustomer)
	e.PUT("/customers/:id", h.UpdateCustomer)
	e.DELETE("/customers/:id", h.DeleteCustomer)
	e.GET("/messages", h.getAllMessage)
	e.POST("/messages", h.CreateMessage)
}

type CustomerHandler struct {
	DB *gorm.DB
}

type Customer struct {
	Id        uint   `gorm:"primary_key" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

type Message struct {
	Id    uint   `gorm:"primary_key" json:"id"`
	Input string `json:"input"`
	Reply string `json:"reply"`
}

type User struct {
	Id       uint   `gorm:"primaty_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *CustomerHandler) Initialize() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres-dev dbname=dev password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Customer{})
	db.AutoMigrate(&Message{})
	db.AutoMigrate(&User{})

	h.DB = db
}

func (hanlders *HTTPCallBackHanlder) Initialize() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres-dev dbname=dev password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Customer{})
	db.AutoMigrate(&Message{})
	db.AutoMigrate(&User{})

	hanlders.DB = db
}

// Callback provides the function to handle request from line
func (handler *HTTPCallBackHanlder) Callback(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	events, err := handler.Bot.ParseRequest(c.Request())
	if err != nil {

		if err == linebot.ErrInvalidSignature {
			c.String(400, linebot.ErrInvalidSignature.Error())
		} else {
			c.String(500, "internal")
		}
	}

	for _, event := range events {

		if event.Type == linebot.EventTypeMessage {

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				messageMedal := Message{}
				content := ""
				if err := handler.DB.Where("Input = ?", message.Text).Find(&messageMedal).Error; err != nil {
					content = "error"
				} else {
					content = messageMedal.Reply
				}
				// messageFromPing := PingService(message.Text, handler.ServicesInfo, time.Second*5)
				if _, err = handler.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
	return c.JSON(200, "")
}

var (
	bankCoreIPAddress = "127.0.0.1"
	bankServiceInfo   = map[string]string{
		"hulk":               "9906",
		"genesis":            "9900",
		"eve":                "9904",
		"ems":                "9901",
		"minio":              "9800",
		"maersk":             "9903",
		"simpleredirectbank": "3001",
		"lawson":             "9902",
		"portainer":          "3002",
		"phpmyadmin":         "3000",
		"I Want Love":        "1234",
	}
)

// PingService provides the function that send the serviceName and services information to match and validate is it online or not.
func PingService(messageServiceName string, servicesInfo *ServicesInfo, timeOut time.Duration) string {
	serviceInfo, err := FindServiceName(messageServiceName, servicesInfo)
	if err != nil {
		return "Sorry, the name did not match to any services in our system."
	}
	if len(serviceInfo.ServiceName) > 0 && len(serviceInfo.Port) > 0 {
		serviceStatus := ping(serviceInfo.ServiceName, serviceInfo.IPAddress, serviceInfo.Port, timeOut)
		message, _ := isServiceOnline(serviceInfo.ServiceName, serviceStatus)
		return message
	}
	return ""
}

func FindServiceName(messageText string, servicesInfo *ServicesInfo) (*ServiceInfo, error) {
	for _, serviceDetail := range *servicesInfo {
		if strings.Contains(strings.ToLower(messageText), strings.ToLower(serviceDetail.ServiceName)) {
			return &serviceDetail, nil
		}
	}
	return nil, errors.New("the name did not match to any services in our system.")
}

// StartPingAllServices provides the function ping to allservice that we send through input.
func StartPingAllServices(servicesInfo *ServicesInfo, timeOut time.Duration) []string {
	var lstServiceDowns []string
	for _, serviceDetail := range *servicesInfo {
		serviceStatus := ping(serviceDetail.ServiceName, serviceDetail.IPAddress, serviceDetail.Port, timeOut)
		if message, isOnline := isServiceOnline(serviceDetail.ServiceName, serviceStatus); !(isOnline) {
			lstServiceDowns = append(lstServiceDowns, message)
		}
	}
	return lstServiceDowns
}

func isServiceOnline(serviceName string, status bool) (string, bool) {
	if status {
		return fmt.Sprintf("%s service is working pretty well.", serviceName), true
	} else {
		return fmt.Sprintf("%s service is down, please contact admin.", serviceName), false
	}

}

func ping(serviceName string, ipAddress string, port string, timeOut time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ipAddress, port), timeOut)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// NewBankCoreServiceInfo provides the all service information that using in BankCore project.
func NewBankCoreServiceInfo() *ServicesInfo {
	bankCoreServices := ServicesInfo{}
	for BankServiceName, BankServicePort := range bankServiceInfo {
		serviceInfo := ServiceInfo{
			ServiceName: BankServiceName,
			IPAddress:   bankCoreIPAddress,
			Port:        BankServicePort}
		bankCoreServices = append(bankCoreServices, serviceInfo)
	}
	return &bankCoreServices
}

func (h *CustomerHandler) GetAllCustomer(c echo.Context) error {
	customers := []Customer{}

	h.DB.Find(&customers)
	fmt.Println("test")
	return c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c echo.Context) error {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	fmt.Println(customer)
	return c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SaveCustomer(c echo.Context) error {
	customer := Customer{}
	fmt.Println("========")

	if err := c.Bind(&customer); err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		fmt.Println(err)

		return c.NoContent(http.StatusInternalServerError)
	}
	fmt.Println(customer)

	return c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c echo.Context) error {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := c.Bind(&customer); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c echo.Context) error {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := h.DB.Delete(&customer).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CustomerHandler) getAllMessage(c echo.Context) error {
	message := []Message{}

	h.DB.Find(&message)
	fmt.Println("test")
	return c.JSON(http.StatusOK, message)
}

func (h *CustomerHandler) CreateMessage(c echo.Context) error {
	message := Message{}
	fmt.Println("========")

	if err := c.Bind(&message); err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.DB.Save(&message).Error; err != nil {
		fmt.Println(err)

		return c.NoContent(http.StatusInternalServerError)
	}
	fmt.Println(message)

	return c.JSON(http.StatusOK, message)
}
