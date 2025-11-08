package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
)

// Execution handlers

func (h *Handler) GetExecution(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/executions/")
	id = strings.Split(id, "/")[0]

	execution, err := h.executionService.GetExecution(r.Context(), id)
	if err != nil {
		if err == domain.ErrExecutionNotFound {
			h.respondError(w, http.StatusNotFound, "Execution not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, execution)
}

func (h *Handler) ListExecutions(w http.ResponseWriter, r *http.Request) {
	filter := domain.ExecutionFilter{
		JobID:      getQueryString(r, "job_id"),
		Status:     domain.ExecutionStatus(getQueryString(r, "status")),
		ProjectID:  getQueryString(r, "project_id"),
		TagIDs:     getQueryStrings(r, "tag_id"),
		Limit:      getQueryInt(r, "limit", 50),
		Offset:     getQueryInt(r, "offset", 0),
		SortBy:     getQueryString(r, "sort_by"),
		SortOrder:  getQueryString(r, "sort_order"),
	}

	if dateFrom := getQueryString(r, "date_from"); dateFrom != "" {
		if t, err := time.Parse(time.RFC3339, dateFrom); err == nil {
			filter.DateFrom = &t
		}
	}
	if dateTo := getQueryString(r, "date_to"); dateTo != "" {
		if t, err := time.Parse(time.RFC3339, dateTo); err == nil {
			filter.DateTo = &t
		}
	}

	executions, err := h.executionService.ListExecutions(r.Context(), filter)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, executions)
}

// Project handlers

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project, err := h.projectService.CreateProject(r.Context(), req.Name, req.Description, req.Color, req.Icon)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, project)
}

func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	id = strings.Split(id, "/")[0]

	project, err := h.projectService.GetProject(r.Context(), id)
	if err != nil {
		if err == domain.ErrProjectNotFound {
			h.respondError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, project)
}

func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	includeArchived := getQueryString(r, "include_archived") == "true"

	projects, err := h.projectService.ListProjects(r.Context(), includeArchived)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, projects)
}

func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	id = strings.Split(id, "/")[0]

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Color       *string `json:"color"`
		Icon        *string `json:"icon"`
		IsArchived  *bool   `json:"is_archived"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project, err := h.projectService.UpdateProject(r.Context(), id, req.Name, req.Description, req.Color, req.Icon, req.IsArchived)
	if err != nil {
		if err == domain.ErrProjectNotFound {
			h.respondError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, project)
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	id = strings.Split(id, "/")[0]

	if err := h.projectService.DeleteProject(r.Context(), id); err != nil {
		if err == domain.ErrProjectNotFound {
			h.respondError(w, http.StatusNotFound, "Project not found")
			return
		}
		if err == domain.ErrCannotDeleteDefault {
			h.respondError(w, http.StatusBadRequest, "Cannot delete default project")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Tag handlers

func (h *Handler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string `json:"name"`
		Color     string `json:"color"`
		IsDefault bool   `json:"is_default"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tag, err := h.tagService.CreateTag(r.Context(), req.Name, req.Color, req.IsDefault)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, tag)
}

func (h *Handler) GetTag(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tags/")
	id = strings.Split(id, "/")[0]

	tag, err := h.tagService.GetTag(r.Context(), id)
	if err != nil {
		if err == domain.ErrTagNotFound {
			h.respondError(w, http.StatusNotFound, "Tag not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, tag)
}

func (h *Handler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.tagService.ListTags(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, tags)
}

func (h *Handler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tags/")
	id = strings.Split(id, "/")[0]

	var req struct {
		Name      *string `json:"name"`
		Color     *string `json:"color"`
		IsDefault *bool   `json:"is_default"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tag, err := h.tagService.UpdateTag(r.Context(), id, req.Name, req.Color, req.IsDefault)
	if err != nil {
		if err == domain.ErrTagNotFound {
			h.respondError(w, http.StatusNotFound, "Tag not found")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, tag)
}

func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tags/")
	id = strings.Split(id, "/")[0]

	if err := h.tagService.DeleteTag(r.Context(), id); err != nil {
		if err == domain.ErrTagNotFound {
			h.respondError(w, http.StatusNotFound, "Tag not found")
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// System handlers

func (h *Handler) GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	stats, err := h.systemService.GetSystemStats(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, stats)
}

func (h *Handler) GetWorkerStatus(w http.ResponseWriter, r *http.Request) {
	status := h.systemService.GetWorkerStatus(r.Context())
	h.respondSuccess(w, http.StatusOK, status)
}

func (h *Handler) GetSystemConfig(w http.ResponseWriter, r *http.Request) {
	config, err := h.systemService.GetConfig(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, config)
}

func (h *Handler) UpdateSystemConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.systemService.UpdateConfig(r.Context(), req.Key, req.Value); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, map[string]string{"message": "Config updated"})
}

func (h *Handler) GetJobTypes(w http.ResponseWriter, r *http.Request) {
	types := h.systemService.GetJobTypes()
	h.respondSuccess(w, http.StatusOK, types)
}
