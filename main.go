package main

import (
	"context"
	"goface-api/handler"
	"goface-api/helper"
	"log"
	"net/http"
	"time"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db   *mongo.Database
	coll *mongo.Collection
)

var (
	rec *face.Recognizer
)

var (
	err error
)

func initDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Panicln(err)
	}

	db = client.Database("db_kita")
	coll = db.Collection("coll_terserah")
}

func initRecognizer() {
	rec, err = face.NewRecognizer(helper.ModelDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
}

func releaseResource(){
	rec.Close()
}

func main() {
	initDB()
	initRecognizer()
	defer releaseResource()

	e := echo.New()

	h := handler.Handler{Rec: rec, Coll: coll}

	e.POST("/register", h.Register)
	e.POST("/find", h.Find)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.Logger.Fatal(e.Start(":8000"))

}
