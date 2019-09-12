package api

import (
	"fmt"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/hb-go/echo-web/middleware/opentracing"
	"github.com/hb-go/echo-web/module/cache"
	"github.com/hb-go/echo-web/module/session"
)

//  Routers API
func Routers() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Context
	e.Use(NewContext())

	// Customization
	if Conf.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("api")
	e.Logger.SetLevel(GetLogLvl())

	// Session
	e.Use(session.Session())

	// OpenTracing
	if !Conf.Opentracing.Disable {
		e.Use(opentracing.OpenTracing("api"))
	}

	// CSRF
	// e.Use(mw.CSRFWithConfig(mw.CSRFConfig{
	// 	TokenLookup: "form:X-XSRF-TOKEN",
	// }))

	// Gzip
	e.Use(mw.GzipWithConfig(mw.GzipConfig{
		Level: 5,
	}))

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Static("/favicon.ico", "./assets/img/favicon.ico")

	// Cache
	e.Use(cache.Cache())

	// e.Use(ec.SiteCache(ec.NewMemcachedStore([]string{conf.MEMCACHED_SERVER}, time.Hour), time.Minute))
	// e.GET("/user/:id", ec.CachePage(ec.NewMemcachedStore([]string{conf.MEMCACHED_SERVER}, time.Hour), time.Minute, UserHandler))

	// Routers
	e.POST("/login", UserLoginHandler)
	e.POST("/register", UserRegisterHandler)
	e.POST("/account", CreateAccount)

	e.GET("/json/encode", handler(JsonEncodeHandler))

	e.POST("/chatchannel", CreateChatChannel)
	e.POST("/chatanswer", CreateChatAnswer)
	e.POST("/product", CreateProduct)

	e.PATCH("/chatchannel/:id", GetChannelAccessToken)
	e.GET("/richmenu/:id", GetRichMenu)
	e.POST("/richmenu/:id", CreateRichMenu)
	e.POST("/richmenu/:id/:richMenuID", UploadImageRichMenu)
	e.PATCH("/richmenu/:id/:richMenuID", AtiveRichMenu)

	e.POST("/richmenu/:id/:richMenuID", UploadImageRichMenu)

	// account := e.Group("/account")
	// {
	// 	account.POST("/", CreateAccount)
	// 	// post.GET("/id/:id", PostHandler)
	// 	// post.GET("/:userId/p/:p/s/:s", PostsHandler)
	// }

	// chatanswer := e.Group("/chatanswer")
	// proDuct := e.Group("/product")
	// richMenu := e.Group("/richmenu")
	// chatChannel := e.Group("/chatchannel")
	// richMenu := e.Group("/richmenu")

	// chatanswer.POST("/:channelID", CreateChatAnswer)

	// JWT
	r := e.Group("")
	fmt.Println(echo.HeaderAuthorization)
	r.Use(mw.JWTWithConfig(mw.JWTConfig{
		SigningKey:  []byte("secret"),
		ContextKey:  "_user",
		TokenLookup: "header:" + echo.HeaderAuthorization,
	}))

	r.GET("/", handler(ApiHandler))

	// curl http://echo.api.localhost:8080/restricted/user -H "Authorization: Bearer XXX"
	r.GET("/user", UserHandler)

	return e
}

type (
	HandlerFunc func(*Context) error
)

// HandlerFunc
func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
