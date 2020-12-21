package integration

import (
	"fmt"
	"net/http"

	"github.com/masagatech/myinvest/module"
	"github.com/masagatech/myinvest/util"
	kiteconnect "github.com/zerodhatech/gokiteconnect"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo/v4"
)

// NewHTTP is ...
func NewHTTP(r *echo.Group) {
	ur := r.Group("/integration")
	ur.POST("/create_session", create_step1)
	ur.POST("/create_order", create_order)
	ur.GET("/back_url", back_url)
	ur.GET("/postback_url", postback_url)

}

func create_step1(c echo.Context) error {
	_uid := c.Param("uid")
	kc := module.GetZerodhaInstance(c)
	login_url := kc.GetLoginURL()
	login_url += login_url + "&q=" + _uid
	return util.Response(c, 1, "", "", bson.M{"login_url": login_url})

}

func back_url(c echo.Context) error {
	_config := c.Get("config").(*util.Configuration)
	request_token := c.QueryParam("request_token")
	redirect_params := c.QueryParam("redirect_params") // user id
	kc := module.GetZerodhaInstance(c)

	data, err := kc.GenerateSession(request_token, _config.Zerodha.ApiSecret)

	if err != nil {
		fmt.Printf("Error: %v", err)
		return c.Render(http.StatusOK, "integration-success.html", map[string]interface{}{
			"name": "",
		})
	}

	// Set access token
	kc.SetAccessToken(data.AccessToken)

	util.SessionSet(c, redirect_params, "access_token", request_token)

	accesstoken := util.SessionGet(c, redirect_params, "access_token")

	return c.Render(http.StatusOK, "integration-success.html", map[string]interface{}{
		"name": accesstoken,
	})
	// return c.HTML(http.StatusOK, "<html><head></head><body><center>  </center></body></html>")
	// return c.Redirect(http.StatusPermanentRedirect, "http://localhost:8080/#/settings/broker")
	// return util.Response(c, 1, "", request_token+" "+redirect_params, bson.M{})
}

func GetSesssion(c echo.Context, uid string) *kiteconnect.Client {
	accesstoken := util.SessionGet(c, uid, "access_token")
	kc := module.GetZerodhaInstance(c)
	kc.SetAccessToken(accesstoken.(string))
	return kc
}

func postback_url(c echo.Context) error {
	_uid := c.Param("uid")
	kc := module.GetZerodhaInstance(c)
	login_url := kc.GetLoginURL()
	login_url += login_url + "&q=" + _uid
	return util.Response(c, 1, "", "", bson.M{"login_url": login_url})
}

func create_order(c echo.Context) error {

	_uid := c.Param("uid")
	kc := module.GetZerodhaInstance(c)
	login_url := kc.GetLoginURL()
	login_url += login_url + "&q=" + _uid
	return util.Response(c, 1, "", "", bson.M{"login_url": login_url})
}
