package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/repository"
)

// APIKeyService handles API key operations
type APIKeyService struct {
	repo repository.Repository
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(repo repository.Repository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

// CreateAPIKey creates a new API key and returns it with the secret
func (s *APIKeyService) CreateAPIKey(ctx context.Context, req domain.CreateAPIKeyRequest) (*domain.APIKeyWithSecret, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Generate random API key (32 bytes = 64 hex chars)
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Format: oneoff_<random_hex>
	rawKey := "oneoff_" + hex.EncodeToString(keyBytes)
	keyPrefix := rawKey[:15] // "oneoff_" + first 8 hex chars

	// Hash the key for storage
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	// Default scopes
	scopes := req.Scopes
	if scopes == "" {
		scopes = "read,write"
	}

	apiKey := &domain.APIKey{
		Name:      req.Name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		Scopes:    scopes,
		ExpiresAt: req.ExpiresAt,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreateAPIKey(ctx, apiKey); err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	// Return with the raw key (only time it's visible)
	return &domain.APIKeyWithSecret{
		APIKey: *apiKey,
		Key:    rawKey,
	}, nil
}

// GetAPIKey retrieves an API key by ID
func (s *APIKeyService) GetAPIKey(ctx context.Context, id string) (*domain.APIKey, error) {
	return s.repo.GetAPIKey(ctx, id)
}

// ListAPIKeys returns all API keys
func (s *APIKeyService) ListAPIKeys(ctx context.Context) ([]*domain.APIKey, error) {
	return s.repo.ListAPIKeys(ctx)
}

// UpdateAPIKey updates an API key
func (s *APIKeyService) UpdateAPIKey(ctx context.Context, id string, req domain.UpdateAPIKeyRequest) (*domain.APIKey, error) {
	// Verify key exists
	existing, err := s.repo.GetAPIKey(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateAPIKey(ctx, id, req); err != nil {
		return nil, err
	}

	// Apply updates to return object
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Scopes != nil {
		existing.Scopes = *req.Scopes
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.ExpiresAt != nil {
		existing.ExpiresAt = req.ExpiresAt
	}

	return existing, nil
}

// DeleteAPIKey deletes an API key
func (s *APIKeyService) DeleteAPIKey(ctx context.Context, id string) error {
	return s.repo.DeleteAPIKey(ctx, id)
}

// ValidateAPIKey validates an API key and returns the key info if valid
func (s *APIKeyService) ValidateAPIKey(ctx context.Context, rawKey string) (*domain.APIKey, error) {
	if rawKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Hash the provided key
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	// Look up by hash
	apiKey, err := s.repo.GetAPIKeyByHash(ctx, keyHash)
	if err != nil {
		return nil, fmt.Errorf("invalid API key")
	}

	// Check if valid
	if !apiKey.IsValid() {
		if !apiKey.IsActive {
			return nil, fmt.Errorf("API key is inactive")
		}
		if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now().UTC()) {
			return nil, fmt.Errorf("API key has expired")
		}
		return nil, fmt.Errorf("API key is invalid")
	}

	// Update last used (async, don't block on errors)
	go func() {
		_ = s.repo.UpdateAPIKeyLastUsed(context.Background(), apiKey.ID)
	}()

	return apiKey, nil
}

// RevokeAPIKey revokes (deactivates) an API key
func (s *APIKeyService) RevokeAPIKey(ctx context.Context, id string) error {
	isActive := false
	return s.repo.UpdateAPIKey(ctx, id, domain.UpdateAPIKeyRequest{
		IsActive: &isActive,
	})
}

// RotateAPIKey creates a new key and revokes the old one
func (s *APIKeyService) RotateAPIKey(ctx context.Context, id string) (*domain.APIKeyWithSecret, error) {
	// Get existing key
	existing, err := s.repo.GetAPIKey(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create new key with same settings
	newKey, err := s.CreateAPIKey(ctx, domain.CreateAPIKeyRequest{
		Name:      existing.Name + " (rotated)",
		Scopes:    existing.Scopes,
		ExpiresAt: existing.ExpiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create rotated key: %w", err)
	}

	// Revoke old key
	if err := s.RevokeAPIKey(ctx, id); err != nil {
		// Log but don't fail - new key was created
		fmt.Printf("warning: failed to revoke old key %s: %v\n", id, err)
	}

	return newKey, nil
}
