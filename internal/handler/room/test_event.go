package room

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"survaive/bus"
	"survaive/internal/handler"
	"survaive/internal/server"

	"github.com/labstack/echo/v4"
)

type startMessageDto struct {
	Id    string `param:"id" validate:"required"`
	Start bool   `form:"start" validate:"required"`
}

func TestEvent(c echo.Context, h *handler.Handler) error {
	dto := new(startMessageDto)
	if err := c.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	game := server.GetGameEngine().GetGame(dto.Id)
	game.State.Running = true

	data, jsonErr := json.Marshal(game.State)
	if jsonErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, jsonErr)
	}

	fmt.Println("set game state 1")

	// TODO: Run in goroutine??
	//h.Broker.Broadcast(dto.Id, "OK")
	err := bus.GetRedisClient().Set(context.Background(), dto.Id, data, 0).Err()

	fmt.Println("set game state 2")

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Start)
	}

	go game.Run()

	fmt.Println("set game state 2")

	return c.JSON(http.StatusOK, dto.Start)
}
