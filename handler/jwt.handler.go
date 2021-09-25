package handler

import (
	"goface-api/models"
	"goface-api/mymiddleware"
	"goface-api/response"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) JWTLogin(c echo.Context) error {
	adminData := models.Admin{}
	err := c.Bind(&adminData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	v := validator.New()
	err = v.Struct(adminData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	repo := h.DBRepo.RepoAdmin
	res, err := repo.FindOneByID(adminData.Username) // username equal _id
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(adminData.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := mymiddleware.Claims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func (h Handler) JWTRegister(c echo.Context) error {
	adminData := models.Admin{}
	err := c.Bind(&adminData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	v := validator.New()
	err = v.Struct(adminData)
	if err != nil {
		return err
	}

	hashed, err := h.Bcrypt.GenerateFromPassword([]byte(adminData.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	adminData.Password = string(hashed)

	repo := h.DBRepo.RepoAdmin
	err = repo.InsertOne(adminData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, response.Response{
		Detail: "Admin created",
	})
}
