package menu

import (
	"fmt"

	"github.com/masagatech/myinvest/util"

	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
)

func NewHTTP(r *echo.Group) {
	ur := r.Group("/menu")
	ur.GET("", get)
}

func get(c echo.Context) error {
	var sysjson, _ = util.BindSysParams(c.Param("cmp"), "", "", "")

	var _params struct {
		Userid     string `json:"userid"`
		Usertype   string `json:"usertype"`
		Cmp        string `json:"cmp"`
		Dispensary string `json:"dispensary"`
		Type       string `json:"type"`
	}

	_params.Cmp = c.QueryParam("cmp")
	_params.Usertype = c.QueryParam("usertype")
	_params.Dispensary = c.QueryParam("dispensary")
	_params.Userid = c.QueryParam("userid")
	_params.Type = c.QueryParam("type")

	sql := `select func.fn_getmenubyuser('` + sysjson + `','` + util.ToJsonString(_params) + `','result');fetch all IN result;`

	var ret []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Url         string `json:"url"`
		Icon        string `json:"icon"`
		Isclickable string `json:"isclickable"`
		Action      string `json:"action"`
		Level       int    `json:"level"`
		Parentid    int    `json:"parentid"`
		Code        string `json:"code"`
		T           int    `json:"t"`
	}
	_db := c.Get("db").(*pg.DB)
	_, err := _db.Query(&ret, sql)

	if err != nil {
		fmt.Println(err)
		return util.CreateResponse(c, 0, ret, "")

	}

	return util.CreateResponse(c, 1, ret, "")
}
