package foil

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"regexp"
)

const (
	legacyTokenVersion         = byte(0x01)
	multiRecipientTokenVersion = byte(0x02)
	tokenNonceBytes            = 12
	tokenTagBytes              = 16
	contentKeyBytes            = 32
	recipientIDBytes           = 32
	maxTokenRecipients         = 256
	v2HeaderBytes              = 1 + 2 + tokenNonceBytes + 4
	v2RecipientBytes           = recipientIDBytes + tokenNonceBytes + contentKeyBytes + tokenTagBytes
)

var (
	v2PayloadAADPrefix = []byte("foil-sealed-results-v2\x00payload\x00")
	v2WrapAADPrefix    = []byte("foil-sealed-results-v2\x00recipient\x00")
)

var hexSecretPattern = regexp.MustCompile(`\A[0-9a-fA-F]{64}\z`)

func resolveSecret(secretKey string) (string, error) {
	if secretKey != "" {
		return secretKey, nil
	}
	if envSecret := os.Getenv("FOIL_SECRET_KEY"); envSecret != "" {
		return envSecret, nil
	}
	return "", &ConfigurationError{
		Message: "Missing Foil secret key. Pass WithSecretKey or set FOIL_SECRET_KEY.",
	}
}

func normalizeSecretMaterial(secretKeyOrHash string) string {
	if hexSecretPattern.MatchString(secretKeyOrHash) {
		return string(bytes.ToLower([]byte(secretKeyOrHash)))
	}
	sum := sha256.Sum256([]byte(secretKeyOrHash))
	return hex.EncodeToString(sum[:])
}

func deriveTokenKey(secretKeyOrHash string) []byte {
	sum := sha256.Sum256([]byte(normalizeSecretMaterial(secretKeyOrHash) + "\x00sealed-results"))
	return sum[:]
}

func decryptTokenGCM(ciphertext, key, nonce, tag, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	sealed := append([]byte{}, ciphertext...)
	sealed = append(sealed, tag...)
	return aead.Open(nil, nonce, sealed, aad)
}

func decryptTokenPayload(raw []byte, secretKey string) ([]byte, error) {
	if raw[0] == legacyTokenVersion {
		return decryptTokenGCM(raw[13:len(raw)-tokenTagBytes], deriveTokenKey(secretKey), raw[1:13], raw[len(raw)-tokenTagBytes:], nil)
	}
	if raw[0] != multiRecipientTokenVersion {
		return nil, &TokenVerificationError{Message: "Unsupported Foil token version."}
	}
	if len(raw) < v2HeaderBytes+tokenTagBytes+v2RecipientBytes {
		return nil, &TokenVerificationError{Message: "Foil token is too short."}
	}
	recipientCount := int(binary.BigEndian.Uint16(raw[1:3]))
	if recipientCount < 1 || recipientCount > maxTokenRecipients {
		return nil, &TokenVerificationError{Message: "Foil token has an invalid recipient count."}
	}
	payloadLength := int(binary.BigEndian.Uint32(raw[15:19]))
	payloadStart := v2HeaderBytes
	payloadTagStart := payloadStart + payloadLength
	recipientsStart := payloadTagStart + tokenTagBytes
	if payloadLength < 1 || recipientsStart+recipientCount*v2RecipientBytes != len(raw) {
		return nil, &TokenVerificationError{Message: "Foil token has an invalid length."}
	}
	expectedIDSum := sha256.Sum256([]byte(normalizeSecretMaterial(secretKey) + "\x00sealed-results-recipient-id"))
	expectedID := expectedIDSum[:]
	recipientIDs := make([]byte, 0, recipientCount*recipientIDBytes)
	for index := 0; index < recipientCount; index++ {
		entryStart := recipientsStart + index*v2RecipientBytes
		recipientIDs = append(recipientIDs, raw[entryStart:entryStart+recipientIDBytes]...)
	}
	var contentKey []byte
	var err error
	for index := 0; index < recipientCount; index++ {
		entryStart := recipientsStart + index*v2RecipientBytes
		id := raw[entryStart : entryStart+recipientIDBytes]
		if subtle.ConstantTimeCompare(id, expectedID) != 1 {
			continue
		}
		nonceStart := entryStart + recipientIDBytes
		wrappedKeyStart := nonceStart + tokenNonceBytes
		tagStart := wrappedKeyStart + contentKeyBytes
		contentKey, err = decryptTokenGCM(
			raw[wrappedKeyStart:tagStart],
			deriveTokenKey(secretKey),
			raw[nonceStart:wrappedKeyStart],
			raw[tagStart:tagStart+tokenTagBytes],
			append(append([]byte{}, v2WrapAADPrefix...), id...),
		)
		if err != nil {
			return nil, err
		}
		break
	}
	if len(contentKey) != contentKeyBytes {
		return nil, &TokenVerificationError{Message: "Secret key is not a recipient of this Foil token."}
	}
	return decryptTokenGCM(
		raw[payloadStart:payloadTagStart],
		contentKey,
		raw[3:15],
		raw[payloadTagStart:recipientsStart],
		append(append(append([]byte{}, v2PayloadAADPrefix...), raw[:v2HeaderBytes]...), recipientIDs...),
	)
}

func VerifyFoilToken(sealedToken string, secretKey string) (*VerifiedFoilToken, error) {
	resolvedSecret, err := resolveSecret(secretKey)
	if err != nil {
		return nil, err
	}

	raw, err := base64.StdEncoding.DecodeString(sealedToken)
	if err != nil {
		return nil, &TokenVerificationError{Message: "Failed to verify Foil token.", Err: err}
	}
	if len(raw) < 29 {
		return nil, &TokenVerificationError{Message: "Foil token is too short."}
	}
	compressed, err := decryptTokenPayload(raw, resolvedSecret)
	if err != nil {
		return nil, &TokenVerificationError{Message: "Failed to verify Foil token.", Err: err}
	}

	reader, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, &TokenVerificationError{Message: "Failed to verify Foil token.", Err: err}
	}
	defer reader.Close()

	payloadBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, &TokenVerificationError{Message: "Failed to verify Foil token.", Err: err}
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, &TokenVerificationError{Message: "Failed to verify Foil token.", Err: err}
	}

	verified := &VerifiedFoilToken{}
	if err := json.Unmarshal(payloadBytes, verified); err != nil {
		return nil, &TokenVerificationError{Message: "Failed to verify Foil token.", Err: err}
	}
	verified.Raw = payload
	return verified, nil
}

func SafeVerifyFoilToken(sealedToken string, secretKey string) VerificationResult {
	verified, err := VerifyFoilToken(sealedToken, secretKey)
	if err != nil {
		return VerificationResult{OK: false, Error: err}
	}
	return VerificationResult{OK: true, Data: verified}
}
