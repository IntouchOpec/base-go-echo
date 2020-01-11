package web

import (
	"fmt"
	"io"
	"net/http"
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

	e.GET("/register/:lineID", LIFFRegisterHandler)
	e.POST("/register/:lineID", LIIFRegisterSaveCustomer)

	e.Use(auth.New())
	e.GET("/", handler(indexHandler))

	e.GET("/login", handler(LoginHandler))
	e.POST("/login", handler(LoginPostHandler))
	e.GET("/logout", handler(LogoutHandler))

	managent := e.Group("/admin")
	managent.Use(auth.LoginRequired())
	{
		managent.GET("/dashboard", handler(DashboardHandler))
		managent.GET("/book", handler(BookingListHandler))
		managent.GET("/customer", handler(CustomerListHandler))
		managent.GET("/customer/:id", handler(CustomerDetailHandler))
		managent.DELETE("/customer/:id", handler(CustomerDeleteHandler))

		managent.GET("/customer_type", handler(CustomerTypeListHandler))
		managent.POST("/customer_type/create", handler(CustomerTypePostHandler))
		managent.GET("/customer_type/create", handler(CustomerTypeCreateHandler))
		managent.GET("/customer_type/:id/edit", handler(CustomerTypeEditViewHandler))
		managent.PUT("/customer_type/:id/edit", handler(CustomerTypeEditPutHandler))
		managent.DELETE("/customer_type/:id", handler(CustomerTypeDeleteHandler))

		managent.GET("/chat_channel", handler(ChatChannelListHandler))
		managent.POST("/chat_channel/create", handler(CustomerTypePostHandler))
		managent.GET("/chat_channel/create", handler(ChatChannelCreateViewHandler))
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
		managent.POST("/promotion/:id", handler(PromotionCreateDetailHandler))
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
		managent.GET("/richmenu/:id/edit", handler(RichMenuCreateViewHandler))
		managent.PATCH("/richmenu/:id/download_image", handler(RichMenuDonwloadImage))
		managent.PATCH("/richmenu/:id/active", handler(RichMenuActiveHandler))
		managent.DELETE("/richmenu/:id", handler(RichMenuDeleteHandler))
		managent.POST("/richmenu/:id/bulk_link", handler(RichMenuAddCustomerHandler))

		managent.GET("/setting", handler(SettingHandler))
		managent.GET("/setting/:id/edit", handler(SettingEditViewHandler))
		managent.PUT("/setting/:id/edit", handler(SettingPutHandler))
		managent.PUT("/setting", handler(SettingPostHandler))
		managent.DELETE("/setting/auth_json_file", handler(RemoaveAuthJSONFile))
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
		managent.POST("/place/:id/place_service/create", handler(PlaceAddSercivePostHandler))
		managent.GET("/place/:id/place_service/create", handler(PlaceAddSerciveViewHandler))

		managent.GET("/employee", handler(EmployeeListHandler))
		managent.GET("/employee/create", handler(EmployeeCreateHandler))
		managent.POST("/employee/create", handler(EmployeePostHandler))
		managent.GET("/employee/:id", handler(EmployeeDetailHandler))
		managent.GET("/employee/:id/edit", handler(EmployeeEditHandler))
		managent.PUT("/employee/:id/edit", handler(EmployeeEditHandler))
		managent.DELETE("/employee/:id", handler(EmployeeDeleteHandler))
		managent.DELETE("/employee/:id/delete_image", handler(EmployeeDeleteImageHandler))

		managent.GET("/transaction", handler(TransactionListHandler))
		managent.GET("/transaction/create", handler(TransactionCreateHandler))
		managent.POST("/transaction/create", handler(TransactionPostHandler))
		managent.GET("/transaction/:id", handler(TransactionDetailHandler))
		managent.PATCH("/transaction/:id", handler(TransactionPatchHandler))
		managent.GET("/transaction/:id/edit", handler(TransactionEditHandler))
		managent.DELETE("/transaction/:id", handler(TransactionDeleteHandler))

		managent.GET("/employee_service/:id/create", handler(EmployeeAddServiceHandler))
		managent.POST("/employee_service/:id/create", handler(EmployeeAddServicePostHandler))
		managent.GET("/employee_service/:prov_id", handler(EmployeeSerciveListHandler))
		// managent.POST("/employee_service/", handler(EmployeeAddServiceHandler))
		managent.GET("/employee_booking/:id/create", handler(EmployeeAddBookingHandler))

		managent.GET("/time_slot/:employee_id/create", handler(TimeSlotCreateHandler))
		managent.POST("/time_slot/:employee_id/create", handler(TimeSlotPostHandler))

		managent.GET("/action_log", handler(ActionLogList))
		managent.GET("/event_log", handler(EventLogList))

		managent.GET("/upload_file", handler(FileListHandler))
		managent.POST("/upload_file", handler(FileCreateHandler))
		managent.DELETE("/upload_file/:id", handler(FileRemoveHandler))
	}

	files := e.Group("/files")
	files.GET("", handler(GetFileGoogleStorageHandler))
	files.PUT("", handler(UploadFileGoogleStorageHandler))

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "public/assets",
	}))

	return e
}

func indexHandler(c *Context) error {
	a := c.Auth()
	if a.User.IsAuthenticated() {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/admin/dashboard"))
		return nil
	}
	csrfValue := c.Get("_csrf")
	return c.Render(http.StatusOK, "login", echo.Map{
		"title":  "login",
		"_csrf":  csrfValue,
		"method": "POST",
		// "redirectParam": auth.RedirectParam,
		"redirect": "",
	})

}

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
