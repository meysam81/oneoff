package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/meysam81/oneoff/internal/domain"
)

// HTTPJob implements JobExecutor for HTTP requests
type HTTPJob struct {
	config *domain.HTTPJobConfig
	client *req.Client
}

// NewHTTPJob creates a new HTTP job
func NewHTTPJob(config string) (domain.JobExecutor, error) {
	cfg, err := domain.ParseHTTPJobConfig(config)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP job config: %w", err)
	}

	job := &HTTPJob{
		config: cfg,
		client: req.C().SetTimeout(30 * time.Second),
	}

	if cfg.Timeout > 0 {
		job.client.SetTimeout(time.Duration(cfg.Timeout) * time.Second)
	}

	return job, nil
}

// Type returns the job type
func (j *HTTPJob) Type() string {
	return "http"
}

// Description returns job description
func (j *HTTPJob) Description() string {
	return fmt.Sprintf("HTTP %s request to %s", j.config.Method, j.config.URL)
}

// Validate validates the job configuration
func (j *HTTPJob) Validate() error {
	if j.config.URL == "" {
		return fmt.Errorf("URL is required")
	}

	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true,
		"DELETE": true, "HEAD": true, "OPTIONS": true,
	}

	method := strings.ToUpper(j.config.Method)
	if method == "" {
		j.config.Method = "GET"
	} else if !validMethods[method] {
		return fmt.Errorf("invalid HTTP method: %s", j.config.Method)
	}

	return nil
}

// Execute executes the HTTP request
func (j *HTTPJob) Execute(ctx context.Context) (*domain.ExecutionResult, error) {
	if err := j.Validate(); err != nil {
		return nil, err
	}

	method := strings.ToUpper(j.config.Method)
	if method == "" {
		method = "GET"
	}

	// Create request
	request := j.client.R().SetContext(ctx)

	// Add headers
	for key, value := range j.config.Headers {
		request.SetHeader(key, value)
	}

	// Add body if present
	if j.config.Body != "" {
		request.SetBodyString(j.config.Body)
	}

	// Execute request with retry
	j.client.SetCommonRetryCount(3).
		SetCommonRetryFixedInterval(2 * time.Second)

	var resp *req.Response
	var err error

	switch method {
	case "GET":
		resp, err = request.Get(j.config.URL)
	case "POST":
		resp, err = request.Post(j.config.URL)
	case "PUT":
		resp, err = request.Put(j.config.URL)
	case "PATCH":
		resp, err = request.Patch(j.config.URL)
	case "DELETE":
		resp, err = request.Delete(j.config.URL)
	case "HEAD":
		resp, err = request.Head(j.config.URL)
	case "OPTIONS":
		resp, err = request.Options(j.config.URL)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return &domain.ExecutionResult{
			Output:   "",
			ExitCode: 1,
			Error:    fmt.Sprintf("HTTP request failed: %v", err),
		}, nil
	}

	statusCode := resp.StatusCode
	output := fmt.Sprintf("Status: %d %s\n\nHeaders:\n", statusCode, resp.Status)

	// Add response headers
	for key, values := range resp.Header {
		output += fmt.Sprintf("%s: %s\n", key, strings.Join(values, ", "))
	}

	// Add response body
	body := resp.String()
	if body != "" {
		output += fmt.Sprintf("\nBody:\n%s", body)
	}

	// Determine exit code based on status
	exitCode := 0
	errorMsg := ""
	if statusCode >= 400 {
		exitCode = 1
		errorMsg = fmt.Sprintf("HTTP request returned error status: %d", statusCode)
	}

	return &domain.ExecutionResult{
		Output:   output,
		ExitCode: exitCode,
		Error:    errorMsg,
	}, nil
}
