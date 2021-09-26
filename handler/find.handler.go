package handler

import (
	"goface-api/helper"
	"goface-api/response"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h Handler) Find(c echo.Context) error {
	start := time.Now()

	content, err := getFileContent(c, "file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = helper.SaveFile(helper.BaseDir, "unknown.jpg", content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	unknownFaces, code, err := helper.RecognizeFile(h.Rec, helper.BaseDir, "unknown.jpg")
	if err != nil {
		return echo.NewHTTPError(code, err.Error())
	}

	repo := h.DBRepo.RepoFace
	samples, err := repo.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	dSlice := response.DetectedSlice{}
	dSlice.FillSortDetectedFromDB(unknownFaces[0].Descriptor, samples, 0.25)

	elapsed := time.Since(start)
	log.Println("Detected:", dSlice, "in", elapsed.String())
	return c.JSON(http.StatusOK, response.Response{
		Detail:       "Sukses mendeteksi wajah",
		Data:         dSlice,
		ResponseTime: elapsed.String(),
	})

}
