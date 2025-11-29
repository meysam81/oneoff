package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/service"
)

// APIKeyHandler handles API key HTTP requests
type APIKeyHandler struct {
	service *service.APIKeyService
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(service *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{service: service}
}

// ListAPIKeys handles GET /api/api-keys
func (h *APIKeyHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check admin scope
	auth := domain.GetAuthContext(ctx)
	if auth != nil && auth.APIKey != nil && !auth.APIKey.HasScope(domain.APIKeyScopeAdmin) {
		respondError(w, http.StatusForbidden, "admin scope required")
		return
	}

	keys, err := h.service.ListAPIKeys(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, keys)
}

// CreateAPIKey handles POST /api/api-keys
func (h *APIKeyHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check admin scope
	auth := domain.GetAuthContext(ctx)
	if auth != nil && auth.APIKey != nil && !auth.APIKey.HasScope(domain.APIKeyScopeAdmin) {
		respondError(w, http.StatusForbidden, "admin scope required")
		return
	}

	var req domain.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	key, err := h.service.CreateAPIKey(ctx, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusCreated, key)
}

// GetAPIKey handles GET /api/api-keys/:id
func (h *APIKeyHandler) GetAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := extractIDFromPath(r.URL.Path, "/api/api-keys/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing API key ID")
		return
	}

	// Check admin scope
	auth := domain.GetAuthContext(ctx)
	if auth != nil && auth.APIKey != nil && !auth.APIKey.HasScope(domain.APIKeyScopeAdmin) {
		respondError(w, http.StatusForbidden, "admin scope required")
		return
	}

	key, err := h.service.GetAPIKey(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "API key not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, key)
}

// UpdateAPIKey handles PATCH /api/api-keys/:id
func (h *APIKeyHandler) UpdateAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := extractIDFromPath(r.URL.Path, "/api/api-keys/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing API key ID")
		return
	}

	// Check admin scope
	auth := domain.GetAuthContext(ctx)
	if auth != nil && auth.APIKey != nil && !auth.APIKey.HasScope(domain.APIKeyScopeAdmin) {
		respondError(w, http.StatusForbidden, "admin scope required")
		return
	}

	var req domain.UpdateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	key, err := h.service.UpdateAPIKey(ctx, id, req)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "API key not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, key)
}

// DeleteAPIKey handles DELETE /api/api-keys/:id
func (h *APIKeyHandler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := extractIDFromPath(r.URL.Path, "/api/api-keys/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing API key ID")
		return
	}

	// Check admin scope
	auth := domain.GetAuthContext(ctx)
	if auth != nil && auth.APIKey != nil && !auth.APIKey.HasScope(domain.APIKeyScopeAdmin) {
		respondError(w, http.StatusForbidden, "admin scope required")
		return
	}

	if err := h.service.DeleteAPIKey(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "API key not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RevokeAPIKey handles POST /api/api-keys/:id/revoke
func (h *APIKeyHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	path := strings.TrimSuffix(r.URL.Path, "/revoke")
	id := extractIDFromPath(path, "/api/api-keys/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing API key ID")
		return
	}

	// Check admin scope
	auth := domain.GetAuthContext(ctx)
	if auth != nil && auth.APIKey != nil && !auth.APIKey.HasScope(domain.APIKeyScopeAdmin) {
		respondError(w, http.StatusForbidden, "admin scope required")
		return
	}

	if err := h.service.RevokeAPIKey(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "API key not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, map[string]string{"status": "revoked"})
}

// RotateAPIKey handles POST /api/api-keys/:id/rotate
func (h *APIKeyHandler) RotateAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	path := strings.TrimSuffix(r.URL.Path, "/rotate")
	id := extractIDFromPath(path, "/api/api-keys/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing API key ID")
		return
	}

	// Check admin scope
	auth := domain.GetAuthContext(ctx)
	if auth != nil && auth.APIKey != nil && !auth.APIKey.HasScope(domain.APIKeyScopeAdmin) {
		respondError(w, http.StatusForbidden, "admin scope required")
		return
	}

	newKey, err := h.service.RotateAPIKey(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "API key not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, newKey)
}

// Helper functions

func extractIDFromPath(path, prefix string) string {
	path = strings.TrimPrefix(path, prefix)
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

type apiKeyErrorResponse struct {
	Error string `json:"error"`
}

type apiKeySuccessResponse struct {
	Data interface{} `json:"data"`
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiKeyErrorResponse{Error: message})
}

func respondSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiKeySuccessResponse{Data: data})
}
