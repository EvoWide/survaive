package handler

import (
	"survaive/sse"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Broker *sse.Broker
}

type FuncHandler func(c echo.Context, h *Handler) error

func NewHandler(broker *sse.Broker) *Handler {
	return &Handler{
		Broker: broker,
	}
}

func (h *Handler) Bind(handler FuncHandler) func(c echo.Context) error {
	return func(c echo.Context) error {
		return handler(c, h)
	}
}
