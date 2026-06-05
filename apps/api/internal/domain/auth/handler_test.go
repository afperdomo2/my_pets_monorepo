package auth

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/config"
	"github.com/my-pets/api/internal/domain/user"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

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
	c.Request = httptest.NewRequest("GET", "/", nil)
	return c, w
}

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret: "test-secret-for-testing",
		GinMode:   "test",
	}
}

func hashPassword(t *testing.T, pwd string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	require.NoError(t, err)
	return string(h)
}

var defaultUser = models.User{ID: "user-1", Name: "Test", Email: "test@test.com"}

func TestLogin_Success(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	pwdHash := hashPassword(t, "password123")
	m.On("GetPasswordByEmail", mock.Anything, "test@test.com").Return(pwdHash, nil)
	m.On("GetByEmail", mock.Anything, "test@test.com").Return(defaultUser, nil)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/api/v1/auth/login",
		bytes.NewBufferString(`{"email":"test@test.com","password":"password123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	require.Equal(t, http.StatusOK, w.Code)
	cookies := w.Result().Cookies()
	cookieMap := make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	require.Contains(t, cookieMap, "access_token")
	require.Contains(t, cookieMap, "refresh_token")
	require.NotEmpty(t, cookieMap["access_token"])
	require.NotEmpty(t, cookieMap["refresh_token"])
	m.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	pwdHash := hashPassword(t, "password123")
	m.On("GetPasswordByEmail", mock.Anything, "test@test.com").Return(pwdHash, nil)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/api/v1/auth/login",
		bytes.NewBufferString(`{"email":"test@test.com","password":"wrong"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	m.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	m.On("GetPasswordByEmail", mock.Anything, "test@test.com").Return("", user.ErrNotFound)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/api/v1/auth/login",
		bytes.NewBufferString(`{"email":"test@test.com","password":"password123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	m.AssertExpectations(t)
}

func TestLogin_ValidationError(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/api/v1/auth/login",
		bytes.NewBufferString(`{"email":"","password":""}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertExpectations(t)
}

func TestLogout(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/api/v1/auth/logout", nil)

	handler.Logout(c)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Header().Get("Set-Cookie"), "access_token")
	m.AssertExpectations(t)
}

func TestMe_Success(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	m.On("GetByID", mock.Anything, "user-1").Return(defaultUser, nil)

	c, w := setupGin()
	c.Set("userID", "user-1")

	handler.Me(c)

	require.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestMe_UserNotFound(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	m.On("GetByID", mock.Anything, "user-1").Return(models.User{}, user.ErrNotFound)

	c, w := setupGin()
	c.Set("userID", "user-1")

	handler.Me(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	m.AssertExpectations(t)
}

func TestUpdateProfile_Success(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	m.On("Update", mock.Anything, "user-1", mock.MatchedBy(func(p user.UpdateUserPayload) bool {
		return p.Name == "Updated" && p.Email == "updated@test.com"
	})).Return(models.User{ID: "user-1", Name: "Updated", Email: "updated@test.com"}, nil)

	c, w := setupGin()
	c.Set("userID", "user-1")
	c.Request = httptest.NewRequest("PUT", "/api/v1/auth/profile",
		bytes.NewBufferString(`{"name":"Updated","email":"updated@test.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateProfile(c)

	require.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestUpdateProfile_Conflict(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	m.On("Update", mock.Anything, "user-1", mock.Anything).Return(models.User{}, user.ErrEmailTaken)

	c, w := setupGin()
	c.Set("userID", "user-1")
	c.Request = httptest.NewRequest("PUT", "/api/v1/auth/profile",
		bytes.NewBufferString(`{"name":"Test","email":"taken@test.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateProfile(c)

	require.Equal(t, http.StatusConflict, w.Code)
	m.AssertExpectations(t)
}

func TestRefresh_Success(t *testing.T) {
	m := new(mockUserRepo)
	cfg := testConfig()
	handler := NewHandler(cfg, m)

	token, err := GenerateRefreshToken(cfg, defaultUser)
	require.NoError(t, err)

	m.On("GetByID", mock.Anything, "user-1").Return(defaultUser, nil)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: token})

	handler.Refresh(c)

	require.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestRefresh_InvalidToken(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: "invalid-token"})

	handler.Refresh(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	m.AssertExpectations(t)
}

func TestRefresh_NoCookie(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	c, w := setupGin()
	c.Request = httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)

	handler.Refresh(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	m.AssertExpectations(t)
}

func TestChangePassword_Success(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	oldHash := hashPassword(t, "old-password")
	m.On("GetPasswordByEmail", mock.Anything, "test@test.com").Return(oldHash, nil)
	m.On("UpdatePassword", mock.Anything, "user-1", mock.Anything).Return(nil)

	c, w := setupGin()
	c.Set("userID", "user-1")
	c.Set("email", "test@test.com")
	c.Request = httptest.NewRequest("PUT", "/api/v1/auth/password",
		bytes.NewBufferString(`{"current_password":"old-password","new_password":"new-password","confirm_password":"new-password"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ChangePassword(c)

	require.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestChangePassword_WrongCurrentPassword(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	oldHash := hashPassword(t, "old-password")
	m.On("GetPasswordByEmail", mock.Anything, "test@test.com").Return(oldHash, nil)

	c, w := setupGin()
	c.Set("userID", "user-1")
	c.Set("email", "test@test.com")
	c.Request = httptest.NewRequest("PUT", "/api/v1/auth/password",
		bytes.NewBufferString(`{"current_password":"wrong","new_password":"new-password","confirm_password":"new-password"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ChangePassword(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	m.AssertExpectations(t)
}

func TestChangePassword_PasswordsDoNotMatch(t *testing.T) {
	m := new(mockUserRepo)
	handler := NewHandler(testConfig(), m)

	c, w := setupGin()
	c.Set("userID", "user-1")
	c.Request = httptest.NewRequest("PUT", "/api/v1/auth/password",
		bytes.NewBufferString(`{"current_password":"old","new_password":"new1","confirm_password":"new2"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ChangePassword(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertExpectations(t)
}
