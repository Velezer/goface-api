package handler_test

import (
	"errors"
	"goface-api/database"
	"goface-api/helper"
	"goface-api/models"
	"goface-api/mymock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)



func TestHandler_Register_Happy(t *testing.T) {

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

	repo := new(mymock.MockRepoFace)
	repo.On("InsertOne", faceData).Return(nil)

	h.DBRepo = &database.DBRepo{
		RepoFace: repo,
	}

	// Assertions
	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
func TestHandler_Register_JpegErr(t *testing.T) {

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
}
func TestHandler_Register_NoFile(t *testing.T) {

	// formfile
	body, writer, err := helper.CreateFormData("error", "../test_noface.png")
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
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}
func TestHandler_Register_ValidationErr(t *testing.T) {

	req := httptest.NewRequest(http.MethodPost, "/api/face/register", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// Assertions
	errHandler := h.Register(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}

func TestHandler_RegisterPatch_Happy(t *testing.T) {

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

	repo := new(mymock.MockRepoFace)
	repo.On("FindById", faceData.Id).Return([]models.Face{faceData}, nil)
	repo.On("PushDescriptor", faceData.Id, faceData.Descriptors[0]).Return(nil)

	h.DBRepo = &database.DBRepo{
		RepoFace: repo,
	}

	// Assertions
	if assert.NoError(t, h.RegisterPatch(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestHandler_RegisterPatch_PushErr(t *testing.T) {
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

	repo := new(mymock.MockRepoFace)
	repo.On("FindById", faceData.Id).Return([]models.Face{faceData}, nil)
	repo.On("PushDescriptor", faceData.Id, faceData.Descriptors[0]).Return(errors.New("PushErr"))

	h.DBRepo = &database.DBRepo{
		RepoFace: repo,
	}
	// Assertions
	errHandler := h.RegisterPatch(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)

}

func TestHandler_RegisterPatch_NotFound(t *testing.T) {

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

	repo := new(mymock.MockRepoFace)
	repo.On("FindById", faceData.Id).Return([]models.Face{}, nil)
	repo.On("PushDescriptor", faceData.Id, faceData.Descriptors[0]).Return(nil)

	h.DBRepo = &database.DBRepo{
		RepoFace: repo,
	}
	// Assertions
	errHandler := h.RegisterPatch(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusNotFound, errHandler.Code)
}
func TestHandler_RegisterPatch_ValidationErr(t *testing.T) {

	req := httptest.NewRequest(http.MethodPut, "/api/face/register", nil)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// Assertions
	errHandler := h.RegisterPatch(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}
