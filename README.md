# Foil Go Library

![Preview](https://img.shields.io/badge/status-preview-111827)
![Go 1.22+](https://img.shields.io/badge/go-1.22%2B-00ADD8?logo=go&logoColor=white)
![License: MIT](https://img.shields.io/badge/license-MIT-0f766e.svg)

The Foil Go library provides convenient access to the Foil API from Go services and applications. It includes a context-aware client for Sessions, visitor fingerprints, Organizations, Organization API key management, webhook endpoints, and sealed token verification.

The library also provides:

- a fast configuration path using `FOIL_SECRET_KEY`
- iterator-style helpers for cursor-based pagination
- structured API errors and built-in sealed token verification
- webhook endpoint management, test sends, event delivery history, and webhook signature verification

## Documentation

See the [Foil docs](https://usefoil.com/docs) and [API reference](https://usefoil.com/docs/api-reference/introduction).

## Installation

You don't need this source code unless you want to modify the module. If you just want to use it, run:

```bash
go get github.com/abxy-labs/foil-server-go
```

## Requirements

- Go 1.22+

## Usage

Use `FOIL_SECRET_KEY` or `WithSecretKey(...)`:

```go
package main

import (
  "context"
  "log"

  foil "github.com/abxy-labs/foil-server-go"
)

func main() {
  client, err := foil.NewClient(foil.WithSecretKey("sk_live_..."))
  if err != nil {
    log.Fatal(err)
  }

  page, err := client.Sessions.List(context.Background(), foil.SessionListParams{
    Verdict: "bot",
    Limit:   25,
  })
  if err != nil {
    log.Fatal(err)
  }

  session, err := client.Sessions.Get(context.Background(), "sid_0123456789abcdefghjkmnpqrs")
  if err != nil {
    log.Fatal(err)
  }

  if _, err := client.Sessions.AttachClientUser(context.Background(), session.ID, "user_123"); err != nil {
    log.Fatal(err)
  }
  if _, err := client.Sessions.ClearClientUser(context.Background(), session.ID); err != nil {
    log.Fatal(err)
  }

  log.Println(page.Items[0].ID, session.Decision.AutomationStatus, session.Highlights[0].Summary)
}
```

### Sealed token verification

```go
result := foil.SafeVerifyFoilToken(sealedToken, "sk_live_...")
if !result.OK {
  log.Fatal(result.Error)
}

log.Println(result.Data.Decision.Verdict, result.Data.Decision.RiskScore)
```

### Pagination

```go
err := client.Sessions.Iter(context.Background(), foil.SessionListParams{Search: "signup"}, func(session foil.SessionSummary) error {
  log.Println(session.ID, session.LatestDecision.Verdict)
  return nil
})
if err != nil {
  log.Fatal(err)
}
```

### Visitor fingerprints

```go
fingerprint, err := client.Fingerprints.Get(context.Background(), "vid_0123456789abcdefghjkmnpqrs")
if err != nil {
  log.Fatal(err)
}

log.Println(fingerprint.ID)
```

### Organizations

```go
organization, err := client.Organizations.Get(context.Background(), "org_0123456789abcdefghjkmnpqrs")
if err != nil {
  log.Fatal(err)
}

updated, err := client.Organizations.Update(context.Background(), "org_0123456789abcdefghjkmnpqrs", foil.UpdateOrganizationParams{
  Name: "New Name",
})
if err != nil {
  log.Fatal(err)
}

_, _ = organization, updated
```

### Organization API keys

```go
created, err := client.Organizations.APIKeys.Create(
  context.Background(),
  "org_0123456789abcdefghjkmnpqrs",
  foil.CreateAPIKeyParams{Name: "Production", Type: "secret", Environment: "live"},
)
if err != nil {
  log.Fatal(err)
}

_, err = client.Organizations.APIKeys.Revoke(context.Background(), "org_0123456789abcdefghjkmnpqrs", created.ID)
if err != nil {
  log.Fatal(err)
}
```

### Webhooks

```go
endpoint, err := client.Webhooks.CreateEndpoint(context.Background(), "org_0123456789abcdefghjkmnpqrs", foil.CreateWebhookEndpointParams{
  Name:       "Production alerts",
  URL:        "https://example.com/foil/webhook",
  EventTypes: []string{"session.result.persisted"},
})
if err != nil {
  log.Fatal(err)
}

events, err := client.Webhooks.ListEvents(context.Background(), "org_0123456789abcdefghjkmnpqrs", foil.EventListParams{
  EndpointID: endpoint.ID,
  Type:       "session.result.persisted",
})
if err != nil {
  log.Fatal(err)
}

log.Println(events.Items[0].WebhookDeliveries[0].Status)
```

#### Verifying webhook deliveries

Every webhook delivery is signed with your endpoint's signing secret. Verify the `X-Foil-Timestamp` and `X-Foil-Signature` headers against the raw request body before trusting the payload:

```go
rawBody, err := io.ReadAll(request.Body)
if err != nil {
  log.Fatal(err)
}

input := foil.VerifyWebhookSignatureInput{
  Secret:    os.Getenv("FOIL_WEBHOOK_SECRET"),
  Timestamp: request.Header.Get("X-Foil-Timestamp"),
  RawBody:   string(rawBody),
  Signature: request.Header.Get("X-Foil-Signature"),
}

if !foil.VerifyWebhookSignature(input) {
  http.Error(writer, "invalid signature", http.StatusUnauthorized)
  return
}

// Verify and parse in one step.
envelope, data, err := foil.VerifyAndParseWebhookEvent(input)
if err != nil {
  log.Fatal(err)
}

if envelope.Type == "session.result.persisted" {
  log.Println(data)
}

// Parse a payload you have already verified.
envelope, data, err = foil.ParseWebhookEvent(rawBody)
```

Signatures older than five minutes are rejected by default. Set `MaxAgeSeconds` to change the tolerance.

### Error handling

```go
_, err := client.Sessions.List(context.Background(), foil.SessionListParams{Limit: 999})
if apiErr, ok := err.(*foil.APIError); ok {
  log.Println(apiErr.Status, apiErr.Code, apiErr.Message)
}
```

## Support

If you need help integrating Foil, start with [usefoil.com/docs](https://usefoil.com/docs).
