package jobs

import "github.com/meysam81/oneoff/internal/domain"

// RegisterJobTypes registers all built-in job types
func RegisterJobTypes(registry *domain.JobRegistry) {
	registry.Register("http", NewHTTPJob)
	registry.Register("shell", NewShellJob)
	registry.Register("docker", NewDockerJob)
}
