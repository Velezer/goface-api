package handler

import (
	"goface-api/helper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
)

func Register(rec *face.Recognizer) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		name := c.FormValue("name")
		if name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status":  "fail",
				"message": "Name field required",
			})
		}
		filename := time.Now().Local().String() + ".jpg"
		filename = strings.Replace(filename,":", "", -1)

		file, _ := c.FormFile("file") //name=file in client html form
		content, _ := file.Open()

		folderSaved := filepath.Join(helper.ImagesDir, name)
		helper.SaveFile(folderSaved, filename, content)

		knownFaces, _ := rec.RecognizeFile(filepath.Join(folderSaved, filename))
		if len(knownFaces) > 1 {
			os.Remove(filepath.Join(folderSaved, filename))
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status":  "fail",
				"message": "Detected more than one faces",
			})
		}
		if len(knownFaces) < 1 {
			os.Remove(filepath.Join(folderSaved, filename))
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status":  "fail",
				"message": "No face detected",
			})
		}
		encFolderSaved := filepath.Join(helper.EncodedDir, name)
		helper.DumpToJson(encFolderSaved, filename, knownFaces[0].Descriptor)

		log.Println("File " + filename + " uploaded")

		log.Println("File " + filename + " uploaded")

		elapsed := time.Since(start)
		return c.JSON(http.StatusOK, map[string]string{
			"status":        "success",
			"message":       "Database wajah "+name+" ditambahkan",
			"response_time": elapsed.String(),
		})
	}
}

func Find(rec *face.Recognizer) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		formParams, _ := c.FormParams()
		excludes := formParams["excludes"]

		samples, cats, labels := helper.GetSamplesCatsLabels(rec, excludes)
		rec.SetSamples(samples, cats)

		log.Println(labels)

		file, err := c.FormFile("file") //name=file in client html form
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  "fail",
				"message": err,
			})
		}
		content, err := file.Open()
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  "fail",
				"message": err,
			})
		}
		helper.SaveFile(helper.DataDir, "unknown.jpg", content)

		unknownFaces, _ := rec.RecognizeFile(filepath.Join(helper.DataDir, "unknown.jpg"))
		if len(unknownFaces) > 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status":  "fail",
				"message": "Detected more than one faces",
			})
		}
		if len(unknownFaces) < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status":  "fail",
				"message": "No face detected",
			})
		}
		catID := rec.ClassifyThreshold(unknownFaces[0].Descriptor, 0.39)

		var detected string
		if catID < 0 {
			detected = "unknown"
		} else {
			detected = labels[catID]
		}
		
		elapsed := time.Since(start)
		log.Println("Detected:", detected, "in", elapsed.String())
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":        "success",
			"detected":      detected,
			"excludes":      excludes,
			"response_time": elapsed.String(),
		})

	}
}
