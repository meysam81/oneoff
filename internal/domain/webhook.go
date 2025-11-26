package domain

import (
	"encoding/json"
	"time"
)

// WebhookEventType represents types of webhook events
type WebhookEventType string

const (
	WebhookEventJobCreated   WebhookEventType = "job.created"
	WebhookEventJobStarted   WebhookEventType = "job.started"
	WebhookEventJobCompleted WebhookEventType = "job.completed"
	WebhookEventJobFailed    WebhookEventType = "job.failed"
	WebhookEventJobCancelled WebhookEventType = "job.cancelled"
)

// AllWebhookEvents returns all available webhook event types
func AllWebhookEvents() []WebhookEventType {
	return []WebhookEventType{
		WebhookEventJobCreated,
		WebhookEventJobStarted,
		WebhookEventJobCompleted,
		WebhookEventJobFailed,
		WebhookEventJobCancelled,
	}
}

// WebhookDeliveryStatus represents the delivery status
type WebhookDeliveryStatus string

const (
	WebhookDeliveryPending WebhookDeliveryStatus = "pending"
	WebhookDeliverySuccess WebhookDeliveryStatus = "success"
	WebhookDeliveryFailed  WebhookDeliveryStatus = "failed"
)

// Webhook represents a webhook configuration
type Webhook struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Secret    string    `json:"secret,omitempty"` // For HMAC signing
	Events    string    `json:"events"`           // Comma-separated event types
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID           string                `json:"id"`
	WebhookID    string                `json:"webhook_id"`
	EventType    string                `json:"event_type"`
	Payload      string                `json:"payload"`
	Status       WebhookDeliveryStatus `json:"status"`
	ResponseCode *int                  `json:"response_code,omitempty"`
	ResponseBody string                `json:"response_body,omitempty"`
	Error        string                `json:"error,omitempty"`
	Attempts     int                   `json:"attempts"`
	NextRetryAt  *time.Time            `json:"next_retry_at,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
}

// CreateWebhookRequest represents a request to create a webhook
type CreateWebhookRequest struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Secret string `json:"secret,omitempty"`
	Events string `json:"events"` // Comma-separated: job.completed,job.failed
}

// UpdateWebhookRequest represents a request to update a webhook
type UpdateWebhookRequest struct {
	Name     *string `json:"name,omitempty"`
	URL      *string `json:"url,omitempty"`
	Secret   *string `json:"secret,omitempty"`
	Events   *string `json:"events,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// WebhookEvent represents an event to be delivered
type WebhookEvent struct {
	Type      WebhookEventType `json:"type"`
	Timestamp time.Time        `json:"timestamp"`
	Job       *Job             `json:"job,omitempty"`
	Execution *JobExecution    `json:"execution,omitempty"`
}

// ToJSON serializes the event to JSON
func (e *WebhookEvent) ToJSON() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// WebhookPayload represents the full webhook payload
type WebhookPayload struct {
	Event     string        `json:"event"`
	Timestamp time.Time     `json:"timestamp"`
	Data      WebhookData   `json:"data"`
}

// WebhookData contains the event data
type WebhookData struct {
	Job       *Job          `json:"job,omitempty"`
	Execution *JobExecution `json:"execution,omitempty"`
}

// HasEvent checks if webhook is subscribed to the event
func (w *Webhook) HasEvent(event WebhookEventType) bool {
	if w.Events == "" {
		return false
	}
	eventList := splitScopes(w.Events) // Reuse the scope splitting logic
	for _, e := range eventList {
		if e == string(event) {
			return true
		}
	}
	return false
}

// WebhookFilter represents filters for listing webhooks
type WebhookFilter struct {
	IsActive *bool
	Limit    int
	Offset   int
}

// WebhookDeliveryFilter represents filters for listing deliveries
type WebhookDeliveryFilter struct {
	WebhookID string
	Status    WebhookDeliveryStatus
	Limit     int
	Offset    int
}
