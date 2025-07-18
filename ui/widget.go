package ui

import (
	"fmt"

	"github.com/iotcore/osk-iotcore/internal/render"
	"github.com/iotcore/osk-iotcore/internal/wayland"
	"github.com/iotcore/osk-iotcore/pkg/keyboard"
)

// Widget represents a UI widget
type Widget interface {
	Render() error
	HandlePointerEvent(event *wayland.PointerEvent) error
	HandleKeyboardEvent(event *wayland.KeyboardEvent) error
	HandleTouchEvent(event *wayland.TouchEvent) error
}

// KeyboardWidget represents the on-screen keyboard widget
type KeyboardWidget struct {
	keyboard *keyboard.Keyboard
	renderer render.Renderer
	x        int
	y        int
	width    int
	height   int
}

// NewKeyboardWidget creates a new keyboard widget
func NewKeyboardWidget(kb *keyboard.Keyboard, renderer render.Renderer) *KeyboardWidget {
	layout := kb.GetLayout()
	return &KeyboardWidget{
		keyboard: kb,
		renderer: renderer,
		x:        0,
		y:        0,
		width:    layout.Width,
		height:   layout.Height,
	}
}

// Render renders the keyboard widget
func (kw *KeyboardWidget) Render() error {
	layout := kw.keyboard.GetLayout()
	theme := kw.keyboard.GetTheme()

	// Render background
	if err := kw.renderBackground(theme); err != nil {
		return fmt.Errorf("failed to render background: %w", err)
	}

	// Render each key
	for _, key := range layout.Keys {
		if err := kw.renderKey(key, theme); err != nil {
			return fmt.Errorf("failed to render key %s: %w", key.ID, err)
		}
	}

	return nil
}

// renderBackground renders the keyboard background
func (kw *KeyboardWidget) renderBackground(theme *keyboard.Theme) error {
	// Background rendering would be implemented here
	// This is a placeholder for the actual rendering logic
	return nil
}

// renderKey renders a single key
func (kw *KeyboardWidget) renderKey(key *keyboard.Key, theme *keyboard.Theme) error {
	// Choose color based on key state
	var color [4]float32
	switch key.State {
	case keyboard.KeyStatePressed:
		color = theme.KeyPressedColor
	default:
		color = theme.KeyColor
	}

	// Render key background (would use actual graphics API)
	// This is a placeholder for the actual rendering logic
	// TODO: Use the color variable when implementing actual rendering
	_ = color

	// Render key label
	textX := key.X + key.Width/2
	textY := key.Y + key.Height/2
	return kw.renderer.RenderText(textX, textY, key.Label, theme.TextColor)
}

// HandlePointerEvent handles pointer events for the keyboard widget
func (kw *KeyboardWidget) HandlePointerEvent(event *wayland.PointerEvent) error {
	if event.Button == 1 { // Left mouse button
		// Find which key was clicked
		key := kw.findKeyAtPosition(int(event.X), int(event.Y))
		if key != nil {
			if event.State == 1 { // Button pressed
				return kw.keyboard.PressKey(key.ID)
			} else { // Button released
				return kw.keyboard.ReleaseKey(key.ID)
			}
		}
	}
	return nil
}

// HandleKeyboardEvent handles keyboard events for the keyboard widget
func (kw *KeyboardWidget) HandleKeyboardEvent(event *wayland.KeyboardEvent) error {
	// This could be used for physical keyboard input to update the virtual keyboard
	// For now, we'll just log it
	fmt.Printf("Keyboard event: key=%d, state=%d\n", event.Key, event.State)
	return nil
}

// HandleTouchEvent handles touch events for the keyboard widget
func (kw *KeyboardWidget) HandleTouchEvent(event *wayland.TouchEvent) error {
	// Find which key was touched
	key := kw.findKeyAtPosition(int(event.X), int(event.Y))
	if key != nil {
		// For touch events, we'll press on touch down and release on touch up
		// This is a simplified implementation
		return kw.keyboard.PressKey(key.ID)
	}
	return nil
}

// findKeyAtPosition finds the key at the given screen coordinates
func (kw *KeyboardWidget) findKeyAtPosition(x, y int) *keyboard.Key {
	layout := kw.keyboard.GetLayout()
	
	// Adjust coordinates relative to widget position
	relX := x - kw.x
	relY := y - kw.y

	for _, key := range layout.Keys {
		if relX >= key.X && relX < key.X+key.Width &&
		   relY >= key.Y && relY < key.Y+key.Height {
			return key
		}
	}
	return nil
}

// SetPosition sets the position of the keyboard widget
func (kw *KeyboardWidget) SetPosition(x, y int) {
	kw.x = x
	kw.y = y
}

// GetPosition returns the position of the keyboard widget
func (kw *KeyboardWidget) GetPosition() (int, int) {
	return kw.x, kw.y
}

// GetSize returns the size of the keyboard widget
func (kw *KeyboardWidget) GetSize() (int, int) {
	return kw.width, kw.height
}
