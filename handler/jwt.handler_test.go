package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_JWTLogin(t *testing.T) {
	reqJSON := `{"_id":"krefa","password":"krefa"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/jwt/login", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{}

	// Assertions
	if assert.NoError(t, h.JWTLogin(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
