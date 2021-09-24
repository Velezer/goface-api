package handler

import (
	"context"
	"goface-api/models"
	"goface-api/mymiddleware"
	"goface-api/response"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) JWTLogin(c echo.Context) error {
	modelAdmin := models.Admin{}
	err := c.Bind(&modelAdmin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	v := validator.New()
	err = v.Struct(modelAdmin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := modelAdmin.FindOneByID(context.Background(), h.DB) // _id equal username
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(modelAdmin.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := mymiddleware.Claims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("rahasia"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func (h Handler) JWTRegister(c echo.Context) error {
	modelAdmin := models.Admin{}
	err := c.Bind(&modelAdmin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	v := validator.New()
	err = v.Struct(modelAdmin)
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(modelAdmin.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	modelAdmin.Password = string(hashed)

	_, err = modelAdmin.InsertOne(context.Background(), h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, response.Response{
		Detail: "Admin created",
	})
}
