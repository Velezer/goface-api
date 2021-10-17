package handler_test

import (
	"goface-api/helper"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Integration_Register(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("Validation Error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api/face/register", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		// Assertions
		errHandler := hInt.Register(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, errHandler.Code)
	})

	t.Run("No File", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api/face/register", nil)
		req.Form = url.Values{} // set field,value of form
		req.Form.Set("id", faceData.Id)
		req.Form.Set("name", faceData.Name)

		req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		// Assertions
		errHandler := hInt.Register(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, errHandler.Code)
	})

	t.Run("No Face", func(t *testing.T) {

		// formfile
		body, writer, err := helper.CreateFormData("file", "../test_noface.png")
		assert.NoError(t, err)
		// end formfile

		req := httptest.NewRequest(http.MethodPost, "/api/face/register", body)
		req.Form = url.Values{} // set field,value of form
		req.Form.Set("id", faceData.Id)
		req.Form.Set("name", faceData.Name)

		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		// Assertions
		errHandler := hInt.Register(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})

	t.Run("Happy File", func(t *testing.T) {

		facerec, err := reco.RecognizeFile("../test_happy.jpg")
		assert.NoError(t, err)
		faceData.Descriptors = []face.Descriptor{facerec[0].Descriptor} // set descriptor

		// formfile
		body, writer, err := helper.CreateFormData("file", "../test_happy.jpg")
		assert.NoError(t, err)
		// end formfile

		req := httptest.NewRequest(http.MethodPost, "/api/face/register", body)
		req.Form = url.Values{} // set field,value of form
		req.Form.Set("id", faceData.Id)
		req.Form.Set("name", faceData.Name)

		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, hInt.Register(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

}

func TestHandler_Integration_RegisterPatch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("Validation Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/face/register", nil)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		// Assertions
		errHandler := hInt.RegisterPatch(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, errHandler.Code)
	})

	t.Run("Happy File", func(t *testing.T) {

		facerec, err := reco.RecognizeFile("../test_happy.jpg")
		assert.NoError(t, err)
		faceData.Descriptors = []face.Descriptor{facerec[0].Descriptor} // set descriptor

		// formfile
		body, writer, err := helper.CreateFormData("file", "../test_happy.jpg")
		assert.NoError(t, err)
		// end formfile

		req := httptest.NewRequest(http.MethodPut, "/api/face/register", body)
		req.Form = url.Values{} // set field,value of form
		req.Form.Set("id", faceData.Id)
		req.Form.Set("name", faceData.Name)

		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		t.Run("Happy", func(t *testing.T) {
			if assert.NoError(t, hInt.RegisterPatch(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
	})

}
