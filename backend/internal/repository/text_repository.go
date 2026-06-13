package repository

import (
	"ancient-text-backend/internal/database"
	"ancient-text-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AncientTextRepository struct {
	db *gorm.DB
}

func NewAncientTextRepository() *AncientTextRepository {
	return &AncientTextRepository{
		db: database.DB,
	}
}

func (r *AncientTextRepository) Create(text *models.AncientText) error {
	return r.db.Create(text).Error
}

func (r *AncientTextRepository) FindByID(id uuid.UUID) (*models.AncientText, error) {
	var text models.AncientText
	err := r.db.Preload("LayoutConfig").First(&text, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &text, nil
}

func (r *AncientTextRepository) FindAll(page, pageSize int) ([]models.AncientText, int64, error) {
	var texts []models.AncientText
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&models.AncientText{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("LayoutConfig").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&texts).Error
	if err != nil {
		return nil, 0, err
	}

	return texts, total, nil
}

func (r *AncientTextRepository) Search(keyword string, page, pageSize int) ([]models.AncientText, int64, error) {
	var texts []models.AncientText
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&models.AncientText{})

	if keyword != "" {
		query = query.Where("title ILIKE ? OR author ILIKE ? OR content ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("LayoutConfig").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&texts).Error
	if err != nil {
		return nil, 0, err
	}

	return texts, total, nil
}

func (r *AncientTextRepository) Update(text *models.AncientText) error {
	return r.db.Save(text).Error
}

func (r *AncientTextRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.AncientText{}, "id = ?", id).Error
}
