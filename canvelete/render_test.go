package canvelete

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderCreate(t *testing.T) {
	expectedData := []byte("image-data")
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		
		var req RenderRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		if req.DesignID != "design-123" {
			t.Errorf("Expected designId 'design-123', got '%s'", req.DesignID)
		}
		
		w.Write(expectedData)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Render.Create(context.Background(), &RenderRequest{
		DesignID: "design-123",
		Format:   "png",
	})
	
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	if string(result) != string(expectedData) {
		t.Errorf("Expected '%s', got '%s'", string(expectedData), string(result))
	}
}

func TestRenderCreateValidation(t *testing.T) {
	client := NewClient("test-key")
	
	_, err := client.Render.Create(context.Background(), &RenderRequest{
		Format: "png",
	})
	
	if err == nil {
		t.Error("Expected error when neither designId nor templateId provided")
	}
}

func TestRenderCreateAsync(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := AsyncRenderResponse{
			JobID:         "job-123",
			Status:        "pending",
			EstimatedTime: 30,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Render.CreateAsync(context.Background(), &RenderRequest{
		DesignID: "design-123",
		Format:   "pdf",
	})
	
	if err != nil {
		t.Fatalf("CreateAsync failed: %v", err)
	}
	
	if result.JobID != "job-123" {
		t.Errorf("Expected jobId 'job-123', got '%s'", result.JobID)
	}
}

func TestRenderGetStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := RenderRecord{
			ID:       "job-123",
			DesignID: "design-123",
			Status:   "processing",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Render.GetStatus(context.Background(), "job-123")
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	
	if result.Status != "processing" {
		t.Errorf("Expected status 'processing', got '%s'", result.Status)
	}
}

func TestRenderList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := RenderHistoryResponse{
			Data: []RenderRecord{
				{ID: "render-1", Status: "completed"},
				{ID: "render-2", Status: "completed"},
			},
			Pagination: PaginatedResponse{Page: 1, Limit: 20, Total: 2, TotalPages: 1},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Render.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	
	if len(result.Data) != 2 {
		t.Errorf("Expected 2 renders, got %d", len(result.Data))
	}
}

func TestRenderBatchCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := BatchRenderResponse{
			BatchID: "batch-123",
			Jobs: []AsyncRenderResponse{
				{JobID: "job-1", Status: "pending"},
				{JobID: "job-2", Status: "pending"},
			},
			TotalJobs: 2,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Render.BatchCreate(context.Background(), &BatchRenderRequest{
		Renders: []RenderRequest{
			{DesignID: "design-1", Format: "png"},
			{DesignID: "design-2", Format: "jpg"},
		},
	})
	
	if err != nil {
		t.Fatalf("BatchCreate failed: %v", err)
	}
	
	if result.TotalJobs != 2 {
		t.Errorf("Expected 2 jobs, got %d", result.TotalJobs)
	}
}

func TestRenderBatchCreateValidation(t *testing.T) {
	client := NewClient("test-key")
	
	_, err := client.Render.BatchCreate(context.Background(), &BatchRenderRequest{
		Renders: []RenderRequest{},
	})
	
	if err == nil {
		t.Error("Expected error for empty batch")
	}
}