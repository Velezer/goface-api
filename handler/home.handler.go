package handler

import "github.com/labstack/echo/v4"

func (h Handler) Home(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"message": "goface-api up",
	})
}
