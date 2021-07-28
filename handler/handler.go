package handler

import (
	"github.com/Kagami/go-face"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Rec  *face.Recognizer
	Coll *mongo.Collection
}
