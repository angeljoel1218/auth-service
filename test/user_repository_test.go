package test

import (
	apierror "auth-service/src/domain/apierrors"
	"auth-service/src/domain/models"
	"auth-service/test/mocks"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestFindByEmail(t *testing.T) {
	t.Helper()

	r, mock := mocks.UserRepository(t)

	var id int64 = 1
	var usr = &models.User{
		ID:       &id,
		Name:     "Angel Joel",
		Lastname: "Alvarado",
		Email:    "angeljoel@gmail.com",
		Password: "12345678",
	}

	cases := map[string]struct {
		arrange func(t *testing.T)
		assert  func(t *testing.T, u *models.User, err error)
	}{
		"Success": {
			arrange: func(t *testing.T) {
				query := `SELECT * FROM "users" WHERE email = $1`
				rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "password"}).
					AddRow(usr.ID, usr.Name, usr.Lastname, usr.Email, usr.Password)
				mock.ExpectQuery(query).WithArgs(usr.Email).WillReturnRows(rows)
			},
			assert: func(t *testing.T, u *models.User, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, u)
			},
		},
		"Error": {
			arrange: func(t *testing.T) {
				query := `SELECT * FROM "users" WHERE email = $1`
				rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "password"})
				mock.ExpectQuery(query).WithArgs(usr.Email).WillReturnRows(rows)
			},
			assert: func(t *testing.T, u *models.User, err error) {
				assert.Nil(t, u)
				assert.Error(t, err)
			},
		},
	}

	for k, tt := range cases {
		t.Run(k, func(t *testing.T) {
			tt.arrange(t)

			u, err := r.FindByEmail(usr.Email)

			tt.assert(t, u, err)
		})
	}
}

func TestCreate(t *testing.T) {
	t.Helper()

	r, mock := mocks.UserRepository(t)

	var id int64 = 1
	var usr = &models.User{
		ID:        &id,
		Name:      "Angel Joel",
		Lastname:  "Alvarado",
		Email:     "angeljoel@gmail.com",
		Password:  "12345678",
		Phone:     "+5198***301",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	cases := map[string]struct {
		arrange func(t *testing.T)
		assert  func(t *testing.T, u *models.User, err error)
	}{
		"Success create": {
			arrange: func(t *testing.T) {
				query := `INSERT INTO "users" ("name","lastname","email","password","phone","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`
				mock.ExpectBegin()
				mock.ExpectQuery(query).
					WithArgs(usr.Name, usr.Lastname, usr.Email, usr.Password, usr.Phone, usr.CreatedAt, usr.UpdatedAt, usr.DeletedAt, usr.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(usr.ID))
				mock.ExpectCommit()
			},
			assert: func(t *testing.T, u *models.User, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, u)
			},
		},
		"Error create": {
			arrange: func(t *testing.T) {
				query := `INSERT INTO "users" ("name","lastname","email","password","phone","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`
				mock.ExpectBegin().WillReturnError(models.ErrorUser(apierror.InternalServer))
				mock.ExpectQuery(query).
					WithArgs(usr.Name, usr.Lastname, usr.Email, usr.Password, usr.Phone, usr.CreatedAt, usr.UpdatedAt, usr.DeletedAt, usr.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectCommit()
			},
			assert: func(t *testing.T, u *models.User, err error) {
				assert.Nil(t, u)
				assert.Error(t, err)
			},
		},
	}

	for k, tt := range cases {
		t.Run(k, func(t *testing.T) {
			tt.arrange(t)

			u, err := r.Create(usr)

			tt.assert(t, u, err)
		})
	}
}
