package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/my-pets/api/internal/models"
)

// gormRepo implements Repository using GORM.
type gormRepo struct {
	db *gorm.DB
}

// NewGormRepo creates a new GORM-backed user repository.
func NewGormRepo(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

func (r *gormRepo) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if result := r.db.WithContext(ctx).Order("id").Find(&users); result.Error != nil {
		return nil, fmt.Errorf("user.GetAll: %w", result.Error)
	}
	return users, nil
}

func (r *gormRepo) GetPaginated(ctx context.Context, page, perPage int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("user.GetPaginated count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := r.db.WithContext(ctx).Order("id").Limit(perPage).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("user.GetPaginated: %w", err)
	}

	return users, total, nil
}

func (r *gormRepo) GetByID(ctx context.Context, id string) (models.User, error) {
	var u models.User
	result := r.db.WithContext(ctx).First(&u, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.User{}, ErrNotFound
	}
	if result.Error != nil {
		return models.User{}, fmt.Errorf("user.GetByID: %w", result.Error)
	}
	return u, nil
}

func (r *gormRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.User{}, ErrNotFound
	}
	if result.Error != nil {
		return models.User{}, fmt.Errorf("user.GetByEmail: %w", result.Error)
	}
	return u, nil
}

func (r *gormRepo) Create(ctx context.Context, payload CreateUserPayload, isSystemUser bool) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("user.Create: hash password: %w", err)
	}

	hashStr := string(hash)
	provider := models.AuthProviderLocal
	u := models.User{
		Name:         payload.Name,
		Email:        payload.Email,
		Password:     &hashStr,
		IsSystemUser: isSystemUser,
		AuthProvider: provider,
	}

	if result := r.db.WithContext(ctx).Create(&u); result.Error != nil {
		if isUniqueViolation(result.Error) {
			return models.User{}, ErrEmailTaken
		}
		return models.User{}, fmt.Errorf("user.Create: %w", result.Error)
	}
	return u, nil
}

func (r *gormRepo) Update(ctx context.Context, id string, payload UpdateUserPayload) (models.User, error) {
	var u models.User
	result := r.db.WithContext(ctx).First(&u, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.User{}, ErrNotFound
	}
	if result.Error != nil {
		return models.User{}, fmt.Errorf("user.Update: %w", result.Error)
	}

	u.Name = payload.Name
	u.Email = payload.Email

	if result := r.db.WithContext(ctx).Save(&u); result.Error != nil {
		if isUniqueViolation(result.Error) {
			return models.User{}, ErrEmailTaken
		}
		return models.User{}, fmt.Errorf("user.Update: %w", result.Error)
	}
	return u, nil
}

func (r *gormRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("user.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormRepo) HasUsers(ctx context.Context) (bool, error) {
	var count int64
	if result := r.db.WithContext(ctx).Model(&models.User{}).Count(&count); result.Error != nil {
		return false, fmt.Errorf("user.HasUsers: %w", result.Error)
	}
	return count > 0, nil
}

func (r *gormRepo) SystemUserExists(ctx context.Context) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("is_system_user = ?", true).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("user.SystemUserExists: %w", result.Error)
	}
	return count > 0, nil
}

func (r *gormRepo) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	var u models.User
	result := r.db.WithContext(ctx).Select("password", "auth_provider").Where("email = ?", email).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", ErrNotFound
	}
	if result.Error != nil {
		return "", fmt.Errorf("user.GetPasswordByEmail: %w", result.Error)
	}
	if u.Password == nil {
		return "", ErrWrongProvider
	}
	return *u.Password, nil
}

func (r *gormRepo) UpsertByGoogleID(ctx context.Context, info GoogleUserInfo, isSystemUser bool) (models.User, error) {
	var u models.User
	provider := models.AuthProviderGoogle

	// Try to find existing user by google_id.
	result := r.db.WithContext(ctx).Where("google_id = ?", info.GoogleID).First(&u)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.User{}, fmt.Errorf("user.UpsertByGoogleID: %w", result.Error)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// New user — create.
		u = models.User{
			Name:         info.Name,
			Email:        info.Email,
			GoogleID:     &info.GoogleID,
			AuthProvider: provider,
			IsSystemUser: isSystemUser,
		}
		if createResult := r.db.WithContext(ctx).Create(&u); createResult.Error != nil {
			return models.User{}, fmt.Errorf("user.UpsertByGoogleID: create: %w", createResult.Error)
		}
	} else {
		// Existing user — update name (email may have changed on Google side).
		u.Name = info.Name
		if saveResult := r.db.WithContext(ctx).Save(&u); saveResult.Error != nil {
			return models.User{}, fmt.Errorf("user.UpsertByGoogleID: update: %w", saveResult.Error)
		}
	}

	return u, nil
}

// isUniqueViolation detects PostgreSQL unique constraint violations (SQLSTATE 23505).
// pgx errors implement the SQLState() method.
type sqlstateErr interface {
	SQLState() string
}

func isUniqueViolation(err error) bool {
	var pgErr sqlstateErr
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
