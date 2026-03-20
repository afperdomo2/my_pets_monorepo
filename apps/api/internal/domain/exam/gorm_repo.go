package exam

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/my-pets/api/internal/models"
)

// gormRepo implementa Repository usando GORM.
type gormRepo struct {
	db *gorm.DB
}

// NewGormRepo crea un nuevo repositorio de exam respaldado por GORM.
func NewGormRepo(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

// parseDateString convierte una cadena "YYYY-MM-DD" en time.Time.
func parseDateString(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// petBelongsToOwner verifica que la mascota existe y pertenece al usuario autenticado.
func (r *gormRepo) petBelongsToOwner(ctx context.Context, petID, ownerID string) error {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Pet{}).
		Where("id = ? AND user_id = ?", petID, ownerID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("exam: verificar mascota: %w", err)
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormRepo) GetAllByOwner(ctx context.Context, ownerID string, page, perPage int) ([]models.Exam, int64, error) {
	var exams []models.Exam
	var total int64

	base := r.db.WithContext(ctx).Model(&models.Exam{}).Where("user_id = ?", ownerID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("exam.GetAllByOwner count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := base.Preload("Pet").Order("created_at DESC").Limit(perPage).Offset(offset).Find(&exams).Error; err != nil {
		return nil, 0, fmt.Errorf("exam.GetAllByOwner: %w", err)
	}

	return exams, total, nil
}

func (r *gormRepo) GetByPetID(ctx context.Context, petID, ownerID string, page, perPage int) ([]models.Exam, int64, error) {
	var exams []models.Exam
	var total int64

	base := r.db.WithContext(ctx).Model(&models.Exam{}).
		Where("pet_id = ? AND user_id = ?", petID, ownerID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("exam.GetByPetID count: %w", err)
	}

	offset := (page - 1) * perPage
	if err := base.Preload("Pet").Order("created_at DESC").Limit(perPage).Offset(offset).Find(&exams).Error; err != nil {
		return nil, 0, fmt.Errorf("exam.GetByPetID: %w", err)
	}

	return exams, total, nil
}

func (r *gormRepo) GetByID(ctx context.Context, id, ownerID string) (models.Exam, error) {
	var exam models.Exam
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, ownerID).
		First(&exam)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Exam{}, ErrNotFound
	}
	if result.Error != nil {
		return models.Exam{}, fmt.Errorf("exam.GetByID: %w", result.Error)
	}
	return exam, nil
}

func (r *gormRepo) GetByIDWithResults(ctx context.Context, id, ownerID string) (models.Exam, []models.ExamResult, error) {
	exam, err := r.GetByID(ctx, id, ownerID)
	if err != nil {
		return models.Exam{}, nil, err
	}

	results, err := r.GetResultsByExamID(ctx, id)
	if err != nil {
		return models.Exam{}, nil, err
	}

	return exam, results, nil
}

func (r *gormRepo) Create(ctx context.Context, exam models.Exam) (models.Exam, error) {
	// Verificar que la mascota pertenece al usuario antes de crear.
	if err := r.petBelongsToOwner(ctx, exam.PetID, exam.UserID); err != nil {
		return models.Exam{}, err
	}

	result := r.db.WithContext(ctx).Create(&exam)
	if result.Error != nil {
		return models.Exam{}, fmt.Errorf("exam.Create: %w", result.Error)
	}
	return exam, nil
}

func (r *gormRepo) Update(ctx context.Context, exam models.Exam) (models.Exam, error) {
	result := r.db.WithContext(ctx).Save(&exam)
	if result.Error != nil {
		return models.Exam{}, fmt.Errorf("exam.Update: %w", result.Error)
	}
	return exam, nil
}

func (r *gormRepo) Delete(ctx context.Context, id, ownerID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, ownerID).
		Delete(&models.Exam{})
	if result.Error != nil {
		return fmt.Errorf("exam.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormRepo) CreateResults(ctx context.Context, results []models.ExamResult) error {
	if len(results) == 0 {
		return nil
	}
	result := r.db.WithContext(ctx).Create(&results)
	if result.Error != nil {
		return fmt.Errorf("exam.CreateResults: %w", result.Error)
	}
	return nil
}

func (r *gormRepo) DeleteResultsByExamID(ctx context.Context, examID string) error {
	result := r.db.WithContext(ctx).
		Where("exam_id = ?", examID).
		Delete(&models.ExamResult{})
	if result.Error != nil {
		return fmt.Errorf("exam.DeleteResultsByExamID: %w", result.Error)
	}
	return nil
}

func (r *gormRepo) GetResultsByExamID(ctx context.Context, examID string) ([]models.ExamResult, error) {
	var results []models.ExamResult
	result := r.db.WithContext(ctx).
		Where("exam_id = ?", examID).
		Find(&results)
	if result.Error != nil {
		return nil, fmt.Errorf("exam.GetResultsByExamID: %w", result.Error)
	}
	return results, nil
}
