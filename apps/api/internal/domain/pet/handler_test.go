package pet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/models"
	"github.com/my-pets/api/internal/domain/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetPaginated(ctx context.Context, ownerID string, page, perPage int) ([]models.Pet, int64, error) {
	args := m.Called(ctx, ownerID, page, perPage)
	return args.Get(0).([]models.Pet), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) GetByID(ctx context.Context, id, ownerID string) (models.Pet, error) {
	args := m.Called(ctx, id, ownerID)
	return args.Get(0).(models.Pet), args.Error(1)
}

func (m *mockRepo) Create(ctx context.Context, ownerID string, payload CreatePetPayload) (models.Pet, error) {
	args := m.Called(ctx, ownerID, payload)
	return args.Get(0).(models.Pet), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, id, ownerID string, payload UpdatePetPayload) (models.Pet, error) {
	args := m.Called(ctx, id, ownerID, payload)
	return args.Get(0).(models.Pet), args.Error(1)
}

func (m *mockRepo) Delete(ctx context.Context, id, ownerID string) error {
	args := m.Called(ctx, id, ownerID)
	return args.Error(0)
}

func (m *mockRepo) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *mockUserRepo) GetPaginated(ctx context.Context, page, perPage int) ([]models.User, int64, error) {
	args := m.Called(ctx, page, perPage)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserRepo) Create(ctx context.Context, payload user.CreateUserPayload, isSystemUser bool) (models.User, error) {
	args := m.Called(ctx, payload, isSystemUser)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, id string, payload user.UpdateUserPayload) (models.User, error) {
	args := m.Called(ctx, id, payload)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserRepo) HasUsers(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepo) SystemUserExists(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepo) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *mockUserRepo) UpsertByGoogleID(ctx context.Context, info user.GoogleUserInfo, isSystemUser bool) (models.User, error) {
	args := m.Called(ctx, info, isSystemUser)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, id string, newHash string) error {
	args := m.Called(ctx, id, newHash)
	return args.Error(0)
}

type petTestCase struct {
	name     string
	method   string
	path     string
	body     string
	ownerID  string
	urlParams gin.Params
	setup    func(*mockRepo, *mockUserRepo)
	check    func(*testing.T, *httptest.ResponseRecorder)
}

func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func routeFromPath(path string) string {
	for i := range path {
		if path[i] == '?' {
			return path[:i]
		}
	}
	return path
}

func (tc *petTestCase) exec(t *testing.T) {
	t.Helper()

	repoMock := new(mockRepo)
	userMock := new(mockUserRepo)
	handler := NewHandler(repoMock, userMock)

	c, w := setupGin()
	c.Request = httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = tc.urlParams
	c.Set("userID", tc.ownerID)

	if tc.setup != nil {
		tc.setup(repoMock, userMock)
	}

	switch tc.method + " " + routeFromPath(tc.path) {
	case "GET /pets":
		handler.GetPets(c)
	case "GET /pets/:id":
		handler.GetPet(c)
	case "POST /pets":
		handler.CreatePet(c)
	case "DELETE /pets/:id":
		handler.DeletePet(c)
	default:
		t.Fatalf("unknown route: %s %s", tc.method, tc.path)
	}

	tc.check(t, w)
	repoMock.AssertExpectations(t)
	userMock.AssertExpectations(t)
}

func TestPetHandler(t *testing.T) {
	tests := []petTestCase{
		{
			name:    "GetPets success",
			method:  "GET",
			path:    "/pets?page=1&per_page=10",
			ownerID: "user-1",
			setup: func(m *mockRepo, _ *mockUserRepo) {
				pets := []models.Pet{
					{ID: "pet-1", Name: "Firulais", Species: "dog"},
					{ID: "pet-2", Name: "Mishi", Species: "cat"},
				}
				m.On("GetPaginated", mock.Anything, "user-1", 1, 10).Return(pets, int64(2), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var resp map[string]any
				json.Unmarshal(w.Body.Bytes(), &resp)
				require.Equal(t, float64(2), resp["total"])
				require.Equal(t, float64(1), resp["page"])
				require.Equal(t, float64(10), resp["per_page"])
				require.Equal(t, float64(1), resp["total_pages"])
			},
		},
		{
			name:    "GetPets default pagination",
			method:  "GET",
			path:    "/pets",
			ownerID: "user-1",
			setup: func(m *mockRepo, _ *mockUserRepo) {
				m.On("GetPaginated", mock.Anything, "user-1", 1, 10).Return([]models.Pet{}, int64(0), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:    "GetPets repo error",
			method:  "GET",
			path:    "/pets",
			ownerID: "user-1",
			setup: func(m *mockRepo, _ *mockUserRepo) {
				m.On("GetPaginated", mock.Anything, "user-1", 1, 10).Return([]models.Pet{}, int64(0), errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:    "GetPet success",
			method:  "GET",
			path:    "/pets/:id",
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "pet-1"}},
			setup: func(m *mockRepo, _ *mockUserRepo) {
				m.On("GetByID", mock.Anything, "pet-1", "user-1").Return(models.Pet{ID: "pet-1", Name: "Firulais"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var resp map[string]models.Pet
				json.Unmarshal(w.Body.Bytes(), &resp)
				require.Equal(t, "Firulais", resp["data"].Name)
			},
		},
		{
			name:    "GetPet not found",
			method:  "GET",
			path:    "/pets/:id",
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "pet-404"}},
			setup: func(m *mockRepo, _ *mockUserRepo) {
				m.On("GetByID", mock.Anything, "pet-404", "user-1").Return(models.Pet{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:    "GetPet repo error",
			method:  "GET",
			path:    "/pets/:id",
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "pet-1"}},
			setup: func(m *mockRepo, _ *mockUserRepo) {
				m.On("GetByID", mock.Anything, "pet-1", "user-1").Return(models.Pet{}, errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:    "CreatePet dog success",
			method:  "POST",
			path:    "/pets",
			ownerID: "user-1",
			body:    `{"name":"Firulais","species":"dog","birth_date":"2023-01-15","size":"small"}`,
			setup: func(m *mockRepo, u *mockUserRepo) {
				u.On("GetByID", mock.Anything, "user-1").Return(models.User{PetLimit: 5}, nil)
				m.On("CountByOwner", mock.Anything, "user-1").Return(int64(1), nil)
				m.On("Create", mock.Anything, "user-1", mock.MatchedBy(func(p CreatePetPayload) bool {
					return p.Name == "Firulais" && p.Species == "dog" && *p.Size == "small"
				})).Return(models.Pet{ID: "pet-1", Name: "Firulais"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			name:    "CreatePet cat without size",
			method:  "POST",
			path:    "/pets",
			ownerID: "user-1",
			body:    `{"name":"Mishi","species":"cat","birth_date":"2022-06-01"}`,
			setup: func(m *mockRepo, u *mockUserRepo) {
				u.On("GetByID", mock.Anything, "user-1").Return(models.User{PetLimit: 5}, nil)
				m.On("CountByOwner", mock.Anything, "user-1").Return(int64(0), nil)
				m.On("Create", mock.Anything, "user-1", mock.MatchedBy(func(p CreatePetPayload) bool {
					return p.Name == "Mishi" && p.Species == "cat" && p.Size == nil
				})).Return(models.Pet{ID: "pet-2", Name: "Mishi"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			name:    "CreatePet validation error",
			method:  "POST",
			path:    "/pets",
			ownerID: "user-1",
			body:    `{"name":"","species":"dog","birth_date":"2023-01-15","size":"small"}`,
			setup:   func(_ *mockRepo, _ *mockUserRepo) {},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:    "CreatePet dog missing size",
			method:  "POST",
			path:    "/pets",
			ownerID: "user-1",
			body:    `{"name":"Firulais","species":"dog","birth_date":"2023-01-15"}`,
			setup:   func(_ *mockRepo, _ *mockUserRepo) {},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:    "CreatePet pet limit reached",
			method:  "POST",
			path:    "/pets",
			ownerID: "user-1",
			body:    `{"name":"Firulais","species":"dog","birth_date":"2023-01-15","size":"small"}`,
			setup: func(m *mockRepo, u *mockUserRepo) {
				u.On("GetByID", mock.Anything, "user-1").Return(models.User{PetLimit: 2}, nil)
				m.On("CountByOwner", mock.Anything, "user-1").Return(int64(2), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:    "DeletePet success",
			method:  "DELETE",
			path:    "/pets/:id",
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "pet-1"}},
			setup: func(m *mockRepo, _ *mockUserRepo) {
				m.On("Delete", mock.Anything, "pet-1", "user-1").Return(nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:    "DeletePet not found",
			method:  "DELETE",
			path:    "/pets/:id",
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "pet-404"}},
			setup: func(m *mockRepo, _ *mockUserRepo) {
				m.On("Delete", mock.Anything, "pet-404", "user-1").Return(ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
