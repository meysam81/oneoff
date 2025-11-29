package jobs

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
)

// DockerJob implements JobExecutor for Docker containers
type DockerJob struct {
	config *domain.DockerJobConfig
}

// NewDockerJob creates a new Docker job
func NewDockerJob(config string) (domain.JobExecutor, error) {
	cfg, err := domain.ParseDockerJobConfig(config)
	if err != nil {
		return nil, fmt.Errorf("invalid docker job config: %w", err)
	}

	return &DockerJob{config: cfg}, nil
}

// Type returns the job type
func (j *DockerJob) Type() string {
	return "docker"
}

// Description returns job description
func (j *DockerJob) Description() string {
	return fmt.Sprintf("Run Docker container: %s", j.config.Image)
}

// Validate validates the job configuration
func (j *DockerJob) Validate() error {
	if j.config.Image == "" {
		return fmt.Errorf("image is required")
	}

	// Check if Docker is available
	cmd := exec.Command("docker", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker is not available or not running: %w", err)
	}

	return nil
}

// Execute executes the Docker container
func (j *DockerJob) Execute(ctx context.Context) (*domain.ExecutionResult, error) {
	if err := j.Validate(); err != nil {
		return nil, err
	}

	// Set timeout if specified
	if j.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(j.config.Timeout)*time.Second)
		defer cancel()
	}

	// Build docker run command
	args := []string{"run"}

	// Auto-remove container after execution
	if j.config.AutoRemove {
		args = append(args, "--rm")
	}

	// Add environment variables
	for key, value := range j.config.Env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// Add volume mounts
	for host, container := range j.config.Volumes {
		args = append(args, "-v", fmt.Sprintf("%s:%s", host, container))
	}

	// Add working directory
	if j.config.WorkDir != "" {
		args = append(args, "-w", j.config.WorkDir)
	}

	// Add image
	args = append(args, j.config.Image)

	// Add command
	if len(j.config.Command) > 0 {
		args = append(args, j.config.Command...)
	}

	// Create command
	cmd := exec.CommandContext(ctx, "docker", args...)

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
			exitCode = 130 // SIGINT exit code
			errorMsg = "Container execution cancelled by user"
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			errorMsg = fmt.Sprintf("Container exited with code %d: %s", exitCode, stderr.String())
		} else if ctx.Err() == context.DeadlineExceeded {
			exitCode = 124
			errorMsg = "Container execution timeout"
		} else {
			exitCode = 1
			errorMsg = fmt.Sprintf("Failed to run container: %v", err)
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

	// Add command info to output
	cmdInfo := fmt.Sprintf("Command: docker %s\n\n", strings.Join(args, " "))
	output = cmdInfo + output

	return &domain.ExecutionResult{
		Output:   output,
		ExitCode: exitCode,
		Error:    errorMsg,
	}, nil
}
