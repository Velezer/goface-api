package handler

import (
	"goface-api/helper"
	"goface-api/models"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var e *echo.Echo
var reco *face.Recognizer
var err error
var h Handler
var faceData models.Face
var adminData models.Admin

func TestMain(m *testing.M) {
	reco, err = face.NewRecognizer(filepath.Join("../", helper.ModelDir))
	if err != nil {
		panic(err)
	}
	defer reco.Close()

	e = echo.New()
	h = Handler{Rec: reco}
	faceData = models.Face{
		Id:          "4871847291721",
		Name:        "myname",
		Descriptors: []face.Descriptor{},
	}
	adminData = models.Admin{
		Username: "krefa",
		Password: "krefa",
	}
	code := m.Run()

	os.Exit(code)
}

func TestHandler_Home(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{}

	// Assertions
	if assert.NoError(t, h.Home(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"message":"goface-api up"}`, strings.TrimSpace(rec.Body.String()))
	}
}
