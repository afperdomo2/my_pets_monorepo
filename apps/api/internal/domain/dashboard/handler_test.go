package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetSummary(ctx context.Context, ownerID string) (DashboardSummary, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(DashboardSummary), args.Error(1)
}

func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	return c, w
}

func TestGetSummary_Success(t *testing.T) {
	m := new(mockRepo)
	handler := NewHandler(m)

	expected := DashboardSummary{TotalPets: 5, HealthyPets: 3, PendingTasks: 2, OverdueTasks: 1}
	m.On("GetSummary", mock.Anything, "user-1").Return(expected, nil)

	c, w := setupGin()
	c.Set("userID", "user-1")

	handler.GetSummary(c)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]DashboardSummary
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, expected, resp["data"])
	m.AssertExpectations(t)
}

func TestGetSummary_RepoError(t *testing.T) {
	m := new(mockRepo)
	handler := NewHandler(m)

	m.On("GetSummary", mock.Anything, "user-1").Return(DashboardSummary{}, errTest)

	c, w := setupGin()
	c.Set("userID", "user-1")

	handler.GetSummary(c)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	m.AssertExpectations(t)
}
