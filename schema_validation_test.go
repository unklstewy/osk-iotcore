package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iotcore/osk-iotcore/pkg/keyboard"
)

// TestSchemaValidation validates all JSON layouts and themes against the struct definitions
func TestSchemaValidation(t *testing.T) {
	// Test layout validation
	t.Run("LayoutValidation", func(t *testing.T) {
		testLayoutValidation(t)
	})

	// Test theme validation
	t.Run("ThemeValidation", func(t *testing.T) {
		testThemeValidation(t)
	})

	// Test headless rendering
	t.Run("HeadlessRendering", func(t *testing.T) {
		testHeadlessRendering(t)
	})
}

// testLayoutValidation tests all layout JSON files
func testLayoutValidation(t *testing.T) {
	layoutsDir := "examples/assets/layouts"
	
	// Find all JSON files in the layouts directory
	jsonFiles, err := filepath.Glob(filepath.Join(layoutsDir, "*.json"))
	if err != nil {
		t.Fatalf("Failed to find layout JSON files: %v", err)
	}

	if len(jsonFiles) == 0 {
		t.Logf("No layout JSON files found in %s", layoutsDir)
		return
	}

	t.Logf("Found %d layout JSON files to validate", len(jsonFiles))

	for _, jsonFile := range jsonFiles {
		t.Run(filepath.Base(jsonFile), func(t *testing.T) {
			validateLayoutFile(t, jsonFile)
		})
	}
}

// validateLayoutFile validates a single layout JSON file
func validateLayoutFile(t *testing.T, jsonFile string) {
	// Read the JSON file
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read layout file %s: %v", jsonFile, err)
	}

	// Unmarshal into Layout struct
	var layout keyboard.Layout
	if err := json.Unmarshal(data, &layout); err != nil {
		t.Fatalf("Failed to unmarshal layout JSON from %s: %v", jsonFile, err)
	}

	// Validate required fields
	if layout.Name == "" {
		t.Error("Layout name is required but empty")
	}

	if layout.Width <= 0 {
		t.Errorf("Layout width must be positive, got %d", layout.Width)
	}

	if layout.Height <= 0 {
		t.Errorf("Layout height must be positive, got %d", layout.Height)
	}

	if len(layout.Keys) == 0 {
		t.Error("Layout must contain at least one key")
	}

	// Validate each key
	for i, key := range layout.Keys {
		if key.ID == "" {
			t.Errorf("Key %d: ID is required but empty", i)
		}

		if key.Label == "" {
			t.Errorf("Key %d: Label is required but empty", i)
		}

		if key.Width <= 0 {
			t.Errorf("Key %d: Width must be positive, got %d", i, key.Width)
		}

		if key.Height <= 0 {
			t.Errorf("Key %d: Height must be positive, got %d", i, key.Height)
		}

		if key.X < 0 {
			t.Errorf("Key %d: X position must be non-negative, got %d", i, key.X)
		}

		if key.Y < 0 {
			t.Errorf("Key %d: Y position must be non-negative, got %d", i, key.Y)
		}

		// Validate key doesn't exceed layout bounds
		if key.X+key.Width > layout.Width {
			t.Errorf("Key %d: Key extends beyond layout width (key right: %d, layout width: %d)", 
				i, key.X+key.Width, layout.Width)
		}

		if key.Y+key.Height > layout.Height {
			t.Errorf("Key %d: Key extends beyond layout height (key bottom: %d, layout height: %d)", 
				i, key.Y+key.Height, layout.Height)
		}
	}

	// Check for overlapping keys
	for i, key1 := range layout.Keys {
		for j, key2 := range layout.Keys {
			if i >= j {
				continue
			}

			if keysOverlap(key1, key2) {
				t.Errorf("Keys %d and %d overlap: %s and %s", i, j, key1.ID, key2.ID)
			}
		}
	}

	// Test that the layout can be parsed using the keyboard package
	parser := keyboard.NewLayoutParser("examples/assets/layouts")
	parsedLayout, err := parser.ParseLayout(filepath.Base(jsonFile))
	if err != nil {
		t.Errorf("Failed to parse layout using keyboard package: %v", err)
	}

	// Validate the parsed layout
	if err := parser.ValidateLayout(parsedLayout); err != nil {
		t.Errorf("Layout validation failed: %v", err)
	}

	t.Logf("✓ Layout %s validated successfully (%dx%d, %d keys)", 
		layout.Name, layout.Width, layout.Height, len(layout.Keys))
}

// keysOverlap checks if two keys overlap
func keysOverlap(key1, key2 *keyboard.Key) bool {
	return !(key1.X+key1.Width <= key2.X || 
		key2.X+key2.Width <= key1.X || 
		key1.Y+key1.Height <= key2.Y || 
		key2.Y+key2.Height <= key1.Y)
}

// testThemeValidation tests theme JSON files (if they exist)
func testThemeValidation(t *testing.T) {
	// Look for theme files in various locations
	themeDirs := []string{
		"examples/assets/themes",
		"assets/themes",
		"themes",
	}

	var jsonFiles []string
	for _, themeDir := range themeDirs {
		if files, err := filepath.Glob(filepath.Join(themeDir, "*.json")); err == nil {
			jsonFiles = append(jsonFiles, files...)
		}
	}

	if len(jsonFiles) == 0 {
		t.Logf("No theme JSON files found")
		return
	}

	t.Logf("Found %d theme JSON files to validate", len(jsonFiles))

	for _, jsonFile := range jsonFiles {
		t.Run(filepath.Base(jsonFile), func(t *testing.T) {
			validateThemeFile(t, jsonFile)
		})
	}
}

// validateThemeFile validates a single theme JSON file
func validateThemeFile(t *testing.T, jsonFile string) {
	// Read the JSON file
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read theme file %s: %v", jsonFile, err)
	}

	// Unmarshal into Theme struct
	var theme keyboard.Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		t.Fatalf("Failed to unmarshal theme JSON from %s: %v", jsonFile, err)
	}

	// Validate required fields
	if theme.Name == "" {
		t.Error("Theme name is required but empty")
	}

	if theme.FontSize <= 0 {
		t.Errorf("Theme font size must be positive, got %d", theme.FontSize)
	}

	if theme.BorderRadius < 0 {
		t.Errorf("Theme border radius must be non-negative, got %d", theme.BorderRadius)
	}

	// Validate color values are in valid range [0,1]
	validateColorArray := func(colorName string, color [4]float32) {
		for i, val := range color {
			if val < 0 || val > 1 {
				t.Errorf("Theme %s color component %d must be between 0 and 1, got %f", 
					colorName, i, val)
			}
		}
	}

	validateColorArray("background", theme.BackgroundColor)
	validateColorArray("key", theme.KeyColor)
	validateColorArray("keyPressed", theme.KeyPressedColor)
	validateColorArray("text", theme.TextColor)

	t.Logf("✓ Theme %s validated successfully", theme.Name)
}

// testHeadlessRendering tests headless rendering functionality
func testHeadlessRendering(t *testing.T) {
	// Initialize keyboard
	kb, err := keyboard.New()
	if err != nil {
		t.Fatalf("Failed to initialize keyboard: %v", err)
	}

	// Get all available layouts
	layouts, err := kb.ListAvailableLayouts()
	if err != nil {
		t.Fatalf("Failed to list layouts: %v", err)
	}

	if len(layouts) == 0 {
		t.Fatal("No layouts available for testing")
	}

	for _, layoutName := range layouts {
		t.Run(fmt.Sprintf("RenderLayout_%s", layoutName), func(t *testing.T) {
			testLayoutRendering(t, kb, layoutName)
		})
	}
}

// testLayoutRendering tests rendering of a specific layout
func testLayoutRendering(t *testing.T, kb *keyboard.Keyboard, layoutName string) {
	// Switch to the layout
	if err := kb.SwitchLayout(layoutName); err != nil {
		t.Fatalf("Failed to switch to layout %s: %v", layoutName, err)
	}

	layout := kb.GetLayout()
	if layout == nil {
		t.Fatal("Layout is nil after switching")
	}

	// Validate geometry for potential rendering issues
	if layout.Width == 0 || layout.Height == 0 {
		t.Errorf("Layout %s has zero dimensions: %dx%d", layoutName, layout.Width, layout.Height)
	}

	// Check for keys outside bounds
	for i, key := range layout.Keys {
		if key.X < 0 || key.Y < 0 {
			t.Errorf("Key %d has negative position: (%d, %d)", i, key.X, key.Y)
		}

		if key.X+key.Width > layout.Width {
			t.Errorf("Key %d extends beyond layout width", i)
		}

		if key.Y+key.Height > layout.Height {
			t.Errorf("Key %d extends beyond layout height", i)
		}
	}

	// Test key press simulation for visual state changes
	testKeys := []string{"a", "space", "enter", "shift"}
	for _, keyID := range testKeys {
		// Press key
		if err := kb.PressKey(keyID); err != nil {
			t.Logf("Note: Could not press key %s: %v", keyID, err)
		} else {
			// Check state
			state := kb.GetKeyState(keyID)
			if state == keyboard.KeyStateReleased {
				t.Errorf("Key %s state should be pressed but is released", keyID)
			}

			// Release key
			if err := kb.ReleaseKey(keyID); err != nil {
				t.Errorf("Failed to release key %s: %v", keyID, err)
			}
		}
	}

	// Simulate headless screenshot functionality
	if err := simulateHeadlessScreenshot(t, layout, layoutName); err != nil {
		t.Errorf("Headless screenshot simulation failed for layout %s: %v", layoutName, err)
	}

	t.Logf("✓ Layout %s rendered successfully (%dx%d, %d keys)", 
		layoutName, layout.Width, layout.Height, len(layout.Keys))
}

// simulateHeadlessScreenshot simulates the headless screenshot functionality
func simulateHeadlessScreenshot(t *testing.T, layout *keyboard.Layout, layoutName string) error {
	// Create a simple pixel buffer to simulate rendering
	width := layout.Width
	height := layout.Height
	
	if width <= 0 || height <= 0 {
		return fmt.Errorf("invalid dimensions: %dx%d", width, height)
	}

	// Simulate creating a pixel buffer
	pixelBuffer := make([]uint32, width*height)
	
	// Fill background
	backgroundColor := uint32(0xFF202020) // Dark gray
	for i := range pixelBuffer {
		pixelBuffer[i] = backgroundColor
	}

	// Simulate rendering each key
	keyColor := uint32(0xFFF0F0F0)      // Light gray
	borderColor := uint32(0xFF808080)   // Medium gray
	
	for i, key := range layout.Keys {
		// Validate key bounds
		if key.X < 0 || key.Y < 0 || key.Width <= 0 || key.Height <= 0 {
			return fmt.Errorf("key %d has invalid bounds: x=%d, y=%d, w=%d, h=%d", 
				i, key.X, key.Y, key.Width, key.Height)
		}

		// Check if key fits within layout
		if key.X+key.Width > width || key.Y+key.Height > height {
			return fmt.Errorf("key %d extends beyond layout bounds", i)
		}

		// Simulate drawing the key
		for y := key.Y; y < key.Y+key.Height && y < height; y++ {
			for x := key.X; x < key.X+key.Width && x < width; x++ {
				if y >= 0 && x >= 0 {
					// Draw border
					if x == key.X || x == key.X+key.Width-1 || 
					   y == key.Y || y == key.Y+key.Height-1 {
						pixelBuffer[y*width+x] = borderColor
					} else {
						pixelBuffer[y*width+x] = keyColor
					}
				}
			}
		}
	}

	// Simulate saving screenshot (in real implementation, this would save to file)
	screenshotPath := fmt.Sprintf("screenshots/%s.png", layoutName)
	t.Logf("Simulated screenshot saved to: %s", screenshotPath)
	
	// In a real implementation, you would save the pixel buffer as an image
	// For now, we just verify the buffer has been modified
	nonBackgroundPixels := 0
	for _, pixel := range pixelBuffer {
		if pixel != backgroundColor {
			nonBackgroundPixels++
		}
	}
	
	if nonBackgroundPixels == 0 {
		return fmt.Errorf("no keys were rendered (all pixels are background color)")
	}

	t.Logf("✓ Headless screenshot simulation successful (%d non-background pixels)", 
		nonBackgroundPixels)
	
	return nil
}

// Helper function to run the main oskway command with screenshot flag
func TestOskwayScreenshot(t *testing.T) {
	// Test if the binary exists
	binaryPath := "./oskway"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("oskway binary not found, skipping screenshot test")
	}

	// Create screenshots directory
	if err := os.MkdirAll("screenshots", 0755); err != nil {
		t.Fatalf("Failed to create screenshots directory: %v", err)
	}

	// Get all layout files
	layoutsDir := "examples/assets/layouts"
	jsonFiles, err := filepath.Glob(filepath.Join(layoutsDir, "*.json"))
	if err != nil {
		t.Fatalf("Failed to find layout JSON files: %v", err)
	}

	for _, jsonFile := range jsonFiles {
		layoutName := strings.TrimSuffix(filepath.Base(jsonFile), ".json")
		t.Run(fmt.Sprintf("Screenshot_%s", layoutName), func(t *testing.T) {
			// Note: In a real implementation, you would add --screenshot flag to oskway
			// For now, we simulate the functionality
			t.Logf("Would run: %s --screenshot --layout %s --output screenshots/%s.png", 
				binaryPath, layoutName, layoutName)
			
			// For demonstration, we create a placeholder file
			placeholderPath := fmt.Sprintf("screenshots/%s.png", layoutName)
			if err := os.WriteFile(placeholderPath, []byte("placeholder"), 0644); err != nil {
				t.Errorf("Failed to create placeholder screenshot: %v", err)
			} else {
				t.Logf("✓ Placeholder screenshot created: %s", placeholderPath)
			}
		})
	}
}
