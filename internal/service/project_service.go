package service

import (
	"context"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/repository"
)

// ProjectService handles business logic for projects
type ProjectService struct {
	repo repository.Repository
}

// NewProjectService creates a new project service
func NewProjectService(repo repository.Repository) *ProjectService {
	return &ProjectService{repo: repo}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, name, description, color, icon string) (*domain.Project, error) {
	if name == "" {
		return nil, domain.ErrMissingRequiredField
	}

	project := &domain.Project{
		Name:        name,
		Description: description,
		Color:       color,
		Icon:        icon,
		IsArchived:  false,
	}

	if err := s.repo.CreateProject(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	return s.repo.GetProject(ctx, id)
}

// ListProjects retrieves all projects
func (s *ProjectService) ListProjects(ctx context.Context, includeArchived bool) ([]*domain.Project, error) {
	return s.repo.ListProjects(ctx, includeArchived)
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(ctx context.Context, id string, name, description, color, icon *string, isArchived *bool) (*domain.Project, error) {
	if err := s.repo.UpdateProject(ctx, id, name, description, color, icon, isArchived); err != nil {
		return nil, err
	}
	return s.repo.GetProject(ctx, id)
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	return s.repo.DeleteProject(ctx, id)
}
