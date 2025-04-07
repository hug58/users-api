package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	pkgUser "github.com/hug58/users-api/pkg/user"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type LogsRouter struct {
	Repository  pkgUser.Repository
	RedisClient *redis.Client // Inyectar el cliente Redis
}

func (ur *LogsRouter) getLogsAll(c echo.Context) error {
	ctx := context.Background()

	listKey := c.QueryParam("key")
	if listKey == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "'key' is required in perimeter",
		})
	}

	startStr := c.QueryParam("start")
	endStr := c.QueryParam("end")

	// Valores por defecto para obtener todos los elementos
	start := 0
	end := -1

	if startStr != "" {
		if val, err := strconv.Atoi(startStr); err == nil {
			start = val
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "The 'start' parameter must be an integer.”",
			})
		}
	}

	if endStr != "" {
		if val, err := strconv.Atoi(endStr); err == nil {
			end = val
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "The 'end' parameter must be an integer.”",
			})
		}
	}

	elements, err := ur.RedisClient.LRange(ctx, listKey, int64(start), int64(end)).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed in get redis logs: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"key":   listKey,
		"start": start,
		"end":   end,
		"count": len(elements),
		"logs":  elements,
	})
}

func (ur *LogsRouter) sendPing(c echo.Context) error {
	ctx := context.Background()

	// Check if Redis client is initialized
	if ur.RedisClient == nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"error":   "Redis client not initialized",
			"details": "The Redis client was not properly configured",
		})
	}

	// Execute Redis PING command
	pong, err := ur.RedisClient.Ping(ctx).Result()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"status":  "error",
			"error":   "Redis connection failed",
			"details": err.Error(),
		})
	}

	// Successful response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Redis connection active",
		"response": map[string]string{
			"command": "PING",
			"result":  pong,
		},
		"metadata": map[string]interface{}{
			"service":  "users-api",
			"endpoint": "/ping",
		},
	})
}
