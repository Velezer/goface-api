package handler

import (
	"context"
	"goface-api/models"
	"goface-api/response"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func (h Handler) JWTLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")


	modelAdmin := models.Admin{Username: username}
	res, err := modelAdmin.FindOneByUsername(context.Background(), h.DB)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
	}


	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
	}

	claims := &jwtCustomClaims{
		"Krefa",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("rahasia"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}