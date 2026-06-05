package health_record

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
	healthCatalog "github.com/my-pets/api/internal/domain/health_catalog"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

func uid(n int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", n)
}

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetAllByOwner(ctx context.Context, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error) {
	args := m.Called(ctx, ownerID, page, perPage)
	return args.Get(0).([]models.HealthRecord), args.Get(1).(int64), args.Error(2)
}
func (m *mockRepo) GetByPetID(ctx context.Context, petID, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error) {
	args := m.Called(ctx, petID, ownerID, page, perPage)
	return args.Get(0).([]models.HealthRecord), args.Get(1).(int64), args.Error(2)
}
func (m *mockRepo) GetByPetIDAndCategory(ctx context.Context, petID, category, ownerID string, page, perPage int) ([]models.HealthRecord, int64, error) {
	args := m.Called(ctx, petID, category, ownerID, page, perPage)
	return args.Get(0).([]models.HealthRecord), args.Get(1).(int64), args.Error(2)
}
func (m *mockRepo) GetByID(ctx context.Context, id, ownerID string) (models.HealthRecord, error) {
	args := m.Called(ctx, id, ownerID)
	return args.Get(0).(models.HealthRecord), args.Error(1)
}
func (m *mockRepo) Create(ctx context.Context, ownerID string, payload CreateHealthRecordPayload) (models.HealthRecord, error) {
	args := m.Called(ctx, ownerID, payload)
	return args.Get(0).(models.HealthRecord), args.Error(1)
}
func (m *mockRepo) Update(ctx context.Context, id, ownerID string, payload UpdateHealthRecordPayload) (models.HealthRecord, error) {
	args := m.Called(ctx, id, ownerID, payload)
	return args.Get(0).(models.HealthRecord), args.Error(1)
}
func (m *mockRepo) Delete(ctx context.Context, id, ownerID string) error {
	args := m.Called(ctx, id, ownerID)
	return args.Error(0)
}
func (m *mockRepo) GetUpcomingRecords(ctx context.Context, ownerID, category string, page, perPage int) ([]models.HealthRecord, int64, error) {
	args := m.Called(ctx, ownerID, category, page, perPage)
	return args.Get(0).([]models.HealthRecord), args.Get(1).(int64), args.Error(2)
}
func (m *mockRepo) UpdateLastDoseDate(ctx context.Context, healthRecordID string, lastDoseDate string) error {
	args := m.Called(ctx, healthRecordID, lastDoseDate)
	return args.Error(0)
}
func (m *mockRepo) IncrementAppliedDosesCount(ctx context.Context, healthRecordID string) error {
	args := m.Called(ctx, healthRecordID)
	return args.Error(0)
}

type mockHealthCatalogRepo struct{ mock.Mock }

func (m *mockHealthCatalogRepo) GetPaginatedByCategory(ctx context.Context, category string, page, perPage int, speciesFilter *string) ([]models.HealthCatalog, int64, error) {
	args := m.Called(ctx, category, page, perPage, speciesFilter)
	return args.Get(0).([]models.HealthCatalog), args.Get(1).(int64), args.Error(2)
}
func (m *mockHealthCatalogRepo) GetBySpeciesAndCategory(ctx context.Context, species, category string) ([]models.HealthCatalog, error) {
	args := m.Called(ctx, species, category)
	return args.Get(0).([]models.HealthCatalog), args.Error(1)
}
func (m *mockHealthCatalogRepo) GetByID(ctx context.Context, id string) (models.HealthCatalog, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.HealthCatalog), args.Error(1)
}
func (m *mockHealthCatalogRepo) Create(ctx context.Context, payload healthCatalog.CreateHealthCatalogPayload) (models.HealthCatalog, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(models.HealthCatalog), args.Error(1)
}
func (m *mockHealthCatalogRepo) Update(ctx context.Context, id string, payload healthCatalog.UpdateHealthCatalogPayload) (models.HealthCatalog, error) {
	args := m.Called(ctx, id, payload)
	return args.Get(0).(models.HealthCatalog), args.Error(1)
}
func (m *mockHealthCatalogRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockVaccineAppRepo struct{ mock.Mock }

func (m *mockVaccineAppRepo) GetByHealthRecordID(ctx context.Context, healthRecordID string) ([]models.VaccineApplication, error) {
	args := m.Called(ctx, healthRecordID)
	return args.Get(0).([]models.VaccineApplication), args.Error(1)
}
func (m *mockVaccineAppRepo) GetByID(ctx context.Context, id string) (models.VaccineApplication, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.VaccineApplication), args.Error(1)
}
func (m *mockVaccineAppRepo) Create(ctx context.Context, app models.VaccineApplication) (models.VaccineApplication, error) {
	args := m.Called(ctx, app)
	return args.Get(0).(models.VaccineApplication), args.Error(1)
}
func (m *mockVaccineAppRepo) Update(ctx context.Context, app models.VaccineApplication) (models.VaccineApplication, error) {
	args := m.Called(ctx, app)
	return args.Get(0).(models.VaccineApplication), args.Error(1)
}
func (m *mockVaccineAppRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockVaccineAppRepo) UpdateHealthRecordAfterApplication(ctx context.Context, healthRecordID, applicationDate string, nextDoseDate *string) error {
	args := m.Called(ctx, healthRecordID, applicationDate, nextDoseDate)
	return args.Error(0)
}

func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

type hrTestCase struct {
	name       string
	method     string
	path       string
	body       string
	ownerID    string
	urlParams  gin.Params
	setupRepo  func(*mockRepo)
	setupHC    func(*mockHealthCatalogRepo)
	setupVA    func(*mockVaccineAppRepo)
	check      func(*testing.T, *httptest.ResponseRecorder)
}

func routeFromPath(path string) string {
	for i := range path {
		if path[i] == '?' {
			return path[:i]
		}
	}
	return path
}

func (tc *hrTestCase) exec(t *testing.T) {
	t.Helper()

	mr := new(mockRepo)
	mhc := new(mockHealthCatalogRepo)
	mva := new(mockVaccineAppRepo)
	handler := NewHandler(mr, mhc, mva)

	c, w := setupGin()
	c.Request = httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = tc.urlParams
	if tc.ownerID != "" {
		c.Set("userID", tc.ownerID)
	}

	if tc.setupRepo != nil {
		tc.setupRepo(mr)
	}
	if tc.setupHC != nil {
		tc.setupHC(mhc)
	}
	if tc.setupVA != nil {
		tc.setupVA(mva)
	}

	switch tc.method + " " + routeFromPath(tc.path) {
	case "GET /health-records":
		handler.GetAllHealthRecords(c)
	case "GET /health-records/pets/:pet_id":
		handler.GetHealthRecordsByPet(c)
	case "GET /health-records/pets/:pet_id/category/:category":
		handler.GetHealthRecordsByPetAndCategory(c)
	case "GET /health-records/:record_id":
		handler.GetHealthRecordByID(c)
	case "POST /health-records":
		handler.CreateHealthRecord(c)
	case "PUT /health-records/:record_id":
		handler.UpdateHealthRecord(c)
	case "DELETE /health-records/:record_id":
		handler.DeleteHealthRecord(c)
	case "GET /health-records/upcoming":
		handler.GetUpcomingRecords(c)
	default:
		t.Fatalf("unknown route: %s %s", tc.method, tc.path)
	}

	tc.check(t, w)
	mr.AssertExpectations(t)
	mhc.AssertExpectations(t)
	mva.AssertExpectations(t)
}

func TestHealthRecordHandler(t *testing.T) {
	petID := uid(1)
	recordID := uid(10)
	catID := uid(20)
	userID := uid(100)

	tests := []hrTestCase{
		// ── GET /health-records ───────────────────────────────
		{
			name: "GetAll success", method: "GET", path: "/health-records?page=1&per_page=10",
			ownerID: userID, setupRepo: func(m *mockRepo) {
				m.On("GetAllByOwner", mock.Anything, userID, 1, 10).
					Return([]models.HealthRecord{{ID: recordID}}, int64(1), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "GetAll repo error", method: "GET", path: "/health-records",
			ownerID: userID, setupRepo: func(m *mockRepo) {
				m.On("GetAllByOwner", mock.Anything, userID, 1, 10).
					Return([]models.HealthRecord{}, int64(0), errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── GET /health-records/pets/:pet_id ───────────────────
		{
			name: "GetByPet success", method: "GET", path: "/health-records/pets/:pet_id",
			ownerID: userID, urlParams: gin.Params{{Key: "pet_id", Value: petID}},
			setupRepo: func(m *mockRepo) {
				m.On("GetByPetID", mock.Anything, petID, userID, 1, 10).
					Return([]models.HealthRecord{{ID: recordID}}, int64(1), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "GetByPet invalid pet_id", method: "GET", path: "/health-records/pets/:pet_id",
			ownerID: userID, check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		// ── GET /health-records/pets/:pet_id/category/:category ─
		{
			name: "GetByPetAndCategory success", method: "GET", path: "/health-records/pets/:pet_id/category/:category",
			ownerID: userID, urlParams: gin.Params{{Key: "pet_id", Value: petID}, {Key: "category", Value: "vaccine"}},
			setupRepo: func(m *mockRepo) {
				m.On("GetByPetIDAndCategory", mock.Anything, petID, "vaccine", userID, 1, 10).
					Return([]models.HealthRecord{{ID: recordID}}, int64(1), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "GetByPetAndCategory invalid pet_id", method: "GET", path: "/health-records/pets/:pet_id/category/:category",
			ownerID: userID, urlParams: gin.Params{{Key: "category", Value: "vaccine"}},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "GetByPetAndCategory invalid category", method: "GET", path: "/health-records/pets/:pet_id/category/:category",
			ownerID: userID, urlParams: gin.Params{{Key: "pet_id", Value: petID}, {Key: "category", Value: "surgery"}},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		// ── GET /health-records/:record_id ─────────────────────
		{
			name: "GetByID success", method: "GET", path: "/health-records/:record_id",
			ownerID: userID, urlParams: gin.Params{{Key: "record_id", Value: recordID}},
			setupRepo: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, recordID, userID).
					Return(models.HealthRecord{ID: recordID, Name: "Vaccine"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var resp map[string]models.HealthRecord
				json.Unmarshal(w.Body.Bytes(), &resp)
				require.Equal(t, "Vaccine", resp["data"].Name)
			},
		},
		{
			name: "GetByID not found", method: "GET", path: "/health-records/:record_id",
			ownerID: userID, urlParams: gin.Params{{Key: "record_id", Value: uid(404)}},
			setupRepo: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, uid(404), userID).
					Return(models.HealthRecord{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		// ── POST /health-records ──────────────────────────────
		{
			name: "Create success", method: "POST", path: "/health-records",
			body:    fmt.Sprintf(`{"pet_id":"%s","category":"vaccine","name":"Rabies","application_date":"2026-01-15"}`, petID),
			ownerID: userID,
			setupRepo: func(m *mockRepo) {
				m.On("Create", mock.Anything, userID, mock.MatchedBy(func(p CreateHealthRecordPayload) bool {
					return p.Name == "Rabies"
				})).Return(models.HealthRecord{ID: recordID, Name: "Rabies"}, nil)
			},
			setupVA: func(m *mockVaccineAppRepo) {
				m.On("Create", mock.Anything, mock.Anything).Return(models.VaccineApplication{}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			name: "Create validation error", method: "POST", path: "/health-records",
			body:    `{"pet_id":"","category":"","name":"","application_date":""}`,
			ownerID: userID,
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Create with catalog", method: "POST", path: "/health-records",
			body:    fmt.Sprintf(`{"pet_id":"%s","health_catalog_id":"%s","application_date":"2026-01-15"}`, petID, catID),
			ownerID: userID,
			setupHC: func(m *mockHealthCatalogRepo) {
				m.On("GetByID", mock.Anything, catID).
					Return(models.HealthCatalog{ID: catID, Name: "CatalogItem", Category: "vaccine"}, nil)
			},
			setupRepo: func(m *mockRepo) {
				m.On("Create", mock.Anything, userID, mock.Anything).Return(models.HealthRecord{ID: recordID}, nil)
			},
			setupVA: func(m *mockVaccineAppRepo) {
				m.On("Create", mock.Anything, mock.Anything).Return(models.VaccineApplication{}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
			},
		},
		// ── PUT /health-records/:record_id ────────────────────
		{
			name: "Update success", method: "PUT", path: "/health-records/:record_id",
			body:    `{"category":"vaccine","name":"Updated"}`,
			ownerID: userID, urlParams: gin.Params{{Key: "record_id", Value: recordID}},
			setupRepo: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, recordID, userID).
					Return(models.HealthRecord{ID: recordID, Name: "Old", Category: "vaccine"}, nil)
				m.On("Update", mock.Anything, recordID, userID, mock.Anything).
					Return(models.HealthRecord{ID: recordID, Name: "Updated"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "Update not found", method: "PUT", path: "/health-records/:record_id",
			body:    `{"category":"vaccine","name":"Updated"}`,
			ownerID: userID, urlParams: gin.Params{{Key: "record_id", Value: uid(404)}},
			setupRepo: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, uid(404), userID).
					Return(models.HealthRecord{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		// ── DELETE /health-records/:record_id ─────────────────
		{
			name: "Delete success", method: "DELETE", path: "/health-records/:record_id",
			ownerID: userID, urlParams: gin.Params{{Key: "record_id", Value: recordID}},
			setupRepo: func(m *mockRepo) {
				m.On("Delete", mock.Anything, recordID, userID).Return(nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "Delete not found", method: "DELETE", path: "/health-records/:record_id",
			ownerID: userID, urlParams: gin.Params{{Key: "record_id", Value: uid(404)}},
			setupRepo: func(m *mockRepo) {
				m.On("Delete", mock.Anything, uid(404), userID).Return(ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		// ── GET /health-records/upcoming ──────────────────────
		{
			name: "GetUpcoming success", method: "GET", path: "/health-records/upcoming?page=1&per_page=10",
			ownerID: userID, setupRepo: func(m *mockRepo) {
				m.On("GetUpcomingRecords", mock.Anything, userID, "", 1, 10).
					Return([]models.HealthRecord{{ID: recordID}}, int64(1), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "GetUpcoming with category filter", method: "GET", path: "/health-records/upcoming?category=vaccine",
			ownerID: userID, setupRepo: func(m *mockRepo) {
				m.On("GetUpcomingRecords", mock.Anything, userID, "vaccine", 1, 10).
					Return([]models.HealthRecord{}, int64(0), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
