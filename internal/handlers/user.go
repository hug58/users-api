package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	pkgUser "github.com/hug58/users-api/pkg/user"
	"github.com/hug58/users-api/pkg/utils"
	"github.com/labstack/echo/v4"
)

type UserRouter struct {
	Repository pkgUser.Repository
}

func (ur *UserRouter) getUsers(c echo.Context) error {

	users, err := ur.Repository.GetAll(c.Request().Context())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, utils.Message{
			Msg:    "Users not found",
			Status: http.StatusNotFound,
		})
	}

	return c.JSON(http.StatusAccepted, users)
}

func (ur *UserRouter) CreateUser(c echo.Context) error {
	var user *pkgUser.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.Message{
			Msg:    "BadRequest",
			Status: http.StatusBadRequest,
		})
	}

	defer c.Request().Body.Close()

	if err := ur.Repository.Create(c.Request().Context(), user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, utils.Message{
			Msg:    "Error create user",
			Status: http.StatusNotFound,
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Message{
		Msg:    "User Created succesfully",
		Status: http.StatusCreated,
	})
}

func (ur *UserRouter) Login(c echo.Context) error {
	var (
		user  *pkgUser.User
		login utils.Login
	)

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.Message{
			Msg:    "BadRequest",
			Status: http.StatusBadRequest,
		})
	}

	defer c.Request().Body.Close()

	user, err := ur.Repository.Login(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.Message{
			Msg:    "Login failed",
			Status: http.StatusNotFound,
			Error:  err.Error(),
		})
	}

	login = utils.Login{User: user}
	login.GenerarToken()

	return c.JSON(http.StatusCreated, login)
}
