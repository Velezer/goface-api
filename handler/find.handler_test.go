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
	"os"
	"path/filepath"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Find_ValidationError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/face/find", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := Handler{}

	// Assertions
	assert.Error(t, h.Find(c))
}

func TestHandler_Find_Happy(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile("file", "filename.jpg") // create empty formFile
	assert.NoError(t, err)

	content, err := os.Open("../test/test_happy.jpg")
	assert.NoError(t, err)

	_, err = io.Copy(formFile, content) // copy content to formFile
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/face/find", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	repo := new(mymock.MockRepoFace)
	repo.On("FindAll").Return([]models.Face{}, nil)

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{Rec: reco, DBRepo: &dbRepo}

	// Assertions
	if assert.NoError(t, h.Find(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
