package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

var channel_secret string = "aadbac180706ab87496224bd4f1a0a74"
var channel_accsss_token string = "QmMegKvZY6n1LJGhxYTRiSCg/F2nuRLrIbafrll5mlPEP9K0ZtOkyuUPJYLN29FjCskLD3n33aBXwaKsLuO7a2oqHPlgd0gCOVFEhw0Th40i5D9Pv0k2StHNGU3AhxtHC+OQP1hzbr9qLRqsrjRpdwdB04t89/1O/w1cDnyilFU="

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
}

// NewServiceHTTPHandler provide the inititail set up service path to handle request
func NewServiceHTTPHandler(e *echo.Echo, linebot *linebot.Client, servicesInfo *ServicesInfo) {

	hanlders := &HTTPCallBackHanlder{Bot: linebot, ServicesInfo: servicesInfo}
	e.GET("/ping", func(c echo.Context) error {

		return c.String(200, "Line boi Service : We are good thank you for asking us.")
	})
	e.POST("/callback", hanlders.Callback)
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
			fmt.Print(event.Type, event.Message)

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				messageFromPing := PingService(message.Text, handler.ServicesInfo, time.Second*5)
				if _, err = handler.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageFromPing)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
	return c.JSON(200, "")
}

var (
	bankCoreIPAddress = "192.168.212.212"
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
