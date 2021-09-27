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
	err := c.Bind(&input)
	if err != nil {
		return input, err
	}

	v := validator.New()
	err = v.Struct(input)
	if err != nil {
		return input, err
	}

	return input, nil
}

func prepFaceData(c echo.Context, h Handler, input inputValidation) (models.Face, int, error) {
	content, err := getFileContent(c, "file")
	if err != nil {
		return models.Face{}, http.StatusBadRequest, err
	}

	filename := time.Now().Local().String() + ".jpg"
	filename = strings.Replace(filename, ":", "_", -1)
	err = helper.SaveFile(helper.BaseDir, filename, content)
	if err != nil {
		return models.Face{}, http.StatusInternalServerError, err
	}
	
	knownFaces, code, err := helper.RecognizeFile(h.Rec, helper.BaseDir, filename)
	if err != nil {
		return models.Face{}, code, err
	}

	faceData := models.Face{
		Id:          input.Id,
		Name:        input.Name,
		Descriptors: []face.Descriptor{knownFaces[0].Descriptor},
	}

	return faceData, 1, nil
}

func (h Handler) Register(c echo.Context) error {
	input, err := prepInputValidation(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log.Println("register input ", input)
	faceData, code, err := prepFaceData(c, h, input)
	if err != nil {
		return echo.NewHTTPError(code, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	faceData, code, err := prepFaceData(c, h, input)
	if err != nil {
		return echo.NewHTTPError(code, err.Error())
	}
	repo := h.DBRepo.RepoFace
	res, err := repo.FindById(faceData.Id)
	if len(res) == 0 || err != nil {
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return echo.NewHTTPError(http.StatusNotFound, "id "+faceData.Id+" not found")
	}

	err = repo.PushDescriptor(faceData.Id, faceData.Descriptors[0])
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	log.Println("Sukses menambahkan descriptor wajah ", faceData.Name, faceData.Id)

	return c.JSON(http.StatusOK, response.Response{
		Detail: "Sukses menambahkan descriptor wajah " + faceData.Name,
		Data:   faceData,
	})
}
