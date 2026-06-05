package vaccine_application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

func hrID(n int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", n)
}

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetByHealthRecordID(ctx context.Context, healthRecordID string) ([]models.VaccineApplication, error) {
	args := m.Called(ctx, healthRecordID)
	return args.Get(0).([]models.VaccineApplication), args.Error(1)
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (models.VaccineApplication, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.VaccineApplication), args.Error(1)
}

func (m *mockRepo) Create(ctx context.Context, app models.VaccineApplication) (models.VaccineApplication, error) {
	args := m.Called(ctx, app)
	return args.Get(0).(models.VaccineApplication), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, app models.VaccineApplication) (models.VaccineApplication, error) {
	args := m.Called(ctx, app)
	return args.Get(0).(models.VaccineApplication), args.Error(1)
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) UpdateHealthRecordAfterApplication(ctx context.Context, healthRecordID string, applicationDate string, nextDoseDate *string) error {
	args := m.Called(ctx, healthRecordID, applicationDate, nextDoseDate)
	return args.Error(0)
}

func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

type appTestCase struct {
	name      string
	method    string
	path      string
	body      string
	urlParams gin.Params
	setup     func(*mockRepo)
	check     func(*testing.T, *httptest.ResponseRecorder)
}

func routeFromPath(path string) string {
	for i := range path {
		if path[i] == '?' {
			return path[:i]
		}
	}
	return path
}

func (tc *appTestCase) exec(t *testing.T) {
	t.Helper()

	m := new(mockRepo)
	handler := NewHandler(m)

	c, w := setupGin()
	c.Request = httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = tc.urlParams

	if tc.setup != nil {
		tc.setup(m)
	}

	switch tc.method + " " + routeFromPath(tc.path) {
	case "GET /health-record":
		handler.GetApplicationsByHealthRecord(c)
	case "GET /:id":
		handler.GetApplicationByID(c)
	case "POST /":
		handler.CreateApplication(c)
	case "PUT /:id":
		handler.UpdateApplication(c)
	case "DELETE /:id":
		handler.DeleteApplication(c)
	default:
		t.Fatalf("unknown route: %s %s", tc.method, tc.path)
	}

	tc.check(t, w)
	m.AssertExpectations(t)
}

func TestVaccineApplicationHandler(t *testing.T) {
	tests := []appTestCase{
		// ── GET /health-record/:id ──────────────────────────────
		{
			name:   "GetApplicationsByHealthRecord success",
			method: "GET", path: "/health-record",
			urlParams: gin.Params{{Key: "id", Value: hrID(1)}},
			setup: func(m *mockRepo) {
				m.On("GetByHealthRecordID", mock.Anything, hrID(1)).
					Return([]models.VaccineApplication{
						{ID: hrID(10), HealthRecordID: hrID(1)},
					}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var resp map[string][]models.VaccineApplication
				json.Unmarshal(w.Body.Bytes(), &resp)
				require.Len(t, resp["data"], 1)
			},
		},
		{
			name:   "GetApplicationsByHealthRecord empty param",
			method: "GET", path: "/health-record",
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "GetApplicationsByHealthRecord repo error",
			method: "GET", path: "/health-record",
			urlParams: gin.Params{{Key: "id", Value: hrID(1)}},
			setup: func(m *mockRepo) {
				m.On("GetByHealthRecordID", mock.Anything, hrID(1)).
					Return([]models.VaccineApplication{}, errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── GET /:id ────────────────────────────────────────────
		{
			name:   "GetApplicationByID success",
			method: "GET", path: "/:id",
			urlParams: gin.Params{{Key: "id", Value: hrID(10)}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, hrID(10)).
					Return(models.VaccineApplication{ID: hrID(10)}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:   "GetApplicationByID not found",
			method: "GET", path: "/:id",
			urlParams: gin.Params{{Key: "id", Value: hrID(404)}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, hrID(404)).
					Return(models.VaccineApplication{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:   "GetApplicationByID invalid id",
			method: "GET", path: "/:id",
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "GetApplicationByID repo error",
			method: "GET", path: "/:id",
			urlParams: gin.Params{{Key: "id", Value: hrID(10)}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, hrID(10)).
					Return(models.VaccineApplication{}, errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── POST / ──────────────────────────────────────────────
		{
			name:   "CreateApplication success",
			method: "POST", path: "/",
			body:   fmt.Sprintf(`{"health_record_id":"%s","application_date":"2026-01-15"}`, hrID(1)),
			setup: func(m *mockRepo) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(a models.VaccineApplication) bool {
					return a.HealthRecordID == hrID(1)
				})).Return(models.VaccineApplication{ID: hrID(20), HealthRecordID: hrID(1)}, nil)
				m.On("UpdateHealthRecordAfterApplication", mock.Anything, hrID(1), "2026-01-15", (*string)(nil)).
					Return(nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			name:   "CreateApplication validation error",
			method: "POST", path: "/",
			body:   `{"health_record_id":"","application_date":""}`,
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "CreateApplication invalid date",
			method: "POST", path: "/",
			body:   fmt.Sprintf(`{"health_record_id":"%s","application_date":"not-a-date"}`, hrID(1)),
			setup:  func(_ *mockRepo) {},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "CreateApplication repo error",
			method: "POST", path: "/",
			body:   fmt.Sprintf(`{"health_record_id":"%s","application_date":"2026-01-15"}`, hrID(1)),
			setup: func(m *mockRepo) {
				m.On("Create", mock.Anything, mock.Anything).
					Return(models.VaccineApplication{}, errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── PUT /:id ────────────────────────────────────────────
		{
			name:   "UpdateApplication success",
			method: "PUT", path: "/:id",
			body:      `{"application_date":"2026-02-01"}`,
			urlParams: gin.Params{{Key: "id", Value: hrID(10)}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, hrID(10)).
					Return(models.VaccineApplication{ID: hrID(10), HealthRecordID: hrID(1)}, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(a models.VaccineApplication) bool {
					return a.ID == hrID(10)
				})).Return(models.VaccineApplication{ID: hrID(10)}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:   "UpdateApplication not found",
			method: "PUT", path: "/:id",
			body:      `{"application_date":"2026-02-01"}`,
			urlParams: gin.Params{{Key: "id", Value: hrID(404)}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, hrID(404)).
					Return(models.VaccineApplication{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:   "UpdateApplication invalid id",
			method: "PUT", path: "/:id",
			body:   `{"application_date":"2026-02-01"}`,
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "UpdateApplication invalid date",
			method: "PUT", path: "/:id",
			body:      `{"application_date":"bad-date"}`,
			urlParams: gin.Params{{Key: "id", Value: hrID(10)}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, hrID(10)).
					Return(models.VaccineApplication{ID: hrID(10), HealthRecordID: hrID(1)}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		// ── DELETE /:id ─────────────────────────────────────────
		{
			name:   "DeleteApplication success",
			method: "DELETE", path: "/:id",
			urlParams: gin.Params{{Key: "id", Value: hrID(10)}},
			setup: func(m *mockRepo) {
				m.On("Delete", mock.Anything, hrID(10)).Return(nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:   "DeleteApplication not found",
			method: "DELETE", path: "/:id",
			urlParams: gin.Params{{Key: "id", Value: hrID(404)}},
			setup: func(m *mockRepo) {
				m.On("Delete", mock.Anything, hrID(404)).Return(ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:   "DeleteApplication invalid id",
			method: "DELETE", path: "/:id",
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "DeleteApplication repo error",
			method: "DELETE", path: "/:id",
			urlParams: gin.Params{{Key: "id", Value: hrID(10)}},
			setup: func(m *mockRepo) {
				m.On("Delete", mock.Anything, hrID(10)).Return(errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
