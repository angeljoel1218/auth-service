package mocks

import (
	irepo "auth-service/src/interface/repository"
	urepo "auth-service/src/usecase/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func UserRepository(t *testing.T) (r urepo.UserRepository, mock sqlmock.Sqlmock) {
	db, mock, _ := NewDBMock(t)

	r = irepo.NewUserRepository(db)

	return r, mock
}
