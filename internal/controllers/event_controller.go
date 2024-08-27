package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"survaive/internal/managers"
	"survaive/internal/models"
)

func EventController(c echo.Context) error {
	sessionId := c.Param("id")

	if sessionId == "" {
		c.JSON(http.StatusBadRequest, "")
		return errors.New("session id required")
	}

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Ensure the response supports flushing
	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return c.String(http.StatusInternalServerError, "Streaming unsupported")
	}
	player := &models.Player{Id: c.RealIP(), Events: make(chan string)}
	gameSession := managers.RetrieveOrCreate(sessionId)

	gameSession.AddPlayerToSession(player)

	// Handle SSE in a separate goroutine to keep this non-blocking

	defer func() {
		gameSession.RemovePlayerFromSession(player.Id) // Clean up on disconnect
		close(player.Events)                           // Close the event channel
	}()

	for {
		select {
		case <-c.Request().Context().Done():
			log.Printf("SSE client disconnected, ip: %v", c.RealIP())
			return nil
		case eventMessage := <-player.Events:
			event := models.Event{
				Data: []byte("message: " + eventMessage),
			}
			if err := event.MarshalTo(w); err != nil {
				log.Printf("Error sending event: %v", err)
				return nil
			}
			flusher.Flush()
		}
	}
}
