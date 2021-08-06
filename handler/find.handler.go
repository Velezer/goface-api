package handler

import (
	"context"
	"goface-api/helper"
	"goface-api/models"
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
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	helper.SaveFile(helper.DataDir, "unknown.jpg", content)

	unknownFaces, err := helper.RecognizeFile(h.Rec, helper.DataDir, "unknown.jpg")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
	}

	modelFace := models.Face{}                                                                          

	samples, err := modelFace.FindAll(context.Background(), h.DB)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err)
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
