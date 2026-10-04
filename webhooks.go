package foil

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type WebhookEndpoint struct {
	Object        string   `json:"object"`
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	URL           string   `json:"url"`
	Status        string   `json:"status"`
	EventTypes    []string `json:"event_types"`
	SigningSecret string   `json:"signing_secret,omitempty"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type WebhookDelivery struct {
	Object         string  `json:"object"`
	ID             string  `json:"id"`
	EventID        string  `json:"event_id"`
	EndpointID     string  `json:"endpoint_id"`
	EventType      string  `json:"event_type"`
	Status         string  `json:"status"`
	Attempts       int     `json:"attempts"`
	ResponseStatus *int    `json:"response_status"`
	ResponseBody   *string `json:"response_body"`
	Error          *string `json:"error"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type WebhookTest struct {
	Object         string           `json:"object"`
	EventID        string           `json:"event_id"`
	DeliveryIDs    []string         `json:"delivery_ids"`
	LatestDelivery *WebhookDelivery `json:"latest_delivery"`
}

type EventSubject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Event struct {
	Object            string            `json:"object"`
	ID                string            `json:"id"`
	Type              string            `json:"type"`
	Subject           EventSubject      `json:"subject"`
	Data              map[string]any    `json:"data"`
	WebhookDeliveries []WebhookDelivery `json:"webhook_deliveries"`
	CreatedAt         string            `json:"created_at"`
}

type WebhookEventEnvelope struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Type    string          `json:"type"`
	Created string          `json:"created"`
	Data    json.RawMessage `json:"data"`
}

type VerifyWebhookSignatureInput struct {
	Secret        string
	Timestamp     string
	RawBody       string
	Signature     string
	MaxAgeSeconds int64
	NowSeconds    int64
}

type CreateWebhookEndpointParams struct {
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	EventTypes []string `json:"event_types"`
}

type UpdateWebhookEndpointParams struct {
	Name       string   `json:"name,omitempty"`
	URL        string   `json:"url,omitempty"`
	Status     string   `json:"status,omitempty"`
	EventTypes []string `json:"event_types,omitempty"`
}

type EventListParams struct {
	EndpointID string
	Type       string
	Limit      int
}

type WebhooksService struct {
	client *Client
}

func (s *WebhooksService) ListEndpoints(ctx context.Context, organizationID string) (ListResult[WebhookEndpoint], error) {
	var envelope resourceListEnvelope[WebhookEndpoint]
	err := s.client.doJSON(ctx, http.MethodGet, "/v1/organizations/"+url.PathEscape(organizationID)+"/webhooks/endpoints", nil, nil, &envelope)
	if err != nil {
		return ListResult[WebhookEndpoint]{}, err
	}
	return normalizeList(envelope), nil
}

func (s *WebhooksService) CreateEndpoint(ctx context.Context, organizationID string, params CreateWebhookEndpointParams) (WebhookEndpoint, error) {
	var envelope resourceEnvelope[WebhookEndpoint]
	err := s.client.doJSON(ctx, http.MethodPost, "/v1/organizations/"+url.PathEscape(organizationID)+"/webhooks/endpoints", nil, params, &envelope)
	if err != nil {
		return WebhookEndpoint{}, err
	}
	return envelope.Data, nil
}

func (s *WebhooksService) UpdateEndpoint(ctx context.Context, organizationID string, endpointID string, params UpdateWebhookEndpointParams) (WebhookEndpoint, error) {
	var envelope resourceEnvelope[WebhookEndpoint]
	err := s.client.doJSON(ctx, http.MethodPatch, "/v1/organizations/"+url.PathEscape(organizationID)+"/webhooks/endpoints/"+url.PathEscape(endpointID), nil, params, &envelope)
	if err != nil {
		return WebhookEndpoint{}, err
	}
	return envelope.Data, nil
}

func (s *WebhooksService) DisableEndpoint(ctx context.Context, organizationID string, endpointID string) (WebhookEndpoint, error) {
	var envelope resourceEnvelope[WebhookEndpoint]
	err := s.client.doJSON(ctx, http.MethodDelete, "/v1/organizations/"+url.PathEscape(organizationID)+"/webhooks/endpoints/"+url.PathEscape(endpointID), nil, nil, &envelope)
	if err != nil {
		return WebhookEndpoint{}, err
	}
	return envelope.Data, nil
}

func (s *WebhooksService) RotateSecret(ctx context.Context, organizationID string, endpointID string) (WebhookEndpoint, error) {
	var envelope resourceEnvelope[WebhookEndpoint]
	err := s.client.doJSON(ctx, http.MethodPost, "/v1/organizations/"+url.PathEscape(organizationID)+"/webhooks/endpoints/"+url.PathEscape(endpointID)+"/rotations", nil, nil, &envelope)
	if err != nil {
		return WebhookEndpoint{}, err
	}
	return envelope.Data, nil
}

func (s *WebhooksService) SendTest(ctx context.Context, organizationID string, endpointID string) (WebhookTest, error) {
	var envelope resourceEnvelope[WebhookTest]
	err := s.client.doJSON(ctx, http.MethodPost, "/v1/organizations/"+url.PathEscape(organizationID)+"/webhooks/endpoints/"+url.PathEscape(endpointID)+"/test", nil, nil, &envelope)
	if err != nil {
		return WebhookTest{}, err
	}
	return envelope.Data, nil
}

func (s *WebhooksService) ListEvents(ctx context.Context, organizationID string, params EventListParams) (ListResult[Event], error) {
	var envelope resourceListEnvelope[Event]
	err := s.client.doJSON(ctx, http.MethodGet, "/v1/organizations/"+url.PathEscape(organizationID)+"/events", map[string]string{
		"endpoint_id": params.EndpointID,
		"type":        params.Type,
		"limit":       intToString(params.Limit),
	}, nil, &envelope)
	if err != nil {
		return ListResult[Event]{}, err
	}
	return normalizeList(envelope), nil
}

func (s *WebhooksService) RetrieveEvent(ctx context.Context, organizationID string, eventID string) (Event, error) {
	var envelope resourceEnvelope[Event]
	err := s.client.doJSON(ctx, http.MethodGet, "/v1/organizations/"+url.PathEscape(organizationID)+"/events/"+url.PathEscape(eventID), nil, nil, &envelope)
	if err != nil {
		return Event{}, err
	}
	return envelope.Data, nil
}

// VerifyWebhookSignature checks the X-Foil-Timestamp and X-Foil-Signature headers
// of a webhook delivery against the raw request body.
func VerifyWebhookSignature(input VerifyWebhookSignatureInput) bool {
	timestamp, err := strconv.ParseInt(input.Timestamp, 10, 64)
	if err != nil {
		return false
	}
	nowSeconds := input.NowSeconds
	if nowSeconds == 0 {
		nowSeconds = time.Now().Unix()
	}
	maxAgeSeconds := input.MaxAgeSeconds
	if maxAgeSeconds == 0 {
		maxAgeSeconds = 5 * 60
	}
	if absInt64(nowSeconds-timestamp) > maxAgeSeconds {
		return false
	}
	expected := hmac.New(sha256.New, []byte(input.Secret))
	expected.Write([]byte(input.Timestamp))
	expected.Write([]byte("."))
	expected.Write([]byte(input.RawBody))
	return subtle.ConstantTimeCompare(
		[]byte(fmt.Sprintf("%x", expected.Sum(nil))),
		[]byte(input.Signature),
	) == 1
}

// ParseWebhookEvent decodes a webhook event envelope and its data object.
func ParseWebhookEvent(rawBody []byte) (*WebhookEventEnvelope, any, error) {
	var envelope WebhookEventEnvelope
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return nil, nil, err
	}
	if envelope.Object != "webhook_event" {
		return nil, nil, errors.New("webhook event object must be webhook_event")
	}
	if envelope.ID == "" {
		return nil, nil, errors.New("webhook event id is required")
	}
	if envelope.Type == "" {
		return nil, nil, errors.New("webhook event type is required")
	}
	if envelope.Created == "" {
		return nil, nil, errors.New("webhook event created timestamp is required")
	}
	if len(envelope.Data) == 0 {
		return nil, nil, errors.New("webhook event data is required")
	}
	var payload map[string]any
	if err := json.Unmarshal(envelope.Data, &payload); err != nil {
		return nil, nil, err
	}
	return &envelope, payload, nil
}

// VerifyAndParseWebhookEvent verifies the delivery signature before parsing the event.
func VerifyAndParseWebhookEvent(input VerifyWebhookSignatureInput) (*WebhookEventEnvelope, any, error) {
	if !VerifyWebhookSignature(input) {
		return nil, nil, errors.New("invalid Foil webhook signature")
	}
	return ParseWebhookEvent([]byte(input.RawBody))
}

func absInt64(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}
