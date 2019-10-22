package channel

import (
	"html/template"
	"io"

	"github.com/IntouchOpec/base-go-echo/router/web"
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
	return t.templates.ExecuteTemplate(w, name, data)
}

// Routers channel.
func Routers() *echo.Echo {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public/assets",
		Browse: true,
	}))
	e.POST("/callback/:account/:ChannelID", HandleWebHookLineAPI)
	e.GET("/register/:lineID", web.LIFFRegisterHandler)
	e.POST("/register/:lineID", web.LIIFRegisterSaveCustomer)
	return e
}
