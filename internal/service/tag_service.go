package service

import (
	"context"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/repository"
)

// TagService handles business logic for tags
type TagService struct {
	repo repository.Repository
}

// NewTagService creates a new tag service
func NewTagService(repo repository.Repository) *TagService {
	return &TagService{repo: repo}
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(ctx context.Context, name, color string, isDefault bool) (*domain.Tag, error) {
	if name == "" {
		return nil, domain.ErrMissingRequiredField
	}

	tag := &domain.Tag{
		Name:      name,
		Color:     color,
		IsDefault: isDefault,
	}

	if err := s.repo.CreateTag(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// GetTag retrieves a tag by ID
func (s *TagService) GetTag(ctx context.Context, id string) (*domain.Tag, error) {
	return s.repo.GetTag(ctx, id)
}

// ListTags retrieves all tags
func (s *TagService) ListTags(ctx context.Context) ([]*domain.Tag, error) {
	return s.repo.ListTags(ctx)
}

// UpdateTag updates a tag
func (s *TagService) UpdateTag(ctx context.Context, id string, name, color *string, isDefault *bool) (*domain.Tag, error) {
	if err := s.repo.UpdateTag(ctx, id, name, color, isDefault); err != nil {
		return nil, err
	}
	return s.repo.GetTag(ctx, id)
}

// DeleteTag deletes a tag
func (s *TagService) DeleteTag(ctx context.Context, id string) error {
	return s.repo.DeleteTag(ctx, id)
}
