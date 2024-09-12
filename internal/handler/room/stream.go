package room

import (
	"fmt"
	"net/http"
	"survaive/internal/handler"
	"survaive/internal/server"
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

	//Send join message
	joinPayload := fmt.Sprintf("%s:join", dto.UserId)
	h.Broker.Broadcast(dto.Id, joinPayload)

	// Subscribe to the room
	dataChan := h.Broker.AddSubscriber(dto.Id, dto.UserId)

	gameEngine := server.GetGameEngine()
	game := gameEngine.GetOrCreateGame(dto.Id)

	fmt.Println(game.State)

	if game.State.Running {
		fmt.Println("game already started, you can't join it")
	}

	// Unsub the room
	defer func() {
		h.Broker.RemoveSubscriber(dto.Id, dto.UserId)
		if game.State.Running {
			game.Stop()
		}
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
