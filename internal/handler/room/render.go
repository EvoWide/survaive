package room

import (
	"net/http"
	"survaive/internal/handler"

	"github.com/labstack/echo/v4"
)

type findRoomDto struct {
	Id string `param:"id" validate:"required"`
}

func Render(c echo.Context, h *handler.Handler) error {
	dto := new(findRoomDto)
	if err := c.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.Render(http.StatusOK, "room.html", map[string]interface{}{"roomId": dto.Id})
}
