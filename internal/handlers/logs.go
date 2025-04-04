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
			"error": "El parámetro 'key' es requerido",
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
				"error": "El parámetro 'start' debe ser un número entero",
			})
		}
	}

	if endStr != "" {
		if val, err := strconv.Atoi(endStr); err == nil {
			end = val
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "El parámetro 'end' debe ser un número entero",
			})
		}
	}

	elements, err := ur.RedisClient.LRange(ctx, listKey, int64(start), int64(end)).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Error al obtener logs de Redis: %v", err),
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
