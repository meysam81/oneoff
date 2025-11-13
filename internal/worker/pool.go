package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/repository"
	"github.com/rs/zerolog/log"
)

// Pool manages a pool of workers for executing jobs
type Pool struct {
	workers      int
	repo         repository.Repository
	registry     *domain.JobRegistry
	jobChan      chan *domain.Job
	stopChan     chan struct{}
	wg           sync.WaitGroup
	runningJobs  map[string]bool
	runningMutex sync.RWMutex
	pollInterval time.Duration
}

// NewPool creates a new worker pool
func NewPool(workers int, repo repository.Repository, registry *domain.JobRegistry) *Pool {
	return &Pool{
		workers:      workers,
		repo:         repo,
		registry:     registry,
		jobChan:      make(chan *domain.Job, workers*2),
		stopChan:     make(chan struct{}),
		runningJobs:  make(map[string]bool),
		pollInterval: 5 * time.Second, // Check for new jobs every 5 seconds
	}
}

// Start starts the worker pool
func (p *Pool) Start(ctx context.Context) error {
	log.Info().Int("workers", p.workers).Msg("Starting worker pool")

	// Start workers
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}

	// Start scheduler
	p.wg.Add(1)
	go p.scheduler(ctx)

	log.Info().Msg("Worker pool started")
	return nil
}

// Stop stops the worker pool gracefully
func (p *Pool) Stop(ctx context.Context) error {
	log.Info().Msg("Stopping worker pool")

	close(p.stopChan)
	close(p.jobChan)

	// Wait for all workers to finish with timeout
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("Worker pool stopped gracefully")
		return nil
	case <-ctx.Done():
		log.Warn().Msg("Worker pool shutdown timeout")
		return fmt.Errorf("worker pool shutdown timeout")
	}
}

// worker processes jobs from the job channel
func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	log.Debug().Int("worker_id", id).Msg("Worker started")

	for {
		select {
		case <-ctx.Done():
			log.Debug().Int("worker_id", id).Msg("Worker context cancelled")
			return
		case <-p.stopChan:
			log.Debug().Int("worker_id", id).Msg("Worker stopped")
			return
		case job, ok := <-p.jobChan:
			if !ok {
				log.Debug().Int("worker_id", id).Msg("Job channel closed")
				return
			}

			if job == nil {
				continue
			}

			log.Info().
				Int("worker_id", id).
				Str("job_id", job.ID).
				Str("job_type", job.Type).
				Str("job_name", job.Name).
				Msg("Executing job")

			p.executeJob(ctx, job)
		}
	}
}

// scheduler polls for scheduled jobs and sends them to workers
func (p *Pool) scheduler(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	log.Debug().Msg("Scheduler started")

	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("Scheduler context cancelled")
			return
		case <-p.stopChan:
			log.Debug().Msg("Scheduler stopped")
			return
		case <-ticker.C:
			p.pollJobs(ctx)
		}
	}
}

// pollJobs checks for jobs that need to be executed
func (p *Pool) pollJobs(ctx context.Context) {
	// Get jobs that are scheduled for now or earlier
	jobs, err := p.repo.GetScheduledJobs(ctx, time.Now().UTC(), p.workers*2)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get scheduled jobs")
		return
	}

	if len(jobs) == 0 {
		return
	}

	log.Debug().Int("count", len(jobs)).Msg("Found scheduled jobs")

	for _, job := range jobs {
		// Skip if already running
		p.runningMutex.RLock()
		isRunning := p.runningJobs[job.ID]
		p.runningMutex.RUnlock()

		if isRunning {
			continue
		}

		// Try to send to job channel (non-blocking)
		select {
		case p.jobChan <- job:
			log.Debug().Str("job_id", job.ID).Msg("Job queued for execution")
		default:
			log.Warn().Str("job_id", job.ID).Msg("Job channel full, will retry on next poll")
		}
	}
}

// executeJob executes a single job
func (p *Pool) executeJob(ctx context.Context, job *domain.Job) {
	// Mark as running
	p.runningMutex.Lock()
	p.runningJobs[job.ID] = true
	p.runningMutex.Unlock()

	defer func() {
		p.runningMutex.Lock()
		delete(p.runningJobs, job.ID)
		p.runningMutex.Unlock()
	}()

	startTime := time.Now()

	// Update job status to running
	if err := p.repo.UpdateJobStatus(ctx, job.ID, domain.JobStatusRunning); err != nil {
		log.Error().Err(err).Str("job_id", job.ID).Msg("Failed to update job status to running")
		return
	}

	// Create execution record
	execution := &domain.JobExecution{
		JobID:     job.ID,
		StartedAt: startTime,
		Status:    domain.ExecutionStatusRunning,
	}

	if err := p.repo.CreateExecution(ctx, execution); err != nil {
		log.Error().Err(err).Str("job_id", job.ID).Msg("Failed to create execution record")
		return
	}

	// Create job executor
	executor, err := p.registry.Create(job.Type, job.Config)
	if err != nil {
		log.Error().Err(err).
			Str("job_id", job.ID).
			Str("job_type", job.Type).
			Msg("Failed to create job executor")

		p.completeExecution(ctx, execution.ID, job.ID, domain.ExecutionStatusFailed, "", fmt.Sprintf("Failed to create executor: %v", err), nil, time.Since(startTime))
		return
	}

	// Execute job
	result, err := executor.Execute(ctx)
	if err != nil {
		log.Error().Err(err).Str("job_id", job.ID).Msg("Job execution error")

		p.completeExecution(ctx, execution.ID, job.ID, domain.ExecutionStatusFailed, "", fmt.Sprintf("Execution error: %v", err), nil, time.Since(startTime))
		return
	}

	// Determine final status
	finalStatus := domain.ExecutionStatusCompleted
	finalJobStatus := domain.JobStatusCompleted
	if result.ExitCode != 0 {
		finalStatus = domain.ExecutionStatusFailed
		finalJobStatus = domain.JobStatusFailed
	}

	// Complete execution
	durationMs := time.Since(startTime).Milliseconds()
	p.completeExecution(ctx, execution.ID, job.ID, finalStatus, result.Output, result.Error, &result.ExitCode, time.Since(startTime))

	// Update job status
	if err := p.repo.UpdateJobStatus(ctx, job.ID, finalJobStatus); err != nil {
		log.Error().Err(err).Str("job_id", job.ID).Msg("Failed to update job final status")
	}

	log.Info().
		Str("job_id", job.ID).
		Str("job_name", job.Name).
		Str("status", string(finalStatus)).
		Int64("duration_ms", durationMs).
		Int("exit_code", result.ExitCode).
		Msg("Job execution completed")
}

// completeExecution marks an execution as complete
func (p *Pool) completeExecution(ctx context.Context, executionID, jobID string, status domain.ExecutionStatus, output, errorMsg string, exitCode *int, duration time.Duration) {
	durationMs := duration.Milliseconds()

	if err := p.repo.CompleteExecution(ctx, executionID, status, output, errorMsg, exitCode, durationMs); err != nil {
		log.Error().Err(err).
			Str("execution_id", executionID).
			Str("job_id", jobID).
			Msg("Failed to complete execution")
	}
}

// GetStatus returns the current status of the worker pool
func (p *Pool) GetStatus(ctx context.Context) *domain.WorkerStatus {
	p.runningMutex.RLock()
	activeWorkers := len(p.runningJobs)
	runningJobIDs := make([]string, 0, len(p.runningJobs))
	for jobID := range p.runningJobs {
		runningJobIDs = append(runningJobIDs, jobID)
	}
	p.runningMutex.RUnlock()

	queuedJobs := len(p.jobChan)

	return &domain.WorkerStatus{
		TotalWorkers:     p.workers,
		ActiveWorkers:    activeWorkers,
		AvailableWorkers: p.workers - activeWorkers,
		QueuedJobs:       queuedJobs,
		RunningJobs:      runningJobIDs,
	}
}

// CancelJob cancels a running job
func (p *Pool) CancelJob(ctx context.Context, jobID string) error {
	p.runningMutex.RLock()
	isRunning := p.runningJobs[jobID]
	p.runningMutex.RUnlock()

	if !isRunning {
		// Job not currently running, just update status
		return p.repo.UpdateJobStatus(ctx, jobID, domain.JobStatusCancelled)
	}

	// For running jobs, we'd need to track contexts to cancel them
	// This is a simplified implementation
	log.Warn().Str("job_id", jobID).Msg("Job cancellation requested but not yet implemented for running jobs")

	return p.repo.UpdateJobStatus(ctx, jobID, domain.JobStatusCancelled)
}
