# Canvelete Go SDK

Official Go client library for the [Canvelete](https://www.canvelete.com) API.

## Installation

```bash
go get github.com/canvelete/canvelete-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/canvelete/canvelete-go/canvelete"
)

func main() {
    // Create client
    client := canvelete.NewClient("cvt_your_api_key")
    ctx := context.Background()
    
    // List designs
    designs, err := client.Designs.List(ctx, &canvelete.ListOptions{
        Limit: 20,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d designs\n", len(designs.Data))
    
    // Create a design
    canvasData := map[string]interface{}{
        "elements": []map[string]interface{}{
            {
                "type": "text",
                "text": "Hello Canvelete!",
                "x":    100,
                "y":    100,
                "fontSize": 48,
            },
        },
    }
    
    design, err := client.Designs.Create(ctx, &canvelete.CreateDesignRequest{
        Name:       "My Design",
        CanvasData: canvasData,
        Width:      1920,
        Height:     1080,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Render to file
    err = client.Render.CreateAndSave(ctx, &canvelete.RenderRequest{
        DesignID: design.ID,
        Format:   "png",
        Quality:  90,
    }, "output.png")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Rendered to output.png")
}
```

## Features

✅ **Zero Dependencies** - Uses only the Go standard library  
✅ **Context Support** - All methods accept `context.Context`  
✅ **Type Safety** - Fully typed request and response structs  
✅ **Idiomatic Go** - Follows Go best practices and conventions  
✅ **Channel-based Iteration** - Efficient pagination with goroutines  
✅ **Error Handling** - Structured error responses

## API Reference

### Client Initialization

```go
// Basic client
client := canvelete.NewClient("cvt_your_api_key")

// With options
client := canvelete.NewClient("cvt_your_api_key",
    canvelete.WithBaseURL("https://custom.canvelete.com"),
    canvelete.WithTimeout(60 * time.Second),
)
```

### Designs

```go
// List designs
designs, err := client.Designs.List(ctx, &canvelete.ListOptions{
    Page:   1,
    Limit:  20,
    Status: "PUBLISHED",
})

// Iterate all designs (auto-pagination)
designChan, errChan := client.Designs.IterateAll(ctx, &canvelete.ListOptions{
    Limit: 100,
})

for design := range designChan {
    fmt.Println(design.Name)
}

if err := <-errChan; err != nil {
    log.Fatal(err)
}

// Create design
design, err := client.Designs.Create(ctx, &canvelete.CreateDesignRequest{
    Name:       "New Design",
    CanvasData: map[string]interface{}{...},
    Width:      1920,
    Height:     1080,
})

// Get design
design, err := client.Designs.Get(ctx, "design_id")

// Update design
name := "Updated Name"
design, err := client.Designs.Update(ctx, "design_id", &canvelete.UpdateDesignRequest{
    Name: &name,
})

// Delete design
err := client.Designs.Delete(ctx, "design_id")
```

### Templates

```go
// List templates
templates, err := client.Templates.List(ctx, &canvelete.TemplateListOptions{
    Search: "certificate",
    Limit:  20,
})

// Iterate all templates
templateChan, errChan := client.Templates.IterateAll(ctx, &canvelete.TemplateListOptions{
    MyOnly: true,
})

for template := range templateChan {
    fmt.Println(template.Name)
}

// Get template
template, err := client.Templates.Get(ctx, "template_id")
```

### Render

```go
// Render to bytes
imageData, err := client.Render.Create(ctx, &canvelete.RenderRequest{
    DesignID: "design_id",
    Format:   "png",
    Quality:  90,
})

// Render and save to file
err := client.Render.CreateAndSave(ctx, &canvelete.RenderRequest{
    TemplateID: "template_id",
    DynamicData: map[string]interface{}{
        "name": "John Doe",
        "date": "2024-01-01",
    },
    Format: "pdf",
}, "certificate.pdf")

// Custom dimensions
imageData, err := client.Render.Create(ctx, &canvelete.RenderRequest{
    DesignID: "design_id",
    Format:   "jpg",
    Width:    1200,
    Height:   630,
})

// List render history
renders, err := client.Render.List(ctx, &canvelete.RenderListOptions{
    Page:  1,
    Limit: 20,
})
```

### API Keys

```go
// List API keys (requires OAuth2)
keys, err := client.APIKeys.List(ctx, &canvelete.APIKeyListOptions{
    Limit: 20,
})

// Create new API key
key, err := client.APIKeys.Create(ctx, &canvelete.CreateAPIKeyRequest{
    Name: "Production Key",
})

// ⚠️ Save the key immediately - it's only shown once!
fmt.Printf("API Key: %s\n", key.Key)
```

## Advanced Usage

### Context with Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

designs, err := client.Designs.List(ctx, nil)
```

### Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 100,
    },
}

client := canvelete.NewClient("cvt_key",
    canvelete.WithHTTPClient(httpClient),
)
```

### Channel-based Iteration

```go
ctx := context.Background()
designChan, errChan := client.Designs.IterateAll(ctx, &canvelete.ListOptions{
    Limit: 100,
})

// Process designs concurrently
for design := range designChan {
    go processDesign(design)
}

// Check for errors
if err := <-errChan; err != nil {
    log.Printf("Error during iteration: %v", err)
}
```

### Error Handling

```go
design, err := client.Designs.Get(ctx, "invalid_id")
if err != nil {
    if apiErr, ok := err.(*canvelete.APIError); ok {
        switch apiErr.StatusCode {
        case 404:
            fmt.Println("Design not found")
        case 401:
            fmt.Println("Invalid API key")
        case 429:
            fmt.Println("Rate limited")
        default:
            fmt.Printf("API error: %v\n", apiErr)
        }
    } else {
        fmt.Printf("Request error: %v\n", err)
    }
}
```

## Examples

See the [examples](./examples) directory for complete working examples:

- `quickstart/` - Basic usage example
- `pagination/` - Channel-based iteration
- `concurrent/` - Concurrent rendering

## Environment Variables

```bash
export CANVELETE_API_KEY="cvt_your_api_key"
export CANVELETE_BASE_URL="https://www.canvelete.com"
```

```go
apiKey := os.Getenv("CANVELETE_API_KEY")
client := canvelete.NewClient(apiKey)
```

## Requirements

- Go 1.21 or higher
- No external dependencies (uses standard library only)

## License

MIT License - see LICENSE file for details.

## Support

- **Documentation**: https://docs.canvelete.com
- **API Reference**: https://docs.canvelete.com/api
- **Issues**: https://github.com/canvelete/canvelete-go/issues
- **Email**: support@canvelete.com
