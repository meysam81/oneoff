package server

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/meysam81/oneoff/internal/config"
	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/handler"
	"github.com/meysam81/oneoff/internal/jobs"
	"github.com/meysam81/oneoff/internal/repository"
	"github.com/meysam81/oneoff/internal/service"
	"github.com/meysam81/oneoff/internal/worker"
	"github.com/rs/zerolog/log"
)

//go:embed dist
var distFS embed.FS

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	config     *config.Config
	repo       repository.Repository
	pool       *worker.Pool
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
	// Initialize database
	repo, err := repository.NewSQLiteRepository(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	// Run migrations
	migrationsPath := filepath.Join(".", "migrations")
	if err := repository.RunMigrations(cfg.DBPath, migrationsPath, "up"); err != nil {
		log.Warn().Err(err).Msg("Failed to run migrations (continuing anyway)")
	}

	// Initialize job registry
	registry := domain.NewJobRegistry()
	jobs.RegisterJobTypes(registry)

	// Initialize worker pool
	pool := worker.NewPool(cfg.WorkersCount, repo, registry)

	// Initialize services
	jobService := service.NewJobService(repo, registry, pool)
	executionService := service.NewExecutionService(repo)
	projectService := service.NewProjectService(repo)
	tagService := service.NewTagService(repo)
	systemService := service.NewSystemService(repo, pool)

	// Initialize handlers
	h := handler.NewHandler(jobService, executionService, projectService, tagService, systemService)

	// Setup router
	mux := http.NewServeMux()
	setupRoutes(mux, h)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    cfg.Address(),
		Handler: corsMiddleware(loggingMiddleware(mux)),
	}

	server := &Server{
		httpServer: httpServer,
		config:     cfg,
		repo:       repo,
		pool:       pool,
	}

	return server, nil
}

// Start starts the server
func (s *Server) Start() error {
	// Start worker pool
	ctx := context.Background()
	if err := s.pool.Start(ctx); err != nil {
		return fmt.Errorf("failed to start worker pool: %w", err)
	}

	// Start HTTP server
	log.Info().Str("address", s.config.Address()).Msg("Starting HTTP server")
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down server")

	// Stop worker pool
	if err := s.pool.Stop(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to stop worker pool")
	}

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown HTTP server")
		return err
	}

	// Close database
	if err := s.repo.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database")
		return err
	}

	return nil
}

// setupRoutes configures all HTTP routes
func setupRoutes(mux *http.ServeMux, h *handler.Handler) {
	// API routes
	// Jobs
	mux.HandleFunc("/api/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListJobs(w, r)
		case http.MethodPost:
			h.CreateJob(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/jobs/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
		parts := strings.Split(path, "/")

		if len(parts) >= 2 {
			// Job actions
			action := parts[1]
			switch action {
			case "execute":
				if r.Method == http.MethodPost {
					h.ExecuteJob(w, r)
					return
				}
			case "clone":
				if r.Method == http.MethodPost {
					h.CloneJob(w, r)
					return
				}
			case "cancel":
				if r.Method == http.MethodPost {
					h.CancelJob(w, r)
					return
				}
			}
		}

		// Single job operations
		switch r.Method {
		case http.MethodGet:
			h.GetJob(w, r)
		case http.MethodPatch:
			h.UpdateJob(w, r)
		case http.MethodDelete:
			h.DeleteJob(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Executions
	mux.HandleFunc("/api/executions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ListExecutions(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/executions/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetExecution(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Projects
	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListProjects(w, r)
		case http.MethodPost:
			h.CreateProject(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetProject(w, r)
		case http.MethodPatch:
			h.UpdateProject(w, r)
		case http.MethodDelete:
			h.DeleteProject(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Tags
	mux.HandleFunc("/api/tags", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListTags(w, r)
		case http.MethodPost:
			h.CreateTag(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/tags/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTag(w, r)
		case http.MethodPatch:
			h.UpdateTag(w, r)
		case http.MethodDelete:
			h.DeleteTag(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// System
	mux.HandleFunc("/api/system/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetSystemStatus(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/system/config", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetSystemConfig(w, r)
		case http.MethodPatch:
			h.UpdateSystemConfig(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/workers/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetWorkerStatus(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/job-types", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetJobTypes(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Serve frontend with SPA fallback
	distFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Warn().Err(err).Msg("Frontend dist not found, will serve placeholder")
		mux.HandleFunc("/", servePlaceholder)
	} else {
		mux.Handle("/", spaHandler(distFS))
	}
}

// spaHandler serves the SPA frontend with fallback to index.html for client-side routing
func spaHandler(fsys fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(fsys))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean the path
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Try to open the file
		file, err := fsys.Open(path)
		if err != nil {
			// File doesn't exist, serve index.html for SPA routing
			indexFile, indexErr := fsys.Open("index.html")
			if indexErr != nil {
				http.Error(w, "index.html not found", http.StatusNotFound)
				return
			}
			defer func() { _ = indexFile.Close() }()

			// Read and serve index.html
			stat, _ := indexFile.Stat()
			http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
			return
		}
		defer func() { _ = file.Close() }()

		// Check if it's a directory
		stat, err := file.Stat()
		if err != nil {
			http.Error(w, "Failed to stat file", http.StatusInternalServerError)
			return
		}

		if stat.IsDir() {
			// For directories, serve index.html (SPA routing)
			indexFile, indexErr := fsys.Open("index.html")
			if indexErr != nil {
				http.Error(w, "index.html not found", http.StatusNotFound)
				return
			}
			defer func() { _ = indexFile.Close() }()

			indexStat, _ := indexFile.Stat()
			http.ServeContent(w, r, "index.html", indexStat.ModTime(), indexFile.(io.ReadSeeker))
			return
		}

		// File exists and is not a directory, serve it normally
		fileServer.ServeHTTP(w, r)
	})
}

// Middleware

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote", r.RemoteAddr).
			Msg("HTTP request")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func servePlaceholder(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>OneOff - Job Scheduler</title>
    <style>
        body { font-family: system-ui; max-width: 800px; margin: 100px auto; padding: 20px; }
        h1 { color: #6366f1; }
        .status { background: #f0fdf4; border: 1px solid #22c55e; padding: 20px; border-radius: 8px; }
    </style>
</head>
<body>
    <h1>🎯 OneOff - Modern Job Scheduler</h1>
    <div class="status">
        <h2>Server Running</h2>
        <p>The OneOff backend is running successfully!</p>
        <p>Frontend will be available once built.</p>
        <h3>Available API Endpoints:</h3>
        <ul>
            <li>GET /api/jobs - List all jobs</li>
            <li>POST /api/jobs - Create a new job</li>
            <li>GET /api/executions - List job executions</li>
            <li>GET /api/projects - List projects</li>
            <li>GET /api/tags - List tags</li>
            <li>GET /api/system/status - Get system status</li>
            <li>GET /api/workers/status - Get worker status</li>
        </ul>
    </div>
</body>
</html>
	`
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(html))
}
