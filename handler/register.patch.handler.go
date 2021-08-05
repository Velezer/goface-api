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
			Error: err,
		})
	}

	content, err := getFileContent(c, "file")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			Error: err,
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
			Error: err,
		})
	}

	modelFace := models.Face{
		Id:          id,
		Name:        name,
		Descriptors: []face.Descriptor{knownFaces[0].Descriptor},
	}

	res, err := modelFace.PushDescriptor(context.Background(), h.Coll, id, knownFaces[0].Descriptor)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response.Response{
			Error: err,
		})
	}

	log.Println(res)

	return c.JSON(http.StatusOK, response.Response{
		Detail: "Sukses menambahkan descriptor wajah " + modelFace.Name,
		Data:   modelFace,
	})
}
