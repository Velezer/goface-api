package handler

import (
	"errors"
	"goface-api/database"
	"goface-api/models"
	"goface-api/mymock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Delete_Happy(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/api/face/:id", nil)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("23124112312")

	repo := new(mymock.MockRepoFace)
	repo.On("DeleteId", "23124112312").Return(nil)

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo}

	// Assertions
	if assert.NoError(t, h.Delete(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestHandler_Delete_DelError(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/api/face/:id", nil)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("23124112312")

	repo := new(mymock.MockRepoFace)
	repo.On("DeleteId", "23124112312").Return(errors.New("delete id error"))

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo}

	// Assertions
	assert.Error(t, h.Delete(c), "delete id error")
}

func TestHandler_FaceAll_Happy(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/face/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	faceData := models.Face{
		Id:          "4871847291721",
		Name:        "myname",
		Descriptors: []face.Descriptor{},
	}
	repo := new(mymock.MockRepoFace)
	repo.On("FindAll").Return([]models.Face{faceData}, nil)

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo}

	// Assertions
	if assert.NoError(t, h.FaceAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestHandler_FaceAll_FindAllErr(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/face/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	faceData := models.Face{
		Id:          "4871847291721",
		Name:        "myname",
		Descriptors: []face.Descriptor{},
	}
	repo := new(mymock.MockRepoFace)
	repo.On("FindAll").Return([]models.Face{faceData}, errors.New("FindAll err"))

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo}

	// Assertions
	assert.Error(t, h.FaceAll(c), "FindAll err")
}

func TestHandler_FaceId_Happy(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/face/:id", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("13213213")

	faceData := models.Face{
		Id:          "4871847291721",
		Name:        "myname",
		Descriptors: []face.Descriptor{},
	}
	repo := new(mymock.MockRepoFace)
	repo.On("FindById", "13213213").Return([]models.Face{faceData}, nil)

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo}

	// Assertions
	if assert.NoError(t, h.FaceId(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestHandler_FaceId_FindByIdErr(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/face/:id", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("13213213")

	faceData := models.Face{
		Id:          "4871847291721",
		Name:        "myname",
		Descriptors: []face.Descriptor{},
	}
	repo := new(mymock.MockRepoFace)
	repo.On("FindById", "13213213").Return([]models.Face{faceData}, errors.New("FindByIdErr"))

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo}

	// Assertions
	assert.Error(t, h.FaceId(c), "FindByIdErr")
}
