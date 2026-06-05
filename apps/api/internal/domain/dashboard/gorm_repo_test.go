package dashboard

import (
	"context"
	"testing"
	"time"

	"github.com/my-pets/api/internal/database/testutil"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func uid(s string) string {
	return "00000000-0000-0000-0000-" + s
}

const (
	userDashID = "00000000-0000-0000-0000-000000000001"
	petHealthy = "00000000-0000-0000-0000-000000000011"
	petPending = "00000000-0000-0000-0000-000000000012"
	petOverdue = "00000000-0000-0000-0000-000000000013"
	petComplet = "00000000-0000-0000-0000-000000000014"
)

func seedDashboardData(t *testing.T, db *gorm.DB) {
	t.Helper()

	now := time.Now().UTC()

	require.NoError(t, db.Create(&models.User{
		ID: userDashID, Name: "Dashboard User", Email: "dash@test.com",
	}).Error)

	for _, p := range []models.Pet{
		{ID: petHealthy, UserID: userDashID, Name: "Healthy Pet", Species: "dog"},
		{ID: petPending, UserID: userDashID, Name: "Pending Pet", Species: "cat"},
		{ID: petOverdue, UserID: userDashID, Name: "Overdue Pet", Species: "dog"},
		{ID: petComplet, UserID: userDashID, Name: "Completed Pet", Species: "cat"},
	} {
		require.NoError(t, db.Create(&p).Error)
	}

	for _, r := range []models.HealthRecord{
		{PetID: petHealthy, UserID: userDashID, Category: "vaccine", Name: "Rabies", ApplicationDate: now.AddDate(0, 0, -30), NextDoseDate: nil, LastDoseDate: &now, TotalDoses: intPtr(3), AppliedDosesCount: 1},
		{PetID: petPending, UserID: userDashID, Category: "vaccine", Name: "Distemper", ApplicationDate: now.AddDate(0, 0, -30), NextDoseDate: datePtr(now.AddDate(0, 0, 15)), LastDoseDate: &now, TotalDoses: intPtr(3), AppliedDosesCount: 1},
		{PetID: petOverdue, UserID: userDashID, Category: "deworming", Name: "Dewormer", ApplicationDate: now.AddDate(0, 0, -60), NextDoseDate: datePtr(now.AddDate(0, 0, -10)), LastDoseDate: &now, TotalDoses: nil, AppliedDosesCount: 2},
		{PetID: petComplet, UserID: userDashID, Category: "vaccine", Name: "Parvo", ApplicationDate: now.AddDate(0, 0, -90), NextDoseDate: datePtr(now.AddDate(0, 0, 30)), LastDoseDate: &now, TotalDoses: intPtr(3), AppliedDosesCount: 3},
	} {
		require.NoError(t, db.Create(&r).Error)
	}
}

func intPtr(n int) *int         { return &n }
func datePtr(t time.Time) *time.Time { return &t }

func TestGormRepo_GetSummary(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	seedDashboardData(t, db)

	repo := NewGormRepo(db)
	summary, err := repo.GetSummary(context.Background(), userDashID)
	require.NoError(t, err)

	require.Equal(t, int64(4), summary.TotalPets)
	require.Equal(t, int64(2), summary.HealthyPets)
	require.Equal(t, int64(2), summary.PendingTasks)
	require.Equal(t, int64(1), summary.OverdueTasks)
}

func TestGormRepo_GetSummary_NoData(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userEmpty := uid("000000000099")
	require.NoError(t, db.Create(&models.User{ID: userEmpty, Name: "Empty", Email: "empty@test.com"}).Error)

	repo := NewGormRepo(db)
	summary, err := repo.GetSummary(context.Background(), userEmpty)
	require.NoError(t, err)

	require.Equal(t, int64(0), summary.TotalPets)
	require.Equal(t, int64(0), summary.HealthyPets)
	require.Equal(t, int64(0), summary.PendingTasks)
	require.Equal(t, int64(0), summary.OverdueTasks)
}

func TestGormRepo_GetSummary_OnlyCompletedRecords(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	now := time.Now().UTC()
	userID := uid("000000000098")
	petID := uid("000000000081")

	require.NoError(t, db.Create(&models.User{ID: userID, Name: "Completed", Email: "completed@test.com"}).Error)
	require.NoError(t, db.Create(&models.Pet{ID: petID, UserID: userID, Name: "All Done", Species: "dog"}).Error)
	require.NoError(t, db.Create(&models.HealthRecord{
		PetID: petID, UserID: userID, Category: "vaccine", Name: "Completed",
		ApplicationDate: now.AddDate(0, 0, -30), NextDoseDate: datePtr(now.AddDate(0, 0, 15)),
		LastDoseDate: &now, TotalDoses: intPtr(2), AppliedDosesCount: 2,
	}).Error)

	repo := NewGormRepo(db)
	summary, err := repo.GetSummary(context.Background(), userID)
	require.NoError(t, err)

	require.Equal(t, int64(1), summary.TotalPets)
	require.Equal(t, int64(1), summary.HealthyPets)
	require.Equal(t, int64(0), summary.PendingTasks)
	require.Equal(t, int64(0), summary.OverdueTasks)
}

func TestGormRepo_GetSummary_OtherUserIsolation(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	now := time.Now().UTC()
	uA := uid("000000000031")
	uB := uid("000000000032")
	pA := uid("000000000041")
	pB := uid("000000000042")

	require.NoError(t, db.Create(&models.User{ID: uA, Name: "User A", Email: "a@test.com"}).Error)
	require.NoError(t, db.Create(&models.User{ID: uB, Name: "User B", Email: "b@test.com"}).Error)
	require.NoError(t, db.Create(&models.Pet{ID: pA, UserID: uA, Name: "A's Pet", Species: "dog"}).Error)
	require.NoError(t, db.Create(&models.Pet{ID: pB, UserID: uB, Name: "B's Pet", Species: "cat"}).Error)
	require.NoError(t, db.Create(&models.HealthRecord{
		PetID: pA, UserID: uA, Category: "vaccine", Name: "A's Record",
		ApplicationDate: now.AddDate(0, 0, -30), NextDoseDate: datePtr(now.AddDate(0, 0, 15)),
		LastDoseDate: &now, TotalDoses: intPtr(3), AppliedDosesCount: 1,
	}).Error)

	repo := NewGormRepo(db)

	summaryA, err := repo.GetSummary(context.Background(), uA)
	require.NoError(t, err)
	require.Equal(t, int64(1), summaryA.TotalPets)
	require.Equal(t, int64(1), summaryA.PendingTasks)

	summaryB, err := repo.GetSummary(context.Background(), uB)
	require.NoError(t, err)
	require.Equal(t, int64(1), summaryB.TotalPets)
	require.Equal(t, int64(0), summaryB.PendingTasks)
}
