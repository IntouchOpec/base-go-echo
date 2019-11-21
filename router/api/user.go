package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// UserHandler is
func UserHandler(c echo.Context) error {
	user := c.Get("_user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		panic(err)
	}

	var User model.User
	u := User.GetUserByID(id)

	c.JSON(http.StatusOK, map[string]interface{}{
		"title":  "User",
		"user":   u,
		"claims": claims,
	})

	return nil
}

// UserLoginHandler login system.
func UserLoginHandler(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	user.GetUserByUserName()
	t, err := getJWTToken(&user)
	if err != nil {
		return err
	}

	c.JSON(200, map[string]interface{}{
		"URI":   "api user login",
		"token": t,
	})

	return nil
}

// UserRegisterHandler is function handler register create user.
func UserRegisterHandler(c echo.Context) error {
	user := model.User{}

	if err := c.Bind(&user); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	user.AddUserWithUserNamePwd()
	c.JSON(200, user)

	return nil
}

func getJWTToken(u *model.User) (t string, e error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["name"] = u.UserName
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	// Generate encoded token and send it as response.
	t, e = token.SignedString([]byte("secret"))
	return
}

func GetUserList(c echo.Context) error {
	users := model.GetUserList()
	return c.JSON(200, users)

}

func UpdateUser(c echo.Context) error {
	id := c.Param("id")

	user, err := model.GetUserDetail(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := user.UpdateUser(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusAccepted, user)
}
