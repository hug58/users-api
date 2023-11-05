package main

import (
	"github.com/hug58/users-api/internal/data"
	"github.com/hug58/users-api/internal/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	group := e.Group("/api")

	handlers.Register(group)

	// Iniciar el servidor
	e.Start(":8080")

	data.DbInstance.Close()
}
