package exam

import (
	"context"
	"fmt"
	"testing"

	"github.com/my-pets/api/internal/database/testutil"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func uid(n int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", n)
}

func setupExamDB(t *testing.T, db *gorm.DB) (userID, petID string) {
	t.Helper()
	userID = uid(1)
	petID = uid(10)

	require.NoError(t, db.Create(&models.User{ID: userID, Name: "Exam User", Email: "exam@test.com"}).Error)
	require.NoError(t, db.Create(&models.Pet{ID: petID, UserID: userID, Name: "Pet", Species: "dog"}).Error)
	return
}

func TestGormRepo_CreateExam(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupExamDB(t, db)
	repo := NewGormRepo(db)

	exam := models.Exam{
		ID: uid(100), PetID: petID, UserID: userID, Name: "Blood Test", Status: "scheduled",
	}

	created, err := repo.Create(context.Background(), exam)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	require.Equal(t, "Blood Test", created.Name)
}

func TestGormRepo_GetByID(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupExamDB(t, db)
	repo := NewGormRepo(db)

	exam := models.Exam{ID: uid(200), PetID: petID, UserID: userID, Name: "X-Ray"}
	require.NoError(t, db.Create(&exam).Error)

	found, err := repo.GetByID(context.Background(), uid(200), userID)
	require.NoError(t, err)
	require.Equal(t, "X-Ray", found.Name)

	_, err = repo.GetByID(context.Background(), uid(200), uid(999))
	require.ErrorIs(t, err, ErrNotFound)

	_, err = repo.GetByID(context.Background(), uid(999), userID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_GetAllByOwner(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupExamDB(t, db)
	repo := NewGormRepo(db)

	for i := 0; i < 3; i++ {
		require.NoError(t, db.Create(&models.Exam{
			ID: uid(300 + i), PetID: petID, UserID: userID, Name: fmt.Sprintf("Exam %d", i+1),
		}).Error)
	}

	result, total, err := repo.GetAllByOwner(context.Background(), userID, 1, 2)
	require.NoError(t, err)
	require.Equal(t, int64(3), total)
	require.Len(t, result, 2)

	result, total, err = repo.GetAllByOwner(context.Background(), uid(999), 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(0), total)
	require.Len(t, result, 0)
}

func TestGormRepo_Update(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupExamDB(t, db)
	repo := NewGormRepo(db)

	exam := models.Exam{ID: uid(400), PetID: petID, UserID: userID, Name: "Old", Status: "scheduled"}
	require.NoError(t, db.Create(&exam).Error)

	exam.Name = "Updated"
	_, err := repo.Update(context.Background(), exam)
	require.NoError(t, err)

	var updated models.Exam
	db.First(&updated, "id = ?", uid(400))
	require.Equal(t, "Updated", updated.Name)
}

func TestGormRepo_Delete(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupExamDB(t, db)
	repo := NewGormRepo(db)

	exam := models.Exam{ID: uid(500), PetID: petID, UserID: userID, Name: "ToDelete"}
	require.NoError(t, db.Create(&exam).Error)

	err := repo.Delete(context.Background(), uid(500), userID)
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), uid(500), userID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_Delete_NotFound(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, _ := setupExamDB(t, db)
	repo := NewGormRepo(db)

	err := repo.Delete(context.Background(), uid(999), userID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_CreateAndGetResults(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupExamDB(t, db)
	repo := NewGormRepo(db)

	exam := models.Exam{ID: uid(600), PetID: petID, UserID: userID, Name: "Blood Test"}
	require.NoError(t, db.Create(&exam).Error)

	results := []models.ExamResult{
		{ExamID: uid(600), ParameterName: "Glucose", Value: "95"},
		{ExamID: uid(600), ParameterName: "WBC", Value: "8.5"},
	}
	err := repo.CreateResults(context.Background(), results)
	require.NoError(t, err)

	fetched, err := repo.GetResultsByExamID(context.Background(), uid(600))
	require.NoError(t, err)
	require.Len(t, fetched, 2)

	exam, fetchedResults, err := repo.GetByIDWithResults(context.Background(), uid(600), userID)
	require.NoError(t, err)
	require.Equal(t, "Blood Test", exam.Name)
	require.Len(t, fetchedResults, 2)
}

func TestGormRepo_DeleteResultsByExamID(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	userID, petID := setupExamDB(t, db)
	repo := NewGormRepo(db)

	exam := models.Exam{ID: uid(700), PetID: petID, UserID: userID, Name: "Test"}
	require.NoError(t, db.Create(&exam).Error)
	require.NoError(t, db.Create(&models.ExamResult{
		ExamID: uid(700), ParameterName: "P1", Value: "v1",
	}).Error)

	err := repo.DeleteResultsByExamID(context.Background(), uid(700))
	require.NoError(t, err)

	results, err := repo.GetResultsByExamID(context.Background(), uid(700))
	require.NoError(t, err)
	require.Len(t, results, 0)
}
