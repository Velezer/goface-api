package main

import (
	"goface-api/database"
	"goface-api/handler"
	"goface-api/helper"
	"goface-api/routes"
	"log"
	"net/http"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	rec *face.Recognizer
)

var (
	err error
)

func initRecognizer() {
	rec, err = face.NewRecognizer(helper.ModelDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
}

func releaseResource() {
	rec.Close()
}

func main() {
	db := database.InitDB()
	initRecognizer()
	defer releaseResource()

	h := handler.Handler{Rec: rec, DB: db}
	e := echo.New()

	routes.Init(e, h)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodOptions, http.MethodConnect, http.MethodTrace, http.MethodPost, http.MethodGet, http.MethodHead, http.MethodPut},
	}))

	e.Logger.Fatal(e.Start(":8000"))

}
