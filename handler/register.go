package handler

import (
	"goface-api/helper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (h Handler) Register(c echo.Context) error {
	start := time.Now()

	name := c.FormValue("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": "Mohon isi nama",
		})
	}
	filename := time.Now().Local().String() + ".jpg"
	filename = strings.Replace(filename, ":", "", -1)

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
	folderSaved := filepath.Join(helper.ImagesDir, name)
	helper.SaveFile(folderSaved, filename, content)

	knownFaces, _ := h.Rec.RecognizeFile(filepath.Join(folderSaved, filename))
	if len(knownFaces) > 1 {
		os.Remove(filepath.Join(folderSaved, filename))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": "Terdeteksi lebih dari satu wajah",
		})
	}
	if len(knownFaces) == 0 {
		os.Remove(filepath.Join(folderSaved, filename))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "fail",
			"detail": "Wajah tidak terdeteksi",
		})
	}
	encFolderSaved := filepath.Join(helper.EncodedDir, name)
	helper.DumpToJson(encFolderSaved, filename, knownFaces[0].Descriptor)

	log.Println("File " + filename + " uploaded")

	elapsed := time.Since(start)
	return c.JSON(http.StatusOK, map[string]string{
		"status":        "success",
		"name":          name,
		"detail":        "Sukses menambahkan wajah",
		"response_time": elapsed.String(),
	})
}
