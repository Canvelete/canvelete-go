package canvelete

import (
	"context"
	"fmt"
	"net/url"
)

// Asset represents a user asset
type Asset struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"` // IMAGE, FONT, VIDEO, AUDIO
	URL       string `json:"url"`
	Format    string `json:"format,omitempty"`
	Size      int64  `json:"size,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// StockImage represents a stock image from Pixabay
type StockImage struct {
	ID            string `json:"id"`
	Tags          string `json:"tags"`
	PreviewURL    string `json:"previewURL"`
	WebformatURL  string `json:"webformatURL"`
	LargeImageURL string `json:"largeImageURL"`
	ImageWidth    int    `json:"imageWidth"`
	ImageHeight   int    `json:"imageHeight"`
}

// Icon represents an icon asset
type Icon struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Tags []string `json:"tags,omitempty"`
}

// Font represents a font
type Font struct {
	Family   string   `json:"family"`
	Variants []string `json:"variants"`
	Category string   `json:"category,omitempty"`
}

// AssetsService handles asset-related operations
type AssetsService struct {
	client *Client
}


// AssetListOptions are options for listing assets
type AssetListOptions struct {
	Page  int
	Limit int
	Type  string // IMAGE, FONT, VIDEO, AUDIO
}

// AssetsListResponse is the response from listing assets
type AssetsListResponse struct {
	Data       []Asset           `json:"data"`
	Pagination PaginatedResponse `json:"pagination"`
}

// List returns a list of user assets
func (s *AssetsService) List(ctx context.Context, opts *AssetListOptions) (*AssetsListResponse, error) {
	endpoint := "/api/assets/library"
	
	if opts != nil {
		params := url.Values{}
		if opts.Page > 0 {
			params.Add("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Limit > 0 {
			params.Add("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Type != "" {
			params.Add("type", opts.Type)
		}
		if len(params) > 0 {
			endpoint += "?" + params.Encode()
		}
	}
	
	var result AssetsListResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Get retrieves an asset by ID
func (s *AssetsService) Get(ctx context.Context, assetID string) (*Asset, error) {
	endpoint := fmt.Sprintf("/api/assets/%s", assetID)
	
	var result struct {
		Data Asset `json:"data"`
	}
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// Delete removes an asset
func (s *AssetsService) Delete(ctx context.Context, assetID string) error {
	endpoint := fmt.Sprintf("/api/assets/%s", assetID)
	return s.client.request(ctx, "DELETE", endpoint, nil, nil)
}

// SearchOptions are options for searching assets
type SearchOptions struct {
	Query   string
	Page    int
	PerPage int
}

// StockImagesResponse is the response from searching stock images
type StockImagesResponse struct {
	Data []StockImage `json:"data"`
}

// SearchStockImages searches for stock images from Pixabay
func (s *AssetsService) SearchStockImages(ctx context.Context, opts *SearchOptions) (*StockImagesResponse, error) {
	params := url.Values{}
	params.Add("query", opts.Query)
	if opts.Page > 0 {
		params.Add("page", fmt.Sprintf("%d", opts.Page))
	}
	if opts.PerPage > 0 {
		params.Add("perPage", fmt.Sprintf("%d", opts.PerPage))
	}
	
	endpoint := "/api/assets/stock-images?" + params.Encode()
	
	var result StockImagesResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// IconsResponse is the response from searching icons
type IconsResponse struct {
	Data []Icon `json:"data"`
}

// SearchIcons searches for icons
func (s *AssetsService) SearchIcons(ctx context.Context, opts *SearchOptions) (*IconsResponse, error) {
	params := url.Values{}
	params.Add("query", opts.Query)
	if opts.Page > 0 {
		params.Add("page", fmt.Sprintf("%d", opts.Page))
	}
	if opts.PerPage > 0 {
		params.Add("perPage", fmt.Sprintf("%d", opts.PerPage))
	}
	
	endpoint := "/api/assets/icons?" + params.Encode()
	
	var result IconsResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// FontsResponse is the response from listing fonts
type FontsResponse struct {
	Data []Font `json:"data"`
}

// ListFonts returns available fonts
func (s *AssetsService) ListFonts(ctx context.Context, category string) (*FontsResponse, error) {
	endpoint := "/api/assets/fonts"
	if category != "" {
		endpoint += "?category=" + url.QueryEscape(category)
	}
	
	var result FontsResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}