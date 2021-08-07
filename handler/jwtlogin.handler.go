package handler

import (
	"context"
	"goface-api/models"
	"goface-api/response"
	"goface-api/mymiddleware"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) JWTLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")


	modelAdmin := models.Admin{Username: username}
	res, err := modelAdmin.FindOneByID(context.Background(), h.DB) // _id equal username
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
	}


	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
	}

	claims := mymiddleware.Claims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("rahasia"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}