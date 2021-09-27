package handler

import (
	"errors"
	"goface-api/database"
	"goface-api/helper"
	"goface-api/models"
	"goface-api/mymock"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Find_NoFile(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/face/find", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := Handler{}

	// Assertions
	errHandler := h.Find(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}

func TestHandler_Find_Happy(t *testing.T) {
	body, writer, err := helper.CreateFormData("file", "../test/test_happy.jpg")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/face/find", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	repoFace:=new(mymock.MockRepoFace)
	repoFace.On("FindAll").Return([]models.Face{}, nil)

	dbRepo := database.DBRepo{
		RepoFace: repoFace,
	}

	h := Handler{Rec: reco, DBRepo: &dbRepo}

	// Assertions
	if assert.NoError(t, h.Find(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestHandler_Find_JpegError(t *testing.T) {
	body, writer, err := helper.CreateFormData("file", "../test/test_noface.png")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/face/find", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	h := Handler{Rec: reco}


	errHandler := h.Find(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}
func TestHandler_Find_FindAllErr(t *testing.T) {
	body, writer, err := helper.CreateFormData("file", "../test/test_happy.jpg")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/face/find", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	repoFace:=new(mymock.MockRepoFace)
	repoFace.On("FindAll").Return([]models.Face{}, errors.New("FindAllErr"))

	dbRepo := database.DBRepo{
		RepoFace: repoFace,
	}

	h := Handler{Rec: reco, DBRepo: &dbRepo}

	// Assertions
	errHandler := h.Find(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}
