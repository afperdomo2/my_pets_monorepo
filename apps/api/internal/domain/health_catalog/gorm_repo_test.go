package health_catalog

import (
	"context"
	"fmt"
	"testing"

	"github.com/my-pets/api/internal/database/testutil"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/require"
)

func uid(n int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", n)
}

func TestGormRepo_Create(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	payload := CreateHealthCatalogPayload{
		Name:        "Rabies",
		Category:    "vaccine",
		Description: "Antirrábica",
		Species:     []string{"dog", "cat"},
		IsMandatory: true,
	}

	created, err := repo.Create(context.Background(), payload)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	require.Equal(t, "Rabies", created.Name)
	require.True(t, created.IsMandatory)
}

func TestGormRepo_GetByID(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	item := models.HealthCatalog{ID: uid(1), Name: "Test", Category: "vaccine", Species: []string{"dog"}}
	require.NoError(t, db.Create(&item).Error)

	found, err := repo.GetByID(context.Background(), uid(1))
	require.NoError(t, err)
	require.Equal(t, "Test", found.Name)

	_, err = repo.GetByID(context.Background(), uid(999))
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_GetPaginatedByCategory(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	for i := 0; i < 3; i++ {
		require.NoError(t, db.Create(&models.HealthCatalog{
			ID: uid(10 + i), Name: fmt.Sprintf("Item %d", i+1),
			Category: "vaccine", Species: []string{"dog"},
		}).Error)
	}
	require.NoError(t, db.Create(&models.HealthCatalog{
		ID: uid(99), Name: "Dewormer", Category: "deworming", Species: []string{"cat"},
	}).Error)

	species := "dog"
	items, total, err := repo.GetPaginatedByCategory(context.Background(), "vaccine", 1, 2, nil)
	require.NoError(t, err)
	require.Equal(t, int64(3), total)
	require.Len(t, items, 2)

	items, total, err = repo.GetPaginatedByCategory(context.Background(), "vaccine", 1, 10, &species)
	require.NoError(t, err)
	require.Equal(t, int64(3), total)

	items, total, err = repo.GetPaginatedByCategory(context.Background(), "deworming", 1, 10, nil)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
}

func TestGormRepo_GetBySpeciesAndCategory(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	require.NoError(t, db.Create(&models.HealthCatalog{
		ID: uid(20), Name: "Rabies", Category: "vaccine", Species: []string{"dog", "cat"}, IsMandatory: true,
	}).Error)
	require.NoError(t, db.Create(&models.HealthCatalog{
		ID: uid(21), Name: "Parvo", Category: "vaccine", Species: []string{"dog"}, IsMandatory: false,
	}).Error)

	items, err := repo.GetBySpeciesAndCategory(context.Background(), "dog", "vaccine")
	require.NoError(t, err)
	require.Len(t, items, 2)
	require.Equal(t, "Rabies", items[0].Name)

	items, err = repo.GetBySpeciesAndCategory(context.Background(), "cat", "vaccine")
	require.NoError(t, err)
	require.Len(t, items, 1)

	items, err = repo.GetBySpeciesAndCategory(context.Background(), "bird", "vaccine")
	require.NoError(t, err)
	require.Len(t, items, 0)
}

func TestGormRepo_Update(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	item := models.HealthCatalog{ID: uid(30), Name: "Old", Category: "vaccine", Species: []string{"dog"}}
	require.NoError(t, db.Create(&item).Error)

	payload := UpdateHealthCatalogPayload{
		Name: "Updated", Category: "vaccine", Species: []string{"dog", "cat"},
	}
	updated, err := repo.Update(context.Background(), uid(30), payload)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Name)
}

func TestGormRepo_Delete(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	item := models.HealthCatalog{ID: uid(40), Name: "Del", Category: "vaccine", Species: []string{"dog"}}
	require.NoError(t, db.Create(&item).Error)

	err := repo.Delete(context.Background(), uid(40))
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), uid(40))
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_Delete_NotFound(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)
	err := repo.Delete(context.Background(), uid(999))
	require.ErrorIs(t, err, ErrNotFound)
}
