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
		log.Println("file error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	helper.SaveFile(helper.DataDir, "unknown.jpg", content)

	unknownFaces, err := helper.RecognizeFile(h.Rec, helper.DataDir, "unknown.jpg")
	if err != nil {
		log.Println("recognize file error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	repo := h.DBRepo.RepoFace
	samples, err := repo.FindAll()
	if err != nil {
		log.Println("db error", err)
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
