package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
)

// CreateJob handles POST /api/jobs
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	job, err := h.jobService.CreateJob(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, job)
}

// GetJob handles GET /api/jobs/:id
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	id = strings.Split(id, "/")[0]

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Job ID is required")
		return
	}

	job, err := h.jobService.GetJob(r.Context(), id)
	if err != nil {
		if err == domain.ErrJobNotFound {
			h.respondError(w, http.StatusNotFound, "Job not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, job)
}

// ListJobs handles GET /api/jobs
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	filter := domain.JobFilter{
		ProjectID: getQueryString(r, "project_id"),
		Status:    domain.JobStatus(getQueryString(r, "status")),
		JobType:   getQueryString(r, "type"),
		Search:    getQueryString(r, "search"),
		TagIDs:    getQueryStrings(r, "tag_id"),
		Limit:     getQueryInt(r, "limit", 50),
		Offset:    getQueryInt(r, "offset", 0),
		SortBy:    getQueryString(r, "sort_by"),
		SortOrder: getQueryString(r, "sort_order"),
	}

	// Parse time filters
	if timeFrom := getQueryString(r, "time_from"); timeFrom != "" {
		if t, err := time.Parse(time.RFC3339, timeFrom); err == nil {
			filter.TimeFrom = &t
		}
	}
	if timeTo := getQueryString(r, "time_to"); timeTo != "" {
		if t, err := time.Parse(time.RFC3339, timeTo); err == nil {
			filter.TimeTo = &t
		}
	}

	jobs, total, err := h.jobService.ListJobs(r.Context(), filter)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondList(w, jobs, total, filter.Limit, filter.Offset)
}

// UpdateJob handles PATCH /api/jobs/:id
func (h *Handler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	id = strings.Split(id, "/")[0]

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Job ID is required")
		return
	}

	var req domain.UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	job, err := h.jobService.UpdateJob(r.Context(), id, req)
	if err != nil {
		if err == domain.ErrJobNotFound {
			h.respondError(w, http.StatusNotFound, "Job not found")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, job)
}

// DeleteJob handles DELETE /api/jobs/:id
func (h *Handler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	id = strings.Split(id, "/")[0]

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Job ID is required")
		return
	}

	if err := h.jobService.DeleteJob(r.Context(), id); err != nil {
		if err == domain.ErrJobNotFound {
			h.respondError(w, http.StatusNotFound, "Job not found")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ExecuteJob handles POST /api/jobs/:id/execute
func (h *Handler) ExecuteJob(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "execute" {
		h.respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}
	id := parts[0]

	if err := h.jobService.ExecuteJobNow(r.Context(), id); err != nil {
		if err == domain.ErrJobNotFound {
			h.respondError(w, http.StatusNotFound, "Job not found")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, map[string]string{"message": "Job scheduled for immediate execution"})
}

// CloneJob handles POST /api/jobs/:id/clone
func (h *Handler) CloneJob(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "clone" {
		h.respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}
	id := parts[0]

	var req struct {
		ScheduledAt string `json:"scheduled_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var scheduledAt time.Time
	var err error

	if req.ScheduledAt == "now" {
		scheduledAt = time.Now()
	} else {
		scheduledAt, err = time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, "Invalid scheduled_at format")
			return
		}
	}

	job, err := h.jobService.CloneJob(r.Context(), id, scheduledAt)
	if err != nil {
		if err == domain.ErrJobNotFound {
			h.respondError(w, http.StatusNotFound, "Job not found")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, job)
}

// CancelJob handles POST /api/jobs/:id/cancel
func (h *Handler) CancelJob(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "cancel" {
		h.respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}
	id := parts[0]

	if err := h.jobService.CancelJob(r.Context(), id); err != nil {
		if err == domain.ErrJobNotFound {
			h.respondError(w, http.StatusNotFound, "Job not found")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, map[string]string{"message": "Job cancelled"})
}
