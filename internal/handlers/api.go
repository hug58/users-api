package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/hug58/users-api/internal/data"
	repos "github.com/hug58/users-api/internal/data/repositories"
	"github.com/redis/go-redis/v9"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Group) {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,                        // O tu dirección de Redis
		Password: os.Getenv("REDIS_PASSWORD"), // Si requiere password
		DB:       0,                           // DB a usar
	})
	// defer redisClient.Close()

	ur := UserRouter{
		Repository:      &repos.UserRepository{Data: data.DbInstance},
		RepositoryToken: &repos.TokenRepository{Data: data.DbInstance},
		RedisClient:     redisClient,
	}

	logRepo := LogsRouter{
		Repository:  &repos.UserRepository{Data: data.DbInstance},
		RedisClient: redisClient,
	}

	ping := ur.RedisClient.Ping(context.TODO())
	fmt.Println("REDIS RESPONSE PING: ")
	fmt.Print(ping)
	fmt.Println(addr)

	group := e.Group("/users")
	middleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
	})

	group.Add("GET", "/:id", ur.getUserByID, middleware)
	group.Add("PUT", "/:id", ur.UpdateUser, middleware)
	group.Add("DELETE", "/:id", ur.DeleteUser, middleware)

	group.Add("GET", "", ur.getUsers)
	group.Add("POST", "", ur.CreateUser)
	group.Add("POST", "/login", ur.Login)

	groupLogs := e.Group("/logs")
	groupLogs.Add("GET", "", logRepo.getLogsAll)
	groupLogs.Add("GET", "/test", logRepo.sendPing)

}
