package iface

import (
	"goface-api/models"

	"github.com/Kagami/go-face"
)

type RepositoryAdminIface interface {
	FindOneByID(id string) (res models.Admin, err error)
	InsertOne(admin models.Admin) (err error)
}

type RepositoryFaceIface interface {
	InsertOne(face models.Face) error
	PushDescriptor(id string, descriptor face.Descriptor) (error)
	FindById(id string) (res []models.Face, err error)
	FindAll() (res []models.Face, err error)
	DeleteId(id string) (error)
}

type BcryptIface interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}