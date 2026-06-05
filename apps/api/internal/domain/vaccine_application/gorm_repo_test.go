//go:build integration

package vaccine_application

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

func setupVAppDB(t *testing.T, db *gorm.DB) (userID, petID, hrID string) {
	t.Helper()
	userID = uid(1)
	petID = uid(10)
	hrID = uid(100)

	require.NoError(t, db.Create(&models.User{ID: userID, Name: "VA User", Email: "va@test.com"}).Error)
	require.NoError(t, db.Create(&models.Pet{ID: petID, UserID: userID, Name: "Pet", Species: "cat"}).Error)
	require.NoError(t, db.Create(&models.HealthRecord{
		ID: hrID, PetID: petID, UserID: userID, Category: "vaccine", Name: "Record",
		ApplicationDate: time.Now().UTC(), TotalDoses: intPtr(3), AppliedDosesCount: 1,
	}).Error)
	return
}

func intPtr(n int) *int { return &n }

func TestGormRepo_Create(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	_, _, hrID := setupVAppDB(t, db)
	repo := NewGormRepo(db)

	app := models.VaccineApplication{
		ID: uid(200), HealthRecordID: hrID,
		ApplicationDate: time.Now().UTC(),
	}

	created, err := repo.Create(context.Background(), app)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
}

func TestGormRepo_GetByID(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	_, _, hrID := setupVAppDB(t, db)
	repo := NewGormRepo(db)

	app := models.VaccineApplication{
		ID: uid(200), HealthRecordID: hrID,
		ApplicationDate: time.Now().UTC(),
	}
	require.NoError(t, db.Create(&app).Error)

	found, err := repo.GetByID(context.Background(), uid(200))
	require.NoError(t, err)
	require.Equal(t, hrID, found.HealthRecordID)

	_, err = repo.GetByID(context.Background(), uid(999))
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_GetByHealthRecordID(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	_, _, hrID := setupVAppDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()
	apps := []models.VaccineApplication{
		{ID: uid(201), HealthRecordID: hrID, ApplicationDate: now.AddDate(0, 0, -10)},
		{ID: uid(202), HealthRecordID: hrID, ApplicationDate: now},
	}
	for _, a := range apps {
		require.NoError(t, db.Create(&a).Error)
	}

	found, err := repo.GetByHealthRecordID(context.Background(), hrID)
	require.NoError(t, err)
	require.Len(t, found, 2)
	require.Equal(t, uid(201), found[0].ID)
}

func TestGormRepo_Update(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	_, _, hrID := setupVAppDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()
	notes := "initial"
	app := models.VaccineApplication{
		ID: uid(300), HealthRecordID: hrID,
		ApplicationDate: now, Notes: &notes,
	}
	require.NoError(t, db.Create(&app).Error)

	updatedNotes := "updated"
	app.Notes = &updatedNotes
	saved, err := repo.Update(context.Background(), app)
	require.NoError(t, err)
	require.Equal(t, "updated", *saved.Notes)
}

func TestGormRepo_Delete(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	_, _, hrID := setupVAppDB(t, db)
	repo := NewGormRepo(db)

	app := models.VaccineApplication{
		ID: uid(400), HealthRecordID: hrID,
		ApplicationDate: time.Now().UTC(),
	}
	require.NoError(t, db.Create(&app).Error)

	err := repo.Delete(context.Background(), uid(400))
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), uid(400))
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_Delete_NotFound(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)
	err := repo.Delete(context.Background(), uid(999))
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_UpdateHealthRecordAfterApplication(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	_, _, hrID := setupVAppDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()
	appDate := now.AddDate(0, 0, -5).Format("2006-01-02")
	nextDate := now.AddDate(0, 0, 30).Format("2006-01-02")

	err := repo.UpdateHealthRecordAfterApplication(context.Background(), hrID, appDate, &nextDate)
	require.NoError(t, err)

	var hr models.HealthRecord
	db.First(&hr, "id = ?", hrID)
	require.Equal(t, appDate, hr.LastDoseDate.Format("2006-01-02"))
	require.Equal(t, 2, hr.AppliedDosesCount)
	require.NotNil(t, hr.NextDoseDate)

	nextDateFormatted := hr.NextDoseDate.Format("2006-01-02")
	require.Equal(t, nextDate, nextDateFormatted)
}

func TestGormRepo_UpdateHealthRecordAfterApplication_NullNextDate(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	_, _, hrID := setupVAppDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()
	appDate := now.AddDate(0, 0, -5).Format("2006-01-02")

	err := repo.UpdateHealthRecordAfterApplication(context.Background(), hrID, appDate, nil)
	require.NoError(t, err)

	var hr models.HealthRecord
	db.First(&hr, "id = ?", hrID)
	require.Equal(t, appDate, hr.LastDoseDate.Format("2006-01-02"))
	require.Equal(t, 2, hr.AppliedDosesCount)
	require.Nil(t, hr.NextDoseDate)
}
