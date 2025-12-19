package canvelete

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// WebhookEvent represents a webhook event
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// WebhookVerifyOptions are options for verifying webhooks
type WebhookVerifyOptions struct {
	Payload   []byte
	Signature string
	Secret    string
	Tolerance time.Duration // Timestamp tolerance
}

// DefaultWebhookTolerance is the default timestamp tolerance
const DefaultWebhookTolerance = 5 * time.Minute

// VerifyWebhookSignature verifies a webhook signature
func VerifyWebhookSignature(opts *WebhookVerifyOptions) bool {
	if opts.Tolerance == 0 {
		opts.Tolerance = DefaultWebhookTolerance
	}
	
	// Parse signature header (format: t=timestamp,v1=signature)
	parts := strings.Split(opts.Signature, ",")
	signatureParts := make(map[string]string)
	
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			signatureParts[kv[0]] = kv[1]
		}
	}
	
	timestamp := signatureParts["t"]
	v1Signature := signatureParts["v1"]
	
	if timestamp == "" || v1Signature == "" {
		return false
	}
	
	// Check timestamp tolerance
	timestampNum, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	
	eventTime := time.Unix(timestampNum, 0)
	if time.Since(eventTime).Abs() > opts.Tolerance {
		return false
	}
	
	// Compute expected signature
	signedPayload := fmt.Sprintf("%s.%s", timestamp, string(opts.Payload))
	mac := hmac.New(sha256.New, []byte(opts.Secret))
	mac.Write([]byte(signedPayload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	
	// Timing-safe comparison
	return hmac.Equal([]byte(v1Signature), []byte(expectedSignature))
}

// ParseWebhookPayload parses a webhook payload
func ParseWebhookPayload(payload []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}
	return &event, nil
}

// ConstructWebhookEvent constructs a webhook event from a verified payload
func ConstructWebhookEvent(payload []byte, signature, secret string) (*WebhookEvent, error) {
	opts := &WebhookVerifyOptions{
		Payload:   payload,
		Signature: signature,
		Secret:    secret,
	}
	
	if !VerifyWebhookSignature(opts) {
		return nil, fmt.Errorf("invalid webhook signature")
	}
	
	return ParseWebhookPayload(payload)
}

// GenerateWebhookSignature generates a webhook signature (for testing)
func GenerateWebhookSignature(payload []byte, secret string) string {
	timestamp := time.Now().Unix()
	signedPayload := fmt.Sprintf("%d.%s", timestamp, string(payload))
	
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	signature := hex.EncodeToString(mac.Sum(nil))
	
	return fmt.Sprintf("t=%d,v1=%s", timestamp, signature)
}