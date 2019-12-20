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
	e.GET("/", handler(LoginHandler))

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
		managent.DELETE("/customer/:id", handler(CustomerDeleteHandler))

		managent.GET("/chat_channel", handler(ChatChannelListHandler))
		managent.GET("/chat_channel/create", handler(ChatChannelCreateViewHandler))
		managent.POST("/chat_channel", handler(ChatChannelCreatePostHandler))
		managent.PATCH("/chat_channel/:id/channel_access_token", handler(ChatChannelGetChannelAccessTokenHandler))
		managent.PATCH("/chat_channel/:id/add_liff_register", handler(ChatChannelAddRegisterLIFF))
		managent.DELETE("/chat_channel/:id", handler(ChatChannelDeleteHandler))

		managent.GET("/chat_channel/:id", handler(ChatChannelDetailHandler))
		managent.GET("/chat_channel/:id/broadcast", handler(ChatChannelBroadcastMessageViewHandler))
		managent.POST("/chat_channel/:id/broadcast", handler(ChatChannelBroadcastMessageHandler))
		managent.GET("/chat_channel/:id/edit", handler(ChatChannelEditHandler))

		managent.GET("/chat_answer", handler(ChatAnswerListHandler))
		managent.GET("/chat_answer/create", handler(ChatAnswerCreateHandler))
		managent.POST("/chat_answer/create", handler(ChatAnswerPostHandler))
		managent.GET("/chat_answer/:id", handler(ChatAnswerDetailHandler))
		managent.GET("/chat_answer/:id/edit", handler(ChatAnswerEditHandler))
		managent.DELETE("/chat_answer/:id", handler(ChatAnswerDeleteHandler))

		managent.GET("/chat_request", handler(ChatRequestListHandler))
		managent.POST("/chat_request/create", handler(ChatRequestPostHandler))
		managent.GET("/chat_request/create", handler(ChatRequestCreateHandler))
		managent.GET("/chat_request/:id", handler(ChatRequestDetailHandler))
		managent.GET("/chat_request/:id/edit", handler(ChatRequestEditHandler))
		managent.DELETE("/chat_request/:id", handler(ChatRequestDeleteHandler))

		managent.GET("/service", handler(ServiceListHandler))
		managent.GET("/service/create", handler(ServiceCreateHandler))
		managent.POST("/service/create", handler(ServicePostHandler))
		managent.PATCH("/service/:id", handler(ServicePostHandler))
		managent.GET("/service/:id", handler(ServiceDetailHandler))
		managent.GET("/service/:id/edit", handler(ServiceEditViewHandler))
		managent.PUT("/service/:id/edit", handler(ServiceEditPutHandler))
		managent.DELETE("/service/:id", handler(ServiceDeleteHandler))
		managent.DELETE("/service/:id/delete_image", handler(ServiceDeleteImageHandler))

		managent.GET("/promotion", handler(PromotionListHandler))
		managent.POST("/promotion/create", handler(PromotionPostHandler))
		managent.GET("/promotion/create", handler(PromotionFormHandler))
		managent.GET("/promotion/:id", handler(PromotionDetailHandler))
		managent.GET("/promotion/:id/edit", handler(PromotionEditHandler))
		managent.PUT("/promotion/:id/edit", handler(PromotionEditPutHandler))
		managent.POST("/promotion/:id/register", handler(PromotionAddRegisterlHandler))
		managent.DELETE("/promotion/:id", handler(PromotionRemoveHandler))
		managent.GET("/promotion_channel/:id/create", handler(PromotionChannelFormHandler))
		managent.POST("/promotion_channel/:id/create", handler(PromotionChannelAddHandler))
		managent.DELETE("/promotion/:id/delete_image", handler(PromotionDeleteImageHandler))

		managent.GET("/users", handler(UserListHandler))
		managent.GET("/users/:id", handler(UserDetailHandler))
		managent.GET("/users/create", handler(UserFormHandler))
		managent.POST("/users/create", handler(UserPostHandler))
		managent.GET("/users/:id/edit", handler(UserEditHandler))
		managent.PUT("/users/:id/edit", handler(UserPutHandler))
		managent.DELETE("/users/:id", handler(UserDeleteHandler))

		managent.GET("/richmenu", handler(RichMenuListHandler))
		managent.GET("/richmenu/create", handler(RichMenuCreateViewHandler))
		managent.POST("/richmenu/create", handler(RichMenuCreateHandler))
		managent.GET("/richmenu/:id", handler(RichMenuDetailHandler))
		managent.PATCH("/richmenu/:id", handler(RichMenuImageHandler))
		managent.PATCH("/richmenu/:id/download_image", handler(RichMenuDonwloadImage))
		managent.PATCH("/richmenu/:id/active", handler(RichMenuActiveHandler))
		managent.DELETE("/richmenu/:id", handler(RichMenuDeleteHandler))

		managent.GET("/setting", handler(SettingHandler))
		// managent.GET("/service/:id/sub_service/create", handler(ServiceSlotCreateHandler))
		// managent.POST("/service/:id/sub_service/create", handler(ServiceSlotPostHandler))
		managent.GET("/service/:id/chatchannel_service/create", handler(ServiceChatChannelViewHandler))
		managent.POST("/service/:id/chatchannel_service/create", handler(ServiceChatChannelPostHandler))

		managent.GET("/LIFF", handler(LIIFListHandler))
		managent.GET("/LIFF/create", handler(LIFFCreateHandler))
		managent.GET("/LIFF/:id", handler(LIIFListHandler))
		managent.POST("/LIFF/create", handler(LIFFPostHandler))
		managent.DELETE("/LIFF/:id", handler(LIFFRemoveHandler))

		managent.GET("/place", handler(PlaceListHandler))
		managent.POST("/place/create", handler(PlacePostHandler))
		managent.GET("/place/create", handler(PlaceCreateHandler))
		managent.GET("/place/:id", handler(PlaceDetailHandler))
		managent.PUT("/place/:id", handler(PlacePutHandler))
		managent.DELETE("/place/:id", handler(PlaceDeleteHandler))
		managent.DELETE("/place/:id/delete_image", handler(PlaceDeleteImageHandler))
		managent.POST("/place/:id/place_chatchannel/create", handler(PlaceAddChatChannelPostHandler))
		managent.GET("/place/:id/place_chatchannel/create", handler(PlaceAddChatChannelViewHandler))

		managent.GET("/provider", handler(ProviderListHandler))
		managent.GET("/provider/create", handler(ProviderCreateHandler))
		managent.POST("/provider/create", handler(ProviderPostHandler))
		managent.GET("/provider/:id", handler(ProviderDetailHandler))
		managent.GET("/provider/:id/edit", handler(ProviderEditHandler))
		managent.PUT("/provider/:id/edit", handler(ProviderEditHandler))
		managent.DELETE("/provider/:id", handler(ProviderDeleteHandler))
		managent.DELETE("/provider/:id/delete_image", handler(ProviderDeleteImageHandler))

		managent.GET("/transaction", handler(TransactionListHandler))
		managent.GET("/transaction/create", handler(TransactionCreateHandler))
		managent.POST("/transaction/create", handler(TransactionPostHandler))
		managent.GET("/transaction/:id", handler(TransactionDetailHandler))
		managent.PATCH("/transaction/:id", handler(TransactionPatchHandler))
		managent.GET("/transaction/:id/edit", handler(TransactionEditHandler))
		managent.DELETE("/transaction/:id", handler(TransactionDeleteHandler))

		managent.GET("/provider_service/:id/create", handler(ProviderAddServiceHandler))
		managent.POST("/provider_service/:id/create", handler(ProviderAddServicePostHandler))
		managent.GET("/provider_service/:prov_id", handler(ProviderSerciveListHandler))
		// managent.POST("/provider_service/", handler(ProviderAddServiceHandler))
		managent.GET("/provider_booking/:id/create", handler(ProviderAddBookingHandler))

		managent.GET("/time_slot/:provider_id/create", handler(TimeSlotCreateHandler))
		managent.POST("/time_slot/:provider_id/create", handler(TimeSlotPostHandler))

		managent.GET("/action_log", handler(ActionLogList))
		managent.GET("/event_log", handler(EventLogList))

		managent.GET("/upload_file", handler(FileListHandler))
		managent.POST("/upload_file", handler(FileCreateHandler))
		managent.DELETE("/upload_file/:id", handler(FileRemoveHandler))
	}

	return e
}

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
