package service

import (
	"errors"

	"ancient-text-backend/internal/cache"
	"ancient-text-backend/internal/models"
	"ancient-text-backend/internal/repository"

	"github.com/google/uuid"
)

type AnnotationService struct {
	repo *repository.AnnotationRepository
}

func NewAnnotationService() *AnnotationService {
	return &AnnotationService{
		repo: repository.NewAnnotationRepository(),
	}
}

func (s *AnnotationService) Create(req *CreateAnnotationRequest) (*models.Annotation, error) {
	textID, err := uuid.Parse(req.TextID)
	if err != nil {
		return nil, errors.New("invalid text id format")
	}

	if req.StartPos < 0 || req.EndPos <= req.StartPos {
		return nil, errors.New("invalid position range")
	}
	if req.Type == "" {
		return nil, errors.New("annotation type is required")
	}

	ann := &models.Annotation{
		TextID:   textID,
		UserID:   req.UserID,
		StartPos: req.StartPos,
		EndPos:   req.EndPos,
		Type:     req.Type,
		Content:  req.Content,
		Color:    req.Color,
		Status:   "draft",
		Version:  1,
	}

	err = s.repo.Create(ann)
	if err != nil {
		return nil, err
	}

	_ = cache.CacheAnnotation(ann)
	_ = cache.DeleteCachedAnnotationList(req.TextID)

	return ann, nil
}

func (s *AnnotationService) GetByID(id string) (*models.Annotation, error) {
	cached, err := cache.GetCachedAnnotation(id)
	if err == nil && cached != nil {
		return cached, nil
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	ann, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	_ = cache.CacheAnnotation(ann)
	return ann, nil
}

func (s *AnnotationService) GetByTextID(textID string) ([]models.Annotation, error) {
	cached, err := cache.GetCachedAnnotationList(textID)
	if err == nil && cached != nil {
		return cached, nil
	}

	uid, err := uuid.Parse(textID)
	if err != nil {
		return nil, errors.New("invalid text id format")
	}

	annotations, err := s.repo.FindByTextID(uid)
	if err != nil {
		return nil, err
	}

	if len(annotations) > 0 {
		_ = cache.CacheAnnotationList(textID, annotations)
	}

	return annotations, nil
}

func (s *AnnotationService) GetByTextIDAndUser(textID, userID string) ([]models.Annotation, error) {
	uid, err := uuid.Parse(textID)
	if err != nil {
		return nil, errors.New("invalid text id format")
	}
	return s.repo.FindByTextIDAndUser(uid, userID)
}

func (s *AnnotationService) Update(id string, req *UpdateAnnotationRequest) (*models.Annotation, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	ann, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, errors.New("annotation not found")
	}

	if req.StartPos >= 0 {
		ann.StartPos = req.StartPos
	}
	if req.EndPos > 0 {
		ann.EndPos = req.EndPos
	}
	if req.Type != "" {
		ann.Type = req.Type
	}
	if req.Content != "" {
		ann.Content = req.Content
	}
	if req.Color != "" {
		ann.Color = req.Color
	}
	if req.Status != "" {
		ann.Status = req.Status
	}
	ann.Version++

	err = s.repo.Update(ann)
	if err != nil {
		return nil, err
	}

	_ = cache.CacheAnnotation(ann)
	_ = cache.DeleteCachedAnnotationList(ann.TextID.String())

	return ann, nil
}

func (s *AnnotationService) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	ann, err := s.repo.FindByID(uid)
	if err == nil {
		_ = cache.DeleteCachedAnnotationList(ann.TextID.String())
	}

	_ = cache.DeleteCachedAnnotation(id)
	return s.repo.Delete(uid)
}

func (s *AnnotationService) BatchSubmit(textID, userID string, reqs []CreateAnnotationRequest) ([]models.Annotation, error) {
	uid, err := uuid.Parse(textID)
	if err != nil {
		return nil, errors.New("invalid text id format")
	}

	annotations := make([]models.Annotation, 0, len(reqs))
	for _, req := range reqs {
		if req.StartPos < 0 || req.EndPos <= req.StartPos || req.Type == "" {
			continue
		}
		ann := models.Annotation{
			TextID:   uid,
			UserID:   userID,
			StartPos: req.StartPos,
			EndPos:   req.EndPos,
			Type:     req.Type,
			Content:  req.Content,
			Color:    req.Color,
			Status:   "submitted",
			Version:  1,
		}
		annotations = append(annotations, ann)
	}

	err = s.repo.BatchCreate(annotations)
	if err != nil {
		return nil, err
	}

	_ = cache.DeleteCachedAnnotationList(textID)

	return annotations, nil
}

type CreateAnnotationRequest struct {
	TextID   string `json:"text_id" binding:"required"`
	UserID   string `json:"user_id"`
	StartPos int    `json:"start_pos" binding:"required"`
	EndPos   int    `json:"end_pos" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Content  string `json:"content"`
	Color    string `json:"color"`
}

type UpdateAnnotationRequest struct {
	StartPos int    `json:"start_pos"`
	EndPos   int    `json:"end_pos"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	Color    string `json:"color"`
	Status   string `json:"status"`
}
