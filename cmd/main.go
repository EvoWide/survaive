package main

import (
	"survaive/internal"

	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// set renderer
	templates := template.Must(template.ParseGlob("./templates/*.html"))
	e.Renderer = internal.NewTemplateRenderer(templates)

	// register new routes
	internal.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":2000"))
}
