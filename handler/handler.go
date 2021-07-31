package handler

import (
	"github.com/Kagami/go-face"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Rec  *face.Recognizer
	Coll *mongo.Collection
}


type inputValidation struct {
	Id   string `validate:"required"`
	Name string `validate:"required"`
}
