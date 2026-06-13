package repository

import (
	"ancient-text-backend/internal/database"
	"ancient-text-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LayoutConfigRepository struct {
	db *gorm.DB
}

func NewLayoutConfigRepository() *LayoutConfigRepository {
	return &LayoutConfigRepository{
		db: database.DB,
	}
}

func (r *LayoutConfigRepository) Create(layout *models.LayoutConfig) error {
	return r.db.Create(layout).Error
}

func (r *LayoutConfigRepository) FindByID(id uuid.UUID) (*models.LayoutConfig, error) {
	var layout models.LayoutConfig
	err := r.db.First(&layout, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &layout, nil
}

func (r *LayoutConfigRepository) FindByTextID(textID uuid.UUID) (*models.LayoutConfig, error) {
	var layout models.LayoutConfig
	err := r.db.First(&layout, "text_id = ?", textID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &layout, nil
}

func (r *LayoutConfigRepository) Update(layout *models.LayoutConfig) error {
	return r.db.Save(layout).Error
}

func (r *LayoutConfigRepository) Upsert(layout *models.LayoutConfig) error {
	existing, err := r.FindByTextID(layout.TextID)
	if err != nil {
		return err
	}

	if existing != nil {
		layout.ID = existing.ID
		return r.db.Save(layout).Error
	}
	return r.db.Create(layout).Error
}

func (r *LayoutConfigRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.LayoutConfig{}, "id = ?", id).Error
}

func (r *LayoutConfigRepository) DeleteByTextID(textID uuid.UUID) error {
	return r.db.Where("text_id = ?", textID).Delete(&models.LayoutConfig{}).Error
}
