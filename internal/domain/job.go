package domain

import (
	"context"
	"encoding/json"
)

// JobExecutor defines the interface that all job types must implement
type JobExecutor interface {
	// Execute runs the job and returns the result
	Execute(ctx context.Context) (*ExecutionResult, error)

	// Validate validates the job configuration
	Validate() error

	// Type returns the job type identifier
	Type() string

	// Description returns a human-readable description of the job
	Description() string
}

// ExecutionResult represents the result of a job execution
type ExecutionResult struct {
	Output   string
	ExitCode int
	Error    string
}

// JobFactory is a function that creates a JobExecutor from a config
type JobFactory func(config string) (JobExecutor, error)

// JobRegistry manages registered job types
type JobRegistry struct {
	factories map[string]JobFactory
}

// NewJobRegistry creates a new job registry
func NewJobRegistry() *JobRegistry {
	return &JobRegistry{
		factories: make(map[string]JobFactory),
	}
}

// Register registers a new job type
func (r *JobRegistry) Register(jobType string, factory JobFactory) {
	r.factories[jobType] = factory
}

// Create creates a JobExecutor for the given job type and config
func (r *JobRegistry) Create(jobType string, config string) (JobExecutor, error) {
	factory, exists := r.factories[jobType]
	if !exists {
		return nil, ErrJobTypeNotFound
	}
	return factory(config)
}

// ListTypes returns all registered job types
func (r *JobRegistry) ListTypes() []string {
	types := make([]string, 0, len(r.factories))
	for jobType := range r.factories {
		types = append(types, jobType)
	}
	return types
}

// HTTPJobConfig represents configuration for HTTP request jobs
type HTTPJobConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Timeout int               `json:"timeout,omitempty"` // seconds
}

// ParseHTTPJobConfig parses HTTP job configuration from JSON
func ParseHTTPJobConfig(config string) (*HTTPJobConfig, error) {
	var cfg HTTPJobConfig
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ShellJobConfig represents configuration for shell script jobs
type ShellJobConfig struct {
	Script  string   `json:"script"`      // Script content or path
	Args    []string `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	WorkDir string   `json:"workdir,omitempty"`
	Timeout int      `json:"timeout,omitempty"` // seconds
	IsPath  bool     `json:"is_path,omitempty"` // true if script is a file path
}

// ParseShellJobConfig parses shell job configuration from JSON
func ParseShellJobConfig(config string) (*ShellJobConfig, error) {
	var cfg ShellJobConfig
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// DockerJobConfig represents configuration for Docker container jobs
type DockerJobConfig struct {
	Image      string            `json:"image"`
	Command    []string          `json:"command,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	Volumes    map[string]string `json:"volumes,omitempty"` // host:container
	WorkDir    string            `json:"workdir,omitempty"`
	AutoRemove bool              `json:"auto_remove"`
	Timeout    int               `json:"timeout,omitempty"` // seconds
}

// ParseDockerJobConfig parses Docker job configuration from JSON
func ParseDockerJobConfig(config string) (*DockerJobConfig, error) {
	var cfg DockerJobConfig
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
