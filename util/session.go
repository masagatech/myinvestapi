package util

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

var ctx = context.Background()

func SessionGet(c echo.Context, userid string, key string) string {
	if userid == "" {
		userid = c.Request().Header.Get("uid")
	}
	rd := GetRedis(c)
	fmt.Println(userid)
	val, err := rd.HGet(ctx, userid, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func SessionSet(c echo.Context, userid string, key string, value string) {
	if userid == "" {
		userid = c.Request().Header.Get("uid")
	}
	fmt.Println(userid)
	rd := GetRedis(c)
	fmt.Println(userid)

	err := rd.HSet(ctx, userid, []string{key, value}).Err()
	if err != nil {
		panic(err)
	}

	// sess, _ := session.Get(userid, c)
	// sess.Values[key] = value
	// err := sess.Save(c.Request(), c.Response())
	// if err != nil {
	// 	http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
	// 	return
	// }

}
