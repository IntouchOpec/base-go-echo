package channel

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/module/cache"
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
	if viewContext, isMap := data.(echo.Map); isMap {
		csrfValue := c.Get("_csrf")
		viewContext["base"] = echo.Map{"title": fmt.Sprintf("%s", viewContext["title"])}
		viewContext["_csrf"] = csrfValue
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

// Routers channel.
func Routers() *echo.Echo {
	e := echo.New()

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		ContextKey:  "_csrf",
		TokenLookup: "form:_csrf",
	}))

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Use(cache.Cache())

	e.Renderer = t

	e.GET("/", indexHandler)

	e.GET("/register/:lineID", web.LIFFRegisterHandler)
	e.POST("/register/:lineID", web.LIIFRegisterSaveCustomer)
	e.POST("/callback/:account/:ChannelID", HandleWebHookLineAPI)
	e.GET("/webhook-facebook", HandleWebHookFacebookAPI)
	e.POST("/webhook-facebook", HandleWebHookFacebookAPI)
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "public/assets",
	}))
	return e
}

func indexHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "welcome line connect "+Conf.Server.DomainLineChannel)
}
