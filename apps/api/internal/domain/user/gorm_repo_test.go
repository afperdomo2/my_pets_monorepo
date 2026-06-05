package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/my-pets/api/internal/database/testutil"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func uid(n int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", n)
}

func TestGormRepo_Create(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	payload := CreateUserPayload{
		Name:     "Test User",
		Email:    "test@test.com",
		Password: "password123",
	}

	created, err := repo.Create(context.Background(), payload, true)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	require.True(t, created.IsSystemUser)
	require.Equal(t, "test@test.com", created.Email)

	_, err = repo.Create(context.Background(), payload, false)
	require.ErrorIs(t, err, ErrEmailTaken)
}

func TestGormRepo_GetByID(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)
	user := models.User{ID: uid(1), Name: "GetTest", Email: "get@test.com"}
	require.NoError(t, db.Create(&user).Error)

	found, err := repo.GetByID(context.Background(), uid(1))
	require.NoError(t, err)
	require.Equal(t, "GetTest", found.Name)

	_, err = repo.GetByID(context.Background(), uid(999))
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_GetByEmail(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)
	user := models.User{ID: uid(2), Name: "EmailTest", Email: "email@test.com"}
	require.NoError(t, db.Create(&user).Error)

	found, err := repo.GetByEmail(context.Background(), "email@test.com")
	require.NoError(t, err)
	require.Equal(t, "EmailTest", found.Name)

	_, err = repo.GetByEmail(context.Background(), "nonexistent@test.com")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_GetPaginated(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)
	for i := 0; i < 3; i++ {
		require.NoError(t, db.Create(&models.User{
			ID: uid(10 + i), Name: fmt.Sprintf("User %d", i+1),
			Email: fmt.Sprintf("user%d@test.com", i+1),
		}).Error)
	}

	result, total, err := repo.GetPaginated(context.Background(), 1, 2)
	require.NoError(t, err)
	require.Equal(t, int64(3), total)
	require.Len(t, result, 2)
}

func TestGormRepo_Update(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)
	user := models.User{ID: uid(3), Name: "Old", Email: "old@test.com"}
	require.NoError(t, db.Create(&user).Error)

	payload := UpdateUserPayload{Name: "Updated", Email: "updated@test.com"}
	updated, err := repo.Update(context.Background(), uid(3), payload)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Name)
	require.Equal(t, "updated@test.com", updated.Email)
}

func TestGormRepo_Update_EmailTaken(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)
	require.NoError(t, db.Create(&models.User{ID: uid(4), Name: "A", Email: "a@test.com"}).Error)
	require.NoError(t, db.Create(&models.User{ID: uid(5), Name: "B", Email: "b@test.com"}).Error)

	payload := UpdateUserPayload{Name: "B", Email: "a@test.com"}
	_, err := repo.Update(context.Background(), uid(5), payload)
	require.ErrorIs(t, err, ErrEmailTaken)
}

func TestGormRepo_Delete(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)
	require.NoError(t, db.Create(&models.User{ID: uid(6), Name: "Del", Email: "del@test.com"}).Error)

	err := repo.Delete(context.Background(), uid(6))
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), uid(6))
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_HasUsers(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	exists, err := repo.HasUsers(context.Background())
	require.NoError(t, err)
	require.False(t, exists)

	require.NoError(t, db.Create(&models.User{ID: uid(7), Name: "First", Email: "first@test.com"}).Error)

	exists, err = repo.HasUsers(context.Background())
	require.NoError(t, err)
	require.True(t, exists)
}

func TestGormRepo_SystemUserExists(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	exists, err := repo.SystemUserExists(context.Background())
	require.NoError(t, err)
	require.False(t, exists)

	require.NoError(t, db.Create(&models.User{ID: uid(8), Name: "Sys", Email: "sys@test.com", IsSystemUser: true}).Error)

	exists, err = repo.SystemUserExists(context.Background())
	require.NoError(t, err)
	require.True(t, exists)
}

func TestGormRepo_GetPasswordByEmail(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	require.NoError(t, err)
	hashStr := string(hash)

	require.NoError(t, db.Create(&models.User{
		ID: uid(9), Name: "Pwd", Email: "pwd@test.com",
		Password: &hashStr, AuthProvider: "local",
	}).Error)

	fetchedHash, err := repo.GetPasswordByEmail(context.Background(), "pwd@test.com")
	require.NoError(t, err)
	require.Equal(t, hashStr, fetchedHash)

	_, err = repo.GetPasswordByEmail(context.Background(), "nonexistent@test.com")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGormRepo_UpdatePassword(t *testing.T) {
	db, cleanup := testutil.NewPostgres(t)
	defer cleanup()

	repo := NewGormRepo(db)

	oldHash, _ := bcrypt.GenerateFromPassword([]byte("old"), bcrypt.DefaultCost)
	oldHashStr := string(oldHash)
	require.NoError(t, db.Create(&models.User{
		ID: uid(10), Name: "PwdUpd", Email: "pwdupd@test.com",
		Password: &oldHashStr, AuthProvider: "local",
	}).Error)

	newHash, _ := bcrypt.GenerateFromPassword([]byte("new"), bcrypt.DefaultCost)
	err := repo.UpdatePassword(context.Background(), uid(10), string(newHash))
	require.NoError(t, err)

	var u models.User
	db.First(&u, "id = ?", uid(10))
	require.NotEqual(t, oldHashStr, *u.Password)
}
