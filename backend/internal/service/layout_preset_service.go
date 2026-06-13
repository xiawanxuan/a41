package service

import (
	"errors"

	"ancient-text-backend/internal/models"
	"ancient-text-backend/internal/repository"

	"github.com/google/uuid"
)

type LayoutPresetService struct {
	repo *repository.LayoutPresetRepository
}

func NewLayoutPresetService() *LayoutPresetService {
	return &LayoutPresetService{
		repo: repository.NewLayoutPresetRepository(),
	}
}

func (s *LayoutPresetService) Create(userID string, req *CreateLayoutPresetRequest) (*models.LayoutPreset, error) {
	if req.Name == "" {
		return nil, errors.New("preset name is required")
	}

	preset := &models.LayoutPreset{
		UserID:         userID,
		Name:           req.Name,
		Description:    req.Description,
		CharsPerColumn: req.CharsPerColumn,
		ColumnGap:      req.ColumnGap,
		FontSize:       req.FontSize,
		FontFamily:     req.FontFamily,
		LineHeight:     req.LineHeight,
		CharSpacing:    req.CharSpacing,
		TextColor:      req.TextColor,
		BackgroundColor: req.BackgroundColor,
		ShowBorder:     req.ShowBorder,
		IsDefault:      req.IsDefault,
		CustomCSS:      req.CustomCSS,
	}

	if preset.CharsPerColumn == 0 {
		preset.CharsPerColumn = 20
	}
	if preset.ColumnGap == 0 {
		preset.ColumnGap = 20
	}
	if preset.FontSize == 0 {
		preset.FontSize = 18
	}
	if preset.FontFamily == "" {
		preset.FontFamily = "SimSun, serif"
	}
	if preset.LineHeight == 0 {
		preset.LineHeight = 1.8
	}
	if preset.TextColor == "" {
		preset.TextColor = "#1a1a1a"
	}
	if preset.BackgroundColor == "" {
		preset.BackgroundColor = "#f5f0e8"
	}

	err := s.repo.Create(preset)
	if err != nil {
		return nil, err
	}

	return preset, nil
}

func (s *LayoutPresetService) GetByID(id string) (*models.LayoutPreset, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	return s.repo.FindByID(uid)
}

func (s *LayoutPresetService) ListByUser(userID string) ([]models.LayoutPreset, error) {
	return s.repo.FindByUserID(userID)
}

func (s *LayoutPresetService) GetDefault(userID string) (*models.LayoutPreset, error) {
	preset, err := s.repo.FindDefaultByUserID(userID)
	if err != nil {
		return nil, err
	}
	return preset, nil
}

func (s *LayoutPresetService) Update(id, userID string, req *UpdateLayoutPresetRequest) (*models.LayoutPreset, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	preset, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, errors.New("preset not found")
	}

	if preset.UserID != userID {
		return nil, errors.New("permission denied")
	}

	if req.Name != "" {
		preset.Name = req.Name
	}
	if req.Description != "" {
		preset.Description = req.Description
	}
	if req.CharsPerColumn > 0 {
		preset.CharsPerColumn = req.CharsPerColumn
	}
	if req.ColumnGap > 0 {
		preset.ColumnGap = req.ColumnGap
	}
	if req.FontSize > 0 {
		preset.FontSize = req.FontSize
	}
	if req.FontFamily != "" {
		preset.FontFamily = req.FontFamily
	}
	if req.LineHeight > 0 {
		preset.LineHeight = req.LineHeight
	}
	if req.CharSpacing >= 0 {
		preset.CharSpacing = req.CharSpacing
	}
	if req.TextColor != "" {
		preset.TextColor = req.TextColor
	}
	if req.BackgroundColor != "" {
		preset.BackgroundColor = req.BackgroundColor
	}
	if req.ShowBorder != nil {
		preset.ShowBorder = *req.ShowBorder
	}
	if req.IsDefault != nil {
		preset.IsDefault = *req.IsDefault
	}
	if req.CustomCSS != "" {
		preset.CustomCSS = req.CustomCSS
	}

	err = s.repo.Update(preset)
	if err != nil {
		return nil, err
	}

	return preset, nil
}

func (s *LayoutPresetService) SetDefault(id, userID string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id format")
	}
	return s.repo.SetDefault(uid, userID)
}

func (s *LayoutPresetService) Delete(id, userID string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	preset, err := s.repo.FindByID(uid)
	if err != nil {
		return errors.New("preset not found")
	}
	if preset.UserID != userID {
		return errors.New("permission denied")
	}

	return s.repo.Delete(uid)
}

type CreateLayoutPresetRequest struct {
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description"`
	CharsPerColumn  int     `json:"chars_per_column"`
	ColumnGap       int     `json:"column_gap"`
	FontSize        int     `json:"font_size"`
	FontFamily      string  `json:"font_family"`
	LineHeight      float64 `json:"line_height"`
	CharSpacing     float64 `json:"char_spacing"`
	TextColor       string  `json:"text_color"`
	BackgroundColor string  `json:"background_color"`
	ShowBorder      bool    `json:"show_border"`
	IsDefault       bool    `json:"is_default"`
	CustomCSS       string  `json:"custom_css"`
}

type UpdateLayoutPresetRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	CharsPerColumn  int      `json:"chars_per_column"`
	ColumnGap       int      `json:"column_gap"`
	FontSize        int      `json:"font_size"`
	FontFamily      string   `json:"font_family"`
	LineHeight      float64  `json:"line_height"`
	CharSpacing     float64  `json:"char_spacing"`
	TextColor       string   `json:"text_color"`
	BackgroundColor string   `json:"background_color"`
	ShowBorder      *bool    `json:"show_border"`
	IsDefault       *bool    `json:"is_default"`
	CustomCSS       string   `json:"custom_css"`
}
