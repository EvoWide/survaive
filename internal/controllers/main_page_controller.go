package controllers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func MainPageController(c echo.Context) error {
	sessionId := c.Param("id")

	if sessionId == "" {
		c.JSON(http.StatusBadRequest, "")
		return errors.New("session id required")
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"channel": sessionId,
	})
}
