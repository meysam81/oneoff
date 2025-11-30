package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/service"
)

// ChainHandler contains chain-related HTTP handlers
type ChainHandler struct {
	chainService *service.ChainService
}

// NewChainHandler creates a new chain handler
func NewChainHandler(chainService *service.ChainService) *ChainHandler {
	return &ChainHandler{
		chainService: chainService,
	}
}

// CreateChain handles POST /api/chains
func (h *ChainHandler) CreateChain(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateChainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	chain, err := h.chainService.CreateChain(r.Context(), req)
	if err != nil {
		if err == domain.ErrChainEmpty {
			respondError(w, http.StatusBadRequest, "Chain must have at least one job")
			return
		}
		if err == domain.ErrChainJobNotFound {
			respondError(w, http.StatusBadRequest, "One or more jobs in chain not found")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusCreated, chain)
}

// GetChain handles GET /api/chains/:id
func (h *ChainHandler) GetChain(w http.ResponseWriter, r *http.Request) {
	id := extractChainID(r.URL.Path)

	if id == "" {
		respondError(w, http.StatusBadRequest, "Chain ID is required")
		return
	}

	chain, err := h.chainService.GetChain(r.Context(), id)
	if err != nil {
		if err == domain.ErrChainNotFound {
			respondError(w, http.StatusNotFound, "Chain not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, chain)
}

// ListChains handles GET /api/chains
func (h *ChainHandler) ListChains(w http.ResponseWriter, r *http.Request) {
	filter := domain.ChainFilter{
		ProjectID: getQueryString(r, "project_id"),
		Limit:     getQueryInt(r, "limit", 50),
		Offset:    getQueryInt(r, "offset", 0),
	}

	chains, total, err := h.chainService.ListChains(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondListChains(w, chains, total, filter.Limit, filter.Offset)
}

// UpdateChain handles PATCH /api/chains/:id
func (h *ChainHandler) UpdateChain(w http.ResponseWriter, r *http.Request) {
	id := extractChainID(r.URL.Path)

	if id == "" {
		respondError(w, http.StatusBadRequest, "Chain ID is required")
		return
	}

	var req domain.UpdateChainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	chain, err := h.chainService.UpdateChain(r.Context(), id, req)
	if err != nil {
		if err == domain.ErrChainNotFound {
			respondError(w, http.StatusNotFound, "Chain not found")
			return
		}
		if err == domain.ErrChainEmpty {
			respondError(w, http.StatusBadRequest, "Chain must have at least one job")
			return
		}
		if err == domain.ErrChainJobNotFound {
			respondError(w, http.StatusBadRequest, "One or more jobs in chain not found")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, chain)
}

// DeleteChain handles DELETE /api/chains/:id
func (h *ChainHandler) DeleteChain(w http.ResponseWriter, r *http.Request) {
	id := extractChainID(r.URL.Path)

	if id == "" {
		respondError(w, http.StatusBadRequest, "Chain ID is required")
		return
	}

	if err := h.chainService.DeleteChain(r.Context(), id); err != nil {
		if err == domain.ErrChainNotFound {
			respondError(w, http.StatusNotFound, "Chain not found")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ExecuteChain handles POST /api/chains/:id/execute
func (h *ChainHandler) ExecuteChain(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/chains/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "execute" {
		respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}
	id := parts[0]

	if err := h.chainService.ExecuteChain(r.Context(), id); err != nil {
		if err == domain.ErrChainNotFound {
			respondError(w, http.StatusNotFound, "Chain not found")
			return
		}
		if err == domain.ErrChainEmpty {
			respondError(w, http.StatusBadRequest, "Chain has no jobs to execute")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusAccepted, map[string]string{"message": "Chain execution started"})
}

// Helper functions for chain handler

func extractChainID(path string) string {
	id := strings.TrimPrefix(path, "/api/chains/")
	id = strings.Split(id, "/")[0]
	return id
}

func respondListChains(w http.ResponseWriter, data interface{}, total int64, limit, offset int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   data,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
