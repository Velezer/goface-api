package handler

import (
	"context"
	"goface-api/models"
	"goface-api/mymiddleware"
	"goface-api/response"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) JWTLogin(c echo.Context) error {

	data := echo.Map{}

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
	}

	username := data["username"].(string)
	password := data["password"].(string)

	log.Println(username, "login jwt")
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
