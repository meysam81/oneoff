package service

import (
	"context"
	"fmt"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/repository"
)

// ChainService handles business logic for job chains
type ChainService struct {
	repo       repository.Repository
	jobService *JobService
}

// NewChainService creates a new chain service
func NewChainService(repo repository.Repository, jobService *JobService) *ChainService {
	return &ChainService{
		repo:       repo,
		jobService: jobService,
	}
}

// CreateChain creates a new job chain
func (s *ChainService) CreateChain(ctx context.Context, req domain.CreateChainRequest) (*domain.JobChain, error) {
	// Validate chain has at least one link
	if len(req.Links) == 0 {
		return nil, domain.ErrChainEmpty
	}

	// Validate all jobs exist
	for _, link := range req.Links {
		if _, err := s.repo.GetJob(ctx, link.JobID); err != nil {
			return nil, domain.ErrChainJobNotFound
		}
	}

	// Set default project
	projectID := req.ProjectID
	if projectID == "" {
		projectID = "default"
	}

	// Validate project exists
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Create chain with links
	chain := &domain.JobChain{
		Name:      req.Name,
		ProjectID: projectID,
	}

	// Convert ChainLinkInput to JobChainLink
	for i, link := range req.Links {
		chain.Links = append(chain.Links, domain.JobChainLink{
			JobID:         link.JobID,
			SequenceOrder: i + 1,
			StopOnFailure: link.StopOnFailure,
		})
	}

	if err := s.repo.CreateChain(ctx, chain); err != nil {
		return nil, fmt.Errorf("failed to create chain: %w", err)
	}

	// Reload chain to get full data
	return s.repo.GetChain(ctx, chain.ID)
}

// GetChain retrieves a chain by ID with all links
func (s *ChainService) GetChain(ctx context.Context, id string) (*domain.JobChain, error) {
	return s.repo.GetChain(ctx, id)
}

// ListChains lists chains with optional filtering
func (s *ChainService) ListChains(ctx context.Context, filter domain.ChainFilter) ([]*domain.JobChain, int64, error) {
	chains, err := s.repo.ListChainsWithFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.CountChains(ctx, filter)
	if err != nil {
		return chains, 0, err
	}

	return chains, count, nil
}

// UpdateChain updates a chain's name and/or links
func (s *ChainService) UpdateChain(ctx context.Context, id string, req domain.UpdateChainRequest) (*domain.JobChain, error) {
	// Verify chain exists
	if _, err := s.repo.GetChain(ctx, id); err != nil {
		return nil, err
	}

	// Update name if provided
	if req.Name != nil {
		if err := s.repo.UpdateChain(ctx, id, req.Name); err != nil {
			return nil, err
		}
	}

	// Update links if provided
	if req.Links != nil {
		// Validate chain has at least one link
		if len(req.Links) == 0 {
			return nil, domain.ErrChainEmpty
		}

		// Validate all jobs exist
		for _, link := range req.Links {
			if _, err := s.repo.GetJob(ctx, link.JobID); err != nil {
				return nil, domain.ErrChainJobNotFound
			}
		}

		if err := s.repo.UpdateChainLinks(ctx, id, req.Links); err != nil {
			return nil, err
		}
	}

	// Return updated chain
	return s.repo.GetChain(ctx, id)
}

// DeleteChain deletes a chain
func (s *ChainService) DeleteChain(ctx context.Context, id string) error {
	return s.repo.DeleteChain(ctx, id)
}

// ExecuteChain executes all jobs in a chain sequentially
func (s *ChainService) ExecuteChain(ctx context.Context, id string) error {
	chain, err := s.repo.GetChain(ctx, id)
	if err != nil {
		return err
	}

	if len(chain.Links) == 0 {
		return domain.ErrChainEmpty
	}

	// Execute jobs in sequence
	for _, link := range chain.Links {
		// Get the job
		job, err := s.repo.GetJob(ctx, link.JobID)
		if err != nil {
			if link.StopOnFailure {
				return fmt.Errorf("job %s not found in chain: %w", link.JobID, err)
			}
			continue
		}

		// Update job to execute immediately
		now := time.Now().UTC().Format(time.RFC3339)
		updates := domain.UpdateJobRequest{
			ScheduledAt: &now,
			Status:      stringPtr("scheduled"),
		}

		if err := s.repo.UpdateJob(ctx, job.ID, updates); err != nil {
			if link.StopOnFailure {
				return fmt.Errorf("failed to schedule job %s: %w", link.JobID, err)
			}
			continue
		}
	}

	return nil
}
