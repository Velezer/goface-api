package handler

import (
	"errors"
	"goface-api/database"
	"goface-api/helper"
	"goface-api/models"
	"goface-api/mymock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var faceData models.Face = models.Face{
	Id:   "2131256312",
	Name: "myname",
}

func TestHandler_Register_Happy(t *testing.T) {
	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	facerec, err := reco.RecognizeFile("../test/test_happy.jpg")
	assert.NoError(t, err)
	faceData.Descriptors = []face.Descriptor{facerec[0].Descriptor} // set descriptor

	// formfile
	body, writer, err := helper.CreateFormData("file", "../test/test_happy.jpg")
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

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo, Rec: reco}
	// Assertions
	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
func TestHandler_Register_JpegErr(t *testing.T) {
	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	// formfile
	body, writer, err := helper.CreateFormData("file", "../test/test_noface.png")
	assert.NoError(t, err)
	// end formfile


	req := httptest.NewRequest(http.MethodPost, "/api/face/register", body)
	req.Form = url.Values{} // set field,value of form
	req.Form.Set("id", faceData.Id)
	req.Form.Set("name", faceData.Name)

	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	h := Handler{Rec: reco}
	// Assertions
	errHandler := h.Register(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}
func TestHandler_Register_NoFile(t *testing.T) {
	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	// formfile
	body, writer, err := helper.CreateFormData("error", "../test/test_noface.png")
	assert.NoError(t, err)
	// end formfile


	req := httptest.NewRequest(http.MethodPost, "/api/face/register", body)
	req.Form = url.Values{} // set field,value of form
	req.Form.Set("id", faceData.Id)
	req.Form.Set("name", faceData.Name)

	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	h := Handler{Rec: reco}
	// Assertions
	errHandler := h.Register(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}
func TestHandler_Register_ValidationErr(t *testing.T) {

	req := httptest.NewRequest(http.MethodPost, "/api/face/register", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	h := Handler{}
	// Assertions
	errHandler := h.Register(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}

func TestHandler_RegisterPatch_Happy(t *testing.T) {
	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	facerec, err := reco.RecognizeFile("../test/test_happy.jpg")
	assert.NoError(t, err)
	faceData.Descriptors = []face.Descriptor{facerec[0].Descriptor} // set descriptor

	// formfile
	body, writer, err := helper.CreateFormData("file", "../test/test_happy.jpg")
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

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo, Rec: reco}
	// Assertions
	if assert.NoError(t, h.RegisterPatch(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
func TestHandler_RegisterPatch_PushErr(t *testing.T) {
	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	facerec, err := reco.RecognizeFile("../test/test_happy.jpg")
	assert.NoError(t, err)
	faceData.Descriptors = []face.Descriptor{facerec[0].Descriptor} // set descriptor

	// formfile
	body, writer, err := helper.CreateFormData("file", "../test/test_happy.jpg")
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

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo, Rec: reco}
	// Assertions
	errHandler := h.RegisterPatch(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)

}

func TestHandler_RegisterPatch_NotFound(t *testing.T) {
	reco, err := face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	assert.NoError(t, err)

	facerec, err := reco.RecognizeFile("../test/test_happy.jpg")
	assert.NoError(t, err)
	faceData.Descriptors = []face.Descriptor{facerec[0].Descriptor} // set descriptor

	// formfile
	body, writer, err := helper.CreateFormData("file", "../test/test_happy.jpg")
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

	dbRepo := database.DBRepo{
		RepoFace: repo,
	}

	h := Handler{DBRepo: &dbRepo, Rec: reco}
	// Assertions
	errHandler := h.RegisterPatch(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusNotFound, errHandler.Code)
}
func TestHandler_RegisterPatch_ValidationErr(t *testing.T) {

	req := httptest.NewRequest(http.MethodPut, "/api/face/register", nil)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	h := Handler{}
	// Assertions
	errHandler := h.RegisterPatch(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}
