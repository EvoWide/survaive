package main

import (
	"survaive/internal"

	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	e := echo.New()

	redis := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// set renderer
	templates := template.Must(template.ParseGlob("./templates/*.html"))
	e.Renderer = internal.NewTemplateRenderer(templates)

	// register new routes
	internal.RegisterRoutes(e, redis)

	e.Logger.Fatal(e.Start(":2000"))
}
