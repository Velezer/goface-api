package handler_test

import (
	"errors"
	"goface-api/database"
	"goface-api/helper"
	"goface-api/models"
	"goface-api/mymock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Find_NoFile(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/face/find", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	errHandler := h.Find(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}

func TestHandler_Find_HappyFile(t *testing.T) {
	body, writer, err := helper.CreateFormData("file", "../test_happy.jpg")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/face/find", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repoFace := new(mymock.MockRepoFace)

	h.DBRepo = &database.DBRepo{
		RepoFace: repoFace,
	}

	t.Run("Find Happy", func(t *testing.T) {
		repoFace.On("FindAll").Return([]models.Face{}, nil).Once()
		if assert.NoError(t, h.Find(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("Find FindAllErr", func(t *testing.T) {
		repoFace.On("FindAll").Return([]models.Face{}, errors.New("FindAllErr"))

		errHandler := h.Find(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})

}
func TestHandler_Find_NoFaceFile(t *testing.T) {
	body, writer, err := helper.CreateFormData("file", "../test_noface.png")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/face/find", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	t.Run("Find FindAllErr", func(t *testing.T) {
		errHandler := h.Find(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})

}