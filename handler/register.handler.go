package handler

import (
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

func prepFaceData(c echo.Context, h Handler, input inputValidation) (models.Face, error) {
	content, err := getFileContent(c, "file")
	if err != nil {
		return models.Face{}, err
	}

	filename := time.Now().Local().String() + ".jpg"
	filename = strings.Replace(filename, ":", "_", -1)
	helper.SaveFile(helper.BaseDir, filename, content)

	knownFaces, err := helper.RecognizeFile(h.Rec, helper.BaseDir, filename)
	if err != nil {
		return models.Face{}, err
	}

	faceData := models.Face{
		Id:          input.Id,
		Name:        input.Name,
		Descriptors: []face.Descriptor{knownFaces[0].Descriptor},
	}

	return faceData, nil
}

func (h Handler) Register(c echo.Context) error {
	input, err := prepInputValidation(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log.Println("register input ", input)
	faceData, err := prepFaceData(c, h, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	repo := h.DBRepo.RepoFace
	err = repo.InsertOne(faceData)
	if mongo.IsDuplicateKeyError(err) {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	log.Println("Insert data success ")

	return c.JSON(http.StatusCreated, response.Response{
		Detail: "Sukses menambahkan wajah",
		Data:   faceData,
	})
}

func (h Handler) RegisterPatch(c echo.Context) error {
	input, err := prepInputValidation(c)
	if err != nil {
		return err
	}

	faceData, err := prepFaceData(c, h, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	repo := h.DBRepo.RepoFace
	res, _ := repo.FindById(faceData.Id)
	if len(res) == 0 || err != nil {
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		log.Println("id " + faceData.Id + " not found")
		return echo.NewHTTPError(http.StatusNotFound, "id "+faceData.Id+" not found")
	}

	_, err = repo.PushDescriptor(faceData.Id, faceData.Descriptors[0])
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	log.Println("Sukses menambahkan descriptor wajah ", faceData.Name, faceData.Id)

	return c.JSON(http.StatusOK, response.Response{
		Detail: "Sukses menambahkan descriptor wajah " + faceData.Name,
		Data:   faceData,
	})
}
