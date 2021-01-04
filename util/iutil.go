package util

import (
	"github.com/dgrijalva/jwt-go"
	redis "github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type TokenData struct {
	Id     string
	Expiry float64
}

func GetTokenData(c echo.Context) TokenData {
	user := c.Get("user")
	_tokenD := TokenData{}
	if user != nil {
		userToken := user.(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		_tokenD.Id = claims["id"].(string)
		_tokenD.Expiry = claims["exp"].(float64)
	}
	return _tokenD
}

func GetRedis(c echo.Context) *redis.Client {
	_db := c.Get("redis").(*redis.Client)
	return _db
}
