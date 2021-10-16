package handler_test

import (
	"errors"
	"goface-api/helper"
	"goface-api/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Register(t *testing.T) {
	t.Run("Validation Error",func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api/face/register", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
	
		rec := httptest.NewRecorder()
	
		c := e.NewContext(req, rec)
	
		// Assertions
		errHandler := h.Register(c).(*echo.HTTPError)
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
		errHandler := h.Register(c).(*echo.HTTPError)
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
		errHandler := h.Register(c).(*echo.HTTPError)
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

		mockRepoFace.On("InsertOne", faceData).Return(nil)

		// Assertions
		if assert.NoError(t, h.Register(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

}

func TestHandler_RegisterPatch(t *testing.T) {
	t.Run("Validation Error",func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/face/register", nil)

		rec := httptest.NewRecorder()
	
		c := e.NewContext(req, rec)
	
		// Assertions
		errHandler := h.RegisterPatch(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, errHandler.Code)
	})

	t.Run("Happy File",func(t *testing.T) {

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
			mockRepoFace.On("FindById", faceData.Id).Return([]models.Face{faceData}, nil).Once()
			mockRepoFace.On("PushDescriptor", faceData.Id, faceData.Descriptors[0]).Return(nil).Once()
			if assert.NoError(t, h.RegisterPatch(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("PushErr", func(t *testing.T) {
			mockRepoFace.On("FindById", faceData.Id).Return([]models.Face{faceData}, nil).Once()
			mockRepoFace.On("PushDescriptor", faceData.Id, faceData.Descriptors[0]).Return(errors.New("PushErr")).Once()
	
			errHandler := h.RegisterPatch(c).(*echo.HTTPError)
			assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
		})
		t.Run("NotFound", func(t *testing.T) {
			mockRepoFace.On("FindById", faceData.Id).Return([]models.Face{}, nil).Once()
			mockRepoFace.On("PushDescriptor", faceData.Id, faceData.Descriptors[0]).Return(nil).Once()
	
			// Assertions
			errHandler := h.RegisterPatch(c).(*echo.HTTPError)
			assert.Equal(t, http.StatusNotFound, errHandler.Code)
		})
	})
	

}
