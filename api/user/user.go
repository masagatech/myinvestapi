package user

import (
	"log"

	"github.com/masagatech/myinvest/db"
	"github.com/masagatech/myinvest/model"
	"github.com/masagatech/myinvest/util"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
)

// NewHTTP is ...
func NewHTTP(r *echo.Group) {
	ur := r.Group("/users")
	ur.POST("", create)
	ur.GET("/id", get)
	ur.GET("/list", getList)
}

func create(c echo.Context) (err error) {
	type _user struct {
		Firstname string `json:"fname"  validate:"required"`
		Lastname  string `json:"lname" validate:"required"`
		Age       int    `json:"age" validate:"required"`
		Id        int    `json:"id" validate:"required"`
	}
	_u := new(_user)
	if err = c.Bind(_u); err != nil {
		return
	}

	if err = c.Validate(_u); err != nil {
		return
	}

	sql := `insert into "user" (firstname,lastname, age, id) values  (?,?,?,?)`
	fmt.Println(sql)
	_db := c.Get("db").(*pg.DB)
	_, err1 := _db.Exec(sql, _u.Firstname, _u.Lastname, _u.Age, _u.Id)

	if err1 != nil {
		fmt.Println(err1)
		return c.JSON(http.StatusBadRequest, err1)

	}

	return c.JSON(http.StatusOK, "Saved Successfully")
}

func get(c echo.Context) error {
	// var cmp = "cmp" + c.Param("cmp") + "."

	var _user []interface{}
	// sql := `SELECT firstname,lastname, age from ` + cmp + `"user"`

	_col, _ := db.GetCollection(c, db.Trainers)
	err := _col.Find(nil).Select(bson.M{"_id": 0, "name": 1}).All(&_user)

	// _, err := _db.QueryOne(&_user, sql)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, _user)

	}

	return util.Response(c, 1, "", "", &_user)
}

func getList(c echo.Context) error {
	var cmp = "cmp" + c.Param("cmp") + "."
	var _users []model.User
	sql := `SELECT id, name, email from ` + cmp + `users`
	_db := c.Get("db").(*pg.DB)
	_, err := _db.Query(&_users, sql)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, _users)
	}

	return c.JSON(http.StatusOK, _users)
}
