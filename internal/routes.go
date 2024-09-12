package internal

import (
	"survaive/internal/handler"
	"survaive/internal/handler/room"
	"survaive/sse"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, broker *sse.Broker) {
	handler := handler.NewHandler(broker)

	e.GET("channels/:id", handler.Bind(room.Render))
	e.GET("channels/:id/stream", handler.Bind(room.Stream))

	// post a new message in the channel
	e.POST("channels/:id", handler.Bind(room.Broadcast))

	//start a game
	e.POST("channels/:id/event", handler.Bind(room.TestEvent))
}
