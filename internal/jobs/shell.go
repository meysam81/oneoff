package jobs

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
)

// ShellJob implements JobExecutor for shell scripts
type ShellJob struct {
	config *domain.ShellJobConfig
}

// NewShellJob creates a new shell job
func NewShellJob(config string) (domain.JobExecutor, error) {
	cfg, err := domain.ParseShellJobConfig(config)
	if err != nil {
		return nil, fmt.Errorf("invalid shell job config: %w", err)
	}

	return &ShellJob{config: cfg}, nil
}

// Type returns the job type
func (j *ShellJob) Type() string {
	return "shell"
}

// Description returns job description
func (j *ShellJob) Description() string {
	if j.config.IsPath {
		return fmt.Sprintf("Execute shell script: %s", j.config.Script)
	}
	scriptPreview := j.config.Script
	if len(scriptPreview) > 50 {
		scriptPreview = scriptPreview[:50] + "..."
	}
	return fmt.Sprintf("Execute shell command: %s", scriptPreview)
}

// Validate validates the job configuration
func (j *ShellJob) Validate() error {
	if j.config.Script == "" {
		return fmt.Errorf("script is required")
	}
	if j.config.IsPath {
		// Expand path to handle ~ as home directory
		expandedPath, err := filepath.Abs(os.ExpandEnv(strings.Replace(j.config.Script, "~", "$HOME", 1)))
		if err != nil {
			return fmt.Errorf("failed to expand path: %w", err)
		}
		j.config.Script = expandedPath

		// Check if file exists
		if _, err := os.Stat(j.config.Script); os.IsNotExist(err) {
			return fmt.Errorf("script file not found: %s", j.config.Script)
		}
	}
	return nil
}

// Execute executes the shell script
func (j *ShellJob) Execute(ctx context.Context) (*domain.ExecutionResult, error) {
	if err := j.Validate(); err != nil {
		return nil, err
	}

	// Set timeout if specified
	if j.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(j.config.Timeout)*time.Second)
		defer cancel()
	}

	var cmd *exec.Cmd

	if j.config.IsPath {
		// Execute script file
		cmd = exec.CommandContext(ctx, "/bin/sh", j.config.Script)
		if len(j.config.Args) > 0 {
			cmd.Args = append(cmd.Args, j.config.Args...)
		}
	} else {
		// Execute inline script
		cmd = exec.CommandContext(ctx, "/bin/sh", "-c", j.config.Script)
		if len(j.config.Args) > 0 {
			// For inline scripts, args need to be passed differently
			// We'll pass them as positional parameters
			script := j.config.Script + " " + strings.Join(j.config.Args, " ")
			cmd = exec.CommandContext(ctx, "/bin/sh", "-c", script)
		}
	}

	// Set working directory
	if j.config.WorkDir != "" {
		cmd.Dir = j.config.WorkDir
	}

	// Set environment variables
	if len(j.config.Env) > 0 {
		cmd.Env = os.Environ()
		for key, value := range j.config.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	setSysProcAttr(cmd)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	err := cmd.Run()

	exitCode := 0
	errorMsg := ""

	if err != nil {
		if ctx.Err() == context.Canceled {
			killProcessGroup(cmd)
			exitCode = 130 // SIGINT exit code
			errorMsg = "Job cancelled by user"
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			errorMsg = fmt.Sprintf("Script exited with code %d: %s", exitCode, stderr.String())
		} else if ctx.Err() == context.DeadlineExceeded {
			exitCode = 124 // Timeout exit code
			errorMsg = "Script execution timeout"
		} else {
			exitCode = 1
			errorMsg = fmt.Sprintf("Failed to execute script: %v", err)
		}
	}

	// Combine stdout and stderr
	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n\n--- STDERR ---\n"
		}
		output += stderr.String()
	}

	return &domain.ExecutionResult{
		Output:   output,
		ExitCode: exitCode,
		Error:    errorMsg,
	}, nil
}
