package util

import (
	"github.com/dgrijalva/jwt-go"
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
