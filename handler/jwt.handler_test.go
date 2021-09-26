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
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminData := models.Admin{
		Username: "krefa",
		Password: "krefa",
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(adminData.Password), bcrypt.DefaultCost)
	adminData.Password = string(hashed)

	repo := new(mymock.MockRepoAdmin)
	repo.On("FindOneByID", "krefa").Return(adminData, nil)

	dbRepo := database.DBRepo{
		RepoAdmin: repo,
	}
	h := Handler{DBRepo: &dbRepo}

	// Assertions
	if assert.NoError(t, h.JWTLogin(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandler_JWTLogin_InsertOneErr(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := new(mymock.MockRepoAdmin)
	repo.On("FindOneByID", "krefa").Return(models.Admin{}, errors.New("InsertOneErr"))

	dbRepo := database.DBRepo{
		RepoAdmin: repo,
	}
	h := Handler{DBRepo: &dbRepo}

	// Assertions
	assert.Error(t, h.JWTLogin(c), "InsertOneErr")
}
func TestHandler_JWTLogin_ValidationError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := Handler{}

	// Assertions
	assert.Error(t, h.JWTLogin(c))
}
func TestHandler_JWTLogin_BindErr_NoContentType(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := Handler{}

	// Assertions
	assert.Error(t, h.JWTLogin(c), 500)
}

func TestHandler_JWTRegister_Happy(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminData := models.Admin{
		Username: "krefa",
		Password: "krefa",
	}

	mockBcrypt := new(mymock.MockBcrypt)
	mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), nil)

	repo := new(mymock.MockRepoAdmin)
	repo.On("InsertOne", adminData).Return(nil)

	dbRepo := database.DBRepo{
		RepoAdmin: repo,
	}
	h := Handler{DBRepo: &dbRepo, Bcrypt: mockBcrypt}

	// Assertions
	if assert.NoError(t, h.JWTRegister(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
func TestHandler_JWTRegister_HashErr(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminData := models.Admin{
		Username: "krefa",
		Password: "krefa",
	}

	mockBcrypt := new(mymock.MockBcrypt)
	mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), errors.New("HashErr"))

	h := Handler{Bcrypt: mockBcrypt}

	// Assertions
	assert.Error(t, h.JWTRegister(c), "HashErr")
}
func TestHandler_JWTRegister_InsertOneErr(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminData := models.Admin{
		Username: "krefa",
		Password: "krefa",
	}

	mockBcrypt := new(mymock.MockBcrypt)
	mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), nil)

	repo := new(mymock.MockRepoAdmin)
	repo.On("InsertOne", adminData).Return(errors.New("InsertOneErr"))

	dbRepo := database.DBRepo{
		RepoAdmin: repo,
	}
	h := Handler{DBRepo: &dbRepo, Bcrypt: mockBcrypt}

	// Assertions
	assert.Error(t, h.JWTRegister(c), "InsertOneErr")
}

func TestHandler_JWTRegister_ValidationError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/register", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := Handler{}

	// Assertions
	assert.Error(t, h.JWTRegister(c))
}
func TestHandler_JWTRegister_BindErr_NoContentType(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := Handler{}

	// Assertions
	assert.Error(t, h.JWTRegister(c))
}
