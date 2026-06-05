//go:build integration

package pet

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/my-pets/api/internal/database/testutil"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func uid(n int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", n)
}

func seedUser(t *testing.T, db *gorm.DB, id, name, email string) {
	t.Helper()
	require.NoError(t, db.Create(&models.User{ID: id, Name: name, Email: email}).Error)
}

func TestGormRepo_Create(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID := uid(1)
	seedUser(t, db, userID, "Creator", "creator@test.com")

	repo := NewGormRepo(db)
	now := time.Now().UTC()
	size := "small"

	payload := CreatePetPayload{
		Name:      "TestPet",
		Species:   "dog",
		Breed:     "Labrador",
		BirthDate: now.AddDate(-2, 0, 0).Format("2006-01-02"),
		Size:      &size,
	}

	created, err := repo.Create(context.Background(), userID, payload)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	require.Equal(t, "TestPet", created.Name)
	require.Equal(t, "dog", created.Species)
	require.Equal(t, userID, created.UserID)
}

func TestGormRepo_GetByID(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID := uid(2)
	petID := uid(11)
	seedUser(t, db, userID, "Getter", "getter@test.com")
	require.NoError(t, db.Create(&models.Pet{ID: petID, UserID: userID, Name: "GetPet", Species: "dog"}).Error)

	repo := NewGormRepo(db)

	found, err := repo.GetByID(context.Background(), petID, userID)
	require.NoError(t, err)
	require.Equal(t, "GetPet", found.Name)

	_, err = repo.GetByID(context.Background(), petID, uid(99))
	require.ErrorIs(t, err, ErrNotFound)

	_, err = repo.GetByID(context.Background(), uid(999), userID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_GetPaginated(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID := uid(3)
	seedUser(t, db, userID, "Pag", "pag@test.com")

	for _, p := range []models.Pet{
		{ID: uid(21), UserID: userID, Name: "Alpha", Species: "dog"},
		{ID: uid(22), UserID: userID, Name: "Beta", Species: "cat"},
		{ID: uid(23), UserID: userID, Name: "Gamma", Species: "bird"},
	} {
		require.NoError(t, db.Create(&p).Error)
	}

	repo := NewGormRepo(db)

	t.Run("first page", func(t *testing.T) {
		result, total, err := repo.GetPaginated(context.Background(), userID, 1, 2)
		require.NoError(t, err)
		require.Equal(t, int64(3), total)
		require.Len(t, result, 2)
	})

	t.Run("second page", func(t *testing.T) {
		result, total, err := repo.GetPaginated(context.Background(), userID, 2, 2)
		require.NoError(t, err)
		require.Equal(t, int64(3), total)
		require.Len(t, result, 1)
	})

	t.Run("other user isolation", func(t *testing.T) {
		otherID := uid(4)
		seedUser(t, db, otherID, "Other", "other@test.com")
		result, total, err := repo.GetPaginated(context.Background(), otherID, 1, 10)
		require.NoError(t, err)
		require.Equal(t, int64(0), total)
		require.Len(t, result, 0)
	})
}

func TestGormRepo_Update(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID := uid(5)
	petID := uid(31)
	seedUser(t, db, userID, "Upd", "upd@test.com")
	require.NoError(t, db.Create(&models.Pet{ID: petID, UserID: userID, Name: "OldName", Species: "dog"}).Error)

	repo := NewGormRepo(db)
	size := "large"
	stage := "adult"

	payload := UpdatePetPayload{
		Name: "NewName", Breed: "Golden", BirthDate: "2020-01-01",
		Size: &size, LifeStage: &stage,
	}

	updated, err := repo.Update(context.Background(), petID, userID, payload)
	require.NoError(t, err)
	require.Equal(t, "NewName", updated.Name)
	require.Equal(t, "Golden", updated.Breed)

	_, err = repo.Update(context.Background(), petID, uid(99), payload)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_Delete(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID := uid(6)
	petID := uid(41)
	seedUser(t, db, userID, "Del", "del@test.com")
	require.NoError(t, db.Create(&models.Pet{ID: petID, UserID: userID, Name: "ToDelete", Species: "dog"}).Error)

	repo := NewGormRepo(db)

	err := repo.Delete(context.Background(), petID, userID)
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), petID, userID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_Delete_NotFound(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID := uid(7)
	seedUser(t, db, userID, "Del2", "del2@test.com")

	repo := NewGormRepo(db)
	err := repo.Delete(context.Background(), uid(999), userID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_Delete_OtherUser(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	ownerID := uid(8)
	intruderID := uid(9)
	petID := uid(51)
	seedUser(t, db, ownerID, "Owner", "owner@test.com")
	seedUser(t, db, intruderID, "Intruder", "intruder@test.com")
	require.NoError(t, db.Create(&models.Pet{ID: petID, UserID: ownerID, Name: "Mine", Species: "cat"}).Error)

	repo := NewGormRepo(db)

	err := repo.Delete(context.Background(), petID, intruderID)
	require.ErrorIs(t, err, ErrNotFound)

	_, err = repo.GetByID(context.Background(), petID, ownerID)
	require.NoError(t, err)
}

func TestGormRepo_CountByOwner(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID := uid(10)
	seedUser(t, db, userID, "Count", "count@test.com")

	for _, p := range []models.Pet{
		{ID: uid(61), UserID: userID, Name: "One", Species: "dog"},
		{ID: uid(62), UserID: userID, Name: "Two", Species: "cat"},
	} {
		require.NoError(t, db.Create(&p).Error)
	}

	repo := NewGormRepo(db)

	count, err := repo.CountByOwner(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, int64(2), count)

	count, err = repo.CountByOwner(context.Background(), uid(99))
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}
