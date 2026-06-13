package service

import (
	"errors"

	"ancient-text-backend/internal/models"
	"ancient-text-backend/internal/repository"

	"github.com/google/uuid"
)

type AncientTextService struct {
	repo     *repository.AncientTextRepository
	layoutRepo *repository.LayoutConfigRepository
}

func NewAncientTextService() *AncientTextService {
	return &AncientTextService{
		repo:       repository.NewAncientTextRepository(),
		layoutRepo: repository.NewLayoutConfigRepository(),
	}
}

func (s *AncientTextService) Create(req *CreateTextRequest) (*models.AncientText, error) {
	if req.Title == "" || req.Content == "" {
		return nil, errors.New("title and content are required")
	}

	text := &models.AncientText{
		Title:       req.Title,
		Author:      req.Author,
		Dynasty:     req.Dynasty,
		Content:     req.Content,
		Description: req.Description,
	}

	err := s.repo.Create(text)
	if err != nil {
		return nil, err
	}

	if req.LayoutConfig != nil {
		layout := &models.LayoutConfig{
			TextID:         text.ID,
			Columns:        req.LayoutConfig.Columns,
			CharsPerColumn: req.LayoutConfig.CharsPerColumn,
			ColumnGap:      req.LayoutConfig.ColumnGap,
			FontSize:       req.LayoutConfig.FontSize,
			FontFamily:     req.LayoutConfig.FontFamily,
			LineHeight:     req.LayoutConfig.LineHeight,
			TextColor:      req.LayoutConfig.TextColor,
			BackgroundColor: req.LayoutConfig.BackgroundColor,
			ShowBorder:     req.LayoutConfig.ShowBorder,
			CustomCSS:      req.LayoutConfig.CustomCSS,
		}
		_ = s.layoutRepo.Create(layout)
		text.LayoutConfig = layout
	}

	return text, nil
}

func (s *AncientTextService) GetByID(id string) (*models.AncientText, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	return s.repo.FindByID(uid)
}

func (s *AncientTextService) List(page, pageSize int) ([]models.AncientText, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.FindAll(page, pageSize)
}

func (s *AncientTextService) Search(keyword string, page, pageSize int) ([]models.AncientText, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.Search(keyword, page, pageSize)
}

func (s *AncientTextService) Update(id string, req *UpdateTextRequest) (*models.AncientText, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	text, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, errors.New("text not found")
	}

	if req.Title != "" {
		text.Title = req.Title
	}
	if req.Author != "" {
		text.Author = req.Author
	}
	if req.Dynasty != "" {
		text.Dynasty = req.Dynasty
	}
	if req.Content != "" {
		text.Content = req.Content
	}
	if req.Description != "" {
		text.Description = req.Description
	}

	err = s.repo.Update(text)
	if err != nil {
		return nil, err
	}

	return text, nil
}

func (s *AncientTextService) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id format")
	}
	return s.repo.Delete(uid)
}

type CreateTextRequest struct {
	Title        string                    `json:"title" binding:"required"`
	Author       string                    `json:"author"`
	Dynasty      string                    `json:"dynasty"`
	Content      string                    `json:"content" binding:"required"`
	Description  string                    `json:"description"`
	LayoutConfig *LayoutConfigRequest      `json:"layout_config"`
}

type UpdateTextRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Dynasty     string `json:"dynasty"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

type LayoutConfigRequest struct {
	Columns         int     `json:"columns"`
	CharsPerColumn  int     `json:"chars_per_column"`
	ColumnGap       int     `json:"column_gap"`
	FontSize        int     `json:"font_size"`
	FontFamily      string  `json:"font_family"`
	LineHeight      float64 `json:"line_height"`
	TextColor       string  `json:"text_color"`
	BackgroundColor string  `json:"background_color"`
	ShowBorder      bool    `json:"show_border"`
	CustomCSS       string  `json:"custom_css"`
}
