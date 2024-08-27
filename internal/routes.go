package internal

import (
	"survaive/internal/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("channels/:id", controllers.MainPageController)
	e.GET("events/:id", controllers.EventController)
	e.POST("message", controllers.SendMessageController)
}
