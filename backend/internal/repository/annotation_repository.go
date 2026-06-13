package repository

import (
	"ancient-text-backend/internal/database"
	"ancient-text-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnnotationRepository struct {
	db *gorm.DB
}

func NewAnnotationRepository() *AnnotationRepository {
	return &AnnotationRepository{
		db: database.DB,
	}
}

func (r *AnnotationRepository) Create(ann *models.Annotation) error {
	return r.db.Create(ann).Error
}

func (r *AnnotationRepository) FindByID(id uuid.UUID) (*models.Annotation, error) {
	var ann models.Annotation
	err := r.db.First(&ann, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ann, nil
}

func (r *AnnotationRepository) FindByTextID(textID uuid.UUID) ([]models.Annotation, error) {
	var annotations []models.Annotation
	err := r.db.Where("text_id = ?", textID).
		Order("start_pos ASC, created_at DESC").
		Find(&annotations).Error
	if err != nil {
		return nil, err
	}
	return annotations, nil
}

func (r *AnnotationRepository) FindByTextIDAndUser(textID uuid.UUID, userID string) ([]models.Annotation, error) {
	var annotations []models.Annotation
	err := r.db.Where("text_id = ? AND user_id = ?", textID, userID).
		Order("start_pos ASC, created_at DESC").
		Find(&annotations).Error
	if err != nil {
		return nil, err
	}
	return annotations, nil
}

func (r *AnnotationRepository) Update(ann *models.Annotation) error {
	return r.db.Save(ann).Error
}

func (r *AnnotationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Annotation{}, "id = ?", id).Error
}

func (r *AnnotationRepository) DeleteByTextID(textID uuid.UUID) error {
	return r.db.Where("text_id = ?", textID).Delete(&models.Annotation{}).Error
}

func (r *AnnotationRepository) BatchCreate(annotations []models.Annotation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, ann := range annotations {
			if err := tx.Create(&ann).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *AnnotationRepository) BatchUpdate(annotations []models.Annotation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, ann := range annotations {
			if err := tx.Save(&ann).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
