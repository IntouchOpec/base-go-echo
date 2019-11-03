package session

import (
	. "github.com/IntouchOpec/base-go-echo/conf"

	es "github.com/IntouchOpec/base-go-echo/middleware/session"
	"github.com/labstack/echo"
)

func Session() echo.MiddlewareFunc {
	switch Conf.SessionStore {
	case REDIS:
		store, err := es.NewRedisStore(10, "tcp", "localhost:6379", Conf.Redis.Pwd, []byte("secret-key"))
		if err != nil {
			panic(err)
		}
		return es.New("sid", store)
	case FILE:
		store := es.NewFilesystemStore("", []byte("secret-key"))
		return es.New("sid", store)
	default:
		store := es.NewCookieStore([]byte("secret-key"))
		return es.New("sid", store)
	}
}
