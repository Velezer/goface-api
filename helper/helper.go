package helper

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/Kagami/go-face"
)

const BaseDir = "."

var (
	ModelDir   = filepath.Join(BaseDir, "dat_models")
	EncodedDir = filepath.Join(BaseDir, "encoded")
)

func RecognizeFile(rec *face.Recognizer, folder string, filename string) ([]face.Face, int, error) {
	knownFaces, err := rec.RecognizeFile(filepath.Join(folder, filename))
	os.Remove(filepath.Join(folder, filename))
	if err != nil {
		return nil, 500, err
	}
	if len(knownFaces) > 1 {
		err = errors.New("detected more than one face, only support one face")
		return nil, 400, err
	}
	if len(knownFaces) == 0 {
		err = errors.New("no face detected")
		return nil, 400, err
	}
	return knownFaces, 200, err
}

func SaveFile(dir string, filename string, content multipart.File) error {
	os.MkdirAll(dir, os.ModeDir)

	destination, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, content)
	if err != nil {
		return err
	}

	return nil
}

// CreateFormData used for testing
func CreateFormData(fieldname string, path string) (*bytes.Buffer, *multipart.Writer, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile(fieldname, "filename.jpg") // create empty formFile
	if err != nil {
		return body, nil, err
	}

	content, err := os.Open(path)
	if err != nil {
		return body, nil, err
	}

	_, err = io.Copy(formFile, content) // copy content to formFile
	if err != nil {
		return body, nil, err
	}
	return body, writer, writer.Close()
}

