package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meysam81/oneoff/internal/logging"
	"github.com/meysam81/oneoff/internal/service"
)

// Handler contains all HTTP handlers
type Handler struct {
	jobService       *service.JobService
	executionService *service.ExecutionService
	projectService   *service.ProjectService
	tagService       *service.TagService
	systemService    *service.SystemService
}

// NewHandler creates a new handler
func NewHandler(
	jobService *service.JobService,
	executionService *service.ExecutionService,
	projectService *service.ProjectService,
	tagService *service.TagService,
	systemService *service.SystemService,
) *Handler {
	return &Handler{
		jobService:       jobService,
		executionService: executionService,
		projectService:   projectService,
		tagService:       tagService,
		systemService:    systemService,
	}
}

// Response helpers

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Data interface{} `json:"data"`
}

type listResponse struct {
	Data   interface{} `json:"data"`
	Total  int64       `json:"total,omitempty"`
	Limit  int         `json:"limit,omitempty"`
	Offset int         `json:"offset,omitempty"`
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logging.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, errorResponse{Error: message})
}

func (h *Handler) respondSuccess(w http.ResponseWriter, status int, data interface{}) {
	h.respondJSON(w, status, successResponse{Data: data})
}

func (h *Handler) respondList(w http.ResponseWriter, data interface{}, total int64, limit, offset int) {
	h.respondJSON(w, http.StatusOK, listResponse{
		Data:   data,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// Helper to parse query parameters
func getQueryInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func getQueryString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func getQueryStrings(r *http.Request, key string) []string {
	return r.URL.Query()[key]
}
