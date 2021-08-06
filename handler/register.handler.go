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
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	content, err := getFileContent(c, "file")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	
	folderSaved := filepath.Join(helper.ImagesDir, name+"_"+id)
	filename := time.Now().Local().String() + ".jpg"
	filename = strings.Replace(filename, ":", "_", -1)
	helper.SaveFile(folderSaved, filename, content)

	knownFaces, err := helper.RecognizeFile(h.Rec, folderSaved, filename)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest,  response.Response{Error: err.Error()})
	}

	modelFace := models.Face{
		Id:          id,
		Name:        name,
		Descriptors: []face.Descriptor{knownFaces[0].Descriptor},
	}

	res, err := modelFace.InsertOne(context.Background(), h.Coll)
	if mongo.IsDuplicateKeyError(err) {
		log.Println(err)
		return c.JSON(http.StatusConflict, err)
	} else if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	log.Println("Insert data success ", res)

	return c.JSON(http.StatusCreated, response.Response{
		Detail: "Sukses menambahkan wajah",
		Data:   modelFace,
	})
}
