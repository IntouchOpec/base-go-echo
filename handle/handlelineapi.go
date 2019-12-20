package Handler

import (
	// "github.com/hb-go/gorm"
	"github.com/hb-go/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ServiceInfo handle set service
type ServiceInfo struct {
	IPAddress   string
	Port        string
	ServiceName string
}

// ServicesInfo array
type ServicesInfo []ServiceInfo

// HTTPCallBackHandler call back webhook
type HTTPCallBackHandler struct {
	Bot          *linebot.Client
	ServicesInfo *ServicesInfo
	DB           *gorm.DB
}

// DBHandler is hanler for connent data base.
type DBHandler struct {
	DB *gorm.DB
}
