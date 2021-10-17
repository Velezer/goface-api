package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)



func TestHandler_Integration_JWTRegister(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("No Content-Type", func(t *testing.T) {
		reqJSON := `{"username":"krefa","password":"krefa"}`

		req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		errHandler := hInt.JWTRegister(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})
	t.Run("Input Not Valid", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/jwt/register", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		errHandler := hInt.JWTRegister(c).(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, errHandler.Code)
	})

	t.Run("Input Valid", func(t *testing.T) {

		reqJSON := `{"username":"krefa","password":"` + adminData.Password + `"}`

		req := httptest.NewRequest(http.MethodPost, "/jwt/register", strings.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		t.Run("Happy", func(t *testing.T) {
			if assert.NoError(t, hInt.JWTRegister(c)) {
				assert.Equal(t, http.StatusCreated, rec.Code)
			}
		})

	})

}

func TestHandler_Integration_JWTLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("No Content-Type", func(t *testing.T) {
		reqJSON := `{"username":"krefa","password":"krefa"}`

		req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		errHandler := hInt.JWTLogin(c).(*echo.HTTPError)
		// Assertions
		assert.Equal(t, http.StatusInternalServerError, errHandler.Code)
	})
	t.Run("Input Not Valid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/jwt/login", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		errHandler := hInt.JWTLogin(c).(*echo.HTTPError)
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
			if assert.NoError(t, hInt.JWTLogin(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})

	})

}
