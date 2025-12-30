// Package canvelete provides a Go client for the Canvelete API.
package canvelete

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://api.canvelete.com"
	defaultTimeout = 30 * time.Second
	sdkVersion     = "2.0.0"
)

// Client is the main Canvelete API client
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	
	// Core resource services
	Designs   *DesignsService
	Templates *TemplatesService
	Render    *RenderService
	APIKeys   *APIKeysService
	
	// New resource services
	Canvas  *CanvasService
	Assets  *AssetsService
	Usage   *UsageService
	Billing *BillingService
}

// NewClient creates a new Canvelete API client
func NewClient(apiKey string, options ...ClientOption) *Client {
	c := &Client{
		APIKey:  apiKey,
		BaseURL: defaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	
	// Apply options
	for _, opt := range options {
		opt(c)
	}
	
	// Initialize core services
	c.Designs = &DesignsService{client: c}
	c.Templates = &TemplatesService{client: c}
	c.Render = &RenderService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	
	// Initialize new services
	c.Canvas = &CanvasService{client: c}
	c.Assets = &AssetsService{client: c}
	c.Usage = &UsageService{client: c}
	c.Billing = &BillingService{client: c}
	
	return c
}

// ClientOption is a functional option for configuring the client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithTimeout sets a custom timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.HTTPClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// request makes an authenticated HTTP request
func (c *Client) request(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	// Build URL
	u, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	
	// Prepare request body
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "canvelete-go/"+sdkVersion)
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	
	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Check for errors
	if resp.StatusCode >= 400 {
		return parseErrorResponse(resp)
	}
	
	// Parse response
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	
	return nil
}

// requestBinary makes a request and returns binary data
func (c *Client) requestBinary(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	u, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return nil, parseErrorResponse(resp)
	}
	
	return io.ReadAll(resp.Body)
}

// parseErrorResponse parses an error response
func parseErrorResponse(resp *http.Response) error {
	bodyBytes, _ := io.ReadAll(resp.Body)
	
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	
	if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
		message := errResp.Error
		if message == "" {
			message = errResp.Message
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    message,
		}
	}
	
	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(bodyBytes)),
	}
}

// APIError represents an API error
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}
