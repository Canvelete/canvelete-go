package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/canvelete/canvelete-go/canvelete"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("CANVELETE_API_KEY")
	if apiKey == "" {
		log.Fatal("CANVELETE_API_KEY environment variable is required")
	}

	// Create client
	client := canvelete.NewClient(apiKey)
	ctx := context.Background()

	fmt.Println("Canvelete Go SDK Quickstart\n")

	// List designs
	fmt.Println("Fetching your designs...")
	designs, err := client.Designs.List(ctx, &canvelete.ListOptions{
		Limit: 5,
	})
	if err != nil {
		log.Fatalf("Failed to list designs: %v", err)
	}

	fmt.Printf("✓ Found %d designs\n", len(designs.Data))
	for _, design := range designs.Data {
		fmt.Printf("  - %s (ID: %s)\n", design.Name, design.ID)
	}

	// Create a new design
	fmt.Println("\nCreating a new design...")
	canvasData := map[string]interface{}{
		"elements": []map[string]interface{}{
			{
				"type":       "text",
				"text":       "Hello from Canvelete Go SDK!",
				"x":          100,
				"y":          100,
				"fontSize":   48,
				"fontFamily": "Arial",
				"fill":       "#000000",
			},
		},
		"background": "#FFFFFF",
	}

	newDesign, err := client.Designs.Create(ctx, &canvelete.CreateDesignRequest{
		Name:        "Go SDK Test Design",
		Description: "Created via Go SDK quickstart",
		CanvasData:  canvasData,
		Width:       1920,
		Height:      1080,
	})
	if err != nil {
		log.Fatalf("Failed to create design: %v", err)
	}

	fmt.Printf("✓ Created design: %s (ID: %s)\n", newDesign.Name, newDesign.ID)

	// Render the design
	fmt.Println("\nRendering design to PNG...")
	outputFile := "quickstart_output.png"

	err = client.Render.CreateAndSave(ctx, &canvelete.RenderRequest{
		DesignID: newDesign.ID,
		Format:   "png",
		Quality:  90,
	}, outputFile)
	if err != nil {
		log.Fatalf("Failed to render design: %v", err)
	}

	fmt.Printf("✓ Rendered and saved to %s\n", outputFile)

	// List templates
	fmt.Println("\nFetching available templates...")
	templates, err := client.Templates.List(ctx, &canvelete.TemplateListOptions{
		Limit: 5,
	})
	if err != nil {
		log.Fatalf("Failed to list templates: %v", err)
	}

	fmt.Printf("✓ Found %d templates\n", len(templates.Data))
	for i, template := range templates.Data {
		if i < 3 {
			fmt.Printf("  - %s\n", template.Name)
		}
	}

	fmt.Println("\n✅ Quickstart complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. View your design at: https://www.canvelete.com/dashboard")
	fmt.Printf("  2. Check the rendered image: %s\n", outputFile)
	fmt.Println("  3. Explore more examples in the examples/ directory")
}
