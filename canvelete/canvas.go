package canvelete

import (
	"context"
	"fmt"
)

// CanvasElement represents an element on the canvas
type CanvasElement struct {
	ID           string                 `json:"id,omitempty"`
	Type         string                 `json:"type"`
	X            float64                `json:"x"`
	Y            float64                `json:"y"`
	Width        float64                `json:"width,omitempty"`
	Height       float64                `json:"height,omitempty"`
	Rotation     float64                `json:"rotation,omitempty"`
	Opacity      float64                `json:"opacity,omitempty"`
	Fill         string                 `json:"fill,omitempty"`
	Stroke       string                 `json:"stroke,omitempty"`
	StrokeWidth  float64                `json:"strokeWidth,omitempty"`
	Text         string                 `json:"text,omitempty"`
	FontSize     float64                `json:"fontSize,omitempty"`
	FontFamily   string                 `json:"fontFamily,omitempty"`
	Src          string                 `json:"src,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
}

// CanvasService handles canvas manipulation operations
type CanvasService struct {
	client *Client
}

// AddElementRequest is the request to add an element
type AddElementRequest struct {
	Element CanvasElement `json:"element"`
}

// AddElementResponse is the response from adding an element
type AddElementResponse struct {
	Data CanvasElement `json:"data"`
}


// AddElement adds an element to a design's canvas
func (s *CanvasService) AddElement(ctx context.Context, designID string, element *CanvasElement) (*CanvasElement, error) {
	endpoint := fmt.Sprintf("/api/designs/%s/elements", designID)
	req := AddElementRequest{Element: *element}
	
	var result AddElementResponse
	if err := s.client.request(ctx, "POST", endpoint, req, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// UpdateElement updates an existing element
func (s *CanvasService) UpdateElement(ctx context.Context, designID, elementID string, updates map[string]interface{}) (*CanvasElement, error) {
	endpoint := fmt.Sprintf("/api/designs/%s/elements/%s", designID, elementID)
	
	var result AddElementResponse
	if err := s.client.request(ctx, "PATCH", endpoint, updates, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// DeleteElement removes an element from the canvas
func (s *CanvasService) DeleteElement(ctx context.Context, designID, elementID string) error {
	endpoint := fmt.Sprintf("/api/designs/%s/elements/%s", designID, elementID)
	return s.client.request(ctx, "DELETE", endpoint, nil, nil)
}

// GetElementsResponse is the response from getting elements
type GetElementsResponse struct {
	Data struct {
		Elements []CanvasElement `json:"elements"`
	} `json:"data"`
}

// GetElements retrieves all elements from a design
func (s *CanvasService) GetElements(ctx context.Context, designID string) ([]CanvasElement, error) {
	endpoint := fmt.Sprintf("/api/designs/%s/canvas", designID)
	
	var result GetElementsResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return result.Data.Elements, nil
}

// ResizeRequest is the request to resize the canvas
type ResizeRequest struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Resize changes the canvas dimensions
func (s *CanvasService) Resize(ctx context.Context, designID string, width, height int) error {
	endpoint := fmt.Sprintf("/api/designs/%s/canvas/resize", designID)
	req := ResizeRequest{Width: width, Height: height}
	return s.client.request(ctx, "PATCH", endpoint, req, nil)
}

// Clear removes all elements from the canvas
func (s *CanvasService) Clear(ctx context.Context, designID string) error {
	endpoint := fmt.Sprintf("/api/designs/%s/canvas/elements", designID)
	return s.client.request(ctx, "DELETE", endpoint, nil, nil)
}

// UpdateBackgroundRequest is the request to update background
type UpdateBackgroundRequest struct {
	Background string `json:"background"`
}

// UpdateBackground changes the canvas background
func (s *CanvasService) UpdateBackground(ctx context.Context, designID, background string) error {
	endpoint := fmt.Sprintf("/api/designs/%s/canvas/background", designID)
	req := UpdateBackgroundRequest{Background: background}
	return s.client.request(ctx, "PATCH", endpoint, req, nil)
}