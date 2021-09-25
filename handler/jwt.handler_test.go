package handler

import (
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

func TestHandler_JWTRegister(t *testing.T) {
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
