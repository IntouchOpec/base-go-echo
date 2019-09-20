package channel

import (
	"io"
	"text/template"

	"github.com/IntouchOpec/base-go-echo/router/web"
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

	e.POST("/callback/:account/:ChannelID", HandleWebHookLineAPI)
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	e.GET("/register/:lineID", web.LIFFRegisterHandler)
	e.POST("/register/:lineID", web.LIIFRegisterSaveCustomer)

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public",
		Browse: true,
	}))
	return e
}
