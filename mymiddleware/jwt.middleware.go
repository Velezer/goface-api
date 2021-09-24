package mymiddleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
)

var Claims = &jwt.StandardClaims{
	ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
}

var JWTAuth = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(os.Getenv("JWT_KEY")),
	Claims:     Claims,
})
