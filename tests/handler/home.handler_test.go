package handler_test

import (
	"context"
	"flag"
	"goface-api/database"
	. "goface-api/handler"
	"goface-api/helper"
	"goface-api/models"
	"goface-api/mymock"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var e *echo.Echo
var reco *face.Recognizer
var err error
var h Handler
var hInt Handler

var mockRepoFace = new(mymock.MockRepoFace)
var mockBcrypt = new(mymock.MockBcrypt)
var mockRepoAdmin = new(mymock.MockRepoAdmin)

var faceData = models.Face{
	Id:          "4871847291721",
	Name:        "myname",
	Descriptors: []face.Descriptor{},
}
var adminData = models.Admin{
	Username: "krefa",
	Password: "krefa",
}

func TestMain(m *testing.M) {
	reco, err = face.NewRecognizer(filepath.Join("../../", helper.ModelDir))
	if err != nil {
		panic(err)
	}
	defer reco.Close()

	e = echo.New()
	h.Rec = reco
	h.Bcrypt = mockBcrypt
	h.DBRepo = &database.DBRepo{
		RepoFace:  mockRepoFace,
		RepoAdmin: mockRepoAdmin,
	}

	flag.Parse()
	if !testing.Short() {
		dbrepo, err := database.InitDB("mongodb://localhost:27017/db_goface_api_echo_test", "db_goface_api_echo_test")
		if err != nil {
			log.Fatalf("Can't init database: %v", err)
		}
		if dbrepo == nil {
			log.Fatalln("dbrepo is nil")
		}

		dbrepo.RepoAdmin.(models.RepoAdmin).Collection.DeleteMany(context.Background(), bson.M{})
		dbrepo.RepoFace.(models.RepoFace).Collection.DeleteMany(context.Background(), bson.M{})

		hInt.Rec = reco
		hInt.Bcrypt = new(RealBcrypt)
		hInt.DBRepo = dbrepo
	}

	code := m.Run()

	os.Exit(code)
}

func TestHandler_Home(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, h.Home(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"message":"goface-api up"}`, strings.TrimSpace(rec.Body.String()))
	}
}
