package orm

import "github.com/labstack/gommon/log"

type Logger struct {
}

// Print format & print log
func (logger Logger) Print(values ...interface{}) {
	log.Debugf("orm log:%v \n", values)
}
