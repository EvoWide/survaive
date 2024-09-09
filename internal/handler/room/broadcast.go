package room

import (
	"net/http"
	"survaive/internal/handler"

	"github.com/labstack/echo/v4"
)

type broadcastMessageDto struct {
	Id string `param:"id" validate:"required"`

	Message string `form:"message" validate:"required"`
}

func Broadcast(c echo.Context, h *handler.Handler) error {
	dto := new(broadcastMessageDto)
	if err := c.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO: Run in goroutine??
	h.Broker.Broadcast(dto.Id, dto.Message)

	return c.JSON(http.StatusOK, dto.Message)
}
