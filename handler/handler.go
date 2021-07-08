package handler

import (
	"goface-api/helper"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
)

func Register(rec *face.Recognizer) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		name := c.FormValue("name")
		filename := name+".jpg"

		file, _ := c.FormFile("file") //name=file in client html form
		content, _ := file.Open()
		helper.SaveFile(helper.ImagesDir, filename, content)

		knownFaces, _ := rec.RecognizeFile(filepath.Join(helper.ImagesDir, filename))
		if len(knownFaces) > 1 {
			os.Remove(filepath.Join(helper.ImagesDir, filename))
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "fail",
				"message": "Detected more than one faces",
			})
		}
		if len(knownFaces) < 1 {
			os.Remove(filepath.Join(helper.ImagesDir, filename))
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "fail",
				"message": "No face detected",
			})
		}
		helper.DumpToJson(helper.EncodedDir, filename, knownFaces[0].Descriptor)

		elapsed := time.Since(start)
		return c.JSON(http.StatusOK,  map[string]string{
			"status": "success",
			"message": "File "+file.Filename+" uploaded",
			"response_time": elapsed.String(),
		})
	}
}

func Find(rec *face.Recognizer) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		
		mapExcludes, _:= c.FormParams()
		excludes := mapExcludes["excludes"]
		
		samples, cats, labels := helper.GetSamplesCatsLabels(rec, excludes)
		rec.SetSamples(samples, cats)

		file, _ := c.FormFile("file") //name=file in client html form
		content, _ := file.Open()
		helper.SaveFile(helper.DataDir, "unknown.jpg", content)

		unknownFaces, _ := rec.RecognizeFile(filepath.Join(helper.DataDir, "unknown.jpg"))
		if len(unknownFaces) > 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "fail",
				"message": "Detected more than one faces",
			})
		}
		if len(unknownFaces) < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "fail",
				"message": "No face detected",
			})
		}
		catID := rec.ClassifyThreshold(unknownFaces[0].Descriptor, 0.4)

		var detected string
		if catID < 0 {
			detected = "unknown"
		} else {
			detected = labels[catID]
		}

		elapsed := time.Since(start)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success",
			"detected": detected,
			"excludes": excludes,
			"response_time": elapsed.String(),
		})

	}
}
