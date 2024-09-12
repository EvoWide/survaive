package main

import (
	"survaive/bus"
	"survaive/internal"
	"survaive/internal/server"
	"survaive/sse"

	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	BROKER_REDIS_CHANNEL = "transport:sse"
)

func main() {
	e := echo.New()

	if connected, err := bus.IsRedisConnected(); !connected {
		e.Logger.Fatal("Failed to connect to Redis:", err)
	}

	broker := sse.NewBroker()
	broker.AttachRedisBus(bus.GetRedisClient(), BROKER_REDIS_CHANNEL)

	/**
	* Game Part
	**/
	server.InitGameEngine(broker)

	/**
	* Web Part
	**/

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// set renderer
	templates := template.Must(template.ParseGlob("./templates/*.html"))
	e.Renderer = internal.NewTemplateRenderer(templates)

	// register new routes
	internal.RegisterRoutes(e, broker)

	e.Logger.Fatal(e.Start(":2000"))

}
