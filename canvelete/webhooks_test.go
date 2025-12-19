package canvelete

import (
	"testing"
	"time"
)

func TestGenerateWebhookSignature(t *testing.T) {
	payload := []byte(`{"event": "test"}`)
	secret := "test-webhook-secret"
	
	signature := GenerateWebhookSignature(payload, secret)
	
	// Check format: t=timestamp,v1=signature
	if len(signature) < 20 {
		t.Error("Signature too short")
	}
	
	if signature[:2] != "t=" {
		t.Error("Signature should start with 't='")
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	payload := []byte(`{"event": "test"}`)
	secret := "test-webhook-secret"
	
	// Generate a valid signature
	signature := GenerateWebhookSignature(payload, secret)
	
	// Verify it
	opts := &WebhookVerifyOptions{
		Payload:   payload,
		Signature: signature,
		Secret:    secret,
		Tolerance: 5 * time.Minute,
	}
	
	if !VerifyWebhookSignature(opts) {
		t.Error("Expected valid signature to verify")
	}
}

func TestVerifyWebhookSignatureInvalid(t *testing.T) {
	payload := []byte(`{"event": "test"}`)
	secret := "test-webhook-secret"
	
	opts := &WebhookVerifyOptions{
		Payload:   payload,
		Signature: "t=123,v1=invalid",
		Secret:    secret,
	}
	
	if VerifyWebhookSignature(opts) {
		t.Error("Expected invalid signature to fail verification")
	}
}

func TestVerifyWebhookSignatureMalformed(t *testing.T) {
	payload := []byte(`{"event": "test"}`)
	secret := "test-webhook-secret"
	
	tests := []struct {
		name      string
		signature string
	}{
		{"Empty", ""},
		{"No timestamp", "v1=abc123"},
		{"No signature", "t=123"},
		{"Invalid format", "invalid"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &WebhookVerifyOptions{
				Payload:   payload,
				Signature: tt.signature,
				Secret:    secret,
			}
			
			if VerifyWebhookSignature(opts) {
				t.Error("Expected malformed signature to fail verification")
			}
		})
	}
}

func TestParseWebhookPayload(t *testing.T) {
	payload := []byte(`{
		"id": "evt-123",
		"type": "render.completed",
		"timestamp": "2024-01-15T10:00:00Z",
		"data": {"designId": "design-123"}
	}`)
	
	event, err := ParseWebhookPayload(payload)
	if err != nil {
		t.Fatalf("ParseWebhookPayload failed: %v", err)
	}
	
	if event.ID != "evt-123" {
		t.Errorf("Expected ID 'evt-123', got '%s'", event.ID)
	}
	
	if event.Type != "render.completed" {
		t.Errorf("Expected type 'render.completed', got '%s'", event.Type)
	}
}

func TestParseWebhookPayloadInvalid(t *testing.T) {
	payload := []byte(`invalid json`)
	
	_, err := ParseWebhookPayload(payload)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestConstructWebhookEvent(t *testing.T) {
	payload := []byte(`{
		"id": "evt-123",
		"type": "render.completed",
		"timestamp": "2024-01-15T10:00:00Z",
		"data": {}
	}`)
	secret := "test-secret"
	
	signature := GenerateWebhookSignature(payload, secret)
	
	event, err := ConstructWebhookEvent(payload, signature, secret)
	if err != nil {
		t.Fatalf("ConstructWebhookEvent failed: %v", err)
	}
	
	if event.Type != "render.completed" {
		t.Errorf("Expected type 'render.completed', got '%s'", event.Type)
	}
}

func TestConstructWebhookEventInvalidSignature(t *testing.T) {
	payload := []byte(`{"event": "test"}`)
	
	_, err := ConstructWebhookEvent(payload, "t=123,v1=invalid", "secret")
	if err == nil {
		t.Error("Expected error for invalid signature")
	}
}