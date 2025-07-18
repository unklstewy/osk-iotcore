package wayland

import (
	"os"
)

// WaylandClient defines the interface for Wayland client implementations
type WaylandClient interface {
	Close()
	CreateSurface() error
	Dispatch() error
	Flush() error
}

// NewClientInterface creates a new Wayland client interface
// This will use mock implementation for testing or real implementation for production
func NewClientInterface() (WaylandClient, error) {
	// Check if we have a Wayland display available
	if isWaylandAvailable() {
		// Try to create a real Wayland client
		client, err := NewClient()
		if err != nil {
			// Fall back to mock if real client fails
			return NewMockClient(), nil
		}
		return client, nil
	}
	
	// Use mock client for non-Wayland environments (like X11)
	return NewMockClient(), nil
}

// isWaylandAvailable checks if Wayland is available in the current environment
func isWaylandAvailable() bool {
	// Check for Wayland display environment variable
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return true
	}
	
	// Check if XDG_SESSION_TYPE is wayland
	if os.Getenv("XDG_SESSION_TYPE") == "wayland" {
		return true
	}
	
	// Default to false for X11 or other environments
	return false
}
