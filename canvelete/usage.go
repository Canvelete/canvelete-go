package canvelete

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// UsageStats represents current usage statistics
type UsageStats struct {
	CreditsUsed      int   `json:"creditsUsed"`
	CreditLimit      int   `json:"creditLimit"`
	CreditsRemaining int   `json:"creditsRemaining"`
	APICalls         int   `json:"apiCalls"`
	APICallLimit     int   `json:"apiCallLimit"`
	Renders          int   `json:"renders"`
	StorageUsed      int64 `json:"storageUsed"`
}

// UsageEvent represents a usage event
type UsageEvent struct {
	Type        string                 `json:"type"`
	CreditsUsed int                    `json:"creditsUsed"`
	Timestamp   string                 `json:"timestamp"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// APIStats represents API usage statistics
type APIStats struct {
	Endpoints  map[string]int `json:"endpoints"`
	TotalCalls int            `json:"totalCalls"`
	Period     string         `json:"period"`
}

// Activity represents a user activity
type Activity struct {
	Action    string                 `json:"action"`
	Timestamp string                 `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// Analytics represents usage analytics
type Analytics struct {
	TotalRenders  int            `json:"totalRenders"`
	AveragePerDay float64        `json:"averagePerDay"`
	PeakDay       string         `json:"peakDay"`
	Trend         string         `json:"trend"`
	Breakdown     map[string]int `json:"breakdown,omitempty"`
}

// UsageService handles usage-related operations
type UsageService struct {
	client *Client
}


// GetStats retrieves current usage statistics
func (s *UsageService) GetStats(ctx context.Context) (*UsageStats, error) {
	var result struct {
		Data UsageStats `json:"data"`
	}
	
	if err := s.client.request(ctx, "GET", "/api/v1/usage/stats", nil, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// UsageHistoryOptions are options for getting usage history
type UsageHistoryOptions struct {
	Page      int
	Limit     int
	StartDate *time.Time
	EndDate   *time.Time
}

// UsageHistoryResponse is the response from getting usage history
type UsageHistoryResponse struct {
	Data       []UsageEvent      `json:"data"`
	Pagination PaginatedResponse `json:"pagination"`
}

// GetHistory retrieves usage history
func (s *UsageService) GetHistory(ctx context.Context, opts *UsageHistoryOptions) (*UsageHistoryResponse, error) {
	params := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			params.Add("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Limit > 0 {
			params.Add("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.StartDate != nil {
			params.Add("startDate", opts.StartDate.Format(time.RFC3339))
		}
		if opts.EndDate != nil {
			params.Add("endDate", opts.EndDate.Format(time.RFC3339))
		}
	}
	
	endpoint := "/api/v1/usage/history"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	
	var result UsageHistoryResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// GetAPIStats retrieves API usage statistics by endpoint
func (s *UsageService) GetAPIStats(ctx context.Context) (*APIStats, error) {
	var result struct {
		Data APIStats `json:"data"`
	}
	
	if err := s.client.request(ctx, "GET", "/api/v1/usage/api-stats", nil, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// ActivityListOptions are options for listing activities
type ActivityListOptions struct {
	Page  int
	Limit int
}

// ActivitiesResponse is the response from getting activities
type ActivitiesResponse struct {
	Data       []Activity        `json:"data"`
	Pagination PaginatedResponse `json:"pagination"`
}

// GetActivities retrieves recent activities
func (s *UsageService) GetActivities(ctx context.Context, opts *ActivityListOptions) (*ActivitiesResponse, error) {
	params := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			params.Add("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Limit > 0 {
			params.Add("limit", fmt.Sprintf("%d", opts.Limit))
		}
	}
	
	endpoint := "/api/usage/activities"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	
	var result ActivitiesResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// GetAnalytics retrieves usage analytics
func (s *UsageService) GetAnalytics(ctx context.Context, period string) (*Analytics, error) {
	if period == "" {
		period = "month"
	}
	
	endpoint := fmt.Sprintf("/api/usage/analytics?period=%s", url.QueryEscape(period))
	
	var result struct {
		Data Analytics `json:"data"`
	}
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}