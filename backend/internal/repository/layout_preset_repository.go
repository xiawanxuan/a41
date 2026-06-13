package repository

import (
	"ancient-text-backend/internal/database"
	"ancient-text-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LayoutPresetRepository struct {
	db *gorm.DB
}

func NewLayoutPresetRepository() *LayoutPresetRepository {
	return &LayoutPresetRepository{
		db: database.DB,
	}
}

func (r *LayoutPresetRepository) Create(preset *models.LayoutPreset) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if preset.IsDefault {
			if err := tx.Model(&models.LayoutPreset{}).
				Where("user_id = ? AND is_default = ?", preset.UserID, true).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(preset).Error
	})
}

func (r *LayoutPresetRepository) FindByID(id uuid.UUID) (*models.LayoutPreset, error) {
	var preset models.LayoutPreset
	err := r.db.First(&preset, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &preset, nil
}

func (r *LayoutPresetRepository) FindByUserID(userID string) ([]models.LayoutPreset, error) {
	var presets []models.LayoutPreset
	err := r.db.Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&presets).Error
	if err != nil {
		return nil, err
	}
	return presets, nil
}

func (r *LayoutPresetRepository) FindDefaultByUserID(userID string) (*models.LayoutPreset, error) {
	var preset models.LayoutPreset
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).
		First(&preset).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &preset, nil
}

func (r *LayoutPresetRepository) Update(preset *models.LayoutPreset) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if preset.IsDefault {
			if err := tx.Model(&models.LayoutPreset{}).
				Where("user_id = ? AND is_default = ? AND id != ?", preset.UserID, true, preset.ID).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Save(preset).Error
	})
}

func (r *LayoutPresetRepository) SetDefault(id uuid.UUID, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.LayoutPreset{}).
			Where("user_id = ? AND is_default = ?", userID, true).
			Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&models.LayoutPreset{}).
			Where("id = ? AND user_id = ?", id, userID).
			Update("is_default", true).Error
	})
}

func (r *LayoutPresetRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.LayoutPreset{}, "id = ?", id).Error
}
