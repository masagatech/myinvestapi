package masters

import (
	"log"

	"github.com/masagatech/myinvest/db"
	"github.com/masagatech/myinvest/util"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewHTTP is ...
func NewHTTP(r *echo.Group) {
	ur := r.Group("/masters")
	ur.GET("/bygroup/:group", getbygroup)
}

func getbygroup(c echo.Context) error {
	// var cmp = "cmp" + c.Param("cmp") + "."
	group := c.Param("group")
	var _result []interface{}
	// sql := `SELECT firstname,lastname, age from ` + cmp + `"user"`

	_col, _ := db.GetCollection(c, db.Masters)
	err := _col.Find(bson.M{"groups": group}).Select(bson.M{"_id": 0, "code": 1, "name": 1}).All(&_result)

	// _, err := _db.QueryOne(&_user, sql)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, nil)

	}

	return util.Response(c, 1, "", "", &_result)
}
