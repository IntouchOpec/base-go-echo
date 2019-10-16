package web

import (
	"fmt"
	"io"
	"text/template"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/IntouchOpec/base-go-echo/module/cache"
	"github.com/IntouchOpec/base-go-echo/module/session"
	"github.com/hb-go/echo-web/middleware/captcha"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Template struct {
	templates *template.Template
}

type BaseTempleRespone struct {
	Title       string `json:"title"`
	AccountName string `json:"account_name"`
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	a := auth.Default(c)
	if viewContext, isMap := data.(echo.Map); isMap {
		acc := a.User.GetAccount()
		csrfValue := c.Get("_csrf")
		viewContext["base"] = echo.Map{"title": fmt.Sprintf("%s", viewContext["title"]), "account": acc}
		viewContext["_csrf"] = csrfValue
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func Routers() *echo.Echo {
	e := echo.New()

	// e.GET("/Book/:lineID", BookingListHandler)
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

	e.Use(cache.Cache())

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = t
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public/assets",
		Browse: true,
	}))
	e.Use(auth.New())

	e.GET("/login", handler(LoginHandler))
	e.POST("/login", handler(LoginPostHandler))
	e.GET("/logout", handler(LogoutHandler))
	e.GET("/register/:lineID", LIFFRegisterHandler)
	e.POST("/register/:lineID", LIIFRegisterSaveCustomer)

	managent := e.Group("/admin")
	managent.Use(auth.LoginRequired())
	{
		managent.GET("/dashboard", handler(DashboardHandler))
		managent.GET("/book", handler(BookingListHandler))
		managent.GET("/customer", handler(CustomerListHandler))
		managent.GET("/customer/:lineID", handler(CustomerListHandler))
		managent.GET("/chat_channel", handler(ChatChannelListHandler))
		managent.GET("/chat_channel/create", handler(ChatChannelCreateViewHandler))
		managent.POST("/chat_channel", handler(ChatChannelCreatePostHandler))
		managent.PATCH("/chat_channel/:id/channel_access_token", handler(ChatChannelGetChannelAccessTokenHandler))
		managent.DELETE("/chat_channel/:id", handler(ChatChannelListHandler))
		managent.GET("/chat_channel/:id", handler(ChatChannelDetailHandler))
		managent.GET("/chat_channel/:id/edit", handler(ChatChannelEditHandler))
		managent.GET("/chat_answer", handler(ChatAnswerListHandler))
		managent.GET("/product/:id", handler(ProductDetailHandler))
		managent.GET("/product", handler(ProductListHandler))
		managent.GET("/promotion", handler(PromotionListHandler))
		managent.GET("/promotion/:lineID", handler(BookingListHandler))
		managent.GET("/user", handler(UserListHandler))
		managent.GET("/setting", handler(SettingHandler))
	}

	return e
}

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
