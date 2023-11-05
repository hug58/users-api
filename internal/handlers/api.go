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
		Repository:      &repos.UserRepository{Data: data.DbInstance},
		RepositoryToken: &repos.TokenRepository{Data: data.DbInstance},
	}
	group := e.Group("/users")
	middleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
	})

	group.Add("GET", "/:id", ur.getUserByID, middleware)
	group.Add("PUT", "/:id", ur.UpdateUser, middleware)

	group.Add("GET", "", ur.getUsers)
	group.Add("POST", "", ur.CreateUser)
	group.Add("POST", "/login", ur.Login)

}
