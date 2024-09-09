package room

import (
	"net/http"
	"survaive/internal/handler"
	"survaive/sse"

	"github.com/labstack/echo/v4"
)

type createStreamDto struct {
	Id     string `param:"id" validate:"required"`
	UserId string `query:"userId" validate:"required"`
}

func Stream(c echo.Context, h *handler.Handler) error {
	dto := new(createStreamDto)
	if err := c.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Subscribe to the room
	dataChan := h.Broker.AddSubscriber(dto.Id, dto.UserId)

	// Unsub the room
	defer func() {
		h.Broker.RemoveSubscriber(dto.Id, dto.UserId)
	}()

	for {
		select {
		// connection lost
		case <-c.Request().Context().Done():
			return nil

		// send new data to the client
		case data := <-dataChan:
			event := sse.Event{
				Data: []byte(data),
			}

			if err := event.MarshalTo(w); err != nil {
				return err
			}
			w.Flush()
		}
	}
}
