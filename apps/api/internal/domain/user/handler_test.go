package user

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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetAll(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *mockRepo) GetPaginated(ctx context.Context, page, perPage int) ([]models.User, int64, error) {
	args := m.Called(ctx, page, perPage)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockRepo) Create(ctx context.Context, payload CreateUserPayload, isSystemUser bool) (models.User, error) {
	args := m.Called(ctx, payload, isSystemUser)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, id string, payload UpdateUserPayload) (models.User, error) {
	args := m.Called(ctx, id, payload)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) HasUsers(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) SystemUserExists(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) UpsertByGoogleID(ctx context.Context, info GoogleUserInfo, isSystemUser bool) (models.User, error) {
	args := m.Called(ctx, info, isSystemUser)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockRepo) UpdatePassword(ctx context.Context, id string, newHash string) error {
	args := m.Called(ctx, id, newHash)
	return args.Error(0)
}

func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	return c, w
}

type userTestCase struct {
	name     string
	method   string
	path     string
	body     string
	ownerID  string
	isSystem bool
	urlParams gin.Params
	setup    func(*mockRepo)
	check    func(*testing.T, *httptest.ResponseRecorder)
}

func routeFromPath(path string) string {
	for i := range path {
		if path[i] == '?' {
			return path[:i]
		}
	}
	return path
}

func (tc *userTestCase) exec(t *testing.T) {
	t.Helper()

	m := new(mockRepo)
	petCountFn := func(_ context.Context, userID string) (int64, error) {
		switch userID {
		case "user-1":
			return 3, nil
		case "user-2":
			return 0, nil
		default:
			return 0, nil
		}
	}
	handler := NewHandler(m, petCountFn)

	c, w := setupGin()
	c.Request = httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = tc.urlParams
	c.Set("userID", tc.ownerID)
	if tc.isSystem {
		c.Set("isSystemUser", true)
	}

	if tc.setup != nil {
		tc.setup(m)
	}

	switch tc.method + " " + routeFromPath(tc.path) {
	case "GET /users":
		handler.GetUsers(c)
	case "GET /users/:id":
		handler.GetUser(c)
	case "POST /users":
		handler.CreateUser(c)
	case "PUT /users/:id":
		handler.UpdateUser(c)
	case "DELETE /users/:id":
		handler.DeleteUser(c)
	default:
		t.Fatalf("unknown route: %s %s", tc.method, tc.path)
	}

	tc.check(t, w)
	m.AssertExpectations(t)
}

var defaultUser = models.User{
	ID:    "user-1",
	Name:  "Test User",
	Email: "test@example.com",
}

func TestUserHandler(t *testing.T) {
	tests := []userTestCase{
		// ── GET /users ───────────────────────────────────────────────
		{
			name:     "GetUsers success",
			method:   "GET",
			path:     "/users?page=1&per_page=10",
			ownerID:  "admin-1",
			isSystem: true,
			setup: func(m *mockRepo) {
				users := []models.User{
					{ID: "user-1", Name: "Alice", Email: "alice@test.com"},
					{ID: "user-2", Name: "Bob", Email: "bob@test.com"},
				}
				m.On("GetPaginated", mock.Anything, 1, 10).Return(users, int64(2), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var resp map[string]any
				json.Unmarshal(w.Body.Bytes(), &resp)
				require.Equal(t, float64(2), resp["total"])
				data := resp["data"].([]any)
				require.Len(t, data, 2)
				alice := data[0].(map[string]any)
				require.Equal(t, float64(3), alice["pet_count"])
			},
		},
		{
			name:     "GetUsers forbidden (non-system)",
			method:   "GET",
			path:     "/users",
			ownerID:  "user-1",
			isSystem: false,
			setup:    func(_ *mockRepo) {},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:     "GetUsers repo error",
			method:   "GET",
			path:     "/users",
			ownerID:  "admin-1",
			isSystem: true,
			setup: func(m *mockRepo) {
				m.On("GetPaginated", mock.Anything, 1, 10).Return([]models.User{}, int64(0), errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── GET /users/:id ───────────────────────────────────────────
		{
			name:     "GetUser success",
			method:   "GET",
			path:     "/users/:id",
			ownerID:  "admin-1",
			isSystem: true,
			urlParams: gin.Params{{Key: "id", Value: "user-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "user-1").Return(defaultUser, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var resp map[string]map[string]any
				json.Unmarshal(w.Body.Bytes(), &resp)
				require.Equal(t, "Test User", resp["data"]["name"])
				require.Equal(t, float64(3), resp["data"]["pet_count"])
			},
		},
		{
			name:     "GetUser forbidden (non-system)",
			method:   "GET",
			path:     "/users/:id",
			ownerID:  "user-1",
			isSystem: false,
			urlParams: gin.Params{{Key: "id", Value: "user-1"}},
			setup:    func(_ *mockRepo) {},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:     "GetUser not found",
			method:   "GET",
			path:     "/users/:id",
			ownerID:  "admin-1",
			isSystem: true,
			urlParams: gin.Params{{Key: "id", Value: "user-404"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "user-404").Return(models.User{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:     "GetUser repo error",
			method:   "GET",
			path:     "/users/:id",
			ownerID:  "admin-1",
			isSystem: true,
			urlParams: gin.Params{{Key: "id", Value: "user-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "user-1").Return(models.User{}, errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── POST /users ──────────────────────────────────────────────
		{
			name:     "CreateUser first user (becomes system)",
			method:   "POST",
			path:     "/users",
			body:     `{"name":"Admin","email":"admin@test.com","password":"pass1234"}`,
			ownerID:  "admin-1",
			isSystem: true,
			setup: func(m *mockRepo) {
				m.On("SystemUserExists", mock.Anything).Return(false, nil)
				m.On("Create", mock.Anything, mock.MatchedBy(func(p CreateUserPayload) bool {
					return p.Name == "Admin" && p.Email == "admin@test.com"
				}), true).Return(models.User{ID: "new-1", Name: "Admin", IsSystemUser: true}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			name:     "CreateUser subsequent user",
			method:   "POST",
			path:     "/users",
			body:     `{"name":"User","email":"user@test.com","password":"pass1234"}`,
			ownerID:  "admin-1",
			isSystem: true,
			setup: func(m *mockRepo) {
				m.On("SystemUserExists", mock.Anything).Return(true, nil)
				m.On("Create", mock.Anything, mock.MatchedBy(func(p CreateUserPayload) bool {
					return p.Name == "User"
				}), false).Return(models.User{ID: "new-2", Name: "User"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			name:     "CreateUser forbidden (non-system)",
			method:   "POST",
			path:     "/users",
			body:     `{"name":"User","email":"user@test.com","password":"pass1234"}`,
			ownerID:  "user-1",
			isSystem: false,
			setup:    func(_ *mockRepo) {},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:     "CreateUser validation error",
			method:   "POST",
			path:     "/users",
			body:     `{"name":"","email":"bad","password":"12"}`,
			ownerID:  "admin-1",
			isSystem: true,
			setup:    func(_ *mockRepo) {},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:     "CreateUser email taken",
			method:   "POST",
			path:     "/users",
			body:     `{"name":"User","email":"taken@test.com","password":"pass1234"}`,
			ownerID:  "admin-1",
			isSystem: true,
			setup: func(m *mockRepo) {
				m.On("SystemUserExists", mock.Anything).Return(true, nil)
				m.On("Create", mock.Anything, mock.Anything, false).Return(models.User{}, ErrEmailTaken)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, w.Code)
			},
		},
		// ── PUT /users/:id ───────────────────────────────────────────
		{
			name:     "UpdateUser own profile",
			method:   "PUT",
			path:     "/users/:id",
			body:     `{"name":"Updated","email":"updated@test.com"}`,
			ownerID:  "user-1",
			isSystem: false,
			urlParams: gin.Params{{Key: "id", Value: "user-1"}},
			setup: func(m *mockRepo) {
				m.On("Update", mock.Anything, "user-1", mock.MatchedBy(func(p UpdateUserPayload) bool {
					return p.Name == "Updated" && p.PetLimit == nil
				})).Return(models.User{ID: "user-1", Name: "Updated"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:     "UpdateUser system updates pet_limit",
			method:   "PUT",
			path:     "/users/:id",
			body:     `{"name":"User","email":"user@test.com","pet_limit":10}`,
			ownerID:  "admin-1",
			isSystem: true,
			urlParams: gin.Params{{Key: "id", Value: "user-2"}},
			setup: func(m *mockRepo) {
				m.On("Update", mock.Anything, "user-2", mock.MatchedBy(func(p UpdateUserPayload) bool {
					return p.PetLimit != nil && *p.PetLimit == 10
				})).Return(models.User{ID: "user-2", Name: "User"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:     "UpdateUser forbidden (different user, not system)",
			method:   "PUT",
			path:     "/users/:id",
			body:     `{"name":"Hacker","email":"hacker@test.com"}`,
			ownerID:  "user-1",
			isSystem: false,
			urlParams: gin.Params{{Key: "id", Value: "user-2"}},
			setup:    func(_ *mockRepo) {},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:     "UpdateUser not found",
			method:   "PUT",
			path:     "/users/:id",
			body:     `{"name":"Ghost","email":"ghost@test.com"}`,
			ownerID:  "user-404",
			isSystem: true,
			urlParams: gin.Params{{Key: "id", Value: "user-404"}},
			setup: func(m *mockRepo) {
				m.On("Update", mock.Anything, "user-404", mock.Anything).Return(models.User{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		// ── DELETE /users/:id ────────────────────────────────────────
		{
			name:     "DeleteUser own profile (non-system)",
			method:   "DELETE",
			path:     "/users/:id",
			ownerID:  "user-1",
			isSystem: false,
			urlParams: gin.Params{{Key: "id", Value: "user-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "user-1").Return(models.User{ID: "user-1", Name: "User", IsSystemUser: false}, nil)
				m.On("Delete", mock.Anything, "user-1").Return(nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:     "DeleteUser system user cannot delete system user",
			method:   "DELETE",
			path:     "/users/:id",
			ownerID:  "admin-1",
			isSystem: true,
			urlParams: gin.Params{{Key: "id", Value: "admin-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "admin-1").Return(models.User{ID: "admin-1", IsSystemUser: true}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:     "DeleteUser forbidden (different user, not system)",
			method:   "DELETE",
			path:     "/users/:id",
			ownerID:  "user-1",
			isSystem: false,
			urlParams: gin.Params{{Key: "id", Value: "user-2"}},
			setup:    func(_ *mockRepo) {},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:     "DeleteUser not found",
			method:   "DELETE",
			path:     "/users/:id",
			ownerID:  "user-404",
			isSystem: false,
			urlParams: gin.Params{{Key: "id", Value: "user-404"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "user-404").Return(models.User{}, ErrNotFound)
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
