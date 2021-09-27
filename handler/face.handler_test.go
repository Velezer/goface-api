package handler

import (
	"errors"
	"goface-api/database"
	"goface-api/models"
	"goface-api/mymock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Delete_Happy(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/face/:id", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("23124112312")

	repo := new(mymock.MockRepoFace)
	repo.On("DeleteId", "23124112312").Return(nil)

	h.DBRepo = &database.DBRepo{
		RepoFace: repo,
	}
	// Assertions
	if assert.NoError(t, h.Delete(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestHandler_Delete_DelError(t *testing.T) {

	req := httptest.NewRequest(http.MethodDelete, "/api/face/:id", nil)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("23124112312")

	repoFace := new(mymock.MockRepoFace)
	repoFace.On("DeleteId", "23124112312").Return(errors.New("delete id error"))

	h.DBRepo = &database.DBRepo{
		RepoFace: repoFace,
	}

	// Assertions
	errHandler := h.Delete(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}

func TestHandler_FaceAll_Happy(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/face/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	repoFace := new(mymock.MockRepoFace)
	repoFace.On("FindAll").Return([]models.Face{faceData}, nil)

	h.DBRepo = &database.DBRepo{
		RepoFace: repoFace,
	}

	// Assertions
	if assert.NoError(t, h.FaceAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestHandler_FaceAll_FindAllErr(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/face/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)


	repoFace := new(mymock.MockRepoFace)
	repoFace.On("FindAll").Return([]models.Face{faceData}, errors.New("FindAll err"))

	h.DBRepo = &database.DBRepo{
		RepoFace: repoFace,
	}

	// Assertions
	errHandler := h.FaceAll(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}

func TestHandler_FaceId_Happy(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/face/:id", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("13213213")


	repoFace := new(mymock.MockRepoFace)
	repoFace.On("FindById", "13213213").Return([]models.Face{faceData}, nil)

	h.DBRepo = &database.DBRepo{
		RepoFace: repoFace,
	}

	// Assertions
	if assert.NoError(t, h.FaceId(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestHandler_FaceId_FindByIdErr(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/face/:id", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("13213213")

	repoFace := new(mymock.MockRepoFace)
	repoFace.On("FindById", "13213213").Return([]models.Face{faceData}, errors.New("FindByIdErr"))

	h.DBRepo = &database.DBRepo{
		RepoFace: repoFace,
	}

	// Assertions
	errHandler := h.FaceId(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}
