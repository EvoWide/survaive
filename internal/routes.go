package internal

import (
	"survaive/internal/handler"
	"survaive/internal/handler/room"
	"survaive/sse"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const (
	BROKER_REDIS_CHANNEL = "transport:sse"
)

func RegisterRoutes(e *echo.Echo, redis *redis.Client) {
	broker := sse.NewBroker()
	broker.AttachRedisBus(redis, BROKER_REDIS_CHANNEL)

	handler := handler.NewHandler(broker)

	e.GET("channels/:id", handler.Bind(room.Render))
	e.GET("channels/:id/stream", handler.Bind(room.Stream))

	// post a new message in the channel
	e.POST("channels/:id", handler.Bind(room.Broadcast))
}
