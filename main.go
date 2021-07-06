package main

import (
	"goface-api/handler"
	"goface-api/helper"
	"log"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
)




func main() {
	rec, err := face.NewRecognizer(helper.ModelDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	defer rec.Close()
	

	e := echo.New()
	
	e.POST("/upload", handler.Upload(rec))
	e.POST("/find", handler.Find(rec))

	e.Logger.Fatal(e.Start(":8081"))

}
