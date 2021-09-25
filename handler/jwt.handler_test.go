package handler

import (
	"goface-api/database"
	"goface-api/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockRepoAdmin struct {
	mock.Mock
}

func (coll *MockRepoAdmin) FindOneByID(id string) (models.Admin, error) {
	args := coll.Called(id)

	return args.Get(0).(models.Admin), args.Error(1) // type cast
}

func (coll *MockRepoAdmin) InsertOne(admin models.Admin) error {
	args := coll.Called(admin)

	return args.Get(0).(error) // type cast
}

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

	repo := new(MockRepoAdmin)
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
