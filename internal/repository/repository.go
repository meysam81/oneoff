package repository

import (
	"context"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
)

// Repository defines the interface for data persistence
type Repository interface {
	// Job operations
	CreateJob(ctx context.Context, job *domain.Job, tagIDs []string) error
	GetJob(ctx context.Context, id string) (*domain.Job, error)
	ListJobs(ctx context.Context, filter domain.JobFilter) ([]*domain.Job, error)
	UpdateJob(ctx context.Context, id string, updates domain.UpdateJobRequest) error
	DeleteJob(ctx context.Context, id string) error
	CountJobs(ctx context.Context, filter domain.JobFilter) (int64, error)

	// Job execution operations
	CreateExecution(ctx context.Context, execution *domain.JobExecution) error
	GetExecution(ctx context.Context, id string) (*domain.JobExecution, error)
	ListExecutions(ctx context.Context, filter domain.ExecutionFilter) ([]*domain.JobExecution, error)
	UpdateExecution(ctx context.Context, id string, status domain.ExecutionStatus, output, error string, exitCode *int) error
	CompleteExecution(ctx context.Context, id string, status domain.ExecutionStatus, output, error string, exitCode *int, durationMs int64) error
	DeleteOldExecutions(ctx context.Context, before time.Time) (int64, error)

	// Project operations
	CreateProject(ctx context.Context, project *domain.Project) error
	GetProject(ctx context.Context, id string) (*domain.Project, error)
	ListProjects(ctx context.Context, includeArchived bool) ([]*domain.Project, error)
	UpdateProject(ctx context.Context, id string, name, description, color, icon *string, isArchived *bool) error
	DeleteProject(ctx context.Context, id string) error

	// Tag operations
	CreateTag(ctx context.Context, tag *domain.Tag) error
	GetTag(ctx context.Context, id string) (*domain.Tag, error)
	GetTagByName(ctx context.Context, name string) (*domain.Tag, error)
	ListTags(ctx context.Context) ([]*domain.Tag, error)
	UpdateTag(ctx context.Context, id string, name, color *string, isDefault *bool) error
	DeleteTag(ctx context.Context, id string) error

	// Job-Tag operations
	AddJobTags(ctx context.Context, jobID string, tagIDs []string) error
	RemoveJobTags(ctx context.Context, jobID string, tagIDs []string) error
	GetJobTags(ctx context.Context, jobID string) ([]domain.Tag, error)

	// Chain operations
	CreateChain(ctx context.Context, chain *domain.JobChain) error
	GetChain(ctx context.Context, id string) (*domain.JobChain, error)
	ListChains(ctx context.Context, projectID string) ([]*domain.JobChain, error)
	DeleteChain(ctx context.Context, id string) error

	// System config operations
	GetConfig(ctx context.Context, key string) (*domain.SystemConfig, error)
	SetConfig(ctx context.Context, key, value string) error
	ListConfig(ctx context.Context) ([]*domain.SystemConfig, error)

	// Stats operations
	GetSystemStats(ctx context.Context) (*domain.SystemStats, error)

	// Scheduler operations
	GetScheduledJobs(ctx context.Context, before time.Time, limit int) ([]*domain.Job, error)
	UpdateJobStatus(ctx context.Context, id string, status domain.JobStatus) error

	// API Key operations
	CreateAPIKey(ctx context.Context, key *domain.APIKey) error
	GetAPIKey(ctx context.Context, id string) (*domain.APIKey, error)
	GetAPIKeyByHash(ctx context.Context, keyHash string) (*domain.APIKey, error)
	ListAPIKeys(ctx context.Context) ([]*domain.APIKey, error)
	UpdateAPIKey(ctx context.Context, id string, updates domain.UpdateAPIKeyRequest) error
	UpdateAPIKeyLastUsed(ctx context.Context, id string) error
	DeleteAPIKey(ctx context.Context, id string) error

	// Webhook operations
	CreateWebhook(ctx context.Context, webhook *domain.Webhook) error
	GetWebhook(ctx context.Context, id string) (*domain.Webhook, error)
	ListWebhooks(ctx context.Context, filter domain.WebhookFilter) ([]*domain.Webhook, error)
	GetActiveWebhooksForEvent(ctx context.Context, event domain.WebhookEventType) ([]*domain.Webhook, error)
	UpdateWebhook(ctx context.Context, id string, updates domain.UpdateWebhookRequest) error
	DeleteWebhook(ctx context.Context, id string) error

	// Webhook delivery operations
	CreateWebhookDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error
	GetWebhookDelivery(ctx context.Context, id string) (*domain.WebhookDelivery, error)
	ListWebhookDeliveries(ctx context.Context, filter domain.WebhookDeliveryFilter) ([]*domain.WebhookDelivery, error)
	UpdateWebhookDelivery(ctx context.Context, id string, status domain.WebhookDeliveryStatus, responseCode *int, responseBody, errMsg string, nextRetry *time.Time) error
	GetPendingWebhookDeliveries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error)
	IncrementDeliveryAttempts(ctx context.Context, id string) error

	// Transaction support
	WithTransaction(ctx context.Context, fn func(Repository) error) error

	// Close closes the database connection
	Close() error
}
