package api

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"

	. "github.com/IntouchOpec/base-go-echo/conf"
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

	e.POST("/user", UserRegisterHandler)
	e.GET("/user", GetUserList)

	e.GET("/account", GetAccontList)
	e.GET("/account/:chatChannelID/list", GetAccontList)
	e.GET("/account/:id", GetAccount)
	e.POST("/account", CreateAccount)
	e.PUT("/account/:id", UpdateAccount)
	e.DELETE("/account/:id", DeleteAccount)

	e.GET("/booking/:chatChannelID/list", GetBookingList)
	e.GET("/booking/:id", GetBookingDetail)
	e.PUT("/booking/:id", UpdateBooking)
	e.DELETE("/booking/:id", DeleteBooking)

	e.GET("/chatchannel/:chatChannelID/list", GetChatChannelList)
	e.GET("/chatchannel/:id", GetChatChannelDetail)
	e.POST("/chatchannel", CreateChatChannel)
	e.POST("/chatchannel/:id/setting", CreateChatChannelSetting)
	e.PATCH("/chatchannel/:id", UpdateChatChannel)
	e.DELETE("/chatchannel/id", DeleteChatChannel)
	e.PATCH("/chatchannel/:id/getaccesstoken", GetChannelAccessToken)
	e.PATCH("/chatchannel/:lineID/activeregisterliff", ActiveRegisterLIFFAPI)

	e.GET("/chatanswer/:chatChannelID/list", GetChatAnswerList)
	e.GET("/chatanswer/:id", GetChatAnswerDetail)
	e.POST("/chatanswer", CreateChatAnswer)
	e.PATCH("/chatanswer/:id", UpdateChatAnswers)
	e.DELETE("/chatanswer/:id", DeleteChatAnswer)

	e.GET("/product/:chatChannelID/list", GetProductList)
	e.GET("/product/:id", GetProduct)
	e.POST("/product", CreateProduct)
	e.PATCH("/product/:id", UpdateProduct)

	// e.GET("/subproduct/:chatChannelID", GetProductList)
	// e.GET("/subproduct/:id", GetProduct)
	// e.POST("/subproduct/:chatChannelID/:id", CreateProduct)
	// e.PATCH("/subproduct/:id", UpdateProduct)
	// e.DELETE("/product/:id", CreateProduct)

	e.GET("/richmenu/:id", GetRichMenu)
	e.POST("/richmenu/:id", CreateRichMenu)
	e.POST("/richmenu/:id/:richMenuID", UploadImageRichMenu)
	e.PATCH("/richmenu/:id/:richMenuID", AtiveRichMenu)
	e.DELETE("/richmenu/:id/:richMenuID", DeleteRichMenu)

	e.GET("/promotion/:chatChannelID/list", GetPromotionList)
	e.GET("/promotion/:id", GetPromotion)
	e.POST("/promotion/:lineID", CreatePromotion)
	e.PATCH("/promotion/:id", UpdatePromotion)
	e.DELETE("/promotion/:id", DeletePromotion)

	e.POST("/liff/:chatChannelID", CreateLIFF)

	e.GET("/eventlog", GetAllEvenLog)
	e.GET("/eventlog/:chatChannelID", GetEvenLogList)

	e.GET("/eventlog", GetAllEvenLog)
	e.GET("/eventlog/:chatChannelID", GetEvenLogList)

	e.GET("/customer/:chatChannelID/list", GetCustomerList)
	e.GET("/customer/:id", GetCustomerDetail)
	e.PATCH("/customer/:id", UpdateCustomer)

	//
	e.GET("/json/encode", handler(JsonEncodeHandler))

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
	// r := e.Group("")
	// fmt.Println(echo.HeaderAuthorization)
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
