package mymock

import (
	"github.com/stretchr/testify/mock"
)

// ------mock bcrypt----------

type MockBcrypt struct {
	mock.Mock
}

func (b MockBcrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	args := b.Called(password, cost)

	return args.Get(0).([]byte), args.Error(1) // type cast
}

// ------end mock bcrypt----------
