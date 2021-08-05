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
		return c.JSON(http.StatusBadRequest, err)
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
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	modelFace := models.Face{
		Id:          id,
		Name:        name,
		Descriptors: []face.Descriptor{knownFaces[0].Descriptor},
	}

	res, err := modelFace.FindById(context.Background(), h.Coll, id)
	if len(res) == 0 {
		log.Println("id not found")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "id not found"})
	}
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	_, err = modelFace.PushDescriptor(context.Background(), h.Coll, id, knownFaces[0].Descriptor)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	log.Println("Sukses menambahkan descriptor wajah ", modelFace.Name, modelFace.Id)

	return c.JSON(http.StatusOK, response.Response{
		Detail: "Sukses menambahkan descriptor wajah " + modelFace.Name,
		Data:   modelFace,
	})
}
