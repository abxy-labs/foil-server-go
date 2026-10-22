package foil

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

const capTestSecret = "sk_live_cap_test"

func sealLegacyToken(t *testing.T, payload map[string]any, secretKey string) string {
	t.Helper()
	plaintext, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	var compressed bytes.Buffer
	writer := zlib.NewWriter(&compressed)
	if _, err := writer.Write(plaintext); err != nil {
		t.Fatalf("compress payload: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close compressor: %v", err)
	}
	block, err := aes.NewCipher(deriveTokenKey(secretKey))
	if err != nil {
		t.Fatalf("create cipher: %v", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("create GCM: %v", err)
	}
	nonce := make([]byte, tokenNonceBytes)
	if _, err := rand.Read(nonce); err != nil {
		t.Fatalf("read nonce: %v", err)
	}
	token := append([]byte{legacyTokenVersion}, nonce...)
	token = aead.Seal(token, nonce, compressed.Bytes(), nil)
	return base64.StdEncoding.EncodeToString(token)
}

func fixturePayloadWithPadding(t *testing.T, padding int) map[string]any {
	t.Helper()
	fixture := loadFixture[struct {
		Payload map[string]any `json:"payload"`
	}](t, "sealed-token/vector.v1.json")
	fixture.Payload["padding"] = strings.Repeat("a", padding)
	return fixture.Payload
}

func TestVerifyFoilTokenFixture(t *testing.T) {
	fixture := loadFixture[struct {
		Token      string         `json:"token"`
		SecretKey  string         `json:"secretKey"`
		SecretHash string         `json:"secretHash"`
		Payload    map[string]any `json:"payload"`
	}](t, "sealed-token/vector.v1.json")

	verified, err := VerifyFoilToken(fixture.Token, fixture.SecretKey)
	if err != nil {
		t.Fatalf("verify token with secret key: %v", err)
	}
	if verified.SessionID != fixture.Payload["session_id"] {
		t.Fatalf("unexpected session id %#v", verified.SessionID)
	}

	verified, err = VerifyFoilToken(fixture.Token, fixture.SecretHash)
	if err != nil {
		t.Fatalf("verify token with secret hash: %v", err)
	}
	decisionRaw, ok := fixture.Payload["decision"].(map[string]any)
	if !ok || verified.Decision.EventID != decisionRaw["event_id"] {
		t.Fatalf("unexpected decision event id %#v", verified.Decision.EventID)
	}
}

func TestVerifyFoilTokenMultiRecipientFixture(t *testing.T) {
	fixture := loadFixture[struct {
		Token        string         `json:"token"`
		SecretKeys   []string       `json:"secretKeys"`
		SecretHashes []string       `json:"secretHashes"`
		Payload      map[string]any `json:"payload"`
	}](t, "sealed-token/vector.v2.json")

	for _, secret := range append(fixture.SecretKeys, fixture.SecretHashes...) {
		verified, err := VerifyFoilToken(fixture.Token, secret)
		if err != nil {
			t.Fatalf("verify multi-recipient token: %v", err)
		}
		if verified.SessionID != fixture.Payload["session_id"] {
			t.Fatalf("unexpected session id %#v", verified.SessionID)
		}
	}
}

func TestSafeVerifyFoilTokenInvalidFixture(t *testing.T) {
	fixture := loadFixture[struct {
		Token string `json:"token"`
	}](t, "sealed-token/invalid.json")

	result := SafeVerifyFoilToken(fixture.Token, "sk_live_fixture_secret")
	if result.OK || result.Error == nil {
		t.Fatal("expected invalid token failure result")
	}
}

func TestVerifyFoilTokenAcceptsPayloadJustUnderSizeCap(t *testing.T) {
	payload := fixturePayloadWithPadding(t, 0)
	encoded, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	payload = fixturePayloadWithPadding(t, maxTokenPayloadBytes-len(encoded)-64)

	verified, err := VerifyFoilToken(sealLegacyToken(t, payload, capTestSecret), capTestSecret)
	if err != nil {
		t.Fatalf("verify token under the cap: %v", err)
	}
	if verified.Raw["padding"] != payload["padding"] {
		t.Fatal("padding did not round-trip")
	}
}

func TestVerifyFoilTokenRejectsPayloadPastSizeCap(t *testing.T) {
	payload := fixturePayloadWithPadding(t, 2*maxTokenPayloadBytes)

	_, err := VerifyFoilToken(sealLegacyToken(t, payload, capTestSecret), capTestSecret)
	var verificationErr *TokenVerificationError
	if !errors.As(err, &verificationErr) || verificationErr.Message != "Foil token payload is too large." {
		t.Fatalf("expected payload too large error, got %v", err)
	}
}

func TestVerifyFoilTokenMissingSecret(t *testing.T) {
	fixture := loadFixture[struct {
		Token string `json:"token"`
	}](t, "sealed-token/vector.v1.json")

	original := os.Getenv("FOIL_SECRET_KEY")
	defer os.Setenv("FOIL_SECRET_KEY", original)
	_ = os.Unsetenv("FOIL_SECRET_KEY")

	if _, err := VerifyFoilToken(fixture.Token, ""); err == nil {
		t.Fatal("expected missing secret error")
	}
}
