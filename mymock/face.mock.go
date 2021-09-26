package mymock

import (
	"goface-api/models"

	"github.com/Kagami/go-face"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockRepoFace struct {
	mock.Mock
}

func (repo MockRepoFace) InsertOne(face models.Face) error {
	args := repo.Called(face)
	return args.Error(0)
}

func (repo MockRepoFace) PushDescriptor(id string, descriptor face.Descriptor) (*mongo.UpdateResult, error) {
	return nil, nil
}

func (repo MockRepoFace) FindById(id string) (res []models.Face, err error) {
	args := repo.Called(id)
	return args.Get(0).([]models.Face), args.Error(1)
}

func (repo MockRepoFace) FindAll() (res []models.Face, err error) {
	args := repo.Called()
	return args.Get(0).([]models.Face), args.Error(1)
}

func (repo MockRepoFace) DeleteId(id string) error {
	args := repo.Called(id)
	return args.Error(0)
}
