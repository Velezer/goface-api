package handler

import (
	"goface-api/database"
	"goface-api/mymock"
	"mime/multipart"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Rec    *face.Recognizer
	DBRepo *database.DBRepo
	Bcrypt mymock.BcryptIface
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
