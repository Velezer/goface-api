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
	"go.mongodb.org/mongo-driver/mongo"
)

func (h Handler) Register(c echo.Context) error {
	id := c.FormValue("id")
	name := c.FormValue("name")

	validate := validator.New()
	input := inputValidation{Id: id, Name: name}
	err := validate.Struct(input)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      err,
		})
	}

	file, err := c.FormFile("file") //name=file in client html form
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      err,
		})
	}

	content, err := file.Open()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      err,
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
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      err,
		})
	}

	dataFace := models.Face{
		Id:          id,
		Name:        name,
		Descriptors: []face.Descriptor{knownFaces[0].Descriptor},
	}

	res, err := dataFace.InsertOne(context.Background(), h.Coll, dataFace)
	if mongo.IsDuplicateKeyError(err) {
		log.Println(err)
		return c.JSON(http.StatusConflict, response.Response{
			StatusCode: http.StatusConflict,
			Status:     http.StatusText(http.StatusConflict),
			Detail:     "_id exist in db",
			Error:      err,
		})
	} else if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      err,
		})
	}

	log.Println("Insert data success ", res)

	return c.JSON(http.StatusCreated, response.Response{
		StatusCode: http.StatusCreated,
		Status:     http.StatusText(http.StatusCreated),
		Detail:     "Sukses menambahkan wajah",
		Data:       dataFace,
	})
}
