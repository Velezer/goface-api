package handler

import (
	"bytes"
	"goface-api/database"
	"goface-api/helper"
	"goface-api/models"
	"goface-api/mymock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Register(t *testing.T) {
	faceData := models.Face{
		Id:   "2131256312",
		Name: "myname",
	}

	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	facerec, err := reco.RecognizeFile("../test/test_happy.jpg")
	assert.NoError(t, err)
	faceData.Descriptors = []face.Descriptor{facerec[0].Descriptor} // set descriptor

	// formfile
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile("file", "filename.jpg") // create empty formFile
	assert.NoError(t, err)

	content, err := os.Open("../test/test_happy.jpg")
	assert.NoError(t, err)

	_, err = io.Copy(formFile, content) // copy content to formFile
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())
	// end formfile

	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/api/face/register", body)
	req.Form = url.Values{} // set field,value of form
	req.Form.Set("id", faceData.Id)
	req.Form.Set("name", faceData.Name)

	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	repo := new(mymock.MockRepoFace)
	repo.On("InsertOne", faceData).Return(nil)

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo, Rec: reco}
	// Assertions
	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
