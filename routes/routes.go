package routes

import (
	"goface-api/handler"
	"goface-api/mymiddleware"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, h handler.Handler) {

	e.GET("api/face/find", h.Find, mymiddleware.JWTAuth)
	e.POST("api/face/register", h.Register, mymiddleware.JWTAuth)
	e.PATCH("api/face/register", h.RegisterPatch, mymiddleware.JWTAuth)

	e.POST("jwt/login", h.JWTLogin)
	e.POST("jwt/register", h.JWTRegister)
}