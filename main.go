package main

import (
	"github.com/hug58/users-api/internal/data"
	"github.com/hug58/users-api/internal/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	data.GetRedisClient()

	group := e.Group("/api/v1")
	handlers.Register(group)

	e.Start(":8080")
	data.DbInstance.Close()
}
