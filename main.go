package main

import (
	"goface-api/database"
	"goface-api/handler"
	"goface-api/helper"
	"goface-api/routes"
	"log"
	"net/http"
	"os"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// setup
	dbrepo, err := database.InitDB()
	if err != nil {
		log.Fatalf("Can't init databse: %v", err)
	}
	if dbrepo == nil {
		log.Fatalln("dbrepo is nil")
	}

	rec, err := face.NewRecognizer(helper.ModelDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	defer rec.Close()
	if rec == nil {
		log.Fatalln("rec is nil")
	}
	log.Println("rec: ", rec)

	// end setup

	h := handler.Handler{Rec: rec, DBRepo: dbrepo, Bcrypt: handler.RealBcrypt{}}
	e := echo.New()

	// e.Debug = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)

		e.DefaultHTTPErrorHandler(err, c)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowOrigin},
		AllowMethods: []string{http.MethodOptions, http.MethodConnect, http.MethodPost, http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete},
	}))

	routes.Init(e, h)

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))

}
