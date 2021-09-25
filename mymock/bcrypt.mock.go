package mymock

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type BcryptIface interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type MockBcrypt struct {
	mock.Mock
}

func (b MockBcrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	args := b.Called(password, cost)

	return args.Get(0).([]byte), args.Error(1) // type cast
}

type RealBcrypt struct{

}

func (b RealBcrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}
