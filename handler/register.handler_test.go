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

func TestHandler_Register(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// idw,err:=writer.CreateFormField("id")
	// assert.NoError(t, err)
	// idw.Write([]byte("21312312"))

	// namew,err:=writer.CreateFormField("name")
	// assert.NoError(t, err)
	// namew.Write([]byte("myname"))

	formFile, err := writer.CreateFormFile("file", "filename.jpg") // create empty formFile
	assert.NoError(t, err)
	content, err := os.Open("../test/test_happy.jpg")
	assert.NoError(t, err)

	_, err = io.Copy(formFile, content) // copy content to formFile
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/api/face/register", body)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	repo := new(mymock.MockRepoFace)
	repo.On("InsertOne", models.Face{}).Return(nil)

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	h := Handler{DBRepo: &dbRepo, Rec: reco}

	// Assertions
	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
