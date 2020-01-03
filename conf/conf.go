package conf

import (
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/labstack/gommon/log"
)

var (
	Conf              Config // holds the global app config.
	defaultConfigFile = "conf/conf.toml"
)

// Config config base
type Config struct {
	ReleaseMode bool   `toml:"release_mode"`
	LogLevel    string `toml:"log_level"`

	SessionStore string `toml:"session_store"`
	CacheStore   string `toml:"cache_store"`

	// name
	App app

	// 模板
	Tmpl tmpl

	Server server

	// prosgret
	DB database `toml:"database"`

	// 静态资源
	Static static

	// Redis
	Redis redis

	// Memcached
	Memcached memcached

	// Opentracing
	Opentracing opentracing

	// Metrics
	Metrics metrics
}

type app struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

type server struct {
	Graceful bool   `toml:"graceful"`
	Addr     string `toml:"addr"`

	Domain            string `toml:"domain"`
	DomainAPI         string `toml:"domain_api"`
	DomainWeb         string `toml:"domain_web"`
	DomainSocket      string `toml:"domain_socket"`
	DomainLineChannel string `toml:"domain_line_channel"`
}

type static struct {
	Type string `toml:"type"`
}

type tmpl struct {
	Type   string `toml:"type"`   // PONGO2,TEMPLATE(TEMPLATE Default)
	Data   string `toml:"data"`   // BINDATA,FILE(FILE Default)
	Dir    string `toml:"dir"`    // PONGO2(template/pongo2),TEMPLATE(template)
	Suffix string `toml:"suffix"` // .html,.tpl
}

type database struct {
	Name     string `toml:"name"`
	UserName string `toml:"user_name"`
	Pwd      string `toml:"pwd"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
}

type redis struct {
	Server string `toml:"server"`
	Pwd    string `toml:"pwd"`
}

type memcached struct {
	Server string `toml:"server"`
}

type opentracing struct {
	Disable     bool   `toml:"disable"`
	Type        string `toml:"type"`
	ServiceName string `toml:"service_name"`
	Address     string `toml:"address"`
}

type metrics struct {
	Disable bool          `toml:"disable"`
	FreqSec time.Duration `toml:"freq_sec"`
	Address string        `toml:"address"`
}

func init() {
}

// InitConfig initializes the app configuration by first setting defaults,
// then overriding settings from the app config file, then overriding
// It returns an error if any.
func InitConfig(configFile string) error {
	if configFile == "" {
		configFile = defaultConfigFile
	}

	// Set defaults.
	Conf = Config{
		ReleaseMode: false,
		LogLevel:    "DEBUG",
	}

	if _, err := os.Stat(configFile); err != nil {
		return errors.New("config file err:" + err.Error())
	} else {
		log.Infof("load config from file:" + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return errors.New("config load err:" + err.Error())
		}
		_, err = toml.Decode(string(configBytes), &Conf)
		if err != nil {
			return errors.New("config decode err:" + err.Error())
		}
	}

	// @TODO 配置检查
	log.Infof("config data:%v", Conf)

	return nil
}

// GetLogLvl state log.
func GetLogLvl() log.Lvl {
	// DEBUG INFO WARN ERROR OFF
	switch Conf.LogLevel {
	case "DEBUG":
		return log.DEBUG
	case "INFO":
		return log.INFO
	case "WARN":
		return log.WARN
	case "ERROR":
		return log.ERROR
	case "OF":
		return log.OFF
	}

	return log.DEBUG
}

const (
	PONGO2    = "PONGO2"
	TEMPLATE  = "TEMPLATE"
	BINDATA   = "BINDATA"
	FILE      = "FILE"
	REDIS     = "REDIS"
	MEMCACHED = "MEMCACHED"
	COOKIE    = "COOKIE"
	IN_MEMORY = "IN_MEMARY"
)
