package canvelete

import (
	"context"
	"fmt"
	"net/url"
)

// Template represents a Canvelete template
type Template Design // Templates are designs with isTemplate=true

// TemplatesService handles template-related API calls
type TemplatesService struct {
	client *Client
}

// TemplateListOptions are options for listing templates
type TemplateListOptions struct {
	Page   int
	Limit  int
	MyOnly bool
	Search string
}

// TemplatesListResponse is the response from listing templates
type TemplatesListResponse struct {
	Data       []Template          `json:"data"`
	Pagination PaginatedResponse   `json:"pagination"`
}

// List returns a list of templates
func (s *TemplatesService) List(ctx context.Context, opts *TemplateListOptions) (*TemplatesListResponse, error) {
	endpoint := "/api/automation/templates"
	
	if opts != nil {
		params := url.Values{}
		if opts.Page > 0 {
			params.Add("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Limit > 0 {
			params.Add("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.MyOnly {
			params.Add("myOnly", "true")
		}
		if opts.Search != "" {
			params.Add("search", opts.Search)
		}
		
		if len(params) > 0 {
			endpoint += "?" + params.Encode()
		}
	}
	
	var result TemplatesListResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Get retrieves a template by ID
func (s *TemplatesService) Get(ctx context.Context, id string) (*Template, error) {
	var result struct {
		Data Template `json:"data"`
	}
	
	endpoint := fmt.Sprintf("/api/automation/designs/%s", id)
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// IterateAll iterates through all templates with automatic pagination
func (s *TemplatesService) IterateAll(ctx context.Context, opts *TemplateListOptions) (<-chan *Template, <-chan error) {
	templateChan := make(chan *Template)
	errChan := make(chan error, 1)
	
	go func() {
		defer close(templateChan)
		defer close(errChan)
		
		page := 1
		limit := 50
		if opts != nil && opts.Limit > 0 {
			limit = opts.Limit
		}
		
		for {
			listOpts := &TemplateListOptions{
				Page:  page,
				Limit: limit,
			}
			if opts != nil {
				listOpts.MyOnly = opts.MyOnly
				listOpts.Search = opts.Search
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
				case templateChan <- &resp.Data[i]:
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
	
	return templateChan, errChan
}


// ApplyTemplateRequest is the request to apply a template
type ApplyTemplateRequest struct {
	TemplateID  string                 `json:"templateId"`
	DynamicData map[string]interface{} `json:"dynamicData,omitempty"`
}

// Apply applies a template to a design
func (s *TemplatesService) Apply(ctx context.Context, designID string, req *ApplyTemplateRequest) (*Design, error) {
	endpoint := fmt.Sprintf("/api/automation/designs/%s/apply-template", designID)
	
	var result struct {
		Data Design `json:"data"`
	}
	
	if err := s.client.request(ctx, "POST", endpoint, req, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// CreateTemplateRequest is the request to create a template from a design
type CreateTemplateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
}

// Create creates a template from an existing design
func (s *TemplatesService) Create(ctx context.Context, designID string, req *CreateTemplateRequest) (*Template, error) {
	endpoint := fmt.Sprintf("/api/automation/designs/%s/save-as-template", designID)
	
	var result struct {
		Data Template `json:"data"`
	}
	
	if err := s.client.request(ctx, "POST", endpoint, req, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}