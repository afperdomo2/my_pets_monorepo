package health_record

import (
	"context"
	"testing"
	"time"

	"github.com/my-pets/api/internal/database/testutil"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupHealthRecordDB(t *testing.T, db *gorm.DB) (userID, petID string) {
	t.Helper()
	userID = uid(1)
	petID = uid(10)
	require.NoError(t, db.Create(&models.User{ID: userID, Name: "HR User", Email: "hr@test.com"}).Error)
	require.NoError(t, db.Create(&models.Pet{ID: petID, UserID: userID, Name: "TestPet", Species: "dog"}).Error)
	return
}

func TestGormRepo_CreateHealthRecord(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupHealthRecordDB(t, db)
	repo := NewGormRepo(db)

	payload := CreateHealthRecordPayload{
		PetID:           petID,
		Category:        "vaccine",
		Name:            "Rabies",
		ApplicationDate: time.Now().UTC().AddDate(0, 0, -30).Format("2006-01-02"),
	}

	created, err := repo.Create(context.Background(), userID, payload)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	require.Equal(t, "vaccine", created.Category)
	require.Equal(t, "Rabies", created.Name)
	require.Equal(t, int(1), created.AppliedDosesCount)
}

func TestGormRepo_GetByID(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupHealthRecordDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()
	rec := models.HealthRecord{
		ID: uid(100), PetID: petID, UserID: userID, Category: "vaccine", Name: "Test",
		ApplicationDate: now, TotalDoses: intPtr(3), AppliedDosesCount: 1,
	}
	require.NoError(t, db.Create(&rec).Error)

	found, err := repo.GetByID(context.Background(), uid(100), userID)
	require.NoError(t, err)
	require.Equal(t, "Test", found.Name)

	_, err = repo.GetByID(context.Background(), uid(100), uid(999))
	require.ErrorIs(t, err, ErrNotFound)

	_, err = repo.GetByID(context.Background(), uid(999), userID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_GetUpcomingRecords(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupHealthRecordDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()

	pending := models.HealthRecord{
		ID: uid(101), PetID: petID, UserID: userID, Category: "vaccine", Name: "Pending",
		ApplicationDate: now, NextDoseDate: datePtr(now.AddDate(0, 0, 15)),
		TotalDoses: intPtr(3), AppliedDosesCount: 1,
	}
	completed := models.HealthRecord{
		ID: uid(102), PetID: petID, UserID: userID, Category: "vaccine", Name: "Done",
		ApplicationDate: now, NextDoseDate: datePtr(now.AddDate(0, 0, 15)),
		TotalDoses: intPtr(3), AppliedDosesCount: 3,
	}
	noDose := models.HealthRecord{
		ID: uid(103), PetID: petID, UserID: userID, Category: "deworming", Name: "NoNext",
		ApplicationDate: now, NextDoseDate: nil, TotalDoses: nil, AppliedDosesCount: 1,
	}
	for _, r := range []models.HealthRecord{pending, completed, noDose} {
		require.NoError(t, db.Create(&r).Error)
	}

	records, total, err := repo.GetUpcomingRecords(context.Background(), userID, "", 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, records, 1)
	require.Equal(t, "Pending", records[0].Name)
}

func TestGormRepo_GetUpcomingRecords_ByCategory(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupHealthRecordDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()

	vaccine := models.HealthRecord{
		ID: uid(201), PetID: petID, UserID: userID, Category: "vaccine", Name: "Vacc",
		ApplicationDate: now, NextDoseDate: datePtr(now.AddDate(0, 0, 15)),
		TotalDoses: nil, AppliedDosesCount: 1,
	}
	deworming := models.HealthRecord{
		ID: uid(202), PetID: petID, UserID: userID, Category: "deworming", Name: "Deworm",
		ApplicationDate: now, NextDoseDate: datePtr(now.AddDate(0, 0, 15)),
		TotalDoses: nil, AppliedDosesCount: 1,
	}
	require.NoError(t, db.Create(&vaccine).Error)
	require.NoError(t, db.Create(&deworming).Error)

	records, total, err := repo.GetUpcomingRecords(context.Background(), userID, "deworming", 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Equal(t, "Deworm", records[0].Name)
}

func TestGormRepo_UpdateHealthRecord(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupHealthRecordDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()
	rec := models.HealthRecord{
		ID: uid(300), PetID: petID, UserID: userID, Category: "vaccine", Name: "Old",
		ApplicationDate: now, TotalDoses: intPtr(3), AppliedDosesCount: 1,
	}
	require.NoError(t, db.Create(&rec).Error)

	payload := UpdateHealthRecordPayload{
		Category: "vaccine",
		Name:     "Updated",
	}

	updated, err := repo.Update(context.Background(), uid(300), userID, payload)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Name)
}

func TestGormRepo_DeleteHealthRecord(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupHealthRecordDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()
	rec := models.HealthRecord{
		ID: uid(400), PetID: petID, UserID: userID, Category: "vaccine", Name: "ToDelete",
		ApplicationDate: now, TotalDoses: intPtr(1), AppliedDosesCount: 1,
	}
	require.NoError(t, db.Create(&rec).Error)

	app := models.VaccineApplication{
		ID: uid(401), HealthRecordID: uid(400), ApplicationDate: now,
	}
	require.NoError(t, db.Create(&app).Error)

	err := repo.Delete(context.Background(), uid(400), userID)
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), uid(400), userID)
	require.ErrorIs(t, err, ErrNotFound)

	var count int64
	db.Model(&models.VaccineApplication{}).Where("health_record_id = ?", uid(400)).Count(&count)
	require.Equal(t, int64(0), count)
}

func TestGormRepo_IncrementAppliedDosesCount(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupHealthRecordDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()
	rec := models.HealthRecord{
		ID: uid(500), PetID: petID, UserID: userID, Category: "vaccine", Name: "CountTest",
		ApplicationDate: now, TotalDoses: intPtr(3), AppliedDosesCount: 1,
	}
	require.NoError(t, db.Create(&rec).Error)

	require.NoError(t, repo.IncrementAppliedDosesCount(context.Background(), uid(500)))

	var updated models.HealthRecord
	db.First(&updated, "id = ?", uid(500))
	require.Equal(t, 2, updated.AppliedDosesCount)
}

func TestGormRepo_UpdateLastDoseDate(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupHealthRecordDB(t, db)
	repo := NewGormRepo(db)

	now := time.Now().UTC()
	rec := models.HealthRecord{
		ID: uid(600), PetID: petID, UserID: userID, Category: "vaccine", Name: "DateTest",
		ApplicationDate: now, LastDoseDate: &now, TotalDoses: intPtr(3), AppliedDosesCount: 1,
	}
	require.NoError(t, db.Create(&rec).Error)

	newDate := now.AddDate(0, 0, 30).Format("2006-01-02")
	require.NoError(t, repo.UpdateLastDoseDate(context.Background(), uid(600), newDate))

	var updated models.HealthRecord
	db.First(&updated, "id = ?", uid(600))
	require.Equal(t, newDate, updated.LastDoseDate.Format("2006-01-02"))
}

func intPtr(n int) *int          { return &n }
func datePtr(t time.Time) *time.Time { return &t }
