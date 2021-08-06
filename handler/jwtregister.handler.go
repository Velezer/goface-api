package handler

import (
	"context"
	"goface-api/models"
	"goface-api/response"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (h Handler) JWTRegister(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	modelAdmin := models.Admin{
		Username: username,
		Password: password,
	}

	validate := validator.New()
	err := validate.Struct(modelAdmin)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
	}
	
	_, err = modelAdmin.InsertOne(context.Background(), h.DB)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, response.Response{
		Detail: "Admin created",
	})
}
