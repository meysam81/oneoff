package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/logging"
)

// CreateWebhook creates a new webhook
func (r *SQLiteRepository) CreateWebhook(ctx context.Context, webhook *domain.Webhook) error {
	query := `
		INSERT INTO webhooks (id, name, url, secret, events, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	if webhook.ID == "" {
		webhook.ID = generateID()
	}
	now := time.Now().UTC()
	if webhook.CreatedAt.IsZero() {
		webhook.CreatedAt = now
	}
	if webhook.UpdatedAt.IsZero() {
		webhook.UpdatedAt = now
	}

	_, err := r.db.ExecContext(ctx, query,
		webhook.ID,
		webhook.Name,
		webhook.URL,
		webhook.Secret,
		webhook.Events,
		webhook.IsActive,
		webhook.CreatedAt,
		webhook.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}

	return nil
}

// GetWebhook retrieves a webhook by ID
func (r *SQLiteRepository) GetWebhook(ctx context.Context, id string) (*domain.Webhook, error) {
	query := `
		SELECT id, name, url, secret, events, is_active, created_at, updated_at
		FROM webhooks
		WHERE id = ?
	`

	webhook := &domain.Webhook{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&webhook.ID,
		&webhook.Name,
		&webhook.URL,
		&webhook.Secret,
		&webhook.Events,
		&webhook.IsActive,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	return webhook, nil
}

// ListWebhooks returns webhooks with optional filters
func (r *SQLiteRepository) ListWebhooks(ctx context.Context, filter domain.WebhookFilter) ([]*domain.Webhook, error) {
	query := `
		SELECT id, name, url, secret, events, is_active, created_at, updated_at
		FROM webhooks
		WHERE 1=1
	`
	var args []interface{}

	if filter.IsActive != nil {
		query += " AND is_active = ?"
		args = append(args, *filter.IsActive)
	}

	query += " ORDER BY created_at DESC"

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
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logging.Error().Err(err).Msg("Failed to close rows in ListWebhooks")
		}
	}()

	var webhooks []*domain.Webhook
	for rows.Next() {
		webhook := &domain.Webhook{}
		if err := rows.Scan(
			&webhook.ID,
			&webhook.Name,
			&webhook.URL,
			&webhook.Secret,
			&webhook.Events,
			&webhook.IsActive,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan webhook: %w", err)
		}
		webhooks = append(webhooks, webhook)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating webhooks: %w", err)
	}

	return webhooks, nil
}

// GetActiveWebhooksForEvent returns active webhooks subscribed to a specific event
func (r *SQLiteRepository) GetActiveWebhooksForEvent(ctx context.Context, event domain.WebhookEventType) ([]*domain.Webhook, error) {
	// SQLite LIKE query for comma-separated events
	query := `
		SELECT id, name, url, secret, events, is_active, created_at, updated_at
		FROM webhooks
		WHERE is_active = 1 AND (
			events LIKE ? OR
			events LIKE ? OR
			events LIKE ? OR
			events = ?
		)
	`

	eventStr := string(event)
	rows, err := r.db.QueryContext(ctx, query,
		eventStr+",%",      // Event at start
		"%,"+eventStr+",%", // Event in middle
		"%,"+eventStr,      // Event at end
		eventStr,           // Event alone
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhooks for event: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logging.Error().Err(err).Msg("Failed to close rows in GetActiveWebhooksForEvent")
		}
	}()

	var webhooks []*domain.Webhook
	for rows.Next() {
		webhook := &domain.Webhook{}
		if err := rows.Scan(
			&webhook.ID,
			&webhook.Name,
			&webhook.URL,
			&webhook.Secret,
			&webhook.Events,
			&webhook.IsActive,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan webhook: %w", err)
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// UpdateWebhook updates a webhook
func (r *SQLiteRepository) UpdateWebhook(ctx context.Context, id string, updates domain.UpdateWebhookRequest) error {
	var setParts []string
	var args []interface{}

	if updates.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *updates.Name)
	}
	if updates.URL != nil {
		setParts = append(setParts, "url = ?")
		args = append(args, *updates.URL)
	}
	if updates.Secret != nil {
		setParts = append(setParts, "secret = ?")
		args = append(args, *updates.Secret)
	}
	if updates.Events != nil {
		setParts = append(setParts, "events = ?")
		args = append(args, *updates.Events)
	}
	if updates.IsActive != nil {
		setParts = append(setParts, "is_active = ?")
		args = append(args, *updates.IsActive)
	}

	if len(setParts) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE webhooks SET %s WHERE id = ?", strings.Join(setParts, ", "))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// DeleteWebhook deletes a webhook
func (r *SQLiteRepository) DeleteWebhook(ctx context.Context, id string) error {
	query := `DELETE FROM webhooks WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// CreateWebhookDelivery creates a webhook delivery record
func (r *SQLiteRepository) CreateWebhookDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error {
	query := `
		INSERT INTO webhook_deliveries (id, webhook_id, event_type, payload, status, attempts, next_retry_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	if delivery.ID == "" {
		delivery.ID = generateID()
	}
	if delivery.CreatedAt.IsZero() {
		delivery.CreatedAt = time.Now().UTC()
	}

	_, err := r.db.ExecContext(ctx, query,
		delivery.ID,
		delivery.WebhookID,
		delivery.EventType,
		delivery.Payload,
		delivery.Status,
		delivery.Attempts,
		delivery.NextRetryAt,
		delivery.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create webhook delivery: %w", err)
	}

	return nil
}

// GetWebhookDelivery retrieves a webhook delivery by ID
func (r *SQLiteRepository) GetWebhookDelivery(ctx context.Context, id string) (*domain.WebhookDelivery, error) {
	query := `
		SELECT id, webhook_id, event_type, payload, status, response_code, response_body, error, attempts, next_retry_at, created_at
		FROM webhook_deliveries
		WHERE id = ?
	`

	delivery := &domain.WebhookDelivery{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&delivery.ID,
		&delivery.WebhookID,
		&delivery.EventType,
		&delivery.Payload,
		&delivery.Status,
		&delivery.ResponseCode,
		&delivery.ResponseBody,
		&delivery.Error,
		&delivery.Attempts,
		&delivery.NextRetryAt,
		&delivery.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook delivery: %w", err)
	}

	return delivery, nil
}

// ListWebhookDeliveries returns webhook deliveries with filters
func (r *SQLiteRepository) ListWebhookDeliveries(ctx context.Context, filter domain.WebhookDeliveryFilter) ([]*domain.WebhookDelivery, error) {
	query := `
		SELECT id, webhook_id, event_type, payload, status, response_code, response_body, error, attempts, next_retry_at, created_at
		FROM webhook_deliveries
		WHERE 1=1
	`
	var args []interface{}

	if filter.WebhookID != "" {
		query += " AND webhook_id = ?"
		args = append(args, filter.WebhookID)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}

	query += " ORDER BY created_at DESC"

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
		return nil, fmt.Errorf("failed to list webhook deliveries: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logging.Error().Err(err).Msg("Failed to close rows in ListWebhookDeliveries")
		}
	}()

	var deliveries []*domain.WebhookDelivery
	for rows.Next() {
		delivery := &domain.WebhookDelivery{}
		if err := rows.Scan(
			&delivery.ID,
			&delivery.WebhookID,
			&delivery.EventType,
			&delivery.Payload,
			&delivery.Status,
			&delivery.ResponseCode,
			&delivery.ResponseBody,
			&delivery.Error,
			&delivery.Attempts,
			&delivery.NextRetryAt,
			&delivery.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan webhook delivery: %w", err)
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, nil
}

// UpdateWebhookDelivery updates a webhook delivery
func (r *SQLiteRepository) UpdateWebhookDelivery(ctx context.Context, id string, status domain.WebhookDeliveryStatus, responseCode *int, responseBody, errMsg string, nextRetry *time.Time) error {
	query := `
		UPDATE webhook_deliveries
		SET status = ?, response_code = ?, response_body = ?, error = ?, next_retry_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, responseCode, responseBody, errMsg, nextRetry, id)
	if err != nil {
		return fmt.Errorf("failed to update webhook delivery: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetPendingWebhookDeliveries returns pending deliveries ready for retry
func (r *SQLiteRepository) GetPendingWebhookDeliveries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error) {
	query := `
		SELECT id, webhook_id, event_type, payload, status, response_code, response_body, error, attempts, next_retry_at, created_at
		FROM webhook_deliveries
		WHERE status = 'pending' AND (next_retry_at IS NULL OR next_retry_at <= ?)
		ORDER BY created_at ASC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, time.Now().UTC(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending deliveries: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logging.Error().Err(err).Msg("Failed to close rows in GetPendingDeliveries")
		}
	}()

	var deliveries []*domain.WebhookDelivery
	for rows.Next() {
		delivery := &domain.WebhookDelivery{}
		if err := rows.Scan(
			&delivery.ID,
			&delivery.WebhookID,
			&delivery.EventType,
			&delivery.Payload,
			&delivery.Status,
			&delivery.ResponseCode,
			&delivery.ResponseBody,
			&delivery.Error,
			&delivery.Attempts,
			&delivery.NextRetryAt,
			&delivery.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan webhook delivery: %w", err)
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, nil
}

// IncrementDeliveryAttempts increments the attempt count
func (r *SQLiteRepository) IncrementDeliveryAttempts(ctx context.Context, id string) error {
	query := `UPDATE webhook_deliveries SET attempts = attempts + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment delivery attempts: %w", err)
	}
	return nil
}
