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

// JobEventCallback is called when a job event occurs
type JobEventCallback func(ctx context.Context, event domain.WebhookEvent)

// MetricsCallback is called to report job metrics
type MetricsCallback func(jobType, status string, duration time.Duration)

// Pool manages a pool of workers for executing jobs
type Pool struct {
	workers          int
	repo             repository.Repository
	registry         *domain.JobRegistry
	jobChan          chan *domain.Job
	stopChan         chan struct{}
	wg               sync.WaitGroup
	runningJobs      map[string]bool
	jobContexts      map[string]context.CancelFunc // Track cancel functions for running jobs
	runningMutex     sync.RWMutex
	pollInterval     time.Duration
	logRetentionDays int              // Days to retain execution logs (0 = no cleanup)
	cleanupInterval  time.Duration    // How often to run cleanup
	onJobEvent       JobEventCallback // Callback for job events (webhooks)
	onMetrics        MetricsCallback  // Callback for metrics
}

// NewPool creates a new worker pool
func NewPool(workers int, repo repository.Repository, registry *domain.JobRegistry) *Pool {
	return &Pool{
		workers:          workers,
		repo:             repo,
		registry:         registry,
		jobChan:          make(chan *domain.Job, workers*2),
		stopChan:         make(chan struct{}),
		runningJobs:      make(map[string]bool),
		jobContexts:      make(map[string]context.CancelFunc),
		pollInterval:     5 * time.Second,  // Check for new jobs every 5 seconds
		logRetentionDays: 90,               // Default: 90 days
		cleanupInterval:  24 * time.Hour,   // Run cleanup daily
	}
}

// SetLogRetention configures log retention cleanup
func (p *Pool) SetLogRetention(days int) {
	p.logRetentionDays = days
}

// SetJobEventCallback sets the callback for job events (used for webhooks)
func (p *Pool) SetJobEventCallback(callback JobEventCallback) {
	p.onJobEvent = callback
}

// SetMetricsCallback sets the callback for job metrics
func (p *Pool) SetMetricsCallback(callback MetricsCallback) {
	p.onMetrics = callback
}

// reportMetrics reports job execution metrics
func (p *Pool) reportMetrics(jobType, status string, duration time.Duration) {
	if p.onMetrics != nil {
		p.onMetrics(jobType, status, duration)
	}
}

// emitJobEvent emits a job event via the callback
func (p *Pool) emitJobEvent(ctx context.Context, eventType domain.WebhookEventType, job *domain.Job, execution *domain.JobExecution) {
	if p.onJobEvent != nil {
		event := domain.WebhookEvent{
			Type:      eventType,
			Timestamp: time.Now().UTC(),
			Job:       job,
			Execution: execution,
		}
		// Run callback in goroutine to not block job execution
		go p.onJobEvent(ctx, event)
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

	// Start log cleanup scheduler if retention is configured
	if p.logRetentionDays > 0 {
		p.wg.Add(1)
		go p.cleanupScheduler(ctx)
		log.Info().Int("retention_days", p.logRetentionDays).Msg("Log retention cleanup enabled")
	}

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
	// Create a cancellable context for this job
	jobCtx, cancel := context.WithCancel(ctx)

	// Mark as running and track the cancel function
	p.runningMutex.Lock()
	p.runningJobs[job.ID] = true
	p.jobContexts[job.ID] = cancel
	p.runningMutex.Unlock()

	defer func() {
		cancel() // Always cancel the context when done
		p.runningMutex.Lock()
		delete(p.runningJobs, job.ID)
		delete(p.jobContexts, job.ID)
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

	// Emit job started event
	p.emitJobEvent(ctx, domain.WebhookEventJobStarted, job, execution)

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

	// Execute job with cancellable context
	result, err := executor.Execute(jobCtx)

	// Check if the job was cancelled
	if jobCtx.Err() == context.Canceled {
		log.Info().Str("job_id", job.ID).Msg("Job was cancelled")
		execution.Status = domain.ExecutionStatusCancelled
		p.completeExecution(ctx, execution.ID, job.ID, domain.ExecutionStatusCancelled, "", "Job cancelled by user", nil, time.Since(startTime))
		// Emit cancelled event
		p.emitJobEvent(ctx, domain.WebhookEventJobCancelled, job, execution)
		// Report metrics
		p.reportMetrics(job.Type, string(domain.ExecutionStatusCancelled), time.Since(startTime))
		// Note: Job status already updated to cancelled by CancelJob
		return
	}

	if err != nil {
		log.Error().Err(err).Str("job_id", job.ID).Msg("Job execution error")

		execution.Status = domain.ExecutionStatusFailed
		execution.Error = fmt.Sprintf("Execution error: %v", err)
		p.completeExecution(ctx, execution.ID, job.ID, domain.ExecutionStatusFailed, "", fmt.Sprintf("Execution error: %v", err), nil, time.Since(startTime))
		// Update job status to failed
		if updateErr := p.repo.UpdateJobStatus(ctx, job.ID, domain.JobStatusFailed); updateErr != nil {
			log.Error().Err(updateErr).Str("job_id", job.ID).Msg("Failed to update job status to failed")
		}
		// Emit failed event
		p.emitJobEvent(ctx, domain.WebhookEventJobFailed, job, execution)
		// Report metrics
		p.reportMetrics(job.Type, string(domain.ExecutionStatusFailed), time.Since(startTime))
		return
	}

	// Determine final status
	finalStatus := domain.ExecutionStatusCompleted
	finalJobStatus := domain.JobStatusCompleted
	webhookEventType := domain.WebhookEventJobCompleted
	if result.ExitCode != 0 {
		finalStatus = domain.ExecutionStatusFailed
		finalJobStatus = domain.JobStatusFailed
		webhookEventType = domain.WebhookEventJobFailed
	}

	// Complete execution
	durationMs := time.Since(startTime).Milliseconds()
	execution.Status = finalStatus
	execution.Output = result.Output
	execution.Error = result.Error
	execution.ExitCode = &result.ExitCode
	p.completeExecution(ctx, execution.ID, job.ID, finalStatus, result.Output, result.Error, &result.ExitCode, time.Since(startTime))

	// Update job status
	if err := p.repo.UpdateJobStatus(ctx, job.ID, finalJobStatus); err != nil {
		log.Error().Err(err).Str("job_id", job.ID).Msg("Failed to update job final status")
	}

	// Emit completion or failure event
	p.emitJobEvent(ctx, webhookEventType, job, execution)

	// Report metrics
	p.reportMetrics(job.Type, string(finalStatus), time.Since(startTime))

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
	p.runningMutex.Lock()
	isRunning := p.runningJobs[jobID]
	cancelFunc := p.jobContexts[jobID]
	p.runningMutex.Unlock()

	// Update job status first
	if err := p.repo.UpdateJobStatus(ctx, jobID, domain.JobStatusCancelled); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	if isRunning && cancelFunc != nil {
		// Cancel the job's context - this will trigger cancellation in the executor
		cancelFunc()
		log.Info().Str("job_id", jobID).Msg("Job cancellation signal sent")
	} else {
		log.Debug().Str("job_id", jobID).Bool("is_running", isRunning).Msg("Job not currently running, status updated")
	}

	return nil
}

// cleanupScheduler periodically cleans up old execution logs
func (p *Pool) cleanupScheduler(ctx context.Context) {
	defer p.wg.Done()

	// Run cleanup immediately on startup
	p.runCleanup(ctx)

	ticker := time.NewTicker(p.cleanupInterval)
	defer ticker.Stop()

	log.Debug().Msg("Cleanup scheduler started")

	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("Cleanup scheduler context cancelled")
			return
		case <-p.stopChan:
			log.Debug().Msg("Cleanup scheduler stopped")
			return
		case <-ticker.C:
			p.runCleanup(ctx)
		}
	}
}

// runCleanup deletes execution logs older than the retention period
func (p *Pool) runCleanup(ctx context.Context) {
	if p.logRetentionDays <= 0 {
		return
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -p.logRetentionDays)

	log.Debug().
		Time("cutoff", cutoff).
		Int("retention_days", p.logRetentionDays).
		Msg("Running execution log cleanup")

	deleted, err := p.repo.DeleteOldExecutions(ctx, cutoff)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup old execution logs")
		return
	}

	if deleted > 0 {
		log.Info().
			Int64("deleted", deleted).
			Int("retention_days", p.logRetentionDays).
			Time("cutoff", cutoff).
			Msg("Cleaned up old execution logs")
	} else {
		log.Debug().Msg("No old execution logs to cleanup")
	}
}
