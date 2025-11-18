package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteRepository implements Repository using SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=on&_journal_mode=WAL", dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &SQLiteRepository{db: db}, nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// CreateJob creates a new job
func (r *SQLiteRepository) CreateJob(ctx context.Context, job *domain.Job, tagIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	query := `
		INSERT INTO jobs (name, type, config, scheduled_at, priority, project_id, timezone, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRowContext(ctx, query,
		job.Name, job.Type, job.Config, job.ScheduledAt.UTC(),
		job.Priority, job.ProjectID, job.Timezone, job.Status,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	// Add tags
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			_, err := tx.ExecContext(ctx, "INSERT INTO job_tags (job_id, tag_id) VALUES (?, ?)", job.ID, tagID)
			if err != nil {
				return fmt.Errorf("failed to add tag: %w", err)
			}
		}
	}

	return tx.Commit()
}

// GetJob retrieves a job by ID
func (r *SQLiteRepository) GetJob(ctx context.Context, id string) (*domain.Job, error) {
	query := `
		SELECT id, name, type, config, scheduled_at, priority, project_id, timezone, status, created_at, updated_at
		FROM jobs
		WHERE id = ?
	`

	job := &domain.Job{}
	var scheduledAt, createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &job.Name, &job.Type, &job.Config, &scheduledAt,
		&job.Priority, &job.ProjectID, &job.Timezone, &job.Status,
		&createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrJobNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	// Parse timestamps
	job.ScheduledAt, _ = time.Parse("2006-01-02 15:04:05", scheduledAt)
	job.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	job.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	// Load tags
	tags, err := r.GetJobTags(ctx, id)
	if err != nil {
		return nil, err
	}
	job.Tags = tags

	return job, nil
}

// ListJobs retrieves jobs based on filter
func (r *SQLiteRepository) ListJobs(ctx context.Context, filter domain.JobFilter) ([]*domain.Job, error) {
	query := `SELECT id, name, type, config, scheduled_at, priority, project_id, timezone, status, created_at, updated_at FROM jobs WHERE 1=1`
	args := []interface{}{}

	if filter.ProjectID != "" {
		query += " AND project_id = ?"
		args = append(args, filter.ProjectID)
	}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}

	if filter.JobType != "" {
		query += " AND type = ?"
		args = append(args, filter.JobType)
	}

	if filter.Search != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+filter.Search+"%")
	}

	if filter.TimeFrom != nil {
		query += " AND scheduled_at >= ?"
		args = append(args, filter.TimeFrom.UTC())
	}

	if filter.TimeTo != nil {
		query += " AND scheduled_at <= ?"
		args = append(args, filter.TimeTo.UTC())
	}

	// Handle tag filtering
	if len(filter.TagIDs) > 0 {
		placeholders := make([]string, len(filter.TagIDs))
		for i, tagID := range filter.TagIDs {
			placeholders[i] = "?"
			args = append(args, tagID)
		}
		query += fmt.Sprintf(" AND id IN (SELECT job_id FROM job_tags WHERE tag_id IN (%s))", strings.Join(placeholders, ","))
	}

	// Sorting
	sortBy := "scheduled_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "ASC"
	if filter.SortOrder == "desc" {
		sortOrder = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Pagination
	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var jobs []*domain.Job
	for rows.Next() {
		job := &domain.Job{}
		var scheduledAt, createdAt, updatedAt string

		err := rows.Scan(
			&job.ID, &job.Name, &job.Type, &job.Config, &scheduledAt,
			&job.Priority, &job.ProjectID, &job.Timezone, &job.Status,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}

		job.ScheduledAt, _ = time.Parse("2006-01-02 15:04:05", scheduledAt)
		job.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		job.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		// Load tags
		tags, err := r.GetJobTags(ctx, job.ID)
		if err != nil {
			return nil, err
		}
		job.Tags = tags

		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

// UpdateJob updates a job
func (r *SQLiteRepository) UpdateJob(ctx context.Context, id string, updates domain.UpdateJobRequest) error {
	sets := []string{}
	args := []interface{}{}

	if updates.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *updates.Name)
	}
	if updates.Config != nil {
		sets = append(sets, "config = ?")
		args = append(args, *updates.Config)
	}
	if updates.ScheduledAt != nil {
		t, err := time.Parse(time.RFC3339, *updates.ScheduledAt)
		if err != nil {
			return fmt.Errorf("invalid scheduled_at format: %w", err)
		}
		sets = append(sets, "scheduled_at = ?")
		args = append(args, t.UTC())
	}
	if updates.Priority != nil {
		sets = append(sets, "priority = ?")
		args = append(args, *updates.Priority)
	}
	if updates.ProjectID != nil {
		sets = append(sets, "project_id = ?")
		args = append(args, *updates.ProjectID)
	}
	if updates.Timezone != nil {
		sets = append(sets, "timezone = ?")
		args = append(args, *updates.Timezone)
	}
	if updates.Status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *updates.Status)
	}

	if len(sets) == 0 && len(updates.TagIDs) == 0 {
		return nil // Nothing to update
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if len(sets) > 0 {
		query := fmt.Sprintf("UPDATE jobs SET %s WHERE id = ?", strings.Join(sets, ", "))
		args = append(args, id)

		result, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to update job: %w", err)
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			return domain.ErrJobNotFound
		}
	}

	// Update tags if provided
	if updates.TagIDs != nil {
		// Remove existing tags
		_, err := tx.ExecContext(ctx, "DELETE FROM job_tags WHERE job_id = ?", id)
		if err != nil {
			return fmt.Errorf("failed to remove old tags: %w", err)
		}

		// Add new tags
		for _, tagID := range updates.TagIDs {
			_, err := tx.ExecContext(ctx, "INSERT INTO job_tags (job_id, tag_id) VALUES (?, ?)", id, tagID)
			if err != nil {
				return fmt.Errorf("failed to add tag: %w", err)
			}
		}
	}

	return tx.Commit()
}

// DeleteJob deletes a job
func (r *SQLiteRepository) DeleteJob(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM jobs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrJobNotFound
	}

	return nil
}

// CountJobs counts jobs matching the filter
func (r *SQLiteRepository) CountJobs(ctx context.Context, filter domain.JobFilter) (int64, error) {
	query := "SELECT COUNT(*) FROM jobs WHERE 1=1"
	args := []interface{}{}

	if filter.ProjectID != "" {
		query += " AND project_id = ?"
		args = append(args, filter.ProjectID)
	}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

// Continue with execution methods...
