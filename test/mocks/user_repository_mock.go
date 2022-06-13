package mocks

import (
	urepo "auth-service/src/app/repository"
	irepo "auth-service/src/infrastructure/database/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func UserRepository(t *testing.T) (r urepo.UserRepository, mock sqlmock.Sqlmock) {
	db, mock, _ := NewDBMock(t)

	r = irepo.NewUserRepository(db)

	return r, mock
}
