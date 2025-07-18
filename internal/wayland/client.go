//go:build !test
// +build !test

package wayland

import (
	"fmt"
)

/*
#cgo pkg-config: wayland-client
#cgo CFLAGS: -DWLR_USE_UNSTABLE
#include <wayland-client.h>
#include <stdlib.h>
*/
import "C"

// Client represents a Wayland client connection
type Client struct {
	display    *C.struct_wl_display
	registry   *C.struct_wl_registry
	compositor *C.struct_wl_compositor
	surface    *C.struct_wl_surface
	shell      *C.struct_wl_shell
	shellSurf  *C.struct_wl_shell_surface
}

// NewClient creates a new Wayland client connection
func NewClient() (*Client, error) {
	display := C.wl_display_connect(nil)
	if display == nil {
		return nil, fmt.Errorf("failed to connect to Wayland display")
	}

	client := &Client{
		display: display,
	}

	// Get the registry
	client.registry = C.wl_display_get_registry(display)
	if client.registry == nil {
		client.Close()
		return nil, fmt.Errorf("failed to get Wayland registry")
	}

	return client, nil
}

// Close cleans up the Wayland client connection
func (c *Client) Close() {
	if c.shellSurf != nil {
		C.wl_shell_surface_destroy(c.shellSurf)
	}
	if c.surface != nil {
		C.wl_surface_destroy(c.surface)
	}
	if c.shell != nil {
		C.wl_shell_destroy(c.shell)
	}
	if c.compositor != nil {
		C.wl_compositor_destroy(c.compositor)
	}
	if c.registry != nil {
		C.wl_registry_destroy(c.registry)
	}
	if c.display != nil {
		C.wl_display_disconnect(c.display)
	}
}

// CreateSurface creates a new Wayland surface
func (c *Client) CreateSurface() error {
	if c.compositor == nil {
		return fmt.Errorf("compositor not available")
	}

	c.surface = C.wl_compositor_create_surface(c.compositor)
	if c.surface == nil {
		return fmt.Errorf("failed to create surface")
	}

	return nil
}

// Dispatch processes pending Wayland events
func (c *Client) Dispatch() error {
	ret := C.wl_display_dispatch(c.display)
	if ret == -1 {
		return fmt.Errorf("failed to dispatch Wayland events")
	}
	return nil
}

// Flush sends buffered requests to the Wayland server
func (c *Client) Flush() error {
	ret := C.wl_display_flush(c.display)
	if ret == -1 {
		return fmt.Errorf("failed to flush Wayland display")
	}
	return nil
}
