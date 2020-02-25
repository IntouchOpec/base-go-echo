package main

import (
	"flag"
	"log"

	"github.com/IntouchOpec/base-go-echo/conf"
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
}

func main() {
	if cmdHelp {
		flag.PrintDefaults()
		return
	}
	if err := conf.InitConfig(confFilePath); err != nil {
		log.Panic(err)
	}
	// model.Initialize()
	// log.Debugf("run with conf:%s", confFilePath)
	router.RunSubdomains(confFilePath)
}
