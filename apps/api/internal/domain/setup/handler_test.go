package setup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/domain/user"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

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

func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestStatus_NeedsSetup(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(m)

	m.On("HasUsers", mock.Anything).Return(false, nil)

	c, w := setupGin()
	c.Request = httptest.NewRequest("GET", "/setup/status", nil)

	handler.Status(c)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]bool
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.True(t, resp["needs_setup"])
	m.AssertExpectations(t)
}

func TestStatus_AlreadySetUp(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(m)

	m.On("HasUsers", mock.Anything).Return(true, nil)

	c, w := setupGin()
	c.Request = httptest.NewRequest("GET", "/setup/status", nil)

	handler.Status(c)

	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]bool
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.False(t, resp["needs_setup"])
	m.AssertExpectations(t)
}

func TestStatus_RepoError(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(m)

	m.On("HasUsers", mock.Anything).Return(false, errTest)

	c, w := setupGin()
	c.Request = httptest.NewRequest("GET", "/setup/status", nil)

	handler.Status(c)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}

func TestCreate_FirstUser(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(m)

	m.On("HasUsers", mock.Anything).Return(false, nil)
	m.On("Create", mock.Anything, mock.MatchedBy(func(p user.CreateUserPayload) bool {
		return p.Name == "Admin" && p.Email == "admin@test.com"
	}), true).Return(models.User{ID: "u1", Name: "Admin", IsSystemUser: true}, nil)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/setup", bytes.NewBufferString(`{"name":"Admin","email":"admin@test.com","password":"secret123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	require.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestCreate_AlreadyInitialized(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(m)

	m.On("HasUsers", mock.Anything).Return(true, nil)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/setup", bytes.NewBufferString(`{"name":"Admin","email":"admin@test.com","password":"secret123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	require.Equal(t, http.StatusConflict, w.Code)
	m.AssertExpectations(t)
}

func TestCreate_ValidationError(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(m)

	m.On("HasUsers", mock.Anything).Return(false, nil)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/setup", bytes.NewBufferString(`{"name":"","email":"bad","password":"12"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertExpectations(t)
}

func TestCreate_EmailTaken(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(m)

	m.On("HasUsers", mock.Anything).Return(false, nil)
	m.On("Create", mock.Anything, mock.Anything, true).Return(models.User{}, user.ErrEmailTaken)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/setup", bytes.NewBufferString(`{"name":"Admin","email":"taken@test.com","password":"secret123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	require.Equal(t, http.StatusConflict, w.Code)
	m.AssertExpectations(t)
}

func TestCreate_HasUsersError(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(m)

	m.On("HasUsers", mock.Anything).Return(false, errTest)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/setup", bytes.NewBufferString(`{"name":"Admin","email":"admin@test.com","password":"secret123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}

func TestCreate_CreateError(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(m)

	m.On("HasUsers", mock.Anything).Return(false, nil)
	m.On("Create", mock.Anything, mock.Anything, true).Return(models.User{}, errTest)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/setup", bytes.NewBufferString(`{"name":"Admin","email":"admin@test.com","password":"secret123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}
