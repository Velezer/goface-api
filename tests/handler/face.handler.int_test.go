package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Integration_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	req := httptest.NewRequest(http.MethodDelete, "/api/face/:id", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(faceData.Id)

	t.Run("Happy", func(t *testing.T) {
		if assert.NoError(t, hInt.Delete(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

}

func TestHandler_Integration_FaceAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	req := httptest.NewRequest(http.MethodGet, "/api/face/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	t.Run("Happy", func(t *testing.T) {
		if assert.NoError(t, hInt.FaceAll(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

}

func TestHandler_Integration_FaceId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	req := httptest.NewRequest(http.MethodGet, "/api/face/:id", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("13213213")

	t.Run("Happy", func(t *testing.T) {
		if assert.NoError(t, hInt.FaceId(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

}
