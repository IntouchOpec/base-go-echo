package web

import (
	"net/http"

	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

// type Template struct {
// 	templates *template.Template
// }

// func (t *Templatem) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.templates.ExecuteTemplate(w, name, data)
// }
type Context struct {
	echo.Context
}

func (ctx *Context) Auth() auth.Auth {
	return auth.Default(ctx)
}

func Routers() *echo.Echo {
	e := echo.New()
	// e.Use(store)
	// Auth
	// e.Use(auth.New(model.GenerateAnonymousUser()))
	// t := &Template{
	// 	templates: template.Must(template.ParseGlob("public/views/*.html")),
	// }
	// e.Renderer = t
	e.GET("/hello", func(c echo.Context) error {
		return c.Render(http.StatusOK, "hello", "World")
	})

	e.GET("/register/:lineID", LIFFRegisterHandler)
	e.POST("/register/:lineID", LIIFRegisterSaveCustomer)

	// e.GET("/Book/:lineID", BookingListHandler)

	return e
}

type (
	HandlerFunc func(*Context) error
)

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
