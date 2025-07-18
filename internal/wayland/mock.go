package wayland

import (
	"fmt"
	"time"
)

// MockClient represents a mock Wayland client for testing
type MockClient struct {
	running bool
}

// NewMockClient creates a new mock Wayland client
func NewMockClient() *MockClient {
	return &MockClient{
		running: true,
	}
}

// Close cleans up the mock client
func (c *MockClient) Close() {
	c.running = false
}

// CreateSurface creates a mock surface
func (c *MockClient) CreateSurface() error {
	fmt.Println("Mock: Created surface")
	return nil
}

// Dispatch processes mock events
func (c *MockClient) Dispatch() error {
	if !c.running {
		return fmt.Errorf("mock client closed")
	}
	
	// Simulate event processing delay
	time.Sleep(16 * time.Millisecond) // ~60 FPS
	return nil
}

// Flush mock flush operation
func (c *MockClient) Flush() error {
	return nil
}

// IsRunning returns whether the mock client is running
func (c *MockClient) IsRunning() bool {
	return c.running
}
