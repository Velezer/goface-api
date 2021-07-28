package handler

import (
	"context"
	"goface-api/database"
	"goface-api/helper"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kagami/go-face"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type inputValidation struct {
	id   string                `validate:"required"`
	name string                `validate:"required"`
	file *multipart.FileHeader `validate:"required"`
}

func (h Handler) Register(c echo.Context) error {
	id := c.FormValue("id")
	name := c.FormValue("name")
	file, err := c.FormFile("file") //name=file in client html form
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": err.Error(),
		})
	}

	validate := validator.New()
	data := inputValidation{id: id, name: name, file: file}
	err = validate.Struct(data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": err.Error(),
		})
	}

	content, err := file.Open()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": err.Error(),
		})
	}
	folderSaved := filepath.Join(helper.ImagesDir, name)
	filename := time.Now().Local().String() + ".jpg"
	filename = strings.Replace(filename, ":", "", -1)
	helper.SaveFile(folderSaved, filename, content)

	knownFaces, err := h.Rec.RecognizeFile(filepath.Join(folderSaved, filename))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": err.Error(),
		})
	}
	if len(knownFaces) > 1 {
		os.Remove(filepath.Join(folderSaved, filename))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": "Terdeteksi lebih dari satu wajah",
		})
	}
	if len(knownFaces) == 0 {
		os.Remove(filepath.Join(folderSaved, filename))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": "Wajah tidak terdeteksi",
		})
	}

	var descriptors []face.Descriptor
	descriptors = append(descriptors, knownFaces[0].Descriptor)

	dataFace := database.Face{
		Id:          id,
		Name:        name,
		Descriptors: descriptors,
	}

	res, err := database.InsertOne(context.Background(), h.Coll, dataFace)
	if mongo.IsDuplicateKeyError(err) {
		_, err := database.PushDescriptor(context.Background(), h.Coll, id, knownFaces[0].Descriptor)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "fail",
				"detail": err.Error(),
			})
		}
	} else if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": err.Error(),
		})
	}

	log.Println(res)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":        "success",
		"data":          dataFace,
		"detail":        "Sukses menambahkan wajah",
	})
}
