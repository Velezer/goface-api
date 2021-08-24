package mymiddleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SkipPreflight(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().Method)
		if c.Request().Method == http.MethodOptions {
			return nil
		}
		return nil

	}

}