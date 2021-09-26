package handler

import (
	"goface-api/response"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Delete(c echo.Context) error {
	id := c.Param("id")

	repo := h.DBRepo.RepoFace
	err := repo.DeleteId(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	log.Println("delete success")

	return c.JSON(http.StatusOK, response.Response{Detail: "deleted"})
}

func (h Handler) FaceAll(c echo.Context) error {
	repo := h.DBRepo.RepoFace
	faces, err := repo.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resFaces := []response.FaceLenDesc{}
	for _, face := range faces {
		resFaces = append(resFaces, response.FaceLenDesc{Id: face.Id, Name: face.Name, Descriptors: len(face.Descriptors)})
	}
	log.Println("FindAll success!")
	return c.JSON(http.StatusOK, response.Response{Data: resFaces})
}

func (h Handler) FaceId(c echo.Context) error {
	id := c.Param("id")
	repo := h.DBRepo.RepoFace
	faces, err := repo.FindById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resFaces := []response.FaceLenDesc{}
	for _, face := range faces {
		resFaces = append(resFaces, response.FaceLenDesc{Id: face.Id, Name: face.Name, Descriptors: len(face.Descriptors)})
	}
	log.Println("FindById success!:", resFaces)
	return c.JSON(http.StatusOK, response.Response{Data: resFaces})
}
