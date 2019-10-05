package channel

import (
	"io"
	"text/template"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/middleware/opentracing"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/IntouchOpec/base-go-echo/module/cache"
	"github.com/IntouchOpec/base-go-echo/module/session"
	"github.com/IntouchOpec/base-go-echo/router/web"
	"github.com/hb-go/echo-web/middleware/captcha"
	"github.com/hb-go/echo-web/module/render"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Routers channel.
func Routers() *echo.Echo {
	e := echo.New()
	// Context
	e.Use(NewContext())
	if Conf.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("web")
	e.Logger.SetLevel(GetLogLvl())

	e.Use(session.Session())

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		ContextKey:  "_csrf",
		TokenLookup: "form:_csrf",
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(captcha.Captcha(captcha.Config{
		CaptchaPath: "/captcha/",
		SkipLogging: true,
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// OpenTracing
	if !Conf.Opentracing.Disable {
		e.Use(opentracing.OpenTracing("web"))
	}

	// 模板
	e.Renderer = render.LoadTemplates()
	e.Use(render.Render())

	// Cache
	e.Use(cache.Cache())

	// Auth
	e.Use(auth.New())
	e.POST("/callback/:account/:ChannelID", HandleWebHookLineAPI)
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.GET("/login", handler(LoginHandler))
	e.POST("/login", handler(LoginPostHandler))
	e.GET("/logout", handler(LogoutHandler))
	e.GET("/register/:lineID", web.LIFFRegisterHandler)
	e.POST("/register/:lineID", web.LIIFRegisterSaveCustomer)

	managent := e.Group("/")
	managent.Use(auth.LoginRequired())
	{
		managent.GET("/dashboard", handler(DashboardHandler))
		managent.GET("book", handler(BookingListHandler))
		managent.GET("book/:lineID", handler(BookingListHandler))
		managent.GET("customer/:id/detail", handler(CustomerDetailHandler))
		managent.GET("customer/:lineID", handler(CustomerListHandler))
		managent.GET("/chat_chennal", handler(BookingListHandler))
		managent.GET("/chat_chennal/:lineID", handler(BookingListHandler))
		managent.GET("/product", handler(BookingListHandler))
		managent.GET("/product/:id", handler(BookingListHandler))
		managent.GET("/promotion", handler(BookingListHandler))
		managent.GET("/promotion/:lineID", handler(BookingListHandler))
		managent.GET("/user", handler(BookingListHandler))
	}
	e.GET("/Book/:lineID", web.BookingListHandler)

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public",
		Browse: true,
	}))
	return e
}

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
