package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

// comment by sothon
// LineLoginRespose is instacne respose line json
type LineLoginRespose struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

// DBHandler is hanler for connent data base.
type DBHandler struct {
	DB *gorm.DB
}

// Customer struct
type Customer struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	User      User
}

// ResposeMessage is property of table message
type ResposeMessage struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Input string `json:"input"`
	Reply string `json:"reply"`
}

// MessageLog soucre user call with line account
type MessageLog struct {
	gorm.Model
	User           User           `gorm:"foreignkey:ID"`
	Message        string         `json:"message"`
	ResposeMessage ResposeMessage `gorm:"foreignkey:ID"`
}

// Source social connet platform
type Source struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Account string `json:"account"`
	Link    string `json:"link"` // webhook
}

// User prepeo use app
type User struct {
	ID       uint   `gorm:"primaty_key" json:"id"`
	Email    string `gorm:"type:varchar(100);unique_index"`
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Birthday string `json:"birth_day"`
	Source   Source `gorm:"foreignkey:ID"`
}

// ServiceInfo handle set service
type ServiceInfo struct {
	IPAddress   string
	Port        string
	ServiceName string
}

// ServicesInfo array
type ServicesInfo []ServiceInfo

// HTTPCallBackHanlder call back webhook
type HTTPCallBackHanlder struct {
	Bot          *linebot.Client
	ServicesInfo *ServicesInfo
	DB           *gorm.DB
}

// Promotion is property of table promotion
type Promotion struct {
	ID       uint   `gorm:"primaty_key" json:"id"`
	Name     string `json:"name"`
	Content  string `json:"connect"`
	Discount int    `json:"discount"`
}

// Setting is property of table setting
type Setting struct {
	ID    uint   `gorm:"primaty_key" json:"id"`
	Key   string `json:"key"`
	Value string `json:""`
}

var channel_secret string = "aadbac180706ab87496224bd4f1a0a74"
var channel_accsss_token string = "UoktnlBJWmbhpjFV+q9s8Cka1DstLx4hkn29lCjYh84ucCkwrNA2gyxCFeLw4AsTCskLD3n33aBXwaKsLuO7a2oqHPlgd0gCOVFEhw0Th43cUZEO9Dc5mKXv5XpgNTws3+JpR9kJbGj2rqep6ny0hAdB04t89/1O/w1cDnyilFU="

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

// NewServiceHTTPHandler provide the inititail set up service path to handle request
func NewServiceHTTPHandler(e *echo.Echo, linebot *linebot.Client, servicesInfo *ServicesInfo) {

	e.GET("/ping", func(c echo.Context) error {

		return c.String(200, "Line boi Service : We are good thank you for asking us.")
	})

	h := DBHandler{}

	h.Initialize()

	hanlders := &HTTPCallBackHanlder{Bot: linebot, ServicesInfo: servicesInfo}
	hanlders.Initialize()

	e.GET("/auth", h.AuthLine)
	e.POST("/callback", hanlders.Callback)
	e.POST("/richMenu", hanlders.CreateRichMenu)
	e.POST("/richMenu/:id", hanlders.UploadRichMenu)
	e.POST("/richMenu/:id", hanlders.UploadRichMenu)
	e.GET("/customers", h.GetAllCustomer)
	e.POST("/customers", h.SaveCustomer)
	e.GET("/customers/:id", h.GetCustomer)
	e.PUT("/customers/:id", h.UpdateCustomer)
	e.DELETE("/customers/:id", h.DeleteCustomer)
	e.GET("/messages", h.getAllMessage)
	e.GET("/register", h.getAllMessage)
	e.POST("/messages", h.CreateMessage)
}

func (h *DBHandler) Initialize() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres-dev dbname=dev password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Customer{})
	db.AutoMigrate(&ResposeMessage{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&LineLoginRespose{})

	h.DB = db
}

func (hanlders *HTTPCallBackHanlder) Initialize() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres-dev dbname=dev password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Customer{})
	db.AutoMigrate(&ResposeMessage{})
	db.AutoMigrate(&User{})

	hanlders.DB = db
}

func (handler *HTTPCallBackHanlder) CreateRichMenu(c echo.Context) error {
	url := "https://api.line.me/v2/bot/richmenu"

	payload := strings.NewReader("{\n  \"size\": {\n    \"width\": 2500,\n    \"height\": 1686\n  },\n  \"selected\": true,\n  \"name\": \"Rich Menu 1\",\n  \"chatBarText\": \"menu\",\n  \"areas\": [\n    {\n      \"bounds\": {\n        \"x\": 0,\n        \"y\": 29,\n        \"width\": 891,\n        \"height\": 821\n      },\n      \"action\": {\n        \"type\": \"message\",\n        \"text\": \"price\"\n      }\n    },\n    {\n      \"bounds\": {\n        \"x\": 895,\n        \"y\": 25,\n        \"width\": 788,\n        \"height\": 829\n      },\n      \"action\": {\n        \"type\": \"message\",\n        \"text\": \"promotion\"\n      }\n    },\n    {\n      \"bounds\": {\n        \"x\": 1691,\n        \"y\": 21,\n        \"width\": 788,\n        \"height\": 829\n      },\n      \"action\": {\n        \"type\": \"message\",\n        \"text\": \"localtion\"\n      }\n    },\n    {\n      \"bounds\": {\n        \"x\": 4,\n        \"y\": 854,\n        \"width\": 883,\n        \"height\": 813\n      },\n      \"action\": {\n        \"type\": \"message\",\n        \"text\": \"Booking\"\n      }\n    },\n    {\n      \"bounds\": {\n        \"x\": 899,\n        \"y\": 862,\n        \"width\": 797,\n        \"height\": 824\n      },\n      \"action\": {\n        \"type\": \"message\",\n        \"text\": \"comment\"\n      }\n    },\n    {\n      \"bounds\": {\n        \"x\": 1708,\n        \"y\": 870,\n        \"width\": 792,\n        \"height\": 816\n      },\n      \"action\": {\n        \"type\": \"message\",\n        \"text\": \"Voucher\"\n      }\n    }\n  ]\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer UoktnlBJWmbhpjFV+q9s8Cka1DstLx4hkn29lCjYh84ucCkwrNA2gyxCFeLw4AsTCskLD3n33aBXwaKsLuO7a2oqHPlgd0gCOVFEhw0Th43cUZEO9Dc5mKXv5XpgNTws3+JpR9kJbGj2rqep6ny0hAdB04t89/1O/w1cDnyilFU=")
	req.Header.Add("User-Agent", "PostmanRuntime/7.15.2")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "3c3c513b-8017-4634-98a9-03140a9b6564,4fa2f3d7-a661-45b8-88bc-b9de3d7bac77")
	req.Header.Add("Host", "api.line.me")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "1335")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	return c.JSON(200, string(body))
}

func (handler *HTTPCallBackHanlder) UploadRichMenu(c echo.Context) error {

	handler.Bot.UploadRichMenuImage("richmenu-7ef5fc568f5546227903a3064466ff2c", "linerichmenu.jpg").Do()

	return c.JSON(200, "state")
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
				fmt.Println(event.Message)
				messageMedal := ResposeMessage{}
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
		} else if event.Type == linebot.EventTypeFollow {
			var contents []linebot.FlexComponent
			text := linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   "taey line bookking plaform",
				Weight: "bold",
				Size:   linebot.FlexTextSizeTypeXl,
				Action: linebot.NewURIAction("register", "https://15e330d8.ngrok.io/register"),
			}
			contents = append(contents, &text)
			// Make Hero
			hero := linebot.ImageComponent{
				Type:        linebot.FlexComponentTypeImage,
				URL:         "https://scontent.fbkk2-7.fna.fbcdn.net/v/t1.0-9/55771768_3311003885591805_86103752003551232_o.jpg?_nc_cat=109&_nc_eui2=AeGPFqTgk7ynFe18QHmR-69H6MogRu5OFJXtXwbMnKDQa2IZeLa57IEayXcXzhyzKDfBKx_tYZevLlEoaJ_bJn6Fl9hCv6mhlWYOOV3ltGoR9Q&_nc_oc=AQkpFLS6szBuMWyOhKz-Ope9I4YkWTFea1DFHE9oNPodtflCUt53bb_kjVd7SVx236w&_nc_ht=scontent.fbkk2-7.fna&oh=62d415b199aaa244c8bea5b9e60dd44b&oe=5DD5122F",
				Size:        "full",
				AspectRatio: linebot.FlexImageAspectRatioType1to1,
				AspectMode:  linebot.FlexImageAspectModeTypeCover,
				Action:      linebot.NewURIAction("register", "https://15e330d8.ngrok.io/register"),
			}
			// Make Body
			body := linebot.BoxComponent{
				Type:     linebot.FlexComponentTypeBox,
				Layout:   linebot.FlexBoxLayoutTypeVertical,
				Contents: contents,
			}
			// Build Container
			bubble := linebot.BubbleContainer{
				Type: linebot.FlexContainerTypeBubble,
				Hero: &hero,
				Body: &body,
			}
			// New Flex Message
			flexMessage := linebot.NewFlexMessage("ขอบคุณที่มาเป็นเพื่อนกันนะ", &bubble)
			if _, err = handler.Bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
				log.Print(err)
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

func (h *DBHandler) GetAllCustomer(c echo.Context) error {
	customers := []Customer{}

	h.DB.Find(&customers)
	fmt.Println("test")
	return c.JSON(http.StatusOK, customers)
}

func (h *DBHandler) GetCustomer(c echo.Context) error {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, customer)
}

func (h *DBHandler) SaveCustomer(c echo.Context) error {
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

func (h *DBHandler) UpdateCustomer(c echo.Context) error {
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

func (h *DBHandler) DeleteCustomer(c echo.Context) error {
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

func (h *DBHandler) getAllMessage(c echo.Context) error {
	message := []ResposeMessage{}

	h.DB.Find(&message)
	fmt.Println("test")
	return c.JSON(http.StatusOK, message)
}

func (h *DBHandler) CreateMessage(c echo.Context) error {
	message := ResposeMessage{}
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

func (h *DBHandler) AuthLine(c echo.Context) error {

	url := "https://api.line.me/oauth2/v2.1/token"
	code := fmt.Sprintf("grant_type=authorization_code&code=%s", c.QueryParam("code"))
	redi := "&redirect_uri=https%3A%2F%2F15e330d8.ngrok.io%2Fauth&client_id=1610710377&client_secret=28a769dc4ec9fc1fdd345ff051827e71"

	payload := strings.NewReader(code + redi)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "PostmanRuntime/7.15.2")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "8b75ec60-fd16-4a9c-89ed-8d8c097a81c0,24b23656-1e23-4227-b90a-b7b2a169276a")
	req.Header.Add("Host", "api.line.me")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "175")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusOK, "")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	lineRespose := LineLoginRespose{}
	json.Unmarshal(body, &lineRespose)
	claims := jwt.MapClaims{}
	if lineRespose.IDToken == "" {
		fmt.Println("====", string(body))
	}
	token, err := jwt.ParseWithClaims(lineRespose.IDToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	fmt.Println(token)
	return c.JSON(http.StatusOK, token)
}
