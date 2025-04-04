package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	pkgToken "github.com/hug58/users-api/pkg/token"
	pkgUser "github.com/hug58/users-api/pkg/user"
	"github.com/hug58/users-api/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type UserRouter struct {
	Repository      pkgUser.Repository
	RepositoryToken pkgToken.Repository
	RedisClient     *redis.Client // Inyectar el cliente Redis
}

func (ur *UserRouter) getUsers(c echo.Context) error {
	println("GET users")

	users, err := ur.Repository.GetAll(c.Request().Context())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, utils.Message{
			Msg:    "Users not found",
			Status: http.StatusNotFound,
		})
	}

	return c.JSON(http.StatusOK, users)
}

func (ur *UserRouter) getUserByID(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Message{
			Msg:    "InvalidID",
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})

	}

	user, err := ur.Repository.GetOne(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Message{
			Msg:    "Failed Get User",
			Error:  err.Error(),
			Status: http.StatusFound,
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (ur *UserRouter) CreateUser(c echo.Context) error {
	var user *pkgUser.User

	println("POST users")

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

	return c.JSON(http.StatusOK, user)
}

func (ur *UserRouter) Login(c echo.Context) error {
	var (
		user  *pkgUser.User
		login utils.Login
	)

	fmt.Sprintln(ur.RedisClient.Ping(context.TODO()))

	// 2. Decodificar el cuerpo de la petición
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		log.Println(err)

		// Log en Redis del fallo
		ur.RedisClient.LPush(c.Request().Context(), "login_errors", fmt.Sprintf("Failed decode at %v: %v", time.Now(), err))

		return c.JSON(http.StatusBadRequest, utils.Message{
			Msg:    "BadRequest",
			Status: http.StatusBadRequest,
		})
	}
	defer c.Request().Body.Close()

	// 3. Registrar intento de login en Redis

	// 4. Procesar login
	user, err := ur.Repository.Login(c.Request().Context(), user)
	if err != nil {
		// Log del fallo en Redis
		ur.RedisClient.LPush(c.Request().Context(), "login_failures", fmt.Sprintf("Failed login for %s at %v: %v", user.Email, time.Now(), err))

		return c.JSON(http.StatusNotFound, utils.Message{
			Msg:    "Login Failed",
			Status: http.StatusNotFound,
			Error:  err.Error(),
		})
	}

	ur.RedisClient.LPush(c.Request().Context(), "login_attempts", fmt.Sprintf("User %s attempting login at %v", user.Email, time.Now()))

	login = utils.Login{User: user}
	login.GenerarToken()

	if err := ur.RepositoryToken.Create(c.Request().Context(), user.ID, login.AccessToken); err != nil {
		ur.RedisClient.LPush(c.Request().Context(), "token_errors", fmt.Sprintf("Error saving token for %s at %v: %v", user.Email, time.Now(), err))

		return c.JSON(http.StatusConflict, utils.Message{
			Msg:    "Error save token",
			Status: http.StatusConflict,
			Error:  err.Error(),
		})
	}

	// 7. Log exitoso en Redis
	ur.RedisClient.LPush(c.Request().Context(), "login_success", fmt.Sprintf("User %s logged in at %v", user.Email, time.Now()))

	// 8. Métricas (opcional)
	ur.RedisClient.Incr(c.Request().Context(), "login_count")
	ur.RedisClient.HIncrBy(c.Request().Context(), "user_logins", user.Email, 1)

	return c.JSON(http.StatusOK, login)
}

func (ur *UserRouter) UpdateUser(c echo.Context) error {
	var user *pkgUser.User

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Message{
			Msg:    "InvalidID",
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})

	}

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.Message{
			Msg:    "BadRequest",
			Status: http.StatusBadRequest,
		})
	}

	defer c.Request().Body.Close()

	if _, err := ur.Repository.Update(c.Request().Context(), uint(id), *user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, utils.Message{
			Msg:    "Error create user",
			Status: http.StatusNotFound,
			Error:  err.Error(),
		})
	}

	if _, err := ur.Repository.ChangePassword(c.Request().Context(), uint(id), user.Password); err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, utils.Message{
			Msg:    "Error create user",
			Status: http.StatusNotFound,
			Error:  err.Error(),
		})
	}

	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

func (ur *UserRouter) DeleteUser(c echo.Context) error {

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Message{
			Msg:    "InvalidID",
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})

	}

	if err := ur.RepositoryToken.DeleteByUserId(c.Request().Context(), uint(id)); err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, utils.Message{
			Msg:    "Error delete token",
			Status: http.StatusNotFound,
			Error:  err.Error(),
		})
	}

	if err := ur.Repository.Delete(c.Request().Context(), uint(id)); err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, utils.Message{
			Msg:    "Error delete user",
			Status: http.StatusNotFound,
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Message{
		Msg:    "User Deleted succesfully",
		Status: http.StatusCreated,
	})
}
