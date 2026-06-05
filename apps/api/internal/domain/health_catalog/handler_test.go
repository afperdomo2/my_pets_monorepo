package health_catalog

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

func (m *mockRepo) GetPaginatedByCategory(ctx context.Context, category string, page, perPage int, speciesFilter *string) ([]models.HealthCatalog, int64, error) {
	args := m.Called(ctx, category, page, perPage, speciesFilter)
	return args.Get(0).([]models.HealthCatalog), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) GetBySpeciesAndCategory(ctx context.Context, species, category string) ([]models.HealthCatalog, error) {
	args := m.Called(ctx, species, category)
	return args.Get(0).([]models.HealthCatalog), args.Error(1)
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (models.HealthCatalog, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.HealthCatalog), args.Error(1)
}

func (m *mockRepo) Create(ctx context.Context, payload CreateHealthCatalogPayload) (models.HealthCatalog, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(models.HealthCatalog), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, id string, payload UpdateHealthCatalogPayload) (models.HealthCatalog, error) {
	args := m.Called(ctx, id, payload)
	return args.Get(0).(models.HealthCatalog), args.Error(1)
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	return c, w
}

type hcTestCase struct {
	name      string
	method    string
	path      string
	body      string
	system    bool
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

func (tc *hcTestCase) exec(t *testing.T) {
	t.Helper()

	m := new(mockRepo)
	handler := NewHandler(m)

	c, w := setupGin()
	c.Request = httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = tc.urlParams
	if tc.system {
		c.Set("isSystemUser", true)
	}

	if tc.setup != nil {
		tc.setup(m)
	}

	switch tc.method + " " + routeFromPath(tc.path) {
	case "GET /category/:category":
		handler.GetHealthCatalogsByCategory(c)
	case "GET /species/:species/category/:category":
		handler.GetHealthCatalogsBySpeciesAndCategory(c)
	case "GET /:id":
		handler.GetHealthCatalogByID(c)
	case "POST /":
		handler.CreateHealthCatalog(c)
	case "PUT /:id":
		handler.UpdateHealthCatalog(c)
	case "DELETE /:id":
		handler.DeleteHealthCatalog(c)
	default:
		t.Fatalf("unknown route: %s %s", tc.method, tc.path)
	}

	tc.check(t, w)
	m.AssertExpectations(t)
}

func TestHealthCatalogHandler(t *testing.T) {
	tests := []hcTestCase{
		// ── GET /category/:category ────────────────────────────
		{
			name:   "GetByCategory success",
			method: "GET", path: "/category/:category",
			urlParams: gin.Params{{Key: "category", Value: "vaccine"}},
			setup: func(m *mockRepo) {
				m.On("GetPaginatedByCategory", mock.Anything, "vaccine", 1, 10, (*string)(nil)).
					Return([]models.HealthCatalog{{ID: "hc-1", Name: "Rabies"}}, int64(1), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:   "GetByCategory invalid category",
			method: "GET", path: "/category/:category",
			urlParams: gin.Params{{Key: "category", Value: "invalid"}},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "GetByCategory with species filter",
			method: "GET", path: "/category/:category?species=dog",
			urlParams: gin.Params{{Key: "category", Value: "vaccine"}},
			setup: func(m *mockRepo) {
				species := "dog"
				m.On("GetPaginatedByCategory", mock.Anything, "vaccine", 1, 10, &species).
					Return([]models.HealthCatalog{}, int64(0), nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:   "GetByCategory invalid species filter",
			method: "GET", path: "/category/:category?species=dragon",
			urlParams: gin.Params{{Key: "category", Value: "vaccine"}},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "GetByCategory repo error",
			method: "GET", path: "/category/:category",
			urlParams: gin.Params{{Key: "category", Value: "vaccine"}},
			setup: func(m *mockRepo) {
				m.On("GetPaginatedByCategory", mock.Anything, "vaccine", 1, 10, (*string)(nil)).
					Return([]models.HealthCatalog{}, int64(0), errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── GET /species/:species/category/:category ───────────
		{
			name:   "GetBySpeciesAndCategory success",
			method: "GET", path: "/species/:species/category/:category",
			urlParams: gin.Params{{Key: "species", Value: "dog"}, {Key: "category", Value: "vaccine"}},
			setup: func(m *mockRepo) {
				m.On("GetBySpeciesAndCategory", mock.Anything, "dog", "vaccine").
					Return([]models.HealthCatalog{{ID: "hc-1"}}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:   "GetBySpeciesAndCategory invalid species",
			method: "GET", path: "/species/:species/category/:category",
			urlParams: gin.Params{{Key: "species", Value: "dragon"}, {Key: "category", Value: "vaccine"}},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "GetBySpeciesAndCategory invalid category",
			method: "GET", path: "/species/:species/category/:category",
			urlParams: gin.Params{{Key: "species", Value: "dog"}, {Key: "category", Value: "invalid"}},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "GetBySpeciesAndCategory repo error",
			method: "GET", path: "/species/:species/category/:category",
			urlParams: gin.Params{{Key: "species", Value: "dog"}, {Key: "category", Value: "vaccine"}},
			setup: func(m *mockRepo) {
				m.On("GetBySpeciesAndCategory", mock.Anything, "dog", "vaccine").
					Return([]models.HealthCatalog{}, errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── GET /:id ───────────────────────────────────────────
		{
			name:   "GetByID success",
			method: "GET", path: "/:id",
			urlParams: gin.Params{{Key: "id", Value: "hc-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "hc-1").
					Return(models.HealthCatalog{ID: "hc-1", Name: "Rabies"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				var resp map[string]models.HealthCatalog
				json.Unmarshal(w.Body.Bytes(), &resp)
				require.Equal(t, "Rabies", resp["data"].Name)
			},
		},
		{
			name:   "GetByID not found",
			method: "GET", path: "/:id",
			urlParams: gin.Params{{Key: "id", Value: "hc-404"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "hc-404").
					Return(models.HealthCatalog{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:   "GetByID repo error",
			method: "GET", path: "/:id",
			urlParams: gin.Params{{Key: "id", Value: "hc-1"}},
			setup: func(m *mockRepo) {
				m.On("GetByID", mock.Anything, "hc-1").
					Return(models.HealthCatalog{}, errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── POST / ─────────────────────────────────────────────
		{
			name:   "Create success",
			method: "POST", path: "/",
			body:   `{"name":"Rabies","category":"vaccine","species":["dog","cat"],"is_mandatory":true}`,
			system: true,
			setup: func(m *mockRepo) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(p CreateHealthCatalogPayload) bool {
					return p.Name == "Rabies" && p.Category == "vaccine"
				})).Return(models.HealthCatalog{ID: "hc-new", Name: "Rabies"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			name:   "Create forbidden (non-system)",
			method: "POST", path: "/",
			body:   `{"name":"Rabies","category":"vaccine","species":["dog"]}`,
			system: false,
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:   "Create validation error",
			method: "POST", path: "/",
			body:   `{"name":"","category":"invalid","species":[]}`,
			system: true,
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:   "Create repo error",
			method: "POST", path: "/",
			body:   `{"name":"Rabies","category":"vaccine","species":["dog"]}`,
			system: true,
			setup: func(m *mockRepo) {
				m.On("Create", mock.Anything, mock.Anything).
					Return(models.HealthCatalog{}, errTest)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		// ── PUT /:id ───────────────────────────────────────────
		{
			name:   "Update success",
			method: "PUT", path: "/:id",
			body:   `{"name":"Updated","category":"vaccine","species":["dog"]}`,
			system: true,
			urlParams: gin.Params{{Key: "id", Value: "hc-1"}},
			setup: func(m *mockRepo) {
				m.On("Update", mock.Anything, "hc-1", mock.Anything).
					Return(models.HealthCatalog{ID: "hc-1", Name: "Updated"}, nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:   "Update forbidden (non-system)",
			method: "PUT", path: "/:id",
			body:   `{"name":"Updated","category":"vaccine","species":["dog"]}`,
			system: false,
			urlParams: gin.Params{{Key: "id", Value: "hc-1"}},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:   "Update not found",
			method: "PUT", path: "/:id",
			body:   `{"name":"Updated","category":"vaccine","species":["dog"]}`,
			system: true,
			urlParams: gin.Params{{Key: "id", Value: "hc-404"}},
			setup: func(m *mockRepo) {
				m.On("Update", mock.Anything, "hc-404", mock.Anything).
					Return(models.HealthCatalog{}, ErrNotFound)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		// ── DELETE /:id ────────────────────────────────────────
		{
			name:   "Delete success",
			method: "DELETE", path: "/:id",
			system: true,
			urlParams: gin.Params{{Key: "id", Value: "hc-1"}},
			setup: func(m *mockRepo) {
				m.On("Delete", mock.Anything, "hc-1").Return(nil)
			},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:   "Delete forbidden (non-system)",
			method: "DELETE", path: "/:id",
			system: false,
			urlParams: gin.Params{{Key: "id", Value: "hc-1"}},
			check: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:   "Delete not found",
			method: "DELETE", path: "/:id",
			system: true,
			urlParams: gin.Params{{Key: "id", Value: "hc-404"}},
			setup: func(m *mockRepo) {
				m.On("Delete", mock.Anything, "hc-404").Return(ErrNotFound)
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
