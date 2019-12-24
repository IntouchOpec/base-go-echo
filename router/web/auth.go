package web

import (
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

type LoginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Response struct {
	Redirect string      `json:"redirect"`
	User     *model.User `json:"user"`
}

func LoginHandler(c *Context) error {
	a := c.Auth()
	if a.User.IsAuthenticated() {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/admin/dashboard"))
		return nil
	}
	csrfValue := c.Get("_csrf")
	err := c.Render(http.StatusOK, "login", echo.Map{
		"title":  "login",
		"_csrf":  csrfValue,
		"method": "POST",
		// "redirectParam": auth.RedirectParam,
		"redirect": "",
	})
	return err
}

func LoginPostHandler(c *Context) error {
	var form LoginForm
	// redirect := c.QueryParam(auth.RedirectParam)
	loginURL := c.Request().RequestURI
	a := c.Auth()
	response := Response{}
	if a.User.IsAuthenticated() {
		// response.Redirect = redirect
		c.JSON(http.StatusOK, response)
		return nil
	}

	if err := c.Bind(&form); err == nil {
		var User model.User
		u := User.GetUserByEmailPwd(form.Email, form.Password)

		if u != nil {
			session := c.Session()
			err := auth.AuthenticateSession(session, u)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
			}
			response.Redirect = "/admin/dashboard"
			response.User = u

			c.JSON(http.StatusMovedPermanently, response)
			return nil
		}
		response.Redirect = loginURL
		c.JSON(http.StatusMovedPermanently, response)
		return nil
	}
	response.Redirect = loginURL
	return c.JSON(http.StatusBadRequest, response)
}

func LogoutHandler(c *Context) error {
	session := c.Session()
	a := c.Auth()
	auth.Logout(session, a.User)

	// redirect := c.QueryParam(auth.RedirectParam)
	// if redirect == "" {
	// 	redirect = "/"
	// }

	c.Redirect(http.StatusMovedPermanently, "/login")

	return nil
}
