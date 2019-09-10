package main

import (
	"flag"
	"log"

	"github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/model"

	"github.com/IntouchOpec/base-go-echo/router"
)

const (
	DefaultConfFilePath = "conf/conf.toml"
)

var (
	confFilePath string
	cmdHelp      bool
)

func init() {
	flag.StringVar(&confFilePath, "c", DefaultConfFilePath, "配置文件路径")
	flag.BoolVar(&cmdHelp, "h", false, "帮助")
	flag.Parse()

}

func main() {
	if cmdHelp {
		flag.PrintDefaults()
		return
	}
	if err := conf.InitConfig(confFilePath); err != nil {
		log.Panic(err)
	}
	// log.Debugf("run with conf:%s", confFilePath)
	model.Initialize()
	router.RunSubdomains(confFilePath)
}
