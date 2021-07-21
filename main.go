package main

import (
	"goface-api/handler"
	"goface-api/helper"
	"log"
	"net/http"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	rec, err := face.NewRecognizer(helper.ModelDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	defer rec.Close()

	e := echo.New()

	e.POST("/register", handler.Register(rec))
	e.POST("/find", handler.Find(rec))
	e.GET("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "success",
		})
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	
	e.Logger.Fatal(e.Start(":8000"))

}
