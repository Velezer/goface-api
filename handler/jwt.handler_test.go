package handler

import (
	"errors"
	"goface-api/database"
	"goface-api/models"
	"goface-api/mymock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHandler_JWTLogin_Happy(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`

	req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(adminData.Password), bcrypt.DefaultCost)
	adminData.Password = string(hashed)

	repo := new(mymock.MockRepoAdmin)
	repo.On("FindOneByID", "krefa").Return(adminData, nil)

	h.DBRepo = &database.DBRepo{
		RepoAdmin: repo,
	}

	// Assertions
	if assert.NoError(t, h.JWTLogin(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandler_JWTLogin_InsertOneErr(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`

	req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := new(mymock.MockRepoAdmin)
	repo.On("FindOneByID", "krefa").Return(models.Admin{}, errors.New("InsertOneErr"))

	h.DBRepo = &database.DBRepo{
		RepoAdmin: repo,
	}

	errHandler := h.JWTLogin(c).(*echo.HTTPError)
	// Assertions
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}
func TestHandler_JWTLogin_ValidationError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/jwt/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := Handler{}

	// Assertions
	errHandler := h.JWTLogin(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}
func TestHandler_JWTLogin_BindErr_NoContentType(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`

	req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errHandler := h.JWTLogin(c).(*echo.HTTPError)
	// Assertions
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}

func TestHandler_JWTRegister_Happy(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`

	req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)


	mockBcrypt := new(mymock.MockBcrypt)
	mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), nil)

	repo := new(mymock.MockRepoAdmin)
	repo.On("InsertOne", adminData).Return(nil)

	h.DBRepo = &database.DBRepo{
		RepoAdmin: repo,
	}
	h.Bcrypt=mockBcrypt

	// Assertions
	if assert.NoError(t, h.JWTRegister(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
func TestHandler_JWTRegister_HashErr(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`

	req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockBcrypt := new(mymock.MockBcrypt)
	mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), errors.New("HashErr"))


	h.Bcrypt=mockBcrypt

	errHandler := h.JWTRegister(c).(*echo.HTTPError)
	// Assertions
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}
func TestHandler_JWTRegister_InsertOneErr(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`

	req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockBcrypt := new(mymock.MockBcrypt)
	mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), nil)

	repo := new(mymock.MockRepoAdmin)
	repo.On("InsertOne", adminData).Return(errors.New("InsertOneErr"))

	h.DBRepo = &database.DBRepo{
		RepoAdmin: repo,
	}
	h.Bcrypt=mockBcrypt

	errHandler := h.JWTRegister(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}

func TestHandler_JWTRegister_ValidationError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/jwt/register", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := Handler{}

	// Assertions
	errHandler := h.JWTRegister(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, errHandler.Code)
}
func TestHandler_JWTRegister_BindErr_NoContentType(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`

	req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	errHandler := h.JWTRegister(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
}
