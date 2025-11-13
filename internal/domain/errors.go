package domain

import "errors"

var (
	// Job errors
	ErrJobNotFound      = errors.New("job not found")
	ErrJobTypeNotFound  = errors.New("job type not found")
	ErrInvalidJobConfig = errors.New("invalid job configuration")
	ErrJobAlreadyExists = errors.New("job already exists")
	ErrJobNotScheduled  = errors.New("job is not in scheduled status")

	// Execution errors
	ErrExecutionNotFound = errors.New("execution not found")
	ErrExecutionFailed   = errors.New("execution failed")
	ErrExecutionTimeout  = errors.New("execution timeout")
	ErrExecutionCancelled = errors.New("execution cancelled")

	// Project errors
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectAlreadyExists = errors.New("project already exists")
	ErrCannotDeleteDefault  = errors.New("cannot delete default project")

	// Tag errors
	ErrTagNotFound      = errors.New("tag not found")
	ErrTagAlreadyExists = errors.New("tag already exists")

	// Chain errors
	ErrChainNotFound = errors.New("job chain not found")
	ErrInvalidChain  = errors.New("invalid job chain")

	// Validation errors
	ErrInvalidPriority     = errors.New("priority must be between 1 and 10")
	ErrInvalidScheduleTime = errors.New("schedule time must be in the future")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrMissingRequiredField = errors.New("missing required field")

	// System errors
	ErrDatabaseError    = errors.New("database error")
	ErrConfigError      = errors.New("configuration error")
	ErrWorkerPoolClosed = errors.New("worker pool is closed")
)
