package handler

import (
	"context"
	"goface-api/database"
	"goface-api/helper"
	"log"
	"net/http"
	"path/filepath"
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

	unknownFaces, err := h.Rec.RecognizeFile(filepath.Join(helper.DataDir, "unknown.jpg"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": err.Error(),
		})
	}
	if len(unknownFaces) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": "Wajah tidak terdeteksi",
		})
	}
	if len(unknownFaces) > 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": "Terdeteksi lebih dari satu wajah",
		})
	}

	samples := database.FindAll(context.Background(), h.Coll)
	log.Println(len(samples))
	if len(samples) == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "fail",
			"detail": "Sampel wajah kosong",
		})
	}

	var dSlice helper.DetectedSlice
	dSlice.FillSortDetectedFromDB(unknownFaces[0].Descriptor, samples, 0.25)

	elapsed := time.Since(start)
	log.Println("Detected:", dSlice, "in", elapsed.String())
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":        "success",
		"data":          dSlice,
		"response_time": elapsed.String(),
	})

}
