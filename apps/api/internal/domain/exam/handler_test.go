package exam

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

func (m *mockRepo) GetAllByOwner(ctx context.Context, ownerID string, page, perPage int) ([]models.Exam, int64, error) {
	args := m.Called(ctx, ownerID, page, perPage)
	return args.Get(0).([]models.Exam), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) GetByPetID(ctx context.Context, petID, ownerID string, page, perPage int) ([]models.Exam, int64, error) {
	args := m.Called(ctx, petID, ownerID, page, perPage)
	return args.Get(0).([]models.Exam), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) GetByID(ctx context.Context, id, ownerID string) (models.Exam, error) {
	args := m.Called(ctx, id, ownerID)
	return args.Get(0).(models.Exam), args.Error(1)
}

func (m *mockRepo) GetByIDWithResults(ctx context.Context, id, ownerID string) (models.Exam, []models.ExamResult, error) {
	args := m.Called(ctx, id, ownerID)
	return args.Get(0).(models.Exam), args.Get(1).([]models.ExamResult), args.Error(2)
}

func (m *mockRepo) Create(ctx context.Context, exam models.Exam) (models.Exam, error) {
	args := m.Called(ctx, exam)
	return args.Get(0).(models.Exam), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, exam models.Exam) (models.Exam, error) {
	args := m.Called(ctx, exam)
	return args.Get(0).(models.Exam), args.Error(1)
}

func (m *mockRepo) Delete(ctx context.Context, id, ownerID string) error {
	args := m.Called(ctx, id, ownerID)
	return args.Error(0)
}

func (m *mockRepo) CreateResults(ctx context.Context, results []models.ExamResult) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

func (m *mockRepo) DeleteResultsByExamID(ctx context.Context, examID string) error {
	args := m.Called(ctx, examID)
	return args.Error(0)
}

func (m *mockRepo) GetResultsByExamID(ctx context.Context, examID string) ([]models.ExamResult, error) {
	args := m.Called(ctx, examID)
	return args.Get(0).([]models.ExamResult), args.Error(1)
}

func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

type examTestCase struct {
	name      string
	method    string
	path      string
	body      string
	ownerID   string
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

func (tc *examTestCase) exec(t *testing.T) {
	t.Helper()

	m := new(mockRepo)
	handler := NewHandler(m)

	c, w := setupGin()
	c.Request = httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = tc.urlParams
	if tc.ownerID != "" {
		c.Set("userID", tc.ownerID)
	}

	if tc.setup != nil {
		tc.setup(m)
	}

	switch tc.method + " " + routeFromPath(tc.path) {
	case "GET /exams":
		handler.GetAllExams(c)
	case "GET /exams/pets/:pet_id":
		handler.GetExamsByPet(c)
	case "GET /exams/:id":
		handler.GetExamByID(c)
	case "POST /exams":
		handler.CreateExam(c)
	case "PUT /exams/:id":
		handler.UpdateExam(c)
	case "PATCH /exams/:id/schedule":
		handler.ScheduleExam(c)
	case "PATCH /exams/:id/complete":
		handler.CompleteExam(c)
	case "DELETE /exams/:id":
		handler.DeleteExam(c)
	default:
		t.Fatalf("unknown route: %s %s", tc.method, tc.path)
	}

	tc.check(t, w)
	m.AssertExpectations(t)
}

func TestExamHandler(t *testing.T) {
	tests := []examTestCase{
		// ── GET /exams ─────────────────────────────────────────
		{
			name: "GetAllExams success", method: "GET", path: "/exams?page=1&per_page=10",
			ownerID: "user-1",
			setup: func(m *mockRepo) {
				m.On("GetAllByOwner", mock.Anything, "user-1", 1, 10).
					Return([]models.Exam{{ID: "ex-1", Name: "Blood Test"}}, int64(1), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var resp map[string]any
				json.Unmarshal(w.Body.Bytes(), &resp)
				require.Equal(t, float64(1), resp["total"])
			},
		},
		{
			name: "GetAllExams repo error", method: "GET", path: "/exams",
			ownerID: "user-1",
			setup: func(m *mockRepo) {
				m.On("GetAllByOwner", mock.Anything, "user-1", 1, 10).
					Return([]models.Exam{}, int64(0), errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── GET /exams/pets/:pet_id ────────────────────────────
		{
			name: "GetExamsByPet success", method: "GET", path: "/exams/pets/:pet_id",
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "pet_id", Value: "pet-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByPetID", mock.Anything, "pet-1", "user-1", 1, 10).
					Return([]models.Exam{{ID: "ex-1"}}, int64(1), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "GetExamsByPet invalid pet_id", method: "GET", path: "/exams/pets/:pet_id",
			ownerID: "user-1",
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		// ── GET /exams/:id ─────────────────────────────────────
		{
			name: "GetExamByID success", method: "GET", path: "/exams/:id",
			ownerID:   "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByIDWithResults", mock.Anything, "ex-1", "user-1").
					Return(models.Exam{ID: "ex-1", Name: "Blood Test"}, []models.ExamResult{}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var resp map[string]json.RawMessage
				json.Unmarshal(w.Body.Bytes(), &resp)
				require.Contains(t, string(resp["data"]), "Blood Test")
			},
		},
		{
			name: "GetExamByID not found", method: "GET", path: "/exams/:id",
			ownerID:   "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-404"}},
			setup: func(m *mockRepo) {
				m.On("GetByIDWithResults", mock.Anything, "ex-404", "user-1").
					Return(models.Exam{}, []models.ExamResult{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "GetExamByID invalid id", method: "GET", path: "/exams/:id",
			ownerID: "user-1",
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		// ── POST /exams ────────────────────────────────────────
		{
			name: "CreateExam success", method: "POST", path: "/exams",
			body:    `{"pet_id":"00000000-0000-0000-0000-000000000001","name":"Blood Test"}`,
			ownerID: "user-1",
			setup: func(m *mockRepo) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(e models.Exam) bool {
					return e.Name == "Blood Test" && e.Status == "scheduled"
				})).Return(models.Exam{ID: "ex-new", Name: "Blood Test"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			name: "CreateExam validation error", method: "POST", path: "/exams",
			body:    `{"pet_id":"","name":""}`,
			ownerID: "user-1",
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "CreateExam pet not found", method: "POST", path: "/exams",
			body:    `{"pet_id":"00000000-0000-0000-0000-000000000001","name":"Blood Test"}`,
			ownerID: "user-1",
			setup: func(m *mockRepo) {
				m.On("Create", mock.Anything, mock.Anything).
					Return(models.Exam{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		// ── PUT /exams/:id ─────────────────────────────────────
		{
			name: "UpdateExam success", method: "PUT", path: "/exams/:id",
			body:    `{"name":"Updated Exam"}`,
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "ex-1", "user-1").
					Return(models.Exam{ID: "ex-1", Name: "Old Name"}, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(e models.Exam) bool {
					return e.Name == "Updated Exam"
				})).Return(models.Exam{ID: "ex-1", Name: "Updated Exam"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "UpdateExam not found", method: "PUT", path: "/exams/:id",
			body:    `{"name":"Updated"}`,
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-404"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "ex-404", "user-1").
					Return(models.Exam{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		// ── PATCH /exams/:id/schedule ──────────────────────────
		{
			name: "ScheduleExam success", method: "PATCH", path: "/exams/:id/schedule",
			body:    `{"scheduled_date":"2026-07-01"}`,
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "ex-1", "user-1").
					Return(models.Exam{ID: "ex-1", Status: "scheduled"}, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(e models.Exam) bool {
					return e.ScheduledDate != nil
				})).Return(models.Exam{ID: "ex-1"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "ScheduleExam not found", method: "PATCH", path: "/exams/:id/schedule",
			body:    `{"scheduled_date":"2026-07-01"}`,
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-404"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "ex-404", "user-1").
					Return(models.Exam{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		// ── PATCH /exams/:id/complete ──────────────────────────
		{
			name: "CompleteExam success", method: "PATCH", path: "/exams/:id/complete",
			body:    `{"completed_date":"2026-06-01"}`,
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "ex-1", "user-1").
					Return(models.Exam{ID: "ex-1", Status: "scheduled"}, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(e models.Exam) bool {
					return e.Status == "completed"
				})).Return(models.Exam{ID: "ex-1", Status: "completed"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "CompleteExam not found", method: "PATCH", path: "/exams/:id/complete",
			body:    `{"completed_date":"2026-06-01"}`,
			ownerID: "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-404"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "ex-404", "user-1").
					Return(models.Exam{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		// ── DELETE /exams/:id ──────────────────────────────────
		{
			name: "DeleteExam success", method: "DELETE", path: "/exams/:id",
			ownerID:   "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-1"}},
			setup: func(m *mockRepo) {
				m.On("Delete", mock.Anything, "ex-1", "user-1").Return(nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "DeleteExam not found", method: "DELETE", path: "/exams/:id",
			ownerID:   "user-1",
			urlParams: gin.Params{{Key: "id", Value: "ex-404"}},
			setup: func(m *mockRepo) {
				m.On("Delete", mock.Anything, "ex-404", "user-1").Return(ErrNotFound)
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
