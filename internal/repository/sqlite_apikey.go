package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/rs/zerolog/log"
)

// CreateAPIKey creates a new API key
func (r *SQLiteRepository) CreateAPIKey(ctx context.Context, key *domain.APIKey) error {
	query := `
		INSERT INTO api_keys (id, name, key_hash, key_prefix, scopes, expires_at, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	if key.ID == "" {
		key.ID = generateID()
	}
	if key.CreatedAt.IsZero() {
		key.CreatedAt = time.Now().UTC()
	}

	_, err := r.db.ExecContext(ctx, query,
		key.ID,
		key.Name,
		key.KeyHash,
		key.KeyPrefix,
		key.Scopes,
		key.ExpiresAt,
		key.IsActive,
		key.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}

	return nil
}

// GetAPIKey retrieves an API key by ID
func (r *SQLiteRepository) GetAPIKey(ctx context.Context, id string) (*domain.APIKey, error) {
	query := `
		SELECT id, name, key_hash, key_prefix, scopes, last_used_at, expires_at, is_active, created_at
		FROM api_keys
		WHERE id = ?
	`

	key := &domain.APIKey{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&key.ID,
		&key.Name,
		&key.KeyHash,
		&key.KeyPrefix,
		&key.Scopes,
		&key.LastUsedAt,
		&key.ExpiresAt,
		&key.IsActive,
		&key.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return key, nil
}

// GetAPIKeyByHash retrieves an API key by its hash (for authentication)
func (r *SQLiteRepository) GetAPIKeyByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	query := `
		SELECT id, name, key_hash, key_prefix, scopes, last_used_at, expires_at, is_active, created_at
		FROM api_keys
		WHERE key_hash = ? AND is_active = 1
	`

	key := &domain.APIKey{}
	err := r.db.QueryRowContext(ctx, query, keyHash).Scan(
		&key.ID,
		&key.Name,
		&key.KeyHash,
		&key.KeyPrefix,
		&key.Scopes,
		&key.LastUsedAt,
		&key.ExpiresAt,
		&key.IsActive,
		&key.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get API key by hash: %w", err)
	}

	return key, nil
}

// ListAPIKeys returns all API keys
func (r *SQLiteRepository) ListAPIKeys(ctx context.Context) ([]*domain.APIKey, error) {
	query := `
		SELECT id, name, key_hash, key_prefix, scopes, last_used_at, expires_at, is_active, created_at
		FROM api_keys
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close rows in ListAPIKeys")
		}
	}()

	var keys []*domain.APIKey
	for rows.Next() {
		key := &domain.APIKey{}
		if err := rows.Scan(
			&key.ID,
			&key.Name,
			&key.KeyHash,
			&key.KeyPrefix,
			&key.Scopes,
			&key.LastUsedAt,
			&key.ExpiresAt,
			&key.IsActive,
			&key.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating API keys: %w", err)
	}

	return keys, nil
}

// UpdateAPIKey updates an API key
func (r *SQLiteRepository) UpdateAPIKey(ctx context.Context, id string, updates domain.UpdateAPIKeyRequest) error {
	var setParts []string
	var args []interface{}

	if updates.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *updates.Name)
	}
	if updates.Scopes != nil {
		setParts = append(setParts, "scopes = ?")
		args = append(args, *updates.Scopes)
	}
	if updates.IsActive != nil {
		setParts = append(setParts, "is_active = ?")
		args = append(args, *updates.IsActive)
	}
	if updates.ExpiresAt != nil {
		setParts = append(setParts, "expires_at = ?")
		args = append(args, *updates.ExpiresAt)
	}

	if len(setParts) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE api_keys SET %s WHERE id = ?", strings.Join(setParts, ", "))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
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

// UpdateAPIKeyLastUsed updates the last_used_at timestamp
func (r *SQLiteRepository) UpdateAPIKeyLastUsed(ctx context.Context, id string) error {
	query := `UPDATE api_keys SET last_used_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update API key last used: %w", err)
	}
	return nil
}

// DeleteAPIKey deletes an API key
func (r *SQLiteRepository) DeleteAPIKey(ctx context.Context, id string) error {
	query := `DELETE FROM api_keys WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
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

// generateID generates a random hex ID
func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
