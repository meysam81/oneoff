package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Collector is the interface for metrics collection
type Collector interface {
	// Job metrics
	IncJobsTotal(jobType, status string)
	ObserveJobDuration(jobType string, duration time.Duration)

	// Worker metrics
	SetActiveWorkers(count int)
	SetTotalWorkers(count int)
	SetQueuedJobs(count int)

	// Request metrics
	IncRequestsTotal(method, path string, status int)
	ObserveRequestDuration(method, path string, duration time.Duration)

	// API key metrics
	IncAPIKeyValidations(valid bool)

	// Webhook metrics
	IncWebhookDeliveries(success bool)

	// Handler returns an HTTP handler for the /metrics endpoint
	Handler() http.Handler
}

// PrometheusCollector implements Prometheus-compatible metrics collection
type PrometheusCollector struct {
	mu sync.RWMutex

	// Job metrics
	jobsTotal    map[string]*int64  // key: type:status
	jobDurations map[string]*bucket // key: type

	// Worker metrics
	activeWorkers int64
	totalWorkers  int64
	queuedJobs    int64

	// Request metrics
	requestsTotal    map[string]*int64  // key: method:path:status
	requestDurations map[string]*bucket // key: method:path

	// API key metrics
	apiKeyValidationsSuccess int64
	apiKeyValidationsFailed  int64

	// Webhook metrics
	webhookDeliveriesSuccess int64
	webhookDeliveriesFailed  int64

	// Start time for uptime
	startTime time.Time
}

// bucket stores histogram data
type bucket struct {
	count int64
	sum   float64
	mu    sync.Mutex
}

// NewCollector creates a new metrics collector
func NewCollector() *PrometheusCollector {
	return &PrometheusCollector{
		jobsTotal:        make(map[string]*int64),
		jobDurations:     make(map[string]*bucket),
		requestsTotal:    make(map[string]*int64),
		requestDurations: make(map[string]*bucket),
		startTime:        time.Now(),
	}
}

// IncJobsTotal increments the jobs total counter
func (c *PrometheusCollector) IncJobsTotal(jobType, status string) {
	key := fmt.Sprintf("%s:%s", jobType, status)
	c.mu.Lock()
	if c.jobsTotal[key] == nil {
		c.jobsTotal[key] = new(int64)
	}
	ptr := c.jobsTotal[key]
	c.mu.Unlock()
	atomic.AddInt64(ptr, 1)
}

// ObserveJobDuration records a job duration
func (c *PrometheusCollector) ObserveJobDuration(jobType string, duration time.Duration) {
	c.mu.Lock()
	if c.jobDurations[jobType] == nil {
		c.jobDurations[jobType] = &bucket{}
	}
	b := c.jobDurations[jobType]
	c.mu.Unlock()

	b.mu.Lock()
	b.count++
	b.sum += duration.Seconds()
	b.mu.Unlock()
}

// SetActiveWorkers sets the current number of active workers
func (c *PrometheusCollector) SetActiveWorkers(count int) {
	atomic.StoreInt64(&c.activeWorkers, int64(count))
}

// SetTotalWorkers sets the total number of workers
func (c *PrometheusCollector) SetTotalWorkers(count int) {
	atomic.StoreInt64(&c.totalWorkers, int64(count))
}

// SetQueuedJobs sets the number of queued jobs
func (c *PrometheusCollector) SetQueuedJobs(count int) {
	atomic.StoreInt64(&c.queuedJobs, int64(count))
}

// IncRequestsTotal increments the requests total counter
func (c *PrometheusCollector) IncRequestsTotal(method, path string, status int) {
	key := fmt.Sprintf("%s:%s:%d", method, path, status)
	c.mu.Lock()
	if c.requestsTotal[key] == nil {
		c.requestsTotal[key] = new(int64)
	}
	ptr := c.requestsTotal[key]
	c.mu.Unlock()
	atomic.AddInt64(ptr, 1)
}

// ObserveRequestDuration records a request duration
func (c *PrometheusCollector) ObserveRequestDuration(method, path string, duration time.Duration) {
	key := fmt.Sprintf("%s:%s", method, path)
	c.mu.Lock()
	if c.requestDurations[key] == nil {
		c.requestDurations[key] = &bucket{}
	}
	b := c.requestDurations[key]
	c.mu.Unlock()

	b.mu.Lock()
	b.count++
	b.sum += duration.Seconds()
	b.mu.Unlock()
}

// IncAPIKeyValidations increments API key validation counter
func (c *PrometheusCollector) IncAPIKeyValidations(valid bool) {
	if valid {
		atomic.AddInt64(&c.apiKeyValidationsSuccess, 1)
	} else {
		atomic.AddInt64(&c.apiKeyValidationsFailed, 1)
	}
}

// IncWebhookDeliveries increments webhook delivery counter
func (c *PrometheusCollector) IncWebhookDeliveries(success bool) {
	if success {
		atomic.AddInt64(&c.webhookDeliveriesSuccess, 1)
	} else {
		atomic.AddInt64(&c.webhookDeliveriesFailed, 1)
	}
}

// Handler returns an HTTP handler for /metrics
func (c *PrometheusCollector) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		var sb strings.Builder

		// Write header comment
		sb.WriteString("# OneOff Metrics\n\n")

		// Uptime
		uptime := time.Since(c.startTime).Seconds()
		sb.WriteString("# HELP oneoff_uptime_seconds Time since service started\n")
		sb.WriteString("# TYPE oneoff_uptime_seconds gauge\n")
		sb.WriteString(fmt.Sprintf("oneoff_uptime_seconds %f\n\n", uptime))

		// Worker metrics
		sb.WriteString("# HELP oneoff_workers_total Total number of workers\n")
		sb.WriteString("# TYPE oneoff_workers_total gauge\n")
		sb.WriteString(fmt.Sprintf("oneoff_workers_total %d\n\n", atomic.LoadInt64(&c.totalWorkers)))

		sb.WriteString("# HELP oneoff_workers_active Number of currently active workers\n")
		sb.WriteString("# TYPE oneoff_workers_active gauge\n")
		sb.WriteString(fmt.Sprintf("oneoff_workers_active %d\n\n", atomic.LoadInt64(&c.activeWorkers)))

		sb.WriteString("# HELP oneoff_jobs_queued Number of jobs in queue\n")
		sb.WriteString("# TYPE oneoff_jobs_queued gauge\n")
		sb.WriteString(fmt.Sprintf("oneoff_jobs_queued %d\n\n", atomic.LoadInt64(&c.queuedJobs)))

		// Job totals
		sb.WriteString("# HELP oneoff_jobs_total Total number of jobs processed\n")
		sb.WriteString("# TYPE oneoff_jobs_total counter\n")
		c.mu.RLock()
		jobKeys := make([]string, 0, len(c.jobsTotal))
		for k := range c.jobsTotal {
			jobKeys = append(jobKeys, k)
		}
		sort.Strings(jobKeys)
		for _, key := range jobKeys {
			parts := strings.SplitN(key, ":", 2)
			if len(parts) == 2 {
				sb.WriteString(fmt.Sprintf("oneoff_jobs_total{type=\"%s\",status=\"%s\"} %d\n",
					parts[0], parts[1], atomic.LoadInt64(c.jobsTotal[key])))
			}
		}
		c.mu.RUnlock()
		sb.WriteString("\n")

		// Job durations
		sb.WriteString("# HELP oneoff_job_duration_seconds Job execution duration in seconds\n")
		sb.WriteString("# TYPE oneoff_job_duration_seconds summary\n")
		c.mu.RLock()
		durationKeys := make([]string, 0, len(c.jobDurations))
		for k := range c.jobDurations {
			durationKeys = append(durationKeys, k)
		}
		sort.Strings(durationKeys)
		for _, jobType := range durationKeys {
			b := c.jobDurations[jobType]
			b.mu.Lock()
			count := b.count
			sum := b.sum
			b.mu.Unlock()
			sb.WriteString(fmt.Sprintf("oneoff_job_duration_seconds_count{type=\"%s\"} %d\n", jobType, count))
			sb.WriteString(fmt.Sprintf("oneoff_job_duration_seconds_sum{type=\"%s\"} %f\n", jobType, sum))
		}
		c.mu.RUnlock()
		sb.WriteString("\n")

		// HTTP request totals
		sb.WriteString("# HELP oneoff_http_requests_total Total number of HTTP requests\n")
		sb.WriteString("# TYPE oneoff_http_requests_total counter\n")
		c.mu.RLock()
		reqKeys := make([]string, 0, len(c.requestsTotal))
		for k := range c.requestsTotal {
			reqKeys = append(reqKeys, k)
		}
		sort.Strings(reqKeys)
		for _, key := range reqKeys {
			parts := strings.SplitN(key, ":", 3)
			if len(parts) == 3 {
				sb.WriteString(fmt.Sprintf("oneoff_http_requests_total{method=\"%s\",path=\"%s\",status=\"%s\"} %d\n",
					parts[0], parts[1], parts[2], atomic.LoadInt64(c.requestsTotal[key])))
			}
		}
		c.mu.RUnlock()
		sb.WriteString("\n")

		// HTTP request durations
		sb.WriteString("# HELP oneoff_http_request_duration_seconds HTTP request duration in seconds\n")
		sb.WriteString("# TYPE oneoff_http_request_duration_seconds summary\n")
		c.mu.RLock()
		reqDurKeys := make([]string, 0, len(c.requestDurations))
		for k := range c.requestDurations {
			reqDurKeys = append(reqDurKeys, k)
		}
		sort.Strings(reqDurKeys)
		for _, key := range reqDurKeys {
			parts := strings.SplitN(key, ":", 2)
			if len(parts) == 2 {
				b := c.requestDurations[key]
				b.mu.Lock()
				count := b.count
				sum := b.sum
				b.mu.Unlock()
				sb.WriteString(fmt.Sprintf("oneoff_http_request_duration_seconds_count{method=\"%s\",path=\"%s\"} %d\n",
					parts[0], parts[1], count))
				sb.WriteString(fmt.Sprintf("oneoff_http_request_duration_seconds_sum{method=\"%s\",path=\"%s\"} %f\n",
					parts[0], parts[1], sum))
			}
		}
		c.mu.RUnlock()
		sb.WriteString("\n")

		// API key validations
		sb.WriteString("# HELP oneoff_apikey_validations_total Total number of API key validations\n")
		sb.WriteString("# TYPE oneoff_apikey_validations_total counter\n")
		sb.WriteString(fmt.Sprintf("oneoff_apikey_validations_total{result=\"success\"} %d\n",
			atomic.LoadInt64(&c.apiKeyValidationsSuccess)))
		sb.WriteString(fmt.Sprintf("oneoff_apikey_validations_total{result=\"failed\"} %d\n\n",
			atomic.LoadInt64(&c.apiKeyValidationsFailed)))

		// Webhook deliveries
		sb.WriteString("# HELP oneoff_webhook_deliveries_total Total number of webhook deliveries\n")
		sb.WriteString("# TYPE oneoff_webhook_deliveries_total counter\n")
		sb.WriteString(fmt.Sprintf("oneoff_webhook_deliveries_total{result=\"success\"} %d\n",
			atomic.LoadInt64(&c.webhookDeliveriesSuccess)))
		sb.WriteString(fmt.Sprintf("oneoff_webhook_deliveries_total{result=\"failed\"} %d\n",
			atomic.LoadInt64(&c.webhookDeliveriesFailed)))

		w.Write([]byte(sb.String())) //nolint:errcheck // HTTP write errors are connection issues
	})
}

// NoopCollector is a no-op metrics collector for when metrics are disabled
type NoopCollector struct{}

func (n *NoopCollector) IncJobsTotal(jobType, status string)                                {}
func (n *NoopCollector) ObserveJobDuration(jobType string, duration time.Duration)          {}
func (n *NoopCollector) SetActiveWorkers(count int)                                         {}
func (n *NoopCollector) SetTotalWorkers(count int)                                          {}
func (n *NoopCollector) SetQueuedJobs(count int)                                            {}
func (n *NoopCollector) IncRequestsTotal(method, path string, status int)                   {}
func (n *NoopCollector) ObserveRequestDuration(method, path string, duration time.Duration) {}
func (n *NoopCollector) IncAPIKeyValidations(valid bool)                                    {}
func (n *NoopCollector) IncWebhookDeliveries(success bool)                                  {}
func (n *NoopCollector) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}
