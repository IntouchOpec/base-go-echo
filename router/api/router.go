package api

import (
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/module/cache"
	"github.com/IntouchOpec/base-go-echo/module/session"
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
	// e.GET("/users/:id", ec.CachePage(ec.NewMemcachedStore([]string{conf.MEMCACHED_SERVER}, time.Hour), time.Minute, UserHandler))

	// Routers
	e.GET("/", indexHandler)
	e.GET("/loaderio-6b18ff04186b994982d4dcf4d621c407", GetAccontList)
	// e.POST("/login", UserLoginHandler)
	// e.POST("/register", UserRegisterHandler)

	// e.POST("/user", UserRegisterHandler)
	// e.GET("/user", GetUserList)
	// e.PUT("/users/:id", UpdateUser)

	e.GET("/account", GetAccontList)
	// e.GET("/account/:chatChannelID/list", GetAccontList)
	// e.GET("/account/:id", GetAccount)
	// e.POST("/account", CreateAccount)
	// e.PUT("/account/:id", UpdateAccount)
	// e.DELETE("/account/:id", DeleteAccount)

	// e.GET("/booking/:chatChannelID/list", GetBookingList)
	// e.GET("/booking/:id", GetBookingDetail)
	// e.PUT("/booking/:id", UpdateBooking)
	// e.DELETE("/booking/:id", DeleteBooking)

	// e.GET("/chatchannel", GetChatChannelList)
	// e.GET("/chatchannel/:id", GetChatChannelDetail)
	// e.POST("/chatchannel", CreateChatChannel)
	// e.POST("/chatchannel/:id/setting", CreateChatChannelSetting)
	// e.PATCH("/chatchannel/:id", UpdateChatChannel)
	// e.DELETE("/chatchannel/:id", DeleteChatChannel)
	// e.PATCH("/chatchannel/:id/getaccesstoken", GetChannelAccessToken)
	// e.PATCH("/chatchannel/:lineID/activeregisterliff", ActiveRegisterLIFFAPI)

	// e.GET("/chatanswer/:chatChannelID/list", GetChatAnswerList)
	// e.GET("/chatanswer/:id", GetChatAnswerDetail)
	// e.POST("/chatanswer", CreateChatAnswer)
	// e.PATCH("/chatanswer/:id", UpdateChatAnswers)
	// e.DELETE("/chatanswer/:id", DeleteChatAnswer)

	// // e.GET("/service/:chatChannelID/list", GetserviceList)
	// e.GET("/service/:id", Getservice)
	// e.POST("/service", Createservice)
	// e.PATCH("/service/:id", UpdateService)

	// // e.GET("/subservice/:chatChannelID", GetserviceList)
	// // e.GET("/subservice/:id", Getservice)
	// // e.POST("/subservice/:chatChannelID/:id", Createservice)
	// // e.PATCH("/subservice/:id", UpdateService)
	// // e.DELETE("/service/:id", Createservice)

	// e.GET("/richmenu/:id", GetRichMenu)
	// e.POST("/richmenu/:id", CreateRichMenu)
	// e.POST("/richmenu/:id/:richMenuID", UploadImageRichMenu)
	// e.PATCH("/richmenu/:id/:richMenuID", AtiveRichMenu)
	// e.DELETE("/richmenu/:id/:richMenuID", DeleteRichMenu)

	// e.GET("/promotion/:chatChannelID/list", GetPromotionList)
	// e.GET("/promotion/:id", GetPromotion)
	// e.POST("/promotion", CreatePromotion)
	// e.PATCH("/promotion/:id", UpdatePromotion)
	// e.DELETE("/promotion/:id", DeletePromotion)

	// e.POST("/liff/:chatChannelID", CreateLIFF)

	// e.GET("/eventlog", GetAllEvenLog)
	// e.GET("/eventlog/:chatChannelID", GetEvenLogList)

	// e.GET("/eventlog", GetAllEvenLog)
	// e.GET("/eventlog/:chatChannelID", GetEvenLogList)

	// e.GET("/customer/:chatChannelID/list", GetCustomerList)
	// e.GET("/customer/:id", GetCustomerDetail)
	// e.PATCH("/customer/:id", UpdateCustomer)

	// e.GET("/place/:chatChannelID/list", GetPlaceList)
	// e.GET("/place/:id", GetPlaceDetail)
	// e.POST("/place", CreatePlace)
	// e.PATCH("/place/:id", UpdatePlace)
	// e.DELETE("/place/:id", DeletePlace)

	// //
	// e.GET("/json/encode", handler(JsonEncodeHandler))

	// account := e.Group("/account")
	// {
	// 	account.POST("/", CreateAccount)
	// 	// post.GET("/id/:id", PostHandler)
	// 	// post.GET("/:userId/p/:p/s/:s", PostsHandler)
	// }

	// chatanswer := e.Group("/chatanswer")
	// Service:= e.Group("/service")
	// richMenu := e.Group("/richmenu")
	// chatChannel := e.Group("/chatchannel")
	// richMenu := e.Group("/richmenu")

	// chatanswer.POST("/:channelID", CreateChatAnswer)

	// JWT
	// r := e.Group("")
	// r.Use(mw.JWTWithConfig(mw.JWTConfig{
	// 	SigningKey:  []byte("secret"),
	// 	ContextKey:  "_user",
	// 	TokenLookup: "header:" + echo.HeaderAuthorization,
	// }))

	// r.GET("/", handler(ApiHandler))

	// curl http://echo.api.localhost:8080/restricted/user -H "Authorization: Bearer XXX"
	// r.GET("/user", UserHandler)

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

func indexHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "api ready on "+Conf.Server.DomainAPI)
}
