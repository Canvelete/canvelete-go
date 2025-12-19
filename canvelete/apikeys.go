package canvelete

import (
	"context"
	"fmt"
)

// APIKey represents an API key
type APIKey struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	KeyPrefix  string   `json:"keyPrefix"`
	Key        string   `json:"key,omitempty"` // Only returned on creation
	Status     string   `json:"status"`
	Scopes     []string `json:"scopes"`
	CreatedAt  string   `json:"createdAt"`
	LastUsedAt string   `json:"lastUsedAt,omitempty"`
	ExpiresAt  string   `json:"expiresAt,omitempty"`
}

// APIKeysService handles API key-related operations
type APIKeysService struct {
	client *Client
}

// APIKeysListResponse is the response from listing API keys
type APIKeysListResponse struct {
	Data       []APIKey            `json:"data"`
	Pagination PaginatedResponse   `json:"pagination"`
}

// APIKeyListOptions are options for listing API keys
type APIKeyListOptions struct {
	Page  int
	Limit int
}

// List returns a list of API keys
func (s *APIKeysService) List(ctx context.Context, opts *APIKeyListOptions) (*APIKeysListResponse, error) {
	endpoint := "/api/automation/api-keys"
	
	if opts != nil {
		page := opts.Page
		if page == 0 {
			page = 1
		}
		limit := opts.Limit
		if limit == 0 {
			limit = 20
		}
		endpoint = fmt.Sprintf("%s?page=%d&limit=%d", endpoint, page, limit)
	}
	
	var result APIKeysListResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// CreateAPIKeyRequest is the request to create an API key
type CreateAPIKeyRequest struct {
	Name      string `json:"name"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// CreateAPIKeyResponse is the response from creating an API key
type CreateAPIKeyResponse struct {
	Data APIKey `json:"data"`
}

// Create creates a new API key
// Note: The raw API key is only returned once!
func (s *APIKeysService) Create(ctx context.Context, req *CreateAPIKeyRequest) (*APIKey, error) {
	var result CreateAPIKeyResponse
	if err := s.client.request(ctx, "POST", "/api/automation/api-keys", req, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}
