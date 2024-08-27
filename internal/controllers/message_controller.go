package controllers

import (
	"log"
	"net/http"
	"survaive/internal/managers"

	"github.com/labstack/echo/v4"
)

type Message struct {
	Channel string `json:"channel" validate:"required"`
	Value   string `json:"value" validate:"required"`
}

func SendMessageController(c echo.Context) (err error) {
	message := new(Message)

	if err = c.Bind(message); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	gameSession := managers.RetrieveOrCreate(message.Channel)

	log.Printf(gameSession.SessionId)
	managers.SendMessageToRoom(gameSession, message.Value)

	return c.JSON(http.StatusOK, message)
}
