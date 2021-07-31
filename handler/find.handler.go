package handler

import (
	"context"
	"goface-api/models"
	"goface-api/helper"
	"goface-api/response"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h Handler) Find(c echo.Context) error {
	start := time.Now()

	file, err := c.FormFile("file") //name=file in client html form
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Error: err,
		})
	}
	content, err := file.Open()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Error: err,
		})
	}
	helper.SaveFile(helper.DataDir, "unknown.jpg", content)


	unknownFaces, err := helper.RecognizeFile(h.Rec, helper.DataDir, "unknown.jpg")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Error: err,
		})
	}
	

	samples := models.Face{}.FindAll(context.Background(), h.Coll)

	var dSlice response.DetectedSlice
	dSlice.FillSortDetectedFromDB(unknownFaces[0].Descriptor, samples, 0.25)

	elapsed := time.Since(start)
	log.Println("Detected:", dSlice, "in", elapsed.String())
	return c.JSON(http.StatusOK, response.Response{
		StatusCode: http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Detail: "Sukses mendeteksi wajah",
		Data: dSlice,
		ResponseTime: elapsed.String(),
	})

}
