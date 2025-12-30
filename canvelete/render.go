package canvelete

import (
	"context"
	"fmt"
	"os"
	"time"
)

// RenderService handles render-related API calls
type RenderService struct {
	client *Client
}

// RenderRequest is the request to render a design
type RenderRequest struct {
	DesignID        string                 `json:"designId,omitempty"`
	TemplateID      string                 `json:"templateId,omitempty"`
	DynamicData     map[string]interface{} `json:"dynamicData,omitempty"`
	DynamicElements map[string]interface{} `json:"dynamicElements,omitempty"`
	Format          string                 `json:"format"` // png, jpg, jpeg, pdf, svg
	Width           int                    `json:"width,omitempty"`
	Height          int                    `json:"height,omitempty"`
	Quality         int                    `json:"quality,omitempty"`
}

// Create renders a design or template to an image/PDF
// Uses the backend API directly at /api/v1/render
func (s *RenderService) Create(ctx context.Context, req *RenderRequest) ([]byte, error) {
	if req.DesignID == "" && req.TemplateID == "" {
		return nil, fmt.Errorf("either DesignID or TemplateID is required")
	}
	
	if req.Format == "" {
		req.Format = "png"
	}
	
	if req.Quality == 0 {
		req.Quality = 90
	}
	
	return s.client.requestBinary(ctx, "POST", "/api/v1/render", req)
}

// CreateAndSave renders a design and saves it to a file
func (s *RenderService) CreateAndSave(ctx context.Context, req *RenderRequest, outputPath string) error {
	data, err := s.Create(ctx, req)
	if err != nil {
		return err
	}
	
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// AsyncRenderResponse is the response from creating an async render
type AsyncRenderResponse struct {
	JobID         string `json:"jobId"`
	Status        string `json:"status"`
	EstimatedTime int    `json:"estimatedTime,omitempty"`
}

// AsyncRenderRequest is the request for async rendering
type AsyncRenderRequest struct {
	DesignID    string                 `json:"designId,omitempty"`
	TemplateID  string                 `json:"templateId,omitempty"`
	DynamicData map[string]interface{} `json:"dynamicData,omitempty"`
	Format      string                 `json:"format"`
	Width       int                    `json:"width,omitempty"`
	Height      int                    `json:"height,omitempty"`
	Quality     int                    `json:"quality,omitempty"`
	Async       bool                   `json:"async"`
}

// CreateAsync creates an asynchronous render job
func (s *RenderService) CreateAsync(ctx context.Context, req *RenderRequest) (*AsyncRenderResponse, error) {
	if req.DesignID == "" && req.TemplateID == "" {
		return nil, fmt.Errorf("either DesignID or TemplateID is required")
	}
	
	asyncReq := AsyncRenderRequest{
		DesignID:    req.DesignID,
		TemplateID:  req.TemplateID,
		DynamicData: req.DynamicData,
		Format:      req.Format,
		Width:       req.Width,
		Height:      req.Height,
		Quality:     req.Quality,
		Async:       true,
	}
	
	if asyncReq.Format == "" {
		asyncReq.Format = "png"
	}
	
	var result AsyncRenderResponse
	if err := s.client.request(ctx, "POST", "/api/v1/render/async", asyncReq, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// RenderRecord represents a render job record
type RenderRecord struct {
	ID          string `json:"id"`
	DesignID    string `json:"designId"`
	Format      string `json:"format"`
	Status      string `json:"status"` // pending, processing, completed, failed
	OutputURL   string `json:"outputUrl,omitempty"`
	FileSize    int    `json:"fileSize"`
	CreatedAt   string `json:"createdAt"`
	CompletedAt string `json:"completedAt,omitempty"`
	Error       string `json:"error,omitempty"`
}

// GetStatus retrieves the status of a render job
func (s *RenderService) GetStatus(ctx context.Context, jobID string) (*RenderRecord, error) {
	endpoint := fmt.Sprintf("/api/v1/render/status/%s", jobID)
	
	var result RenderRecord
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// WaitOptions are options for waiting on render completion
type WaitOptions struct {
	Timeout      time.Duration
	PollInterval time.Duration
}

// WaitForCompletion waits for a render job to complete
func (s *RenderService) WaitForCompletion(ctx context.Context, jobID string, opts *WaitOptions) (*RenderRecord, error) {
	timeout := 5 * time.Minute
	pollInterval := 2 * time.Second
	
	if opts != nil {
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
		if opts.PollInterval > 0 {
			pollInterval = opts.PollInterval
		}
	}
	
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		status, err := s.GetStatus(ctx, jobID)
		if err != nil {
			return nil, err
		}
		
		if status.Status == "completed" {
			return status, nil
		}
		
		if status.Status == "failed" {
			errMsg := status.Error
			if errMsg == "" {
				errMsg = "Unknown error"
			}
			return nil, fmt.Errorf("render job failed: %s", errMsg)
		}
		
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}
	
	return nil, fmt.Errorf("render job timed out after %v", timeout)
}

// RenderHistoryResponse is the response from listing renders
type RenderHistoryResponse struct {
	Data       []RenderRecord    `json:"data"`
	Pagination PaginatedResponse `json:"pagination"`
}

// RenderListOptions are options for listing renders
type RenderListOptions struct {
	Page  int
	Limit int
}

// List returns a list of render history
func (s *RenderService) List(ctx context.Context, opts *RenderListOptions) (*RenderHistoryResponse, error) {
	endpoint := "/api/v1/render/history"
	
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
	
	var result RenderHistoryResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// GetHistory retrieves render history
func (s *RenderService) GetHistory(ctx context.Context, opts *RenderListOptions) (*RenderHistoryResponse, error) {
	endpoint := "/api/v1/render/history"
	
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
	
	var result RenderHistoryResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// BatchRenderRequest is the request for batch rendering
type BatchRenderRequest struct {
	Renders []RenderRequest `json:"renders"`
	Webhook string          `json:"webhook,omitempty"`
}

// BatchRenderResponse is the response from batch rendering
type BatchRenderResponse struct {
	BatchID   string                `json:"batchId"`
	Jobs      []AsyncRenderResponse `json:"jobs"`
	TotalJobs int                   `json:"totalJobs"`
}

// BatchCreate creates multiple render jobs
func (s *RenderService) BatchCreate(ctx context.Context, req *BatchRenderRequest) (*BatchRenderResponse, error) {
	if len(req.Renders) == 0 {
		return nil, fmt.Errorf("at least one render is required")
	}
	
	var result BatchRenderResponse
	if err := s.client.request(ctx, "POST", "/api/v1/render/batch", req, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// BatchStatusResponse is the response from getting batch status
type BatchStatusResponse struct {
	Jobs      []RenderRecord `json:"jobs"`
	Completed bool           `json:"completed"`
}

// WaitForBatch waits for all jobs in a batch to complete
func (s *RenderService) WaitForBatch(ctx context.Context, batchID string, opts *WaitOptions) ([]RenderRecord, error) {
	timeout := 10 * time.Minute
	pollInterval := 5 * time.Second
	
	if opts != nil {
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
		if opts.PollInterval > 0 {
			pollInterval = opts.PollInterval
		}
	}
	
	deadline := time.Now().Add(timeout)
	endpoint := fmt.Sprintf("/api/v1/render/batch/%s/status", batchID)
	
	for time.Now().Before(deadline) {
		var result BatchStatusResponse
		if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
			return nil, err
		}
		
		if result.Completed {
			return result.Jobs, nil
		}
		
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}
	
	return nil, fmt.Errorf("batch render timed out after %v", timeout)
}
