package hanlder

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

// HTTPCallBackHanlder call back webhook
type HTTPCallBackHanlder struct {
	Bot          *linebot.Client
	ServicesInfo *ServicesInfo
	DB           *gorm.DB
}

// DBHandler is hanler for connent data base.
type DBHandler struct {
	DB *gorm.DB
}
