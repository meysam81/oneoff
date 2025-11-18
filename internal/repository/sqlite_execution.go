package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
)

// CreateExecution creates a new job execution
func (r *SQLiteRepository) CreateExecution(ctx context.Context, execution *domain.JobExecution) error {
	query := `
		INSERT INTO job_executions (job_id, started_at, status)
		VALUES (?, ?, ?)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(ctx, query,
		execution.JobID,
		execution.StartedAt.UTC(),
		execution.Status,
	).Scan(&execution.ID, &execution.CreatedAt)
}

// GetExecution retrieves an execution by ID
func (r *SQLiteRepository) GetExecution(ctx context.Context, id string) (*domain.JobExecution, error) {
	query := `
		SELECT id, job_id, started_at, completed_at, status, output, exit_code, error, duration_ms, created_at
		FROM job_executions
		WHERE id = ?
	`

	execution := &domain.JobExecution{}
	var startedAt, createdAt string
	var completedAt sql.NullString
	var output, error sql.NullString
	var exitCode sql.NullInt64
	var durationMs sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&execution.ID,
		&execution.JobID,
		&startedAt,
		&completedAt,
		&execution.Status,
		&output,
		&exitCode,
		&error,
		&durationMs,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrExecutionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}

	execution.StartedAt, _ = time.Parse("2006-01-02 15:04:05", startedAt)
	execution.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

	if completedAt.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", completedAt.String)
		execution.CompletedAt = &t
	}
	if output.Valid {
		execution.Output = output.String
	}
	if error.Valid {
		execution.Error = error.String
	}
	if exitCode.Valid {
		code := int(exitCode.Int64)
		execution.ExitCode = &code
	}
	if durationMs.Valid {
		duration := durationMs.Int64
		execution.DurationMs = &duration
	}

	return execution, nil
}

// ListExecutions retrieves executions based on filter
func (r *SQLiteRepository) ListExecutions(ctx context.Context, filter domain.ExecutionFilter) ([]*domain.JobExecution, error) {
	query := `
		SELECT e.id, e.job_id, e.started_at, e.completed_at, e.status, e.output, e.exit_code, e.error, e.duration_ms, e.created_at
		FROM job_executions e
		WHERE 1=1
	`
	args := []interface{}{}

	if filter.JobID != "" {
		query += " AND e.job_id = ?"
		args = append(args, filter.JobID)
	}

	if filter.Status != "" {
		query += " AND e.status = ?"
		args = append(args, filter.Status)
	}

	if filter.ProjectID != "" {
		query += " AND e.job_id IN (SELECT id FROM jobs WHERE project_id = ?)"
		args = append(args, filter.ProjectID)
	}

	if filter.DateFrom != nil {
		query += " AND e.started_at >= ?"
		args = append(args, filter.DateFrom.UTC())
	}

	if filter.DateTo != nil {
		query += " AND e.started_at <= ?"
		args = append(args, filter.DateTo.UTC())
	}

	// Sorting
	sortBy := "started_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY e.%s %s", sortBy, sortOrder)

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
		return nil, fmt.Errorf("failed to list executions: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var executions []*domain.JobExecution
	for rows.Next() {
		execution := &domain.JobExecution{}
		var startedAt, createdAt string
		var completedAt sql.NullString
		var output, error sql.NullString
		var exitCode sql.NullInt64
		var durationMs sql.NullInt64

		err := rows.Scan(
			&execution.ID,
			&execution.JobID,
			&startedAt,
			&completedAt,
			&execution.Status,
			&output,
			&exitCode,
			&error,
			&durationMs,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan execution: %w", err)
		}

		execution.StartedAt, _ = time.Parse("2006-01-02 15:04:05", startedAt)
		execution.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

		if completedAt.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", completedAt.String)
			execution.CompletedAt = &t
		}
		if output.Valid {
			execution.Output = output.String
		}
		if error.Valid {
			execution.Error = error.String
		}
		if exitCode.Valid {
			code := int(exitCode.Int64)
			execution.ExitCode = &code
		}
		if durationMs.Valid {
			duration := durationMs.Int64
			execution.DurationMs = &duration
		}

		executions = append(executions, execution)
	}

	return executions, rows.Err()
}

// UpdateExecution updates an execution's status and output
func (r *SQLiteRepository) UpdateExecution(ctx context.Context, id string, status domain.ExecutionStatus, output, errorMsg string, exitCode *int) error {
	query := `UPDATE job_executions SET status = ?, output = ?, error = ?, exit_code = ? WHERE id = ?`

	var exitCodeVal sql.NullInt64
	if exitCode != nil {
		exitCodeVal = sql.NullInt64{Int64: int64(*exitCode), Valid: true}
	}

	result, err := r.db.ExecContext(ctx, query, status, output, errorMsg, exitCodeVal, id)
	if err != nil {
		return fmt.Errorf("failed to update execution: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrExecutionNotFound
	}

	return nil
}

// CompleteExecution marks an execution as complete
func (r *SQLiteRepository) CompleteExecution(ctx context.Context, id string, status domain.ExecutionStatus, output, errorMsg string, exitCode *int, durationMs int64) error {
	query := `
		UPDATE job_executions
		SET status = ?, output = ?, error = ?, exit_code = ?, completed_at = datetime('now', 'utc'), duration_ms = ?
		WHERE id = ?
	`

	var exitCodeVal sql.NullInt64
	if exitCode != nil {
		exitCodeVal = sql.NullInt64{Int64: int64(*exitCode), Valid: true}
	}

	result, err := r.db.ExecContext(ctx, query, status, output, errorMsg, exitCodeVal, durationMs, id)
	if err != nil {
		return fmt.Errorf("failed to complete execution: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrExecutionNotFound
	}

	return nil
}

// DeleteOldExecutions deletes executions older than the specified date
func (r *SQLiteRepository) DeleteOldExecutions(ctx context.Context, before time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx, "DELETE FROM job_executions WHERE created_at < ?", before.UTC())
	if err != nil {
		return 0, fmt.Errorf("failed to delete old executions: %w", err)
	}

	rows, _ := result.RowsAffected()
	return rows, nil
}
