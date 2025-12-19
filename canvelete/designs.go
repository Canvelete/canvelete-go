package canvelete

import (
	"context"
	"fmt"
	"net/url"
)

// Design represents a Canvelete design
type Design struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	CanvasData   map[string]interface{} `json:"canvasData"`
	Width        int                    `json:"width"`
	Height       int                    `json:"height"`
	Status       string                 `json:"status"`
	Visibility   string                 `json:"visibility"`
	IsTemplate   bool                   `json:"isTemplate"`
	ThumbnailURL string                 `json:"thumbnailUrl,omitempty"`
	CreatedAt    string                 `json:"createdAt"`
	UpdatedAt    string                 `json:"updatedAt"`
}

// DesignsService handles design-related API calls
type DesignsService struct {
	client *Client
}

// ListOptions are options for listing designs
type ListOptions struct {
	Page       int
	Limit      int
	IsTemplate *bool
	Status     string
}

// DesignsListResponse is the response from listing designs
type DesignsListResponse struct {
	Data       []Design            `json:"data"`
	Pagination PaginatedResponse   `json:"pagination"`
}

// List returns a list of designs
func (s *DesignsService) List(ctx context.Context, opts *ListOptions) (*DesignsListResponse, error) {
	endpoint := "/api/automation/designs"
	
	if opts != nil {
		params := url.Values{}
		if opts.Page > 0 {
			params.Add("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Limit > 0 {
			params.Add("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.IsTemplate != nil {
			if *opts.IsTemplate {
				params.Add("isTemplate", "true")
			} else {
				params.Add("isTemplate", "false")
			}
		}
		if opts.Status != "" {
			params.Add("status", opts.Status)
		}
		
		if len(params) > 0 {
			endpoint += "?" + params.Encode()
		}
	}
	
	var result DesignsListResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// CreateDesignRequest is the request to create a design
type CreateDesignRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	CanvasData  map[string]interface{} `json:"canvasData"`
	Width       int                    `json:"width"`
	Height      int                    `json:"height"`
	IsTemplate  bool                   `json:"isTemplate,omitempty"`
	Visibility  string                 `json:"visibility,omitempty"`
}

// CreateDesignResponse is the response from creating a design
type CreateDesignResponse struct {
	Success bool   `json:"success"`
	Data    Design `json:"data"`
}

// Create creates a new design
func (s *DesignsService) Create(ctx context.Context, req *CreateDesignRequest) (*Design, error) {
	var result CreateDesignResponse
	if err := s.client.request(ctx, "POST", "/api/automation/designs", req, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// Get retrieves a design by ID
func (s *DesignsService) Get(ctx context.Context, id string) (*Design, error) {
	var result struct {
		Data Design `json:"data"`
	}
	
	endpoint := fmt.Sprintf("/api/automation/designs/%s", id)
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// UpdateDesignRequest is the request to update a design
type UpdateDesignRequest struct {
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	CanvasData  map[string]interface{}  `json:"canvasData,omitempty"`
	Status      *string                 `json:"status,omitempty"`
	Visibility  *string                 `json:"visibility,omitempty"`
}

// Update updates a design
func (s *DesignsService) Update(ctx context.Context, id string, req *UpdateDesignRequest) (*Design, error) {
	var result struct {
		Data Design `json:"data"`
	}
	
	endpoint := fmt.Sprintf("/api/automation/designs/%s", id)
	if err := s.client.request(ctx, "PATCH", endpoint, req, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// Delete deletes a design
func (s *DesignsService) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/api/automation/designs/%s", id)
	return s.client.request(ctx, "DELETE", endpoint, nil, nil)
}

// IterateAll iterates through all designs with automatic pagination
func (s *DesignsService) IterateAll(ctx context.Context, opts *ListOptions) (<-chan *Design, <-chan error) {
	designChan := make(chan *Design)
	errChan := make(chan error, 1)
	
	go func() {
		defer close(designChan)
		defer close(errChan)
		
		page := 1
		limit := 50
		if opts != nil && opts.Limit > 0 {
			limit = opts.Limit
		}
		
		for {
			listOpts := &ListOptions{
				Page:  page,
				Limit: limit,
			}
			if opts != nil {
				listOpts.IsTemplate = opts.IsTemplate
				listOpts.Status = opts.Status
			}
			
			resp, err := s.List(ctx, listOpts)
			if err != nil {
				errChan <- err
				return
			}
			
			if len(resp.Data) == 0 {
				return
			}
			
			for i := range resp.Data {
				select {
				case designChan <- &resp.Data[i]:
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				}
			}
			
			if page >= resp.Pagination.TotalPages {
				return
			}
			page++
		}
	}()
	
	return designChan, errChan
}
