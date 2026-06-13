package service

import (
	"errors"

	"ancient-text-backend/internal/cache"
	"ancient-text-backend/internal/models"
	"ancient-text-backend/internal/repository"

	"github.com/google/uuid"
)

type LayoutConfigService struct {
	repo *repository.LayoutConfigRepository
}

func NewLayoutConfigService() *LayoutConfigService {
	return &LayoutConfigService{
		repo: repository.NewLayoutConfigRepository(),
	}
}

func (s *LayoutConfigService) GetByTextID(textID string) (*models.LayoutConfig, error) {
	cached, err := cache.GetCachedLayoutConfig(textID)
	if err == nil && cached != nil {
		return cached, nil
	}

	uid, err := uuid.Parse(textID)
	if err != nil {
		return nil, errors.New("invalid text id format")
	}

	layout, err := s.repo.FindByTextID(uid)
	if err != nil {
		return nil, err
	}

	if layout == nil {
		layout = &models.LayoutConfig{
			TextID:          uid,
			Columns:         12,
			CharsPerColumn:  20,
			ColumnGap:       20,
			FontSize:        18,
			FontFamily:      "SimSun, serif",
			LineHeight:      1.8,
			CharSpacing:     0,
			TextColor:       "#1a1a1a",
			BackgroundColor: "#f5f0e8",
			ShowBorder:      true,
		}
	}

	if layout.ID != uuid.Nil {
		_ = cache.CacheLayoutConfig(layout)
	}

	return layout, nil
}

func (s *LayoutConfigService) Save(textID string, req *SaveLayoutConfigRequest) (*models.LayoutConfig, error) {
	uid, err := uuid.Parse(textID)
	if err != nil {
		return nil, errors.New("invalid text id format")
	}

	layout := &models.LayoutConfig{
		TextID:          uid,
		Columns:         req.Columns,
		CharsPerColumn:  req.CharsPerColumn,
		ColumnGap:       req.ColumnGap,
		FontSize:        req.FontSize,
		FontFamily:      req.FontFamily,
		LineHeight:      req.LineHeight,
		CharSpacing:     req.CharSpacing,
		TextColor:       req.TextColor,
		BackgroundColor: req.BackgroundColor,
		ShowBorder:      req.ShowBorder,
		CustomCSS:       req.CustomCSS,
	}

	err = s.repo.Upsert(layout)
	if err != nil {
		return nil, err
	}

	_ = cache.CacheLayoutConfig(layout)

	return layout, nil
}

func (s *LayoutConfigService) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	layout, err := s.repo.FindByID(uid)
	if err == nil && layout != nil {
		_ = cache.DeleteCachedLayoutConfig(layout.TextID.String())
	}

	return s.repo.Delete(uid)
}

type SaveLayoutConfigRequest struct {
	Columns         int     `json:"columns"`
	CharsPerColumn  int     `json:"chars_per_column"`
	ColumnGap       int     `json:"column_gap"`
	FontSize        int     `json:"font_size"`
	FontFamily      string  `json:"font_family"`
	LineHeight      float64 `json:"line_height"`
	CharSpacing     float64 `json:"char_spacing"`
	TextColor       string  `json:"text_color"`
	BackgroundColor string  `json:"background_color"`
	ShowBorder      bool    `json:"show_border"`
	CustomCSS       string  `json:"custom_css"`
}
