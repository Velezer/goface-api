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
	helper.SaveFile(helper.DataDir, "unknown.jpg", content)


	unknownFaces, err := helper.RecognizeFile(h.Rec, helper.DataDir, "unknown.jpg")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": err.Error(),
		})
	}
	

	samples := models.Face{}.FindAll(context.Background(), h.Coll)

	var dSlice response.DetectedSlice
	dSlice.FillSortDetectedFromDB(unknownFaces[0].Descriptor, samples, 0.25)

	elapsed := time.Since(start)
	log.Println("Detected:", dSlice, "in", elapsed.String())
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":        "success",
		"data":          dSlice,
		"response_time": elapsed.String(),
	})

}
