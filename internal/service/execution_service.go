package service

import (
	"context"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/repository"
)

// ExecutionService handles business logic for job executions
type ExecutionService struct {
	repo repository.Repository
}

// NewExecutionService creates a new execution service
func NewExecutionService(repo repository.Repository) *ExecutionService {
	return &ExecutionService{repo: repo}
}

// GetExecution retrieves an execution by ID
func (s *ExecutionService) GetExecution(ctx context.Context, id string) (*domain.JobExecution, error) {
	return s.repo.GetExecution(ctx, id)
}

// ListExecutions retrieves executions based on filter
func (s *ExecutionService) ListExecutions(ctx context.Context, filter domain.ExecutionFilter) ([]*domain.JobExecution, error) {
	return s.repo.ListExecutions(ctx, filter)
}

// GetJobExecutions retrieves all executions for a job
func (s *ExecutionService) GetJobExecutions(ctx context.Context, jobID string, limit, offset int) ([]*domain.JobExecution, error) {
	filter := domain.ExecutionFilter{
		JobID:     jobID,
		Limit:     limit,
		Offset:    offset,
		SortBy:    "started_at",
		SortOrder: "desc",
	}
	return s.repo.ListExecutions(ctx, filter)
}
