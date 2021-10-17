package handler

import (
	"goface-api/database"
	"goface-api/iface"
	"mime/multipart"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Rec    *face.Recognizer
	DBRepo *database.DBRepo
	Bcrypt iface.BcryptIface
}



type inputValidation struct {
	Id   string `form:"id" json:"id" validate:"required,gte=9,lte=16"`
	Name string `form:"name" json:"name" validate:"required"`
}

func getFileContent(c echo.Context, fieldName string) (multipart.File, error) {
	file, err := c.FormFile(fieldName) //name=file in client html form
	if err != nil {
		return nil, err
	}

	content, err := file.Open()
	if err != nil {
		return nil, err
	}

	return content, nil
}


// ------real bcrypt----------

type RealBcrypt struct {
}

func (b RealBcrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}
func (b RealBcrypt) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

// ------end real bcrypt----------
