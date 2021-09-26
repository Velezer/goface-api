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
	path := "../test/test_happy.jpg"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", path)
	assert.NoError(t, err)

	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
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

	h := Handler{Rec: reco,DBRepo: &dbRepo}

	// Assertions
	assert.NoError(t, h.Find(c))
}
