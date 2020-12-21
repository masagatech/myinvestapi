package model

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	DbCtx struct {
		echo.Context
		*mongo.Database
	}
)
