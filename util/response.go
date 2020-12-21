package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Response(c echo.Context, status int, errcode string, message string, data interface{}) error {
	m := new(ResponseModel)
	m.ErrorCode = errcode
	m.Message = message
	m.ResultKey = status
	m.ResultValue = data

	return c.JSON(http.StatusOK, &m)
}
