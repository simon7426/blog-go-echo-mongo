package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Alive(c echo.Context) error {
	return c.JSON(http.StatusOK, &echo.Map{
		"message": "alive",
	})
}
