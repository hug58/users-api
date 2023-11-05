package handlers

import (
	"os"

	"github.com/hug58/users-api/internal/data"
	repos "github.com/hug58/users-api/internal/data/repositories"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Group) {

	ur := UserRouter{
		Repository: &repos.UserRepository{Data: data.DbInstance},
	}

	group := e.Group("/users")
	groupToken := e.Group(("/token"))

	group.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
	}))

	group.Add("GET", "/", ur.getUsers)

	groupToken.Add("POST", "/user/", ur.CreateUser)
	groupToken.Add("POST", "/", ur.Login)

}
