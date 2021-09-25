package mymock

import (
	"goface-api/models"

	"github.com/stretchr/testify/mock"
)

type MockRepoAdmin struct {
	mock.Mock
}

func (coll *MockRepoAdmin) FindOneByID(id string) (models.Admin, error) {
	args := coll.Called(id)

	return args.Get(0).(models.Admin), args.Error(1) // type cast
}

func (coll *MockRepoAdmin) InsertOne(admin models.Admin) error {
	args := coll.Called(admin)

	return args.Error(0)
}
