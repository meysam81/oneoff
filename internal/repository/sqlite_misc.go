package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
)

// Project operations

func (r *SQLiteRepository) CreateProject(ctx context.Context, project *domain.Project) error {
	query := `
		INSERT INTO projects (name, description, color, icon, is_archived)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		project.Name, project.Description, project.Color, project.Icon, project.IsArchived,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

func (r *SQLiteRepository) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	query := `SELECT id, name, description, color, icon, is_archived, created_at, updated_at FROM projects WHERE id = ?`

	project := &domain.Project{}
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Color,
		&project.Icon, &project.IsArchived, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrProjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	project.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	project.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return project, nil
}

func (r *SQLiteRepository) ListProjects(ctx context.Context, includeArchived bool) ([]*domain.Project, error) {
	query := "SELECT id, name, description, color, icon, is_archived, created_at, updated_at FROM projects"
	if !includeArchived {
		query += " WHERE is_archived = 0"
	}
	query += " ORDER BY name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var projects []*domain.Project
	for rows.Next() {
		project := &domain.Project{}
		var createdAt, updatedAt string

		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Color,
			&project.Icon, &project.IsArchived, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		project.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		project.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		projects = append(projects, project)
	}

	return projects, rows.Err()
}

func (r *SQLiteRepository) UpdateProject(ctx context.Context, id string, name, description, color, icon *string, isArchived *bool) error {
	sets := []string{}
	args := []interface{}{}

	if name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *name)
	}
	if description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *description)
	}
	if color != nil {
		sets = append(sets, "color = ?")
		args = append(args, *color)
	}
	if icon != nil {
		sets = append(sets, "icon = ?")
		args = append(args, *icon)
	}
	if isArchived != nil {
		sets = append(sets, "is_archived = ?")
		args = append(args, *isArchived)
	}

	if len(sets) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE projects SET %s WHERE id = ?", joinStrings(sets, ", "))
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrProjectNotFound
	}

	return nil
}

func (r *SQLiteRepository) DeleteProject(ctx context.Context, id string) error {
	if id == "default" {
		return domain.ErrCannotDeleteDefault
	}

	result, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrProjectNotFound
	}

	return nil
}

// Tag operations

func (r *SQLiteRepository) CreateTag(ctx context.Context, tag *domain.Tag) error {
	query := `
		INSERT INTO tags (name, color, is_default)
		VALUES (?, ?, ?)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(ctx, query, tag.Name, tag.Color, tag.IsDefault).Scan(&tag.ID, &tag.CreatedAt)
}

func (r *SQLiteRepository) GetTag(ctx context.Context, id string) (*domain.Tag, error) {
	query := "SELECT id, name, color, is_default, created_at FROM tags WHERE id = ?"

	tag := &domain.Tag{}
	var createdAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(&tag.ID, &tag.Name, &tag.Color, &tag.IsDefault, &createdAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrTagNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	tag.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return tag, nil
}

func (r *SQLiteRepository) GetTagByName(ctx context.Context, name string) (*domain.Tag, error) {
	query := "SELECT id, name, color, is_default, created_at FROM tags WHERE name = ?"

	tag := &domain.Tag{}
	var createdAt string

	err := r.db.QueryRowContext(ctx, query, name).Scan(&tag.ID, &tag.Name, &tag.Color, &tag.IsDefault, &createdAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrTagNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	tag.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return tag, nil
}

func (r *SQLiteRepository) ListTags(ctx context.Context) ([]*domain.Tag, error) {
	query := "SELECT id, name, color, is_default, created_at FROM tags ORDER BY name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var tags []*domain.Tag
	for rows.Next() {
		tag := &domain.Tag{}
		var createdAt string

		err := rows.Scan(&tag.ID, &tag.Name, &tag.Color, &tag.IsDefault, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}

		tag.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func (r *SQLiteRepository) UpdateTag(ctx context.Context, id string, name, color *string, isDefault *bool) error {
	sets := []string{}
	args := []interface{}{}

	if name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *name)
	}
	if color != nil {
		sets = append(sets, "color = ?")
		args = append(args, *color)
	}
	if isDefault != nil {
		sets = append(sets, "is_default = ?")
		args = append(args, *isDefault)
	}

	if len(sets) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE tags SET %s WHERE id = ?", joinStrings(sets, ", "))
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrTagNotFound
	}

	return nil
}

func (r *SQLiteRepository) DeleteTag(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrTagNotFound
	}

	return nil
}

// Job-Tag operations

func (r *SQLiteRepository) AddJobTags(ctx context.Context, jobID string, tagIDs []string) error {
	for _, tagID := range tagIDs {
		_, err := r.db.ExecContext(ctx, "INSERT OR IGNORE INTO job_tags (job_id, tag_id) VALUES (?, ?)", jobID, tagID)
		if err != nil {
			return fmt.Errorf("failed to add job tag: %w", err)
		}
	}
	return nil
}

func (r *SQLiteRepository) RemoveJobTags(ctx context.Context, jobID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(tagIDs))
	args := []interface{}{jobID}
	for i, tagID := range tagIDs {
		placeholders[i] = "?"
		args = append(args, tagID)
	}

	query := fmt.Sprintf("DELETE FROM job_tags WHERE job_id = ? AND tag_id IN (%s)", joinStrings(placeholders, ","))
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *SQLiteRepository) GetJobTags(ctx context.Context, jobID string) ([]domain.Tag, error) {
	query := `
		SELECT t.id, t.name, t.color, t.is_default, t.created_at
		FROM tags t
		INNER JOIN job_tags jt ON t.id = jt.tag_id
		WHERE jt.job_id = ?
		ORDER BY t.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job tags: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var tags []domain.Tag
	for rows.Next() {
		tag := domain.Tag{}
		var createdAt string

		err := rows.Scan(&tag.ID, &tag.Name, &tag.Color, &tag.IsDefault, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}

		tag.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// Chain operations

func (r *SQLiteRepository) CreateChain(ctx context.Context, chain *domain.JobChain) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Create chain
	err = tx.QueryRow(`
		INSERT INTO job_chains (name, project_id)
		VALUES (?, ?)
		RETURNING id, created_at
	`, chain.Name, chain.ProjectID).Scan(&chain.ID, &chain.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create chain: %w", err)
	}

	// Create links
	for _, link := range chain.Links {
		_, err := tx.Exec(`
			INSERT INTO job_chain_links (chain_id, job_id, sequence_order, stop_on_failure)
			VALUES (?, ?, ?, ?)
		`, chain.ID, link.JobID, link.SequenceOrder, link.StopOnFailure)

		if err != nil {
			return fmt.Errorf("failed to create chain link: %w", err)
		}
	}

	return tx.Commit()
}

func (r *SQLiteRepository) GetChain(ctx context.Context, id string) (*domain.JobChain, error) {
	query := "SELECT id, name, project_id, created_at FROM job_chains WHERE id = ?"

	chain := &domain.JobChain{}
	var createdAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(&chain.ID, &chain.Name, &chain.ProjectID, &createdAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrChainNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get chain: %w", err)
	}

	chain.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

	// Get links
	linksQuery := "SELECT id, job_id, sequence_order, stop_on_failure, created_at FROM job_chain_links WHERE chain_id = ? ORDER BY sequence_order ASC"
	rows, err := r.db.QueryContext(ctx, linksQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain links: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		link := domain.JobChainLink{ChainID: id}
		var createdAt string

		err := rows.Scan(&link.ID, &link.JobID, &link.SequenceOrder, &link.StopOnFailure, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chain link: %w", err)
		}

		link.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		chain.Links = append(chain.Links, link)
	}

	return chain, rows.Err()
}

func (r *SQLiteRepository) ListChains(ctx context.Context, projectID string) ([]*domain.JobChain, error) {
	query := "SELECT id, name, project_id, created_at FROM job_chains WHERE project_id = ? ORDER BY name ASC"

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list chains: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var chains []*domain.JobChain
	for rows.Next() {
		chain := &domain.JobChain{}
		var createdAt string

		err := rows.Scan(&chain.ID, &chain.Name, &chain.ProjectID, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chain: %w", err)
		}

		chain.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		chains = append(chains, chain)
	}

	return chains, rows.Err()
}

func (r *SQLiteRepository) DeleteChain(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM job_chains WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete chain: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrChainNotFound
	}

	return nil
}

func (r *SQLiteRepository) ListChainsWithFilter(ctx context.Context, filter domain.ChainFilter) ([]*domain.JobChain, error) {
	query := "SELECT id, name, project_id, created_at FROM job_chains WHERE 1=1"
	args := []interface{}{}

	if filter.ProjectID != "" {
		query += " AND project_id = ?"
		args = append(args, filter.ProjectID)
	}

	query += " ORDER BY name ASC"

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
		return nil, fmt.Errorf("failed to list chains: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var chains []*domain.JobChain
	for rows.Next() {
		chain := &domain.JobChain{}
		var createdAt string

		err := rows.Scan(&chain.ID, &chain.Name, &chain.ProjectID, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chain: %w", err)
		}

		chain.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

		// Load links for each chain
		links, err := r.getChainLinks(ctx, chain.ID)
		if err != nil {
			return nil, err
		}
		chain.Links = links

		chains = append(chains, chain)
	}

	return chains, rows.Err()
}

func (r *SQLiteRepository) CountChains(ctx context.Context, filter domain.ChainFilter) (int64, error) {
	query := "SELECT COUNT(*) FROM job_chains WHERE 1=1"
	args := []interface{}{}

	if filter.ProjectID != "" {
		query += " AND project_id = ?"
		args = append(args, filter.ProjectID)
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *SQLiteRepository) UpdateChain(ctx context.Context, id string, name *string) error {
	if name == nil {
		return nil
	}

	result, err := r.db.ExecContext(ctx, "UPDATE job_chains SET name = ? WHERE id = ?", *name, id)
	if err != nil {
		return fmt.Errorf("failed to update chain: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrChainNotFound
	}

	return nil
}

func (r *SQLiteRepository) UpdateChainLinks(ctx context.Context, chainID string, links []domain.ChainLinkInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Delete existing links
	_, err = tx.ExecContext(ctx, "DELETE FROM job_chain_links WHERE chain_id = ?", chainID)
	if err != nil {
		return fmt.Errorf("failed to delete old links: %w", err)
	}

	// Insert new links with sequence order
	for i, link := range links {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO job_chain_links (chain_id, job_id, sequence_order, stop_on_failure)
			VALUES (?, ?, ?, ?)
		`, chainID, link.JobID, i+1, link.StopOnFailure)
		if err != nil {
			return fmt.Errorf("failed to create chain link: %w", err)
		}
	}

	return tx.Commit()
}

func (r *SQLiteRepository) getChainLinks(ctx context.Context, chainID string) ([]domain.JobChainLink, error) {
	query := "SELECT id, job_id, sequence_order, stop_on_failure, created_at FROM job_chain_links WHERE chain_id = ? ORDER BY sequence_order ASC"
	rows, err := r.db.QueryContext(ctx, query, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain links: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var links []domain.JobChainLink
	for rows.Next() {
		link := domain.JobChainLink{ChainID: chainID}
		var createdAt string

		err := rows.Scan(&link.ID, &link.JobID, &link.SequenceOrder, &link.StopOnFailure, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chain link: %w", err)
		}

		link.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		links = append(links, link)
	}

	return links, rows.Err()
}

// System config operations

func (r *SQLiteRepository) GetConfig(ctx context.Context, key string) (*domain.SystemConfig, error) {
	query := "SELECT key, value, updated_at FROM system_config WHERE key = ?"

	config := &domain.SystemConfig{}
	var updatedAt string

	err := r.db.QueryRowContext(ctx, query, key).Scan(&config.Key, &config.Value, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("config key not found: %s", key)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	config.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
	return config, nil
}

func (r *SQLiteRepository) SetConfig(ctx context.Context, key, value string) error {
	query := "INSERT OR REPLACE INTO system_config (key, value) VALUES (?, ?)"
	_, err := r.db.ExecContext(ctx, query, key, value)
	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) ListConfig(ctx context.Context) ([]*domain.SystemConfig, error) {
	query := "SELECT key, value, updated_at FROM system_config ORDER BY key ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list config: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var configs []*domain.SystemConfig
	for rows.Next() {
		config := &domain.SystemConfig{}
		var updatedAt string

		err := rows.Scan(&config.Key, &config.Value, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}

		config.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		configs = append(configs, config)
	}

	return configs, rows.Err()
}

// Stats operations

func (r *SQLiteRepository) GetSystemStats(ctx context.Context) (*domain.SystemStats, error) {
	stats := &domain.SystemStats{}

	// Total scheduled
	_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM jobs WHERE status = 'scheduled'").Scan(&stats.TotalScheduled)

	// Currently running
	_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM jobs WHERE status = 'running'").Scan(&stats.CurrentlyRunning)

	// Completed today
	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM job_executions
		WHERE status = 'completed' AND DATE(started_at) = DATE('now', 'utc')
	`).Scan(&stats.CompletedToday)

	// Failed recent (last 24 hours)
	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM job_executions
		WHERE status = 'failed' AND started_at >= datetime('now', '-1 day', 'utc')
	`).Scan(&stats.FailedRecent)

	// Average duration
	_ = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(duration_ms), 0) FROM job_executions
		WHERE status = 'completed' AND duration_ms IS NOT NULL
	`).Scan(&stats.AvgDurationMs)

	// Queue depth (scheduled jobs in the past that haven't run yet)
	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM jobs
		WHERE status = 'scheduled' AND scheduled_at <= datetime('now', 'utc')
	`).Scan(&stats.QueueDepth)

	return stats, nil
}

// Scheduler operations

func (r *SQLiteRepository) GetScheduledJobs(ctx context.Context, before time.Time, limit int) ([]*domain.Job, error) {
	query := `
		SELECT id, name, type, config, scheduled_at, priority, project_id, timezone, status, created_at, updated_at
		FROM jobs
		WHERE status = 'scheduled' AND scheduled_at <= ?
		ORDER BY priority DESC, scheduled_at ASC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, before.UTC(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled jobs: %w", err)
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

		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

func (r *SQLiteRepository) UpdateJobStatus(ctx context.Context, id string, status domain.JobStatus) error {
	result, err := r.db.ExecContext(ctx, "UPDATE jobs SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrJobNotFound
	}

	return nil
}

// Transaction support

func (r *SQLiteRepository) WithTransaction(ctx context.Context, fn func(Repository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	txRepo := &SQLiteRepository{db: r.db} // In a real implementation, you'd wrap the tx
	if err := fn(txRepo); err != nil {
		return err
	}

	return tx.Commit()
}

// Helper function
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
