package models

import (
	"github.com/Kagami/go-face"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryAdminIface interface {
	FindOneByID(id string) (res Admin, err error)
	InsertOne(admin Admin) error
}

type RepositoryFaceIface interface {
	InsertOne(face Face) error
	PushDescriptor(id string, descriptor face.Descriptor) (*mongo.UpdateResult, error)
	FindById(id string) (res []Face, err error)
	FindAll() (res []Face, err error)
	DeleteId(id string) (*mongo.DeleteResult, error)
}