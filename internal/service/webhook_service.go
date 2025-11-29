package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/meysam81/oneoff/internal/domain"
	"github.com/meysam81/oneoff/internal/repository"
	"github.com/rs/zerolog/log"
)

// WebhookService handles webhook operations and delivery
type WebhookService struct {
	repo       repository.Repository
	httpClient *http.Client
	queue      chan *webhookDeliveryTask
	maxRetries int
	wg         sync.WaitGroup
	stopChan   chan struct{}
	workers    int
}

// webhookDeliveryTask represents a webhook to be delivered
type webhookDeliveryTask struct {
	webhook  *domain.Webhook
	delivery *domain.WebhookDelivery
	event    *domain.WebhookEvent
}

// NewWebhookService creates a new webhook service
func NewWebhookService(repo repository.Repository) *WebhookService {
	return &WebhookService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		queue:      make(chan *webhookDeliveryTask, 100),
		maxRetries: 5,
		stopChan:   make(chan struct{}),
		workers:    3,
	}
}

// Start starts the webhook delivery workers
func (s *WebhookService) Start(ctx context.Context) {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(ctx, i)
	}

	// Start retry processor
	s.wg.Add(1)
	go s.retryProcessor(ctx)

	log.Info().Int("workers", s.workers).Msg("Webhook service started")
}

// Stop stops the webhook service
func (s *WebhookService) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	log.Info().Msg("Webhook service stopped")
}

// worker processes webhook deliveries
func (s *WebhookService) worker(ctx context.Context, id int) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case task, ok := <-s.queue:
			if !ok {
				return
			}
			s.deliver(ctx, task)
		}
	}
}

// retryProcessor periodically checks for failed deliveries to retry
func (s *WebhookService) retryProcessor(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.processPendingRetries(ctx)
		}
	}
}

// processPendingRetries fetches and requeues pending deliveries
func (s *WebhookService) processPendingRetries(ctx context.Context) {
	deliveries, err := s.repo.GetPendingWebhookDeliveries(ctx, 50)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get pending webhook deliveries")
		return
	}

	for _, delivery := range deliveries {
		webhook, err := s.repo.GetWebhook(ctx, delivery.WebhookID)
		if err != nil {
			log.Error().Err(err).Str("webhook_id", delivery.WebhookID).Msg("Failed to get webhook for retry")
			continue
		}

		if !webhook.IsActive {
			// Mark as failed if webhook is disabled
			_ = s.repo.UpdateWebhookDelivery(ctx, delivery.ID, domain.WebhookDeliveryFailed, nil, "", "Webhook disabled", nil)
			continue
		}

		// Requeue for delivery
		select {
		case s.queue <- &webhookDeliveryTask{webhook: webhook, delivery: delivery}:
		default:
			log.Warn().Str("delivery_id", delivery.ID).Msg("Webhook queue full, will retry later")
		}
	}
}

// Dispatch dispatches a webhook event to all subscribed webhooks
func (s *WebhookService) Dispatch(ctx context.Context, event domain.WebhookEvent) {
	webhooks, err := s.repo.GetActiveWebhooksForEvent(ctx, event.Type)
	if err != nil {
		log.Error().Err(err).Str("event", string(event.Type)).Msg("Failed to get webhooks for event")
		return
	}

	if len(webhooks) == 0 {
		return
	}

	log.Debug().
		Str("event", string(event.Type)).
		Int("webhooks", len(webhooks)).
		Msg("Dispatching webhook event")

	// Build payload
	payload := domain.WebhookPayload{
		Event:     string(event.Type),
		Timestamp: event.Timestamp,
		Data: domain.WebhookData{
			Job:       event.Job,
			Execution: event.Execution,
		},
	}
	payloadJSON, _ := json.Marshal(payload)

	for _, webhook := range webhooks {
		delivery := &domain.WebhookDelivery{
			WebhookID: webhook.ID,
			EventType: string(event.Type),
			Payload:   string(payloadJSON),
			Status:    domain.WebhookDeliveryPending,
			Attempts:  0,
			CreatedAt: time.Now().UTC(),
		}

		// Create delivery record
		if err := s.repo.CreateWebhookDelivery(ctx, delivery); err != nil {
			log.Error().Err(err).Str("webhook_id", webhook.ID).Msg("Failed to create webhook delivery")
			continue
		}

		// Queue for delivery
		select {
		case s.queue <- &webhookDeliveryTask{
			webhook:  webhook,
			delivery: delivery,
			event:    &event,
		}:
		default:
			log.Warn().Str("webhook_id", webhook.ID).Msg("Webhook queue full, delivery will be retried")
		}
	}
}

// deliver attempts to deliver a webhook
func (s *WebhookService) deliver(ctx context.Context, task *webhookDeliveryTask) {
	webhook := task.webhook
	delivery := task.delivery

	// Increment attempt count
	_ = s.repo.IncrementDeliveryAttempts(ctx, delivery.ID)
	delivery.Attempts++

	// Build request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook.URL, strings.NewReader(delivery.Payload))
	if err != nil {
		s.handleDeliveryFailure(ctx, delivery, fmt.Sprintf("Failed to create request: %v", err))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "OneOff-Webhook/1.0")
	req.Header.Set("X-OneOff-Event", delivery.EventType)
	req.Header.Set("X-OneOff-Delivery", delivery.ID)

	// Add HMAC signature if secret is configured
	if webhook.Secret != "" {
		sig := computeHMAC(delivery.Payload, webhook.Secret)
		req.Header.Set("X-OneOff-Signature", "sha256="+sig)
	}

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.handleDeliveryFailure(ctx, delivery, fmt.Sprintf("Request failed: %v", err))
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Str("delivery_id", delivery.ID).Msg("Failed to close response body")
		}
	}()

	// Read response (limit to 1KB)
	responseBody := make([]byte, 1024)
	n, _ := resp.Body.Read(responseBody)
	responseBodyStr := string(responseBody[:n])

	// Check response status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Success
		statusCode := resp.StatusCode
		if err := s.repo.UpdateWebhookDelivery(ctx, delivery.ID, domain.WebhookDeliverySuccess, &statusCode, responseBodyStr, "", nil); err != nil {
			log.Error().Err(err).Str("delivery_id", delivery.ID).Msg("Failed to update delivery status")
		}
		log.Debug().
			Str("webhook_id", webhook.ID).
			Str("delivery_id", delivery.ID).
			Int("status_code", resp.StatusCode).
			Msg("Webhook delivered successfully")
	} else {
		// Failure
		statusCode := resp.StatusCode
		errMsg := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, responseBodyStr)
		s.handleDeliveryFailureWithResponse(ctx, delivery, errMsg, &statusCode, responseBodyStr)
	}
}

// handleDeliveryFailure handles a failed delivery
func (s *WebhookService) handleDeliveryFailure(ctx context.Context, delivery *domain.WebhookDelivery, errMsg string) {
	s.handleDeliveryFailureWithResponse(ctx, delivery, errMsg, nil, "")
}

// handleDeliveryFailureWithResponse handles a failed delivery with response details
func (s *WebhookService) handleDeliveryFailureWithResponse(ctx context.Context, delivery *domain.WebhookDelivery, errMsg string, statusCode *int, responseBody string) {
	if delivery.Attempts >= s.maxRetries {
		// Max retries exceeded, mark as failed
		if err := s.repo.UpdateWebhookDelivery(ctx, delivery.ID, domain.WebhookDeliveryFailed, statusCode, responseBody, errMsg, nil); err != nil {
			log.Error().Err(err).Str("delivery_id", delivery.ID).Msg("Failed to update delivery status")
		}
		log.Warn().
			Str("delivery_id", delivery.ID).
			Int("attempts", delivery.Attempts).
			Str("error", errMsg).
			Msg("Webhook delivery failed after max retries")
		return
	}

	// Schedule retry with exponential backoff
	// Cap the exponent to prevent overflow (max 30, i.e., ~17 minutes, but will be capped to 5 minutes below)
	maxExponent := 30
	exponent := delivery.Attempts
	if exponent > maxExponent {
		exponent = maxExponent
	}
	backoff := time.Duration(1<<uint(exponent)) * time.Second
	if backoff > 5*time.Minute {
		backoff = 5 * time.Minute
	}
	nextRetry := time.Now().UTC().Add(backoff)

	if err := s.repo.UpdateWebhookDelivery(ctx, delivery.ID, domain.WebhookDeliveryPending, statusCode, responseBody, errMsg, &nextRetry); err != nil {
		log.Error().Err(err).Str("delivery_id", delivery.ID).Msg("Failed to schedule retry")
	}

	log.Debug().
		Str("delivery_id", delivery.ID).
		Int("attempts", delivery.Attempts).
		Time("next_retry", nextRetry).
		Str("error", errMsg).
		Msg("Webhook delivery failed, scheduled retry")
}

// CreateWebhook creates a new webhook
func (s *WebhookService) CreateWebhook(ctx context.Context, req domain.CreateWebhookRequest) (*domain.Webhook, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}
	if req.Events == "" {
		return nil, fmt.Errorf("at least one event is required")
	}

	// Validate events
	events := strings.Split(req.Events, ",")
	validEvents := make(map[string]bool)
	for _, e := range domain.AllWebhookEvents() {
		validEvents[string(e)] = true
	}
	for _, e := range events {
		if !validEvents[strings.TrimSpace(e)] {
			return nil, fmt.Errorf("invalid event type: %s", e)
		}
	}

	webhook := &domain.Webhook{
		Name:      req.Name,
		URL:       req.URL,
		Secret:    req.Secret,
		Events:    req.Events,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreateWebhook(ctx, webhook); err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return webhook, nil
}

// GetWebhook retrieves a webhook by ID
func (s *WebhookService) GetWebhook(ctx context.Context, id string) (*domain.Webhook, error) {
	return s.repo.GetWebhook(ctx, id)
}

// ListWebhooks returns all webhooks
func (s *WebhookService) ListWebhooks(ctx context.Context, filter domain.WebhookFilter) ([]*domain.Webhook, error) {
	return s.repo.ListWebhooks(ctx, filter)
}

// UpdateWebhook updates a webhook
func (s *WebhookService) UpdateWebhook(ctx context.Context, id string, req domain.UpdateWebhookRequest) (*domain.Webhook, error) {
	// Validate events if provided
	if req.Events != nil {
		events := strings.Split(*req.Events, ",")
		validEvents := make(map[string]bool)
		for _, e := range domain.AllWebhookEvents() {
			validEvents[string(e)] = true
		}
		for _, e := range events {
			if !validEvents[strings.TrimSpace(e)] {
				return nil, fmt.Errorf("invalid event type: %s", e)
			}
		}
	}

	if err := s.repo.UpdateWebhook(ctx, id, req); err != nil {
		return nil, err
	}

	return s.repo.GetWebhook(ctx, id)
}

// DeleteWebhook deletes a webhook
func (s *WebhookService) DeleteWebhook(ctx context.Context, id string) error {
	return s.repo.DeleteWebhook(ctx, id)
}

// ListDeliveries returns webhook deliveries
func (s *WebhookService) ListDeliveries(ctx context.Context, filter domain.WebhookDeliveryFilter) ([]*domain.WebhookDelivery, error) {
	return s.repo.ListWebhookDeliveries(ctx, filter)
}

// TestWebhook sends a test event to a webhook
func (s *WebhookService) TestWebhook(ctx context.Context, id string) error {
	webhook, err := s.repo.GetWebhook(ctx, id)
	if err != nil {
		return err
	}

	testEvent := domain.WebhookEvent{
		Type:      "test",
		Timestamp: time.Now().UTC(),
	}

	payload := domain.WebhookPayload{
		Event:     "test",
		Timestamp: testEvent.Timestamp,
		Data:      domain.WebhookData{},
	}
	payloadJSON, _ := json.Marshal(payload)

	delivery := &domain.WebhookDelivery{
		WebhookID: webhook.ID,
		EventType: "test",
		Payload:   string(payloadJSON),
		Status:    domain.WebhookDeliveryPending,
		Attempts:  0,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreateWebhookDelivery(ctx, delivery); err != nil {
		return fmt.Errorf("failed to create test delivery: %w", err)
	}

	// Deliver synchronously for test
	s.deliver(ctx, &webhookDeliveryTask{
		webhook:  webhook,
		delivery: delivery,
		event:    &testEvent,
	})

	return nil
}

// computeHMAC computes HMAC-SHA256 signature
func computeHMAC(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
