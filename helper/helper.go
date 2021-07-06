package helper

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kagami/go-face"
)

const DataDir = "."

var (
	ModelDir   = filepath.Join(DataDir, "models")
	ImagesDir  = filepath.Join(DataDir, "images")
	EncodedDir = filepath.Join(DataDir, "encoded")
)


func contains(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

func GetSamplesCatsLabels(rec *face.Recognizer, excludes []string) (samples []face.Descriptor, cats []int32, labels []string) {
	files := OSReadDir(ImagesDir, ".jpg")
	for i, file := range files {
		label := strings.Split(file, ".")[0]
		if contains(excludes, label) {
			continue
		}
		encodedFile := filepath.Join(EncodedDir, file)
		_, err := os.Stat(encodedFile)
		if os.IsNotExist(err) {
			knownFace, _ := rec.RecognizeSingleFile(filepath.Join(ImagesDir, file))
			DumpToJson(EncodedDir, file, knownFace.Descriptor)
		}

		descriptor := DecodeFromJson(EncodedDir, file)
		samples = append(samples, descriptor)
		cats = append(cats, int32(i))
		labels = append(labels, label)
	}
	return
}

func SaveFile(dir string, filename string, content multipart.File) {
	destination, _ := os.Create(filepath.Join(dir, filename))
	defer destination.Close()
	io.Copy(destination, content)
}

func OSReadDir(root string, extension string) []string {
	var files []string
	f, _ := os.Open(root)
	defer f.Close()

	fileInfo, _ := f.Readdir(-1)

	for _, file := range fileInfo {
		if strings.Contains(file.Name(), extension) {
			files = append(files, file.Name())
		}
	}
	return files
}

func DumpToJson(dir string, filename string, object face.Descriptor) {
	file, _ := os.Create(filepath.Join(dir, filename))
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.Encode(object)

}

func DecodeFromJson(dir string, filename string) face.Descriptor {
	file, _ := os.Open(filepath.Join(dir, filename))
	defer file.Close()

	dec := json.NewDecoder(file)
	var descriptor face.Descriptor
	dec.Decode(&descriptor)

	return descriptor
}
