package database

import "github.com/Kagami/go-face"

type Face struct {
	Id          string            `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string            `json:"name,omitempty" bson:"name,omitempty"`
	Descriptors []face.Descriptor `json:"descriptors,omitempty" bson:"descriptors,omitempty"`
}
