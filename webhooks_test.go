package foil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

type webhookSignatureFixture struct {
	Secret           string `json:"secret"`
	Timestamp        string `json:"timestamp"`
	ExpiredTimestamp string `json:"expired_timestamp"`
	NowSeconds       int64  `json:"now_seconds"`
	RawBody          string `json:"raw_body"`
	Signature        string `json:"signature"`
	InvalidSignature string `json:"invalid_signature"`
}

func (f webhookSignatureFixture) input() VerifyWebhookSignatureInput {
	return VerifyWebhookSignatureInput{
		Secret:     f.Secret,
		Timestamp:  f.Timestamp,
		RawBody:    f.RawBody,
		Signature:  f.Signature,
		NowSeconds: f.NowSeconds,
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	fixture := loadFixture[webhookSignatureFixture](t, "webhooks/signature.json")

	if !VerifyWebhookSignature(fixture.input()) {
		t.Fatal("expected valid signature fixture to verify")
	}

	tampered := fixture.input()
	tampered.Signature = fixture.InvalidSignature
	if VerifyWebhookSignature(tampered) {
		t.Fatal("expected invalid signature fixture to fail")
	}

	tampered = fixture.input()
	tampered.Signature = "short"
	if VerifyWebhookSignature(tampered) {
		t.Fatal("expected truncated signature to fail")
	}

	tampered = fixture.input()
	tampered.RawBody += " "
	if VerifyWebhookSignature(tampered) {
		t.Fatal("expected modified body to fail")
	}

	tampered = fixture.input()
	tampered.Secret = "whsec_other"
	if VerifyWebhookSignature(tampered) {
		t.Fatal("expected wrong secret to fail")
	}

	expired := fixture.input()
	expired.Timestamp = fixture.ExpiredTimestamp
	if VerifyWebhookSignature(expired) {
		t.Fatal("expected expired timestamp to fail")
	}

	malformed := fixture.input()
	malformed.Timestamp = "not-a-timestamp"
	if VerifyWebhookSignature(malformed) {
		t.Fatal("expected malformed timestamp to fail")
	}

	empty := fixture.input()
	empty.Secret = ""
	emptyKeyMAC := hmac.New(sha256.New, nil)
	emptyKeyMAC.Write([]byte(fixture.Timestamp + "." + fixture.RawBody))
	empty.Signature = hex.EncodeToString(emptyKeyMAC.Sum(nil))
	if VerifyWebhookSignature(empty) {
		t.Fatal("expected empty secret to fail")
	}

	relaxed := fixture.input()
	relaxed.NowSeconds = fixture.NowSeconds + 600
	relaxed.MaxAgeSeconds = 900
	if !VerifyWebhookSignature(relaxed) {
		t.Fatal("expected custom max age to be honored")
	}
}

func TestParseWebhookEvent(t *testing.T) {
	fixture := loadFixture[webhookSignatureFixture](t, "webhooks/signature.json")

	envelope, payload, err := ParseWebhookEvent([]byte(fixture.RawBody))
	if err != nil {
		t.Fatalf("parse webhook event: %v", err)
	}
	if envelope.Type != "session.result.persisted" || envelope.Object != "webhook_event" {
		t.Fatalf("unexpected envelope %#v", envelope)
	}
	data, ok := payload.(map[string]any)
	if !ok || data["object"] != "session_result" {
		t.Fatalf("unexpected payload %#v", payload)
	}

	envelope, _, err = ParseWebhookEvent([]byte(`{"id":"wevt_1","object":"webhook_event","type":"webhook.test","created":"2026-04-27T00:00:00.000Z","data":{}}`))
	if err != nil || envelope.Type != "webhook.test" {
		t.Fatalf("parse webhook.test event: %#v %v", envelope, err)
	}

	for name, body := range map[string]string{
		"unknown type": `{"id":"wevt_1","object":"webhook_event","type":"unknown.event","created":"2026-04-27T00:00:00.000Z","data":{}}`,
		"retired type": `{"id":"wevt_1","object":"webhook_event","type":"session.fingerprint.calculated","created":"2026-04-27T00:00:00.000Z","data":{}}`,
		"wrong object": `{"id":"wevt_1","object":"event","type":"webhook.test","created":"2026-04-27T00:00:00.000Z","data":{}}`,
		"missing id":   `{"object":"webhook_event","type":"webhook.test","created":"2026-04-27T00:00:00.000Z","data":{}}`,
		"missing type": `{"id":"wevt_1","object":"webhook_event","created":"2026-04-27T00:00:00.000Z","data":{}}`,
		"missing data": `{"id":"wevt_1","object":"webhook_event","type":"webhook.test","created":"2026-04-27T00:00:00.000Z"}`,
		"not json":     `nope`,
	} {
		if _, _, err := ParseWebhookEvent([]byte(body)); err == nil {
			t.Fatalf("expected %s to be rejected", name)
		}
	}
}

func TestVerifyAndParseWebhookEvent(t *testing.T) {
	fixture := loadFixture[webhookSignatureFixture](t, "webhooks/signature.json")

	envelope, _, err := VerifyAndParseWebhookEvent(fixture.input())
	if err != nil || envelope.Type != "session.result.persisted" {
		t.Fatalf("verify and parse webhook event: %#v err=%v", envelope, err)
	}

	invalid := fixture.input()
	invalid.Signature = fixture.InvalidSignature
	if _, _, err := VerifyAndParseWebhookEvent(invalid); err == nil {
		t.Fatal("expected invalid signature to be rejected")
	}
}
