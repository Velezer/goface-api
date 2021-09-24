package handler

import (
	"context"
	"goface-api/helper"
	"goface-api/models"
	"goface-api/response"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Kagami/go-face"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func prepInputValidation(c echo.Context) (inputValidation, error) {
	input := inputValidation{}
	c.Bind(&input)

	v := validator.New()
	err := v.Struct(input)
	if err != nil {
		return input, err
	}

	return input, nil
}

func prepModelFace(c echo.Context, h Handler, input inputValidation) (models.Face, error) {
	content, err := getFileContent(c, "file")
	if err != nil {
		return models.Face{}, err
	}

	filename := time.Now().Local().String() + ".jpg"
	filename = strings.Replace(filename, ":", "_", -1)
	helper.SaveFile(helper.DataDir, filename, content)

	knownFaces, err := helper.RecognizeFile(h.Rec, helper.DataDir, filename)
	if err != nil {
		return models.Face{}, err
	}

	modelFace := models.Face{
		Id:          input.Id,
		Name:        input.Name,
		Descriptors: []face.Descriptor{knownFaces[0].Descriptor},
	}

	return modelFace, nil
}

func (h Handler) Register(c echo.Context) error {
	input, err := prepInputValidation(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log.Println("register input ", input)
	modelFace, err := prepModelFace(c, h, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := modelFace.InsertOne(context.Background(), h.DB)
	if mongo.IsDuplicateKeyError(err) {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	log.Println("Insert data success ", res)

	return c.JSON(http.StatusCreated, response.Response{
		Detail: "Sukses menambahkan wajah",
		Data:   modelFace,
	})
}

func (h Handler) RegisterPatch(c echo.Context) error {
	input, err := prepInputValidation(c)
	if err != nil {
		return err
	}

	modelFace, err := prepModelFace(c, h, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, _ := modelFace.FindById(context.Background(), h.DB)
	if len(res) == 0 || err != nil {
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		log.Println("id " + modelFace.Id + " not found")
		return echo.NewHTTPError(http.StatusNotFound, "id "+modelFace.Id+" not found")
	}

	_, err = modelFace.PushDescriptor(context.Background(), h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	log.Println("Sukses menambahkan descriptor wajah ", modelFace.Name, modelFace.Id)

	return c.JSON(http.StatusOK, response.Response{
		Detail: "Sukses menambahkan descriptor wajah " + modelFace.Name,
		Data:   modelFace,
	})
}
