package service

import (
	"context"
	"fmt"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/repository"
	"github.com/meysam81/oneoff/internal/worker"
)

// JobService handles business logic for jobs
type JobService struct {
	repo     repository.Repository
	registry *domain.JobRegistry
	pool     *worker.Pool
}

// NewJobService creates a new job service
func NewJobService(repo repository.Repository, registry *domain.JobRegistry, pool *worker.Pool) *JobService {
	return &JobService{
		repo:     repo,
		registry: registry,
		pool:     pool,
	}
}

// CreateJob creates a new job
func (s *JobService) CreateJob(ctx context.Context, req domain.CreateJobRequest) (*domain.Job, error) {
	if _, err := s.registry.Create(req.Type, req.Config); err != nil {
		return nil, fmt.Errorf("invalid job type or config: %w", err)
	}

	var scheduledAt time.Time

	if req.Immediate || req.ScheduledAt == "now" {
		scheduledAt = time.Now().UTC()
	} else {
		if req.ScheduledAt == "" {
			return nil, fmt.Errorf("scheduled_at is required when immediate is false")
		}

		parsed, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			return nil, fmt.Errorf("invalid scheduled_at format (use RFC3339): %w", err)
		}

		if parsed.Before(time.Now().UTC()) {
			return nil, domain.ErrInvalidScheduleTime
		}

		scheduledAt = parsed.UTC()
	}

	priority := req.Priority
	if priority == 0 {
		priority = 5
	}
	if priority < 1 || priority > 10 {
		return nil, domain.ErrInvalidPriority
	}

	projectID := req.ProjectID
	if projectID == "" {
		projectID = "default"
	}

	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	timezone := req.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	job := &domain.Job{
		Name:        req.Name,
		Type:        req.Type,
		Config:      req.Config,
		ScheduledAt: scheduledAt,
		Priority:    priority,
		ProjectID:   projectID,
		Timezone:    timezone,
		Status:      domain.JobStatusScheduled,
	}

	if err := s.repo.CreateJob(ctx, job, req.TagIDs); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	tags, _ := s.repo.GetJobTags(ctx, job.ID)
	job.Tags = tags

	return job, nil
}

// GetJob retrieves a job by ID
func (s *JobService) GetJob(ctx context.Context, id string) (*domain.Job, error) {
	return s.repo.GetJob(ctx, id)
}

// ListJobs retrieves jobs based on filter
func (s *JobService) ListJobs(ctx context.Context, filter domain.JobFilter) ([]*domain.Job, int64, error) {
	jobs, err := s.repo.ListJobs(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.CountJobs(ctx, filter)
	if err != nil {
		return jobs, 0, err
	}

	return jobs, count, nil
}

// UpdateJob updates a job
func (s *JobService) UpdateJob(ctx context.Context, id string, updates domain.UpdateJobRequest) (*domain.Job, error) {
	// Get existing job
	job, err := s.repo.GetJob(ctx, id)
	if err != nil {
		return nil, err
	}

	// Only allow updates to scheduled jobs
	if job.Status != domain.JobStatusScheduled {
		return nil, fmt.Errorf("cannot update job in status: %s", job.Status)
	}

	// Validate updates
	if updates.ScheduledAt != nil {
		scheduledAt, err := time.Parse(time.RFC3339, *updates.ScheduledAt)
		if err != nil {
			return nil, fmt.Errorf("invalid scheduled_at format: %w", err)
		}
		if scheduledAt.Before(time.Now().UTC()) {
			return nil, domain.ErrInvalidScheduleTime
		}
	}

	if updates.Priority != nil {
		if *updates.Priority < 1 || *updates.Priority > 10 {
			return nil, domain.ErrInvalidPriority
		}
	}

	if updates.ProjectID != nil {
		if _, err := s.repo.GetProject(ctx, *updates.ProjectID); err != nil {
			return nil, fmt.Errorf("project not found: %w", err)
		}
	}

	// Update job
	if err := s.repo.UpdateJob(ctx, id, updates); err != nil {
		return nil, err
	}

	return s.repo.GetJob(ctx, id)
}

// DeleteJob deletes a job
func (s *JobService) DeleteJob(ctx context.Context, id string) error {
	job, err := s.repo.GetJob(ctx, id)
	if err != nil {
		return err
	}

	// Only allow deletion of scheduled or failed/completed jobs
	if job.Status == domain.JobStatusRunning {
		return fmt.Errorf("cannot delete running job, cancel it first")
	}

	return s.repo.DeleteJob(ctx, id)
}

// CancelJob cancels a job
func (s *JobService) CancelJob(ctx context.Context, id string) error {
	job, err := s.repo.GetJob(ctx, id)
	if err != nil {
		return err
	}

	if job.Status == domain.JobStatusCancelled {
		return fmt.Errorf("job already cancelled")
	}

	if job.Status == domain.JobStatusCompleted || job.Status == domain.JobStatusFailed {
		return fmt.Errorf("cannot cancel job in status: %s", job.Status)
	}

	// If running, ask pool to cancel
	if job.Status == domain.JobStatusRunning {
		return s.pool.CancelJob(ctx, id)
	}

	// Otherwise just update status
	return s.repo.UpdateJobStatus(ctx, id, domain.JobStatusCancelled)
}

// ExecuteJobNow executes a job immediately
func (s *JobService) ExecuteJobNow(ctx context.Context, id string) error {
	job, err := s.repo.GetJob(ctx, id)
	if err != nil {
		return err
	}

	if job.Status == domain.JobStatusRunning {
		return fmt.Errorf("job is already running")
	}

	// Update scheduled time to now
	now := time.Now().UTC().Format(time.RFC3339)
	updates := domain.UpdateJobRequest{
		ScheduledAt: &now,
		Status:      stringPtr("scheduled"),
	}

	return s.repo.UpdateJob(ctx, id, updates)
}

// CloneJob creates a copy of an existing job
func (s *JobService) CloneJob(ctx context.Context, id string, newScheduledAt time.Time) (*domain.Job, error) {
	original, err := s.repo.GetJob(ctx, id)
	if err != nil {
		return nil, err
	}

	tagIDs := make([]string, len(original.Tags))
	for i, tag := range original.Tags {
		tagIDs[i] = tag.ID
	}

	req := domain.CreateJobRequest{
		Name:        original.Name + " (clone)",
		Type:        original.Type,
		Config:      original.Config,
		ScheduledAt: newScheduledAt.Format(time.RFC3339),
		Priority:    original.Priority,
		ProjectID:   original.ProjectID,
		Timezone:    original.Timezone,
		TagIDs:      tagIDs,
	}

	return s.CreateJob(ctx, req)
}

func stringPtr(s string) *string {
	return &s
}
