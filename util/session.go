package util

import (
	"fmt"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionGet(c echo.Context, userid string, key interface{}) interface{} {
	if userid == "" {
		userid = c.Request().Header.Get("uid")
	}

	fmt.Println(userid)
	sess, _ := session.Get(userid, c)
	return sess.Values[key]

}

func SessionSet(c echo.Context, userid string, key interface{}, value interface{}) {
	if userid == "" {
		userid = c.Request().Header.Get("uid")
	}
	fmt.Println(userid)
	sess, _ := session.Get(userid, c)
	sess.Values[key] = value
	sess.Save(c.Request(), c.Response())

}
