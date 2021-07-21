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



func Contains(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

// 			encodedFile := filepath.Join(encFolder, file)
// 			_, err := os.Stat(encodedFile)
// 			if os.IsNotExist(err) {
// 				knownFace, _ := rec.RecognizeSingleFile(filepath.Join(dir, file))
// 				DumpToJson(encFolder, file, knownFace.Descriptor)
// 			}

func GetSamplesCatsLabels(rec *face.Recognizer) (samples []face.Descriptor, cats []int32, labels []string) {
	dirs, _ := OSReadDir(EncodedDir, ".jpg")
	for _, dir := range dirs {
		encFolder := filepath.Join(EncodedDir, dir)
		_, files := OSReadDir(encFolder, ".jpg")
		for i, file := range files {
			label := dir

			descriptor := DecodeFromJson(encFolder, file)
			samples = append(samples, descriptor)
			cats = append(cats, int32(i))
			labels = append(labels, label)
		}
	}

	return
}

func SaveFile(dir string, filename string, content multipart.File) {
	os.Mkdir(dir, os.ModeDir)
	destination, _ := os.Create(filepath.Join(dir, filename))
	defer destination.Close()
	io.Copy(destination, content)
}

func OSReadDir(root string, extension string) (dirs []string, files []string) {
	f, _ := os.Open(root)
	defer f.Close()

	fileInfo, _ := f.Readdir(-1)

	for _, file := range fileInfo {
		if strings.Contains(file.Name(), extension) {
			files = append(files, file.Name())
		}
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	return
}

func DumpToJson(dir string, filename string, object face.Descriptor) {
	os.Mkdir(dir, os.ModeDir)
	file, _ := os.Create(filepath.Join(dir, filename))
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(object)

}

func DecodeFromJson(dir string, filename string) (descriptor face.Descriptor) {
	file, _ := os.Open(filepath.Join(dir, filename))
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.Decode(&descriptor)

	return descriptor
}
