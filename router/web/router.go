package web

import (
	"fmt"
	"io"
	"text/template"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/IntouchOpec/base-go-echo/module/cache"
	"github.com/IntouchOpec/base-go-echo/module/session"
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
		managent.GET("/customer/:id", handler(CustomerDetailHandler))

		managent.GET("/chat_channel", handler(ChatChannelListHandler))
		managent.GET("/chat_channel/create", handler(ChatChannelCreateViewHandler))
		managent.POST("/chat_channel", handler(ChatChannelCreatePostHandler))
		managent.PATCH("/chat_channel/:id/channel_access_token", handler(ChatChannelGetChannelAccessTokenHandler))
		managent.PATCH("/chat_channel/:id/add_liff_register", handler(ChatChannelAddRegisterLIFF))
		managent.DELETE("/chat_channel/:id", handler(ChatChannelListHandler))
		managent.GET("/chat_channel/:id", handler(ChatChannelDetailHandler))
		managent.GET("/chat_channel/:id/edit", handler(ChatChannelEditHandler))

		managent.GET("/chat_answer", handler(ChatAnswerListHandler))
		managent.GET("/chat_answer/create", handler(ChatAnswerCreateHandler))
		managent.GET("/chat_answer/:id", handler(ChatAnswerDetailHandler))
		managent.POST("/chat_answer", handler(ChatAnswerPostHandler))
		managent.GET("/chat_answer/:id/edit", handler(ChatAnswerEditHandler))
		managent.DELETE("/chat_answer/:id", handler(ChatAnswerDeleteHandler))

		managent.GET("/chat_request", handler(ChatRequestListHandler))
		managent.GET("/chat_request/create", handler(ChatRequestCreateHandler))
		managent.GET("/chat_request/:id", handler(ChatRequestDetailHandler))
		managent.POST("/chat_request", handler(ChatRequestPostHandler))
		managent.GET("/chat_request/:id/edit", handler(ChatRequestEditHandler))
		managent.DELETE("/chat_request/:id", handler(ChatRequestDeleteHandler))

		managent.GET("/service", handler(ServiceListHandler))
		managent.GET("/service/create", handler(ServiceCreateHandler))
		managent.GET("/service/:id", handler(ServiceDetailHandler))
		managent.POST("/service", handler(ServicePostHandler))

		managent.GET("/promotion", handler(PromotionListHandler))
		managent.POST("/promotion", handler(PromotionPostHandler))
		managent.GET("/promotion/create", handler(PromotionFormHandler))
		managent.GET("/promotion/:id", handler(PromotionDetailHandler))
		managent.GET("/promotion/:id", handler(PromotionDetailHandler))
		managent.DELETE("/promotion/:id", handler(PromotionDetailHandler))
		managent.GET("/promotion_channel/:id/create", handler(PromotionChannelFormHandler))
		managent.POST("/promotion_channel/:id/create", handler(PromotionChannelAddHandler))

		managent.GET("/user", handler(UserListHandler))
		managent.GET("/user/:id", handler(UserDetailHandler))
		managent.GET("/user/create", handler(UserFormHamdeler))
		managent.GET("/user/:id/edit", handler(UserEditHamdeler))
		managent.DELETE("/user/:id", handler(UserDeleteHandler))

		managent.GET("/richmenu", handler(RichMenuListHandler))
		managent.GET("/richmenu/create", handler(RichMenuCreateViewHandler))
		managent.POST("/richmenu/create", handler(RichMenuCreateHandler))
		managent.GET("/richmenu/:id", handler(RichMenuDetailHandler))
		managent.PATCH("/richmenu/:id", handler(RichMenuImageHandler))
		managent.PATCH("/richmenu/:id/active", handler(RichMenuActiveHandler))
		managent.GET("/setting", handler(SettingHandler))
		// managent.GET("/service/:id/sub_service/create", handler(ServiceSlotCreateHandler))
		// managent.POST("/service/:id/sub_service/create", handler(ServiceSlotPostHandler))
		managent.GET("/service/:id/chatchannel_service/create", handler(ServiceChatChannelViewHandler))
		managent.POST("/service/:id/chatchannel_service/create", handler(ServiceChatChannelPostHandler))

		managent.GET("/LIFF", handler(LIIFListHandler))
		managent.GET("/LIFF/create", handler(LIFFCreateHandler))
		managent.GET("/LIFF/:id", handler(LIIFListHandler))
		managent.POST("/LIFF/create", handler(LIFFPostHanlder))
		managent.DELETE("/LIFF/:id", handler(LIFFRemoveHanlder))

		managent.GET("/place", handler(PlaceListHandler))
		managent.POST("/place", handler(PlaceCreateHandler))
		managent.GET("/place/create", handler(PromotionFormHandler))
		managent.GET("/place/:id", handler(PromotionDetailHandler))
		managent.DELETE("/place/:id", handler(PromotionDetailHandler))

		managent.GET("/provider", handler(ProviderListHandler))
		managent.POST("/provider/create", handler(ProviderPostHandler))
		managent.GET("/provider/create", handler(ProviderCreateHandler))
		managent.GET("/provider/:id", handler(ProviderDetailHandler))
		managent.GET("/provider/:id/edit", handler(ProviderDetailHandler))
		managent.DELETE("/provider/:id", handler(ProviderDetailHandler))

		managent.GET("/provider_service/:id/create", handler(ProviderAddServiceHandler))
		managent.GET("/provider_service/:prov_id", handler(ProviderSerciveListHandler))
		managent.POST("/provider_service/", handler(ProviderAddServiceHandler))
		managent.POST("/provider_service/:id/create", handler(ProviderAddServicePostHandler))
		managent.GET("/provider_booking/:id/create", handler(ProviderAddBookingHandler))

		managent.GET("/time_slot/:provider_id/create", handler(TimeSlotCreateHandler))
		managent.POST("/time_slot/:provider_id/create", handler(TimeSlotPostHandler))
	}

	return e
}

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
