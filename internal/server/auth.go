package server

import (
	"container/list"
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/logging"
	"github.com/meysam81/oneoff/internal/service"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Enabled   bool
	CacheSize int
	CacheTTL  time.Duration
	SkipPaths []string // Paths that don't require auth
}

// DefaultAuthConfig returns default auth configuration
func DefaultAuthConfig() AuthConfig {
	return AuthConfig{
		Enabled:   true,
		CacheSize: 128,
		CacheTTL:  5 * time.Minute,
		SkipPaths: []string{
			"/",        // Frontend
			"/assets/", // Static assets
			"/favicon.ico",
		},
	}
}

// authCacheEntry represents a cached API key
type authCacheEntry struct {
	key       *domain.APIKey
	expiresAt time.Time
}

// authCache is a thread-safe LRU cache for API keys
type authCache struct {
	mu      sync.RWMutex
	items   map[string]*list.Element
	order   *list.List
	maxSize int
	ttl     time.Duration
}

type cacheItem struct {
	key   string
	entry authCacheEntry
}

// newAuthCache creates a new auth cache
func newAuthCache(maxSize int, ttl time.Duration) *authCache {
	return &authCache{
		items:   make(map[string]*list.Element),
		order:   list.New(),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// Get retrieves a cached API key
func (c *authCache) Get(keyHash string) (*domain.APIKey, bool) {
	c.mu.RLock()
	elem, exists := c.items[keyHash]
	c.mu.RUnlock()

	if !exists {
		return nil, false
	}

	item := elem.Value.(*cacheItem)
	if time.Now().After(item.entry.expiresAt) {
		// Expired, remove it
		c.mu.Lock()
		c.order.Remove(elem)
		delete(c.items, keyHash)
		c.mu.Unlock()
		return nil, false
	}

	// Move to front (most recently used)
	c.mu.Lock()
	c.order.MoveToFront(elem)
	c.mu.Unlock()

	return item.entry.key, true
}

// Set adds or updates a cached API key
func (c *authCache) Set(keyHash string, key *domain.APIKey) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if already exists
	if elem, exists := c.items[keyHash]; exists {
		c.order.MoveToFront(elem)
		item := elem.Value.(*cacheItem)
		item.entry = authCacheEntry{
			key:       key,
			expiresAt: time.Now().Add(c.ttl),
		}
		return
	}

	// Evict if at capacity
	if c.order.Len() >= c.maxSize {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.items, oldest.Value.(*cacheItem).key)
		}
	}

	// Add new entry
	item := &cacheItem{
		key: keyHash,
		entry: authCacheEntry{
			key:       key,
			expiresAt: time.Now().Add(c.ttl),
		},
	}
	elem := c.order.PushFront(item)
	c.items[keyHash] = elem
}

// Invalidate removes a key from cache
func (c *authCache) Invalidate(keyHash string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[keyHash]; exists {
		c.order.Remove(elem)
		delete(c.items, keyHash)
	}
}

// AuthMiddleware creates authentication middleware
type AuthMiddleware struct {
	apiKeyService *service.APIKeyService
	config        AuthConfig
	cache         *authCache
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(apiKeyService *service.APIKeyService, config AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		apiKeyService: apiKeyService,
		config:        config,
		cache:         newAuthCache(config.CacheSize, config.CacheTTL),
	}
}

// Middleware returns the HTTP middleware function
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if auth is disabled
		if !m.config.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Skip auth for certain paths
		if m.shouldSkipAuth(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract API key from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.respondUnauthorized(w, "missing Authorization header")
			return
		}

		// Expect "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			m.respondUnauthorized(w, "invalid Authorization header format")
			return
		}

		rawKey := parts[1]
		if rawKey == "" {
			m.respondUnauthorized(w, "missing API key")
			return
		}

		// Validate API key
		apiKey, err := m.validateKey(r.Context(), rawKey)
		if err != nil {
			logging.Debug().Err(err).Str("key_prefix", safeKeyPrefix(rawKey)).Msg("API key validation failed")
			m.respondUnauthorized(w, "invalid or expired API key")
			return
		}

		// Check required scope based on HTTP method
		requiredScope := m.getRequiredScope(r.Method)
		if !apiKey.HasScope(requiredScope) {
			m.respondForbidden(w, "insufficient permissions")
			return
		}

		// Add auth context to request
		authCtx := &domain.AuthContext{
			APIKey: apiKey,
			KeyID:  apiKey.ID,
		}
		ctx := domain.WithAuthContext(r.Context(), authCtx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateKey validates an API key, using cache when possible
func (m *AuthMiddleware) validateKey(ctx context.Context, rawKey string) (*domain.APIKey, error) {
	// Check cache first (using a hash of the key as cache key)
	// We hash here to avoid storing raw keys in cache
	cacheKey := hashKey(rawKey)
	if cached, ok := m.cache.Get(cacheKey); ok {
		return cached, nil
	}

	// Validate against database
	apiKey, err := m.apiKeyService.ValidateAPIKey(ctx, rawKey)
	if err != nil {
		return nil, err
	}

	// Cache the result
	m.cache.Set(cacheKey, apiKey)

	return apiKey, nil
}

// shouldSkipAuth checks if a path should skip authentication
func (m *AuthMiddleware) shouldSkipAuth(path string) bool {
	// API paths always require auth
	if strings.HasPrefix(path, "/api/") {
		return false
	}

	// Check skip paths
	for _, skipPath := range m.config.SkipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	// Non-API paths (frontend) don't require auth
	return true
}

// getRequiredScope returns the required scope for an HTTP method
func (m *AuthMiddleware) getRequiredScope(method string) domain.APIKeyScope {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return domain.APIKeyScopeRead
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return domain.APIKeyScopeWrite
	default:
		return domain.APIKeyScopeWrite
	}
}

// respondUnauthorized sends a 401 response
func (m *AuthMiddleware) respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"` + message + `"}`))
}

// respondForbidden sends a 403 response
func (m *AuthMiddleware) respondForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":"` + message + `"}`))
}

// hashKey creates a simple hash for cache key lookup
func hashKey(key string) string {
	// Use first 32 chars as a quick lookup key
	// The actual validation will use the full hash
	if len(key) > 32 {
		return key[:32]
	}
	return key
}

// safeKeyPrefix returns a safe prefix for logging
func safeKeyPrefix(key string) string {
	if len(key) > 15 {
		return key[:15] + "..."
	}
	return "***"
}
