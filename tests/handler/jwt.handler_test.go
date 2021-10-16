package handler_test

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

func TestHandler_JWTLogin(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"krefa"}`

	req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(adminData.Password), bcrypt.DefaultCost)
	adminData.Password = string(hashed)

	repo := new(mymock.MockRepoAdmin)

	h.DBRepo = &database.DBRepo{
		RepoAdmin: repo,
	}

	t.Run("JWTLogin Happy", func(t *testing.T) {
		repo.On("FindOneByID", "krefa").Return(adminData, nil).Once()
		if assert.NoError(t, h.JWTLogin(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("JWTLogin InsertOneErr", func(t *testing.T) {
		repo.On("FindOneByID", "krefa").Return(models.Admin{}, errors.New("InsertOneErr")).Once()

		errHandler := h.JWTLogin(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})

}

func TestHandler_JWTLogin_ValidationError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/jwt/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

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

func TestHandler_JWTRegister(t *testing.T) {
	reqJSON := `{"username":"krefa","password":"` + adminData.Password + `"}`

	req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockBcrypt := new(mymock.MockBcrypt)
	repo := new(mymock.MockRepoAdmin)
	h.DBRepo = &database.DBRepo{
		RepoAdmin: repo,
	}
	h.Bcrypt = mockBcrypt

	t.Run("JWTRegister Happy", func(t *testing.T) {
		mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), nil).Once()
		repo.On("InsertOne", adminData).Return(nil).Once()
		if assert.NoError(t, h.JWTRegister(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})
	t.Run("JWTRegister HashErr", func(t *testing.T) {
		mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), errors.New("HashErr")).Once()
		errHandler := h.JWTRegister(c).(*echo.HTTPError)
		// Assertions
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})
	t.Run("JWTRegister InsertOneErr", func(t *testing.T) {
		mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), nil).Once()
		repo.On("InsertOne", adminData).Return(errors.New("InsertOneErr")).Once()

		errHandler := h.JWTRegister(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})

}

func TestHandler_JWTRegister_ValidationError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/jwt/register", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

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
