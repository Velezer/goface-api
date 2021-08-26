package handler

import (
	"context"
	"goface-api/models"
	"goface-api/response"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Delete(c echo.Context) error {
	id := c.Param("id")

	modelFace := models.Face{Id: id}
	err := modelFace.Delete(context.Background(), h.DB)
	if err != nil {
		log.Println("delete error:", err)
		return c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, response.Response{Detail: "deleted"})
}