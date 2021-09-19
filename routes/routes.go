package routes

import (
	"goface-api/handler"
	"goface-api/mymiddleware"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, h handler.Handler) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "see?",
		})
	})

	e.POST("api/face/find", h.Find, mymiddleware.JWTAuth)
	e.POST("api/face/register", h.Register, mymiddleware.JWTAuth)
	e.PUT("api/face/register", h.RegisterPatch, mymiddleware.JWTAuth)

	e.GET("api/face", h.FaceAll, mymiddleware.JWTAuth)
	e.GET("api/face/:id", h.FaceId, mymiddleware.JWTAuth)
	e.DELETE("api/face/:id", h.Delete, mymiddleware.JWTAuth)

	e.POST("jwt/login", h.JWTLogin)
	e.POST("jwt/register", h.JWTRegister)
}
