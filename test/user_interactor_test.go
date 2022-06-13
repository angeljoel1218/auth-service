package test

import (
	uinter "auth-service/src/app/interactor"
	apierror "auth-service/src/domain/apierrors"
	fixture "auth-service/src/domain/fixtures"
	"auth-service/src/domain/models"
	"auth-service/test/mocks"
	"log"
	"testing"
	"time"

	common "auth-service/config"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	t.Helper()

	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("error loading environment variables:", err)
	}

	err = common.LoadConfig()
	if err != nil {
		log.Fatal("error loading environment variables:", err)
	}

	r, mock := mocks.UserRepository(t)
	i := uinter.NewUserInteractor(r)

	hashed, _ := fixture.HashAndSalt("12345678")

	var id int64 = 1
	var usr = &models.User{
		ID:       &id,
		Name:     "Angel Joel",
		Lastname: "Alvarado",
		Email:    "angeljoel@gmail.com",
		Password: hashed,
		Phone:    "+5198***301",
	}

	cases := map[string]struct {
		arrange func(t *testing.T)
		assert  func(t *testing.T, u *models.User, err error)
	}{
		"Success": {
			arrange: func(t *testing.T) {
				query := `SELECT * FROM "users" WHERE email = $1`
				rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "password", "phone"}).
					AddRow(usr.ID, usr.Name, usr.Lastname, usr.Email, usr.Password, usr.Phone)
				mock.ExpectQuery(query).WithArgs(usr.Email).WillReturnRows(rows)
			},
			assert: func(t *testing.T, u *models.User, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, u)
				assert.NotEmpty(t, u.Token)
			},
		},
		"Invalid password": {
			arrange: func(t *testing.T) {
				hashed, _ := fixture.HashAndSalt("1234567800")
				usr.Password = hashed
				query := `SELECT * FROM "users" WHERE email = $1`
				rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "password", "phone"}).
					AddRow(usr.ID, usr.Name, usr.Lastname, usr.Email, usr.Password, usr.Phone)
				mock.ExpectQuery(query).WithArgs(usr.Email).WillReturnRows(rows)
			},
			assert: func(t *testing.T, u *models.User, err error) {
				assert.Nil(t, u)
				assert.Error(t, err)
				assert.EqualError(t, err, models.ErrorUser(apierror.Unauthorized).Error())
			},
		},
		"Username does not exist": {
			arrange: func(t *testing.T) {
				query := `SELECT * FROM "users" WHERE email = $1`
				rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "password", "phone"})
				mock.ExpectQuery(query).WithArgs(usr.Email).WillReturnRows(rows)
			},
			assert: func(t *testing.T, u *models.User, err error) {
				assert.Nil(t, u)
				assert.Error(t, err)
				assert.EqualError(t, err, models.ErrorUser(apierror.NotFound).Error())
			},
		},
	}

	var usrl = &models.User{
		Email:    "angeljoel@gmail.com",
		Password: "12345678",
	}

	for k, tt := range cases {
		t.Run(k, func(t *testing.T) {
			tt.arrange(t)

			u, err := i.SignIn(usrl)

			tt.assert(t, u, err)
		})
	}
}

func TestSignUp(t *testing.T) {
	t.Helper()

	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("error loading environment variables:", err)
	}

	err = common.LoadConfig()
	if err != nil {
		log.Fatal("error loading environment variables:", err)
	}

	r, mock := mocks.UserRepository(t)
	i := uinter.NewUserInteractor(r)

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
				rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "password"})
				mock.ExpectQuery(`SELECT * FROM "users" WHERE email = $1`).
					WithArgs(usr.Email).WillReturnRows(rows)

				query := `INSERT INTO "users" ("name","lastname","email","password","phone","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`
				mock.ExpectBegin()

				mock.ExpectQuery(query).
					WithArgs(usr.Name, usr.Lastname, usr.Email, sqlmock.AnyArg(), usr.Phone, usr.CreatedAt, usr.UpdatedAt, usr.DeletedAt, usr.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(usr.ID))

				mock.ExpectCommit()
			},
			assert: func(t *testing.T, u *models.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, u)
				assert.NotEmpty(t, u.Token)
			},
		},
		"Error create": {
			arrange: func(t *testing.T) {
				rows := sqlmock.NewRows([]string{"name", "lastname", "email", "phone", "created_at", "updated_at", "deleted_at", "id"}).
					AddRow(usr.Name, usr.Lastname, usr.Email, usr.Phone, usr.CreatedAt, usr.UpdatedAt, usr.DeletedAt, usr.ID)
				mock.ExpectQuery(`SELECT * FROM "users" WHERE email = $1`).
					WithArgs(usr.Email).WillReturnRows(rows)

				query := `INSERT INTO "users" ("name","lastname","email","password","phone","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`
				mock.ExpectBegin()

				mock.ExpectQuery(query).
					WithArgs(usr.Name, usr.Lastname, usr.Email, sqlmock.AnyArg(), usr.Phone, usr.CreatedAt, usr.UpdatedAt, usr.DeletedAt, usr.ID).
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

			u, err := i.SignUp(usr)

			tt.assert(t, u, err)
		})
	}
}
