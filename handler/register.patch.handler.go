package handler

import (
	"context"
	"goface-api/helper"
	"goface-api/models"
	"goface-api/response"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kagami/go-face"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (h Handler) RegisterPatch(c echo.Context) error {
	id := c.FormValue("id")
	name := c.FormValue("name")

	validate := validator.New()
	err := validate.Struct(inputValidation{Id: id, Name: name})
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Error: err,
		})
	}

	file, err := c.FormFile("file") //name=file in client html form
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Error: err,
		})
	}

	content, err := file.Open()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Detail: err.Error(),
		})
	}

	folderSaved := filepath.Join(helper.ImagesDir, name+"_"+id)
	filename := time.Now().Local().String() + ".jpg"
	filename = strings.Replace(filename, ":", "_", -1)
	helper.SaveFile(folderSaved, filename, content)

	knownFaces, err := helper.RecognizeFile(h.Rec, folderSaved, filename)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Error: err,
		})
	}

	dataFace := models.Face{
		Id:          id,
		Name:        name,
		Descriptors: []face.Descriptor{knownFaces[0].Descriptor},
	}

	res, err := dataFace.PushDescriptor(context.Background(), h.Coll, id, knownFaces[0].Descriptor)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response.Response{
			StatusCode: http.StatusInternalServerError,
			Status: http.StatusText(http.StatusInternalServerError),
			Error: err,
		})
	}

	log.Println(res)

	return c.JSON(http.StatusOK, response.Response{
		StatusCode: http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Detail: "Sukses menambahkan wajah",
		Data: dataFace,
	})
}
