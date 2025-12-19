package canvelete

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")
	
	if client.APIKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", client.APIKey)
	}
	
	if client.BaseURL != defaultBaseURL {
		t.Errorf("Expected base URL '%s', got '%s'", defaultBaseURL, client.BaseURL)
	}
	
	if client.Designs == nil {
		t.Error("Expected Designs service to be initialized")
	}
	
	if client.Canvas == nil {
		t.Error("Expected Canvas service to be initialized")
	}
	
	if client.Assets == nil {
		t.Error("Expected Assets service to be initialized")
	}
	
	if client.Usage == nil {
		t.Error("Expected Usage service to be initialized")
	}
	
	if client.Billing == nil {
		t.Error("Expected Billing service to be initialized")
	}
}

func TestClientWithOptions(t *testing.T) {
	customURL := "https://custom.canvelete.com"
	customTimeout := 60 * time.Second
	
	client := NewClient("test-key",
		WithBaseURL(customURL),
		WithTimeout(customTimeout),
	)
	
	if client.BaseURL != customURL {
		t.Errorf("Expected base URL '%s', got '%s'", customURL, client.BaseURL)
	}
	
	if client.HTTPClient.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.HTTPClient.Timeout)
	}
}

func TestClientRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-api-key" {
			t.Errorf("Expected Authorization 'Bearer test-api-key', got '%s'", auth)
		}
		
		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}
		
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{"id": "test-123"},
		})
	}))
	defer server.Close()
	
	client := NewClient("test-api-key", WithBaseURL(server.URL))
	
	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	
	err := client.request(context.Background(), "GET", "/test", nil, &result)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	
	if result.Data.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got '%s'", result.Data.ID)
	}
}

func TestClientErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
	}{
		{"Success", 200, `{"data": {}}`, false},
		{"Unauthorized", 401, `{"error": "Invalid API key"}`, true},
		{"Not Found", 404, `{"error": "Design not found"}`, true},
		{"Rate Limited", 429, `{"error": "Rate limit exceeded"}`, true},
		{"Server Error", 500, `{"error": "Internal server error"}`, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()
			
			client := NewClient("test-key", WithBaseURL(server.URL))
			
			var result map[string]interface{}
			err := client.request(context.Background(), "GET", "/test", nil, &result)
			
			if tt.wantErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}