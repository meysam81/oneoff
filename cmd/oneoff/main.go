package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/meysam81/oneoff/internal/config"
	"github.com/meysam81/oneoff/internal/repository"
	"github.com/meysam81/oneoff/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	cmd := &cli.Command{
		Name:    "oneoff",
		Usage:   "OneOff - Modern Job Scheduler",
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Start the OneOff server",
				Action: func(ctx context.Context, c *cli.Command) error {
					return runServer()
				},
			},
			{
				Name:  "migrate",
				Usage: "Run database migrations",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "direction",
						Usage: "Migration direction (up or down)",
						Value: "up",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return runMigrations(c.String("direction"))
				},
			},
		},
		DefaultCommand: "serve",
		Action: func(ctx context.Context, c *cli.Command) error {
			return runServer()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("Application failed")
	}
}

func runServer() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup logging
	setupLogging(cfg.LogLevel)

	log.Info().
		Str("version", version).
		Str("commit", commit).
		Str("environment", cfg.Environment).
		Msg("Starting OneOff")

	// Create server
	srv, err := server.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		log.Info().Msg("Received shutdown signal")
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Info().Msg("Server stopped successfully")
	return nil
}

func runMigrations(direction string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup logging
	setupLogging(cfg.LogLevel)

	log.Info().
		Str("direction", direction).
		Str("db_path", cfg.DBPath).
		Msg("Running migrations")

	// Get migrations path
	migrationsPath := filepath.Join(".", "migrations")

	// Validate direction
	if direction != "up" && direction != "down" {
		return fmt.Errorf("invalid migration direction: %s (must be 'up' or 'down')", direction)
	}

	// Run migrations
	if err := repository.RunMigrations(cfg.DBPath, migrationsPath, direction); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Info().Msg("Migrations completed successfully")
	return nil
}

func setupLogging(level string) {
	// Set up zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Pretty logging for development
	if level == "debug" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
}
