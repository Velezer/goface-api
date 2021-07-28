package database

import "github.com/Kagami/go-face"

type Face struct {
	Id          string `bson:"_id,omitempty"`
	Name        string `bson:"name,omitempty"`
	Descriptors []face.Descriptor `bson:"descriptors,omitempty"`
}