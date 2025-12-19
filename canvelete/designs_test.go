package canvelete

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDesignsList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		
		response := DesignsListResponse{
			Data: []Design{
				{ID: "design-1", Name: "Test Design 1"},
				{ID: "design-2", Name: "Test Design 2"},
			},
			Pagination: PaginatedResponse{Page: 1, Limit: 20, Total: 2, TotalPages: 1},
		}
		
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Designs.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	
	if len(result.Data) != 2 {
		t.Errorf("Expected 2 designs, got %d", len(result.Data))
	}
}

func TestDesignsCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		
		var req CreateDesignRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		if req.Name != "New Design" {
			t.Errorf("Expected name 'New Design', got '%s'", req.Name)
		}
		
		response := CreateDesignResponse{
			Success: true,
			Data: Design{
				ID:     "new-design-id",
				Name:   req.Name,
				Width:  req.Width,
				Height: req.Height,
			},
		}
		
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Designs.Create(context.Background(), &CreateDesignRequest{
		Name:   "New Design",
		Width:  1920,
		Height: 1080,
	})
	
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	if result.Name != "New Design" {
		t.Errorf("Expected name 'New Design', got '%s'", result.Name)
	}
}

func TestDesignsGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Data Design `json:"data"`
		}{
			Data: Design{ID: "design-123", Name: "Test Design"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Designs.Get(context.Background(), "design-123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	
	if result.ID != "design-123" {
		t.Errorf("Expected ID 'design-123', got '%s'", result.ID)
	}
}

func TestDesignsUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH, got %s", r.Method)
		}
		
		response := struct {
			Data Design `json:"data"`
		}{
			Data: Design{ID: "design-123", Name: "Updated Design"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	name := "Updated Design"
	result, err := client.Designs.Update(context.Background(), "design-123", &UpdateDesignRequest{
		Name: &name,
	})
	
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	
	if result.Name != "Updated Design" {
		t.Errorf("Expected name 'Updated Design', got '%s'", result.Name)
	}
}

func TestDesignsDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	err := client.Designs.Delete(context.Background(), "design-123")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}