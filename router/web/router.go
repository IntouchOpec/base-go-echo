package web

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Routers() *echo.Echo {
	e := echo.New()
	fmt.Println("test")
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.GET("/hello", func(c echo.Context) error {
		return c.Render(http.StatusOK, "hello", "World")
	})
	e.GET("/register/:lineID", LIFFRegisterHandler)

	return e
}
