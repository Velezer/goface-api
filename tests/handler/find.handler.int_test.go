package handler_test

import (
	"goface-api/helper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Integration_Find(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("No File", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/face/find", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		errHandler := hInt.Find(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, errHandler.Code)
	})


	t.Run("Happy file", func(t *testing.T) {
		body, writer, err := helper.CreateFormData("file", "../test_happy.jpg")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/face/find", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)


		t.Run("Happy", func(t *testing.T) {
			if assert.NoError(t, hInt.Find(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})

	})

}