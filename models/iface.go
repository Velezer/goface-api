package models

import (
	"github.com/Kagami/go-face"
)

type RepositoryAdminIface interface {
	FindOneByID(id string) (res Admin, err error)
	InsertOne(admin Admin) (err error)
}

type RepositoryFaceIface interface {
	InsertOne(face Face) error
	PushDescriptor(id string, descriptor face.Descriptor) (error)
	FindById(id string) (res []Face, err error)
	FindAll() (res []Face, err error)
	DeleteId(id string) (error)
}