package canvelete

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCanvasAddElement(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		
		response := AddElementResponse{
			Data: CanvasElement{
				ID:     "element-1",
				Type:   "rectangle",
				X:      100,
				Y:      100,
				Width:  200,
				Height: 150,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Canvas.AddElement(context.Background(), "design-123", &CanvasElement{
		Type:   "rectangle",
		X:      100,
		Y:      100,
		Width:  200,
		Height: 150,
	})
	
	if err != nil {
		t.Fatalf("AddElement failed: %v", err)
	}
	
	if result.Type != "rectangle" {
		t.Errorf("Expected type 'rectangle', got '%s'", result.Type)
	}
}

func TestCanvasUpdateElement(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH, got %s", r.Method)
		}
		
		response := AddElementResponse{
			Data: CanvasElement{
				ID: "element-1",
				X:  200,
				Y:  200,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Canvas.UpdateElement(context.Background(), "design-123", "element-1", map[string]interface{}{
		"x": 200,
		"y": 200,
	})
	
	if err != nil {
		t.Fatalf("UpdateElement failed: %v", err)
	}
	
	if result.X != 200 {
		t.Errorf("Expected X 200, got %f", result.X)
	}
}

func TestCanvasDeleteElement(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	err := client.Canvas.DeleteElement(context.Background(), "design-123", "element-1")
	if err != nil {
		t.Fatalf("DeleteElement failed: %v", err)
	}
}

func TestCanvasGetElements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := GetElementsResponse{
			Data: struct {
				Elements []CanvasElement `json:"elements"`
			}{
				Elements: []CanvasElement{
					{ID: "el-1", Type: "rectangle"},
					{ID: "el-2", Type: "text"},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	result, err := client.Canvas.GetElements(context.Background(), "design-123")
	if err != nil {
		t.Fatalf("GetElements failed: %v", err)
	}
	
	if len(result) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(result))
	}
}

func TestCanvasResize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH, got %s", r.Method)
		}
		
		var req ResizeRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		if req.Width != 1920 || req.Height != 1080 {
			t.Errorf("Expected 1920x1080, got %dx%d", req.Width, req.Height)
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	err := client.Canvas.Resize(context.Background(), "design-123", 1920, 1080)
	if err != nil {
		t.Fatalf("Resize failed: %v", err)
	}
}

func TestCanvasClear(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	err := client.Canvas.Clear(context.Background(), "design-123")
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
}

func TestCanvasUpdateBackground(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH, got %s", r.Method)
		}
		
		var req UpdateBackgroundRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		if req.Background != "#ffffff" {
			t.Errorf("Expected background '#ffffff', got '%s'", req.Background)
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	client := NewClient("test-key", WithBaseURL(server.URL))
	
	err := client.Canvas.UpdateBackground(context.Background(), "design-123", "#ffffff")
	if err != nil {
		t.Fatalf("UpdateBackground failed: %v", err)
	}
}