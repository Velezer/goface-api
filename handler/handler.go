package handler

import (
	"mime/multipart"

	"github.com/Kagami/go-face"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Rec  *face.Recognizer
	Coll *mongo.Collection
	DB *mongo.Database
}


type inputValidation struct {
	Id   string `form:"id" validate:"required,gte=9,lte=16"`
	Name string `form:"name" validate:"required"`
}

func getFileContent(c echo.Context, fieldName string) (multipart.File, error){
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
