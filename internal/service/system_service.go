package service

import (
	"context"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/repository"
	"github.com/meysam81/oneoff/internal/worker"
)

// SystemService handles business logic for system operations
type SystemService struct {
	repo repository.Repository
	pool *worker.Pool
}

// NewSystemService creates a new system service
func NewSystemService(repo repository.Repository, pool *worker.Pool) *SystemService {
	return &SystemService{
		repo: repo,
		pool: pool,
	}
}

// GetSystemStats retrieves system-wide statistics
func (s *SystemService) GetSystemStats(ctx context.Context) (*domain.SystemStats, error) {
	return s.repo.GetSystemStats(ctx)
}

// GetWorkerStatus retrieves worker pool status
func (s *SystemService) GetWorkerStatus(ctx context.Context) *domain.WorkerStatus {
	return s.pool.GetStatus(ctx)
}

// GetConfig retrieves system configuration
func (s *SystemService) GetConfig(ctx context.Context) ([]*domain.SystemConfig, error) {
	return s.repo.ListConfig(ctx)
}

// UpdateConfig updates system configuration
func (s *SystemService) UpdateConfig(ctx context.Context, key, value string) error {
	return s.repo.SetConfig(ctx, key, value)
}

// GetJobTypes retrieves all available job types
func (s *SystemService) GetJobTypes() []string {
	// This would come from the registry
	return []string{"http", "shell", "docker"}
}
