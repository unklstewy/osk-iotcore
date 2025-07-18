package keyboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// LayoutParser handles parsing of keyboard layout files
type LayoutParser struct {
	layoutDir string
}

// NewLayoutParser creates a new layout parser
func NewLayoutParser(layoutDir string) *LayoutParser {
	return &LayoutParser{
		layoutDir: layoutDir,
	}
}

// ParseLayout parses a keyboard layout from a JSON file
func (p *LayoutParser) ParseLayout(filename string) (*Layout, error) {
	fullPath := filepath.Join(p.layoutDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("layout file %s does not exist", fullPath)
	}

	// Read file content
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read layout file %s: %w", fullPath, err)
	}

	// Parse JSON
	var layout Layout
	if err := json.Unmarshal(data, &layout); err != nil {
		return nil, fmt.Errorf("failed to parse layout JSON: %w", err)
	}

	return &layout, nil
}

// SaveLayout saves a keyboard layout to a JSON file
func (p *LayoutParser) SaveLayout(layout *Layout, filename string) error {
	fullPath := filepath.Join(p.layoutDir, filename)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(p.layoutDir, 0755); err != nil {
		return fmt.Errorf("failed to create layout directory: %w", err)
	}

	// Marshal layout to JSON
	data, err := json.MarshalIndent(layout, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal layout to JSON: %w", err)
	}

	// Write to file
	if err := ioutil.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write layout file: %w", err)
	}

	return nil
}

// ListLayouts returns a list of available layout files
func (p *LayoutParser) ListLayouts() ([]string, error) {
	if _, err := os.Stat(p.layoutDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	files, err := ioutil.ReadDir(p.layoutDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read layout directory: %w", err)
	}

	var layouts []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			layouts = append(layouts, file.Name())
		}
	}

	return layouts, nil
}

// ValidateLayout validates a keyboard layout structure
func (p *LayoutParser) ValidateLayout(layout *Layout) error {
	if layout.Name == "" {
		return fmt.Errorf("layout name cannot be empty")
	}

	if layout.Width <= 0 || layout.Height <= 0 {
		return fmt.Errorf("layout dimensions must be positive")
	}

	if len(layout.Keys) == 0 {
		return fmt.Errorf("layout must contain at least one key")
	}

	// Validate each key
	for i, key := range layout.Keys {
		if err := p.validateKey(key, i); err != nil {
			return fmt.Errorf("key %d is invalid: %w", i, err)
		}
	}

	return nil
}

// validateKey validates a single key structure
func (p *LayoutParser) validateKey(key *Key, index int) error {
	if key.ID == "" {
		return fmt.Errorf("key ID cannot be empty")
	}

	if key.Label == "" {
		return fmt.Errorf("key label cannot be empty")
	}

	if key.Width <= 0 || key.Height <= 0 {
		return fmt.Errorf("key dimensions must be positive")
	}

	if key.X < 0 || key.Y < 0 {
		return fmt.Errorf("key position must be non-negative")
	}

	return nil
}
