package handler_test

import (
	"errors"
	"goface-api/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHandler_JWTLogin(t *testing.T) {
	t.Run("No Content-Type", func(t *testing.T) {
		reqJSON := `{"username":"krefa","password":"krefa"}`

		req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		errHandler := h.JWTLogin(c).(*echo.HTTPError)
		// Assertions
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})
	t.Run("Input Not Valid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/jwt/login", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		errHandler := h.JWTLogin(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, errHandler.Code)
	})
	t.Run("Input valid", func(t *testing.T) {
		reqJSON := `{"username":"krefa","password":"krefa"}`

		req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		hashed, _ := bcrypt.GenerateFromPassword([]byte(adminData.Password), bcrypt.DefaultCost)
		adminData.Password = string(hashed)

		t.Run("Happy", func(t *testing.T) {
			mockRepoAdmin.On("FindOneByID", "krefa").Return(adminData, nil).Once()
			if assert.NoError(t, h.JWTLogin(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})

		t.Run("InsertOneErr", func(t *testing.T) {
			mockRepoAdmin.On("FindOneByID", "krefa").Return(models.Admin{}, errors.New("InsertOneErr")).Once()

			errHandler := h.JWTLogin(c).(*echo.HTTPError)
			assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
		})

	})

}

func TestHandler_JWTRegister(t *testing.T) {

	t.Run("No Content-Type", func(t *testing.T) {
		reqJSON := `{"username":"krefa","password":"krefa"}`

		req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		errHandler := h.JWTRegister(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})
	t.Run("Input Not Valid", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/jwt/register", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		errHandler := h.JWTRegister(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, errHandler.Code)
	})

	t.Run("Input Valid", func(t *testing.T) {

		reqJSON := `{"username":"krefa","password":"` + adminData.Password + `"}`

		req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)


		t.Run("Happy", func(t *testing.T) {
			mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), nil).Once()
			mockRepoAdmin.On("InsertOne", adminData).Return(nil).Once()
			if assert.NoError(t, h.JWTRegister(c)) {
				assert.Equal(t, http.StatusCreated, rec.Code)
			}
		})
		t.Run("HashErr", func(t *testing.T) {
			mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), errors.New("HashErr")).Once()
			errHandler := h.JWTRegister(c).(*echo.HTTPError)
			// Assertions
			assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
		})
		t.Run("InsertOneErr", func(t *testing.T) {
			mockBcrypt.On("GenerateFromPassword", []byte(adminData.Password), bcrypt.DefaultCost).Return([]byte(adminData.Password), nil).Once()
			mockRepoAdmin.On("InsertOne", adminData).Return(errors.New("InsertOneErr")).Once()

			errHandler := h.JWTRegister(c).(*echo.HTTPError)
			assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
		})

	})

}
