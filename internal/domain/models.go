package domain

import "time"

// JobStatus represents the current status of a job
type JobStatus string

const (
	JobStatusScheduled JobStatus = "scheduled"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// ExecutionStatus represents the status of a job execution
type ExecutionStatus string

const (
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
)

// Job represents a scheduled one-time job
type Job struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Config      string    `json:"config"` // JSON string for job-specific config
	ScheduledAt time.Time `json:"scheduled_at"`
	Priority    int       `json:"priority"` // 1-10
	ProjectID   string    `json:"project_id"`
	Timezone    string    `json:"timezone"`
	Status      JobStatus `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []Tag     `json:"tags,omitempty"`
}

// JobExecution represents an execution instance of a job
type JobExecution struct {
	ID          string          `json:"id"`
	JobID       string          `json:"job_id"`
	StartedAt   time.Time       `json:"started_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Status      ExecutionStatus `json:"status"`
	Output      string          `json:"output,omitempty"`
	ExitCode    *int            `json:"exit_code,omitempty"`
	Error       string          `json:"error,omitempty"`
	DurationMs  *int64          `json:"duration_ms,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// Project represents a project for organizing jobs
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	IsArchived  bool      `json:"is_archived"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag represents a tag for categorizing jobs
type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color,omitempty"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

// JobChain represents a sequence of jobs to execute
type JobChain struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	ProjectID string         `json:"project_id"`
	CreatedAt time.Time      `json:"created_at"`
	Links     []JobChainLink `json:"links,omitempty"`
}

// JobChainLink represents a link in a job chain
type JobChainLink struct {
	ID            string    `json:"id"`
	ChainID       string    `json:"chain_id"`
	JobID         string    `json:"job_id"`
	SequenceOrder int       `json:"sequence_order"`
	StopOnFailure bool      `json:"stop_on_failure"`
	CreatedAt     time.Time `json:"created_at"`
}

// SystemConfig represents system-wide configuration
type SystemConfig struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"` // JSON string
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateJobRequest represents a request to create a new job
type CreateJobRequest struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Config      string   `json:"config"`
	ScheduledAt string   `json:"scheduled_at,omitempty"`
	Immediate   bool     `json:"immediate,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	ProjectID   string   `json:"project_id,omitempty"`
	Timezone    string   `json:"timezone,omitempty"`
	TagIDs      []string `json:"tag_ids,omitempty"`
}

// UpdateJobRequest represents a request to update a job
type UpdateJobRequest struct {
	Name        *string  `json:"name,omitempty"`
	Config      *string  `json:"config,omitempty"`
	ScheduledAt *string  `json:"scheduled_at,omitempty"`
	Priority    *int     `json:"priority,omitempty"`
	ProjectID   *string  `json:"project_id,omitempty"`
	Timezone    *string  `json:"timezone,omitempty"`
	Status      *string  `json:"status,omitempty"`
	TagIDs      []string `json:"tag_ids,omitempty"`
}

// JobFilter represents filters for querying jobs
type JobFilter struct {
	ProjectID string
	Status    JobStatus
	TagIDs    []string
	JobType   string
	Search    string
	TimeFrom  *time.Time
	TimeTo    *time.Time
	Limit     int
	Offset    int
	SortBy    string // "scheduled_at", "priority", "created_at"
	SortOrder string // "asc", "desc"
}

// ExecutionFilter represents filters for querying executions
type ExecutionFilter struct {
	JobID     string
	Status    ExecutionStatus
	ProjectID string
	TagIDs    []string
	DateFrom  *time.Time
	DateTo    *time.Time
	Limit     int
	Offset    int
	SortBy    string // "started_at", "duration_ms"
	SortOrder string // "asc", "desc"
}

// SystemStats represents system-wide statistics
type SystemStats struct {
	TotalScheduled   int64   `json:"total_scheduled"`
	CurrentlyRunning int64   `json:"currently_running"`
	CompletedToday   int64   `json:"completed_today"`
	FailedRecent     int64   `json:"failed_recent"`
	AvgDurationMs    float64 `json:"avg_duration_ms"`
	QueueDepth       int64   `json:"queue_depth"`
}

// WorkerStatus represents the status of workers
type WorkerStatus struct {
	TotalWorkers     int      `json:"total_workers"`
	ActiveWorkers    int      `json:"active_workers"`
	AvailableWorkers int      `json:"available_workers"`
	QueuedJobs       int      `json:"queued_jobs"`
	RunningJobs      []string `json:"running_jobs"` // Job IDs
}

// CreateChainRequest represents a request to create a new job chain
type CreateChainRequest struct {
	Name      string           `json:"name"`
	ProjectID string           `json:"project_id,omitempty"`
	Links     []ChainLinkInput `json:"links"`
}

// ChainLinkInput represents input for a chain link
type ChainLinkInput struct {
	JobID         string `json:"job_id"`
	StopOnFailure bool   `json:"stop_on_failure"`
}

// UpdateChainRequest represents a request to update a chain
type UpdateChainRequest struct {
	Name  *string          `json:"name,omitempty"`
	Links []ChainLinkInput `json:"links,omitempty"`
}

// ChainFilter represents filters for querying chains
type ChainFilter struct {
	ProjectID string
	Limit     int
	Offset    int
}
