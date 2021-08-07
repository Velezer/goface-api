package routes

import (
	"goface-api/handler"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, h handler.Handler) {
	middlewareIsAuth := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("rahasia"),
	})
	
	e.GET("api/face/find", h.Find)
	e.POST("api/face/register", h.Register, middlewareIsAuth)
	e.PATCH("api/face/register", h.RegisterPatch, middlewareIsAuth)

	e.POST("jwt/login", h.JWTLogin)
	e.POST("jwt/register", h.JWTRegister)
}