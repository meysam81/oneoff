package domain

import (
	"context"
	"time"
)

// APIKeyScope represents permission scopes for API keys
type APIKeyScope string

const (
	APIKeyScopeRead  APIKeyScope = "read"
	APIKeyScopeWrite APIKeyScope = "write"
	APIKeyScopeAdmin APIKeyScope = "admin"
)

// APIKey represents an API key for authentication
type APIKey struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	KeyHash    string     `json:"-"` // Never expose hash
	KeyPrefix  string     `json:"key_prefix"`
	Scopes     string     `json:"scopes"` // Comma-separated: read,write,admin
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
}

// APIKeyWithSecret is returned only when creating a new key
type APIKeyWithSecret struct {
	APIKey
	Key string `json:"key"` // Only returned on creation
}

// CreateAPIKeyRequest represents a request to create a new API key
type CreateAPIKeyRequest struct {
	Name      string     `json:"name"`
	Scopes    string     `json:"scopes,omitempty"`    // Default: "read,write"
	ExpiresAt *time.Time `json:"expires_at,omitempty"` // Optional expiration
}

// UpdateAPIKeyRequest represents a request to update an API key
type UpdateAPIKeyRequest struct {
	Name      *string    `json:"name,omitempty"`
	Scopes    *string    `json:"scopes,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// HasScope checks if the API key has the specified scope
func (k *APIKey) HasScope(scope APIKeyScope) bool {
	if k.Scopes == "" {
		return false
	}
	// Admin scope has all permissions
	if containsScope(k.Scopes, APIKeyScopeAdmin) {
		return true
	}
	// Read scope is implied by write scope
	if scope == APIKeyScopeRead && containsScope(k.Scopes, APIKeyScopeWrite) {
		return true
	}
	return containsScope(k.Scopes, scope)
}

// IsValid checks if the API key is valid (active and not expired)
func (k *APIKey) IsValid() bool {
	if !k.IsActive {
		return false
	}
	if k.ExpiresAt != nil && k.ExpiresAt.Before(time.Now().UTC()) {
		return false
	}
	return true
}

// containsScope checks if scopes string contains the given scope
func containsScope(scopes string, scope APIKeyScope) bool {
	scopeList := splitScopes(scopes)
	for _, s := range scopeList {
		if s == string(scope) {
			return true
		}
	}
	return false
}

// splitScopes splits a comma-separated scopes string
func splitScopes(scopes string) []string {
	if scopes == "" {
		return nil
	}
	result := make([]string, 0)
	start := 0
	for i := 0; i <= len(scopes); i++ {
		if i == len(scopes) || scopes[i] == ',' {
			if i > start {
				result = append(result, scopes[start:i])
			}
			start = i + 1
		}
	}
	return result
}

// AuthContext holds authentication context for a request
type AuthContext struct {
	APIKey   *APIKey
	KeyID    string
	IsSystem bool // For internal system calls
}

// AuthContextKey is the context key for auth context
type authContextKey struct{}

// AuthContextKeyValue is used to store/retrieve auth context
var AuthContextKeyValue = authContextKey{}

// GetAuthContext retrieves auth context from context
func GetAuthContext(ctx context.Context) *AuthContext {
	if auth, ok := ctx.Value(AuthContextKeyValue).(*AuthContext); ok {
		return auth
	}
	return nil
}

// WithAuthContext adds auth context to context
func WithAuthContext(ctx context.Context, auth *AuthContext) context.Context {
	return context.WithValue(ctx, AuthContextKeyValue, auth)
}
