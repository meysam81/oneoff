package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/service"
)

// WebhookHandler handles webhook HTTP requests
type WebhookHandler struct {
	service *service.WebhookService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(service *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{service: service}
}

// ListWebhooks handles GET /api/webhooks
func (h *WebhookHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := domain.WebhookFilter{
		Limit:  getQueryInt(r, "limit", 100),
		Offset: getQueryInt(r, "offset", 0),
	}

	webhooks, err := h.service.ListWebhooks(ctx, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, webhooks)
}

// CreateWebhook handles POST /api/webhooks
func (h *WebhookHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req domain.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	webhook, err := h.service.CreateWebhook(ctx, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusCreated, webhook)
}

// GetWebhook handles GET /api/webhooks/:id
func (h *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := extractIDFromPath(r.URL.Path, "/api/webhooks/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing webhook ID")
		return
	}

	webhook, err := h.service.GetWebhook(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "webhook not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, webhook)
}

// UpdateWebhook handles PATCH /api/webhooks/:id
func (h *WebhookHandler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := extractIDFromPath(r.URL.Path, "/api/webhooks/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing webhook ID")
		return
	}

	var req domain.UpdateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	webhook, err := h.service.UpdateWebhook(ctx, id, req)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "webhook not found")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, webhook)
}

// DeleteWebhook handles DELETE /api/webhooks/:id
func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := extractIDFromPath(r.URL.Path, "/api/webhooks/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing webhook ID")
		return
	}

	if err := h.service.DeleteWebhook(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "webhook not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// TestWebhook handles POST /api/webhooks/:id/test
func (h *WebhookHandler) TestWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	path := strings.TrimSuffix(r.URL.Path, "/test")
	id := extractIDFromPath(path, "/api/webhooks/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing webhook ID")
		return
	}

	if err := h.service.TestWebhook(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "webhook not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, map[string]string{"status": "test sent"})
}

// ListDeliveries handles GET /api/webhooks/:id/deliveries
func (h *WebhookHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	path := strings.TrimSuffix(r.URL.Path, "/deliveries")
	id := extractIDFromPath(path, "/api/webhooks/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing webhook ID")
		return
	}

	filter := domain.WebhookDeliveryFilter{
		WebhookID: id,
		Limit:     getQueryInt(r, "limit", 50),
		Offset:    getQueryInt(r, "offset", 0),
	}

	deliveries, err := h.service.ListDeliveries(ctx, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, deliveries)
}

// GetWebhookEvents handles GET /api/webhook-events
func (h *WebhookHandler) GetWebhookEvents(w http.ResponseWriter, r *http.Request) {
	events := domain.AllWebhookEvents()
	eventStrings := make([]string, len(events))
	for i, e := range events {
		eventStrings[i] = string(e)
	}
	respondSuccess(w, http.StatusOK, eventStrings)
}
