package test

import (
	fixture "auth-service/src/domain/fixtures"
	"auth-service/src/registry"
	"auth-service/test/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignUpAction(t *testing.T) {
	t.Helper()

	db, mock, _ := mocks.NewDBMock(t)

	cases := map[string]struct {
		arrange func(t *testing.T) *httptest.ResponseRecorder
		assert  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		"Success create": {
			arrange: func(t *testing.T) *httptest.ResponseRecorder {
				usr := gin.H{
					"name":     "Angel Joel",
					"lastname": "Alvarado",
					"email":    "angeljoel@gmail.com",
					"password": "12345678",
					"phone":    "+5198***301",
				}

				rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "password"})
				mock.ExpectQuery(`SELECT * FROM "users" WHERE email = $1`).
					WithArgs(usr["email"]).WillReturnRows(rows)

				query := `INSERT INTO "users" ("name","lastname","email","password","phone","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`
				mock.ExpectBegin()

				mock.ExpectQuery(query).
					WithArgs(usr["name"], usr["lastname"], usr["email"], sqlmock.AnyArg(), usr["phone"], sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()

				reg := registry.NewRegistry(db)
				router, _ := reg.StartHandlerInterface()

				rr := httptest.NewRecorder()

				reqBody, _ := json.Marshal(usr)

				request, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(reqBody))
				request.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(rr, request)

				return rr
			},
			assert: func(t *testing.T, rr *httptest.ResponseRecorder) {
				response := make(map[string]interface{})
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, 200, rr.Code)
				assert.Equal(t, float64(0), response["code"])
				assert.NotNil(t, response["data"])
				if response["data"] != nil {
					data := response["data"].(map[string]interface{})
					assert.NotEmpty(t, data["token"].(string))
				}
			},
		},
		"Invalid email": {
			arrange: func(t *testing.T) *httptest.ResponseRecorder {
				usr := gin.H{
					"name":     "Angel Joel",
					"lastname": "Alvarado",
					"email":    "angeljoel",
					"password": "12345678",
					"phone":    "+5198***301",
				}

				reg := registry.NewRegistry(db)
				router, _ := reg.StartHandlerInterface()

				rr := httptest.NewRecorder()

				reqBody, _ := json.Marshal(usr)

				request, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(reqBody))
				request.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(rr, request)

				return rr
			},
			assert: func(t *testing.T, rr *httptest.ResponseRecorder) {
				response := make(map[string]interface{})
				json.NewDecoder(rr.Body).Decode(&response)
				assert.Equal(t, response["message"], "Invalid e-mail")
				assert.Equal(t, 400, rr.Code)
				assert.NotEqual(t, float64(0), response["code"])
			},
		},
		"Invalid password": {
			arrange: func(t *testing.T) *httptest.ResponseRecorder {
				usr := gin.H{
					"name":     "Angel Joel",
					"lastname": "Alvarado",
					"email":    "angeljoel@gmail.com",
					"password": "12345",
					"phone":    "+5198***301",
				}

				reg := registry.NewRegistry(db)
				router, _ := reg.StartHandlerInterface()

				rr := httptest.NewRecorder()

				reqBody, _ := json.Marshal(usr)

				request, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(reqBody))
				request.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(rr, request)

				return rr
			},
			assert: func(t *testing.T, rr *httptest.ResponseRecorder) {
				response := make(map[string]interface{})
				json.NewDecoder(rr.Body).Decode(&response)
				assert.Equal(t, response["message"], "The password field must be greater than or equal to 6 characters")
				assert.Equal(t, 400, rr.Code)
				assert.NotEqual(t, float64(0), response["code"])
			},
		},
		"Create error": {
			arrange: func(t *testing.T) *httptest.ResponseRecorder {
				usr := gin.H{}

				reg := registry.NewRegistry(db)
				router, _ := reg.StartHandlerInterface()

				rr := httptest.NewRecorder()

				reqBody, _ := json.Marshal(usr)

				request, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(reqBody))
				request.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(rr, request)

				return rr
			},
			assert: func(t *testing.T, rr *httptest.ResponseRecorder) {
				response := make(map[string]interface{})
				json.NewDecoder(rr.Body).Decode(&response)
				assert.Equal(t, 400, rr.Code)
				assert.NotEqual(t, float64(0), response["code"])
			},
		},
	}

	for k, tt := range cases {
		t.Run(k, func(t *testing.T) {
			rr := tt.arrange(t)
			tt.assert(t, rr)
		})
	}
}

func TestSignInAction(t *testing.T) {
	t.Helper()

	db, mock, _ := mocks.NewDBMock(t)

	cases := map[string]struct {
		arrange func(t *testing.T) *httptest.ResponseRecorder
		assert  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		"Successful login": {
			arrange: func(t *testing.T) *httptest.ResponseRecorder {
				usr := gin.H{
					"email":    "angeljoel@gmail.com",
					"password": "12345678",
				}

				hashed, _ := fixture.HashAndSalt("12345678")

				query := `SELECT * FROM "users" WHERE email = $1`
				rows := sqlmock.NewRows([]string{"password", "id"}).AddRow(hashed, 1)
				mock.ExpectQuery(query).WithArgs(usr["email"]).WillReturnRows(rows)

				reg := registry.NewRegistry(db)
				router, _ := reg.StartHandlerInterface()

				rr := httptest.NewRecorder()

				reqBody, _ := json.Marshal(usr)

				request, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
				request.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(rr, request)

				return rr
			},
			assert: func(t *testing.T, rr *httptest.ResponseRecorder) {
				response := make(map[string]interface{})
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, 200, rr.Code)
				assert.Equal(t, float64(0), response["code"])
				assert.NotNil(t, response["data"])
				if response["data"] != nil {
					data := response["data"].(map[string]interface{})
					assert.NotEmpty(t, data["token"].(string))
				}
			},
		},
		"Invalid password": {
			arrange: func(t *testing.T) *httptest.ResponseRecorder {
				usr := gin.H{
					"email":    "angeljoel@gmail.com",
					"password": "12345678",
				}

				hashed, _ := fixture.HashAndSalt("12345678dsds")

				query := `SELECT * FROM "users" WHERE email = $1`
				rows := sqlmock.NewRows([]string{"password", "id"}).AddRow(hashed, 1)
				mock.ExpectQuery(query).WithArgs(usr["email"]).WillReturnRows(rows)

				reg := registry.NewRegistry(db)
				router, _ := reg.StartHandlerInterface()

				rr := httptest.NewRecorder()

				reqBody, _ := json.Marshal(usr)

				request, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
				request.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(rr, request)

				return rr
			},
			assert: func(t *testing.T, rr *httptest.ResponseRecorder) {
				response := make(map[string]interface{})
				json.NewDecoder(rr.Body).Decode(&response)
				assert.Equal(t, 401, rr.Code)
				assert.NotEqual(t, float64(0), response["code"])
			},
		},
		"Login error": {
			arrange: func(t *testing.T) *httptest.ResponseRecorder {
				usr := gin.H{}

				reg := registry.NewRegistry(db)
				router, _ := reg.StartHandlerInterface()

				rr := httptest.NewRecorder()

				reqBody, _ := json.Marshal(usr)

				request, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
				request.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(rr, request)

				return rr
			},
			assert: func(t *testing.T, rr *httptest.ResponseRecorder) {
				response := make(map[string]interface{})
				json.NewDecoder(rr.Body).Decode(&response)
				assert.Equal(t, 400, rr.Code)
				assert.NotEqual(t, float64(0), response["code"])
			},
		},
	}

	for k, tt := range cases {
		t.Run(k, func(t *testing.T) {
			rr := tt.arrange(t)
			tt.assert(t, rr)
		})
	}
}

func TestAuthAction(t *testing.T) {
	t.Helper()

	db, _, _ := mocks.NewDBMock(t)

	cases := map[string]struct {
		arrange func(t *testing.T) *httptest.ResponseRecorder
		assert  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		"Valid token": {
			arrange: func(t *testing.T) *httptest.ResponseRecorder {
				usr := gin.H{}

				token, _ := fixture.CreateToken(int64(1))

				reg := registry.NewRegistry(db)
				router, _ := reg.StartHandlerInterface()

				rr := httptest.NewRecorder()

				reqBody, _ := json.Marshal(usr)

				request, _ := http.NewRequest(http.MethodPost, "/auth/authenticated", bytes.NewBuffer(reqBody))
				request.Header.Set("Authorization", "Bearer "+token)

				router.ServeHTTP(rr, request)

				return rr
			},
			assert: func(t *testing.T, rr *httptest.ResponseRecorder) {
				response := make(map[string]interface{})
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, 200, rr.Code)
				assert.Equal(t, float64(0), response["code"])
				assert.Equal(t, response["message"], "Authenticated user")
			},
		},
		"Invalid Token": {
			arrange: func(t *testing.T) *httptest.ResponseRecorder {
				usr := gin.H{}

				token := "test"

				reg := registry.NewRegistry(db)
				router, _ := reg.StartHandlerInterface()

				rr := httptest.NewRecorder()

				reqBody, _ := json.Marshal(usr)

				request, _ := http.NewRequest(http.MethodPost, "/auth/authenticated", bytes.NewBuffer(reqBody))
				request.Header.Set("Authorization", "Bearer "+token)

				router.ServeHTTP(rr, request)

				return rr
			},
			assert: func(t *testing.T, rr *httptest.ResponseRecorder) {
				response := make(map[string]interface{})
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, 401, rr.Code)
			},
		},
	}

	for k, tt := range cases {
		t.Run(k, func(t *testing.T) {
			rr := tt.arrange(t)
			tt.assert(t, rr)
		})
	}
}
