package log

import (
	l "github.com/labstack/gommon/log"
)

type Logger struct {
	*l.Logger
}

var (
	global        = l.New("log")
	defaultHeader = `{"time":"${time_rfc3339}","level":"${level}","prefix":"${prefix}",` + `"file":"${long_file}","line":"${line}"}`
)

func init() {
	l.SetHeader(defaultHeader)
	l.SetLevel(l.DEBUG)

	global.SetHeader(defaultHeader)
	global.SetLevel(l.DEBUG)
}

func SetLevel(v l.Lvl) {
	l.SetLevel(v)
	global.SetLevel(v)
}

func Debug(i ...interface{}) {
	global.Debug(i)
}

func Debugf(format string, values ...interface{}) {
	global.Debugf(format, values...)
}

func Info(i ...interface{}) {
	global.Info(i)
}

func Infof(format string, values ...interface{}) {
	global.Infof(format, values...)
}

func Warn(i ...interface{}) {
	global.Warn(i)
}

func Warnf(format string, values ...interface{}) {
	global.Warnf(format, values...)
}

func Error(i ...interface{}) {
	global.Error(i)
}

func Errorf(format string, values ...interface{}) {
	global.Errorf(format, values...)
}

func Fatal(i ...interface{}) {
	global.Fatal(i)
}

func Fatalf(format string, values ...interface{}) {
	global.Fatalf(format, values...)
}

func Panic(i ...interface{}) {
	global.Panic(i)
}

func Panicf(format string, args ...interface{}) {
	global.Panicf(format, args)
}
