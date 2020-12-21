package api

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/masagatech/myinvest/api/auth"
	"github.com/masagatech/myinvest/api/broker"
	"github.com/masagatech/myinvest/api/instrument"
	"github.com/masagatech/myinvest/api/integration"
	"github.com/masagatech/myinvest/api/masters"
	"github.com/masagatech/myinvest/api/menu"
	"github.com/masagatech/myinvest/api/strategy"
	"github.com/masagatech/myinvest/api/user"
	"github.com/masagatech/myinvest/util"
)

// RegisterRoute register route
func RegisterRoute(cfg *util.Configuration, e *echo.Echo) {

	api := e.Group("/api")
	o := api.Group("/o")
	oroute(o)

	midware := api.Group("/v1")
	midware.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(cfg.JWT.Secret),

		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "/auth") {
				return true
			}
			return false
		},
	}))

	midlewareroute(midware)

}

func oroute(cmpgrp *echo.Group) {
	integration.NewHTTP(cmpgrp)

}

func midlewareroute(midlew *echo.Group) {
	auth.NewHTTP(midlew)
	user.NewHTTP(midlew)
	menu.NewHTTP(midlew)
	instrument.NewHTTP(midlew)
	masters.NewHTTP(midlew)
	strategy.NewHTTP(midlew)
	broker.NewHTTP(midlew)
}
