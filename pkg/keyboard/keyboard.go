// Package keyboard provides keyboard layout parsing, key state management,
// and theming functionality for the on-screen keyboard.
package keyboard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// KeyState represents the state of a key
type KeyState int

const (
	KeyStateReleased KeyState = iota
	KeyStatePressed
	KeyStateRepeating
)

// Key represents a single key on the keyboard
type Key struct {
	ID       string   `json:"id"`
	Label    string   `json:"label"`
	Code     int32    `json:"code"`  // Changed to int32 to handle negative codes for special keys
	X        int      `json:"x"`
	Y        int      `json:"y"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	State    KeyState `json:"-"`       // Not serialized
	Modifier bool     `json:"modifier,omitempty"`
}

// Layout represents a keyboard layout
type Layout struct {
	Name        string `json:"name"`
	Keys        []*Key `json:"keys"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Description string `json:"description,omitempty"`
}

// Theme represents keyboard visual theme
type Theme struct {
	Name            string     `json:"name"`
	BackgroundColor [4]float32 `json:"background_color"`
	KeyColor        [4]float32 `json:"key_color"`
	KeyPressedColor [4]float32 `json:"key_pressed_color"`
	TextColor       [4]float32 `json:"text_color"`
	FontSize        int        `json:"font_size"`
	BorderRadius    int        `json:"border_radius"`
}

// Keyboard manages keyboard state and layout
type Keyboard struct {
	layout     *Layout
	theme      *Theme
	keyStates  map[string]KeyState
	mutex      sync.RWMutex
	callbacks  map[string]func(*Key)
}

// New creates a new keyboard instance
func New() (*Keyboard, error) {
	kb := &Keyboard{
		keyStates: make(map[string]KeyState),
		callbacks: make(map[string]func(*Key)),
	}

	// Load default layout
	if err := kb.LoadLayout("qwerty"); err != nil {
		return nil, fmt.Errorf("failed to load default layout: %w", err)
	}

	// Load glass theme
	if err := kb.LoadTheme("glass"); err != nil {
		return nil, fmt.Errorf("failed to load glass theme: %w", err)
	}

	return kb, nil
}

// LoadLayout loads a keyboard layout by name from the layout.d directory
func (kb *Keyboard) LoadLayout(name string) error {
	kb.mutex.Lock()
	defer kb.mutex.Unlock()

	// Try to load from layout.d directory first
	parser := NewLayoutParser("assets/layouts")
	layout, err := parser.ParseLayout(name + ".json")
	if err != nil {
		// Fallback to built-in layout if file not found
		if name == "qwerty" {
			layout = &Layout{
				Name:   name,
				Width:  900,
				Height: 400,
				Keys:   createQWERTYLayout(),
			}
		} else {
			return fmt.Errorf("layout %s not found: %w", name, err)
		}
	}

	// Validate the loaded layout
	if err := parser.ValidateLayout(layout); err != nil {
		return fmt.Errorf("invalid layout %s: %w", name, err)
	}

	kb.layout = layout
	return nil
}

// LoadTheme loads a visual theme by name
func (kb *Keyboard) LoadTheme(name string) error {
	kb.mutex.Lock()
	defer kb.mutex.Unlock()

	// Try to load theme from JSON file
	theme, err := loadThemeFromFile(name)
	if err != nil {
		// Fallback to default theme if file not found
		theme = &Theme{
			Name:            name,
			BackgroundColor: [4]float32{0.2, 0.2, 0.2, 1.0},
			KeyColor:        [4]float32{0.8, 0.8, 0.8, 1.0},
			KeyPressedColor: [4]float32{0.6, 0.6, 0.6, 1.0},
			TextColor:       [4]float32{0.0, 0.0, 0.0, 1.0},
			FontSize:        16,
			BorderRadius:    4,
		}
	}

	kb.theme = theme
	return nil
}

// GetLayout returns the current keyboard layout
func (kb *Keyboard) GetLayout() *Layout {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()
	return kb.layout
}

// GetTheme returns the current keyboard theme
func (kb *Keyboard) GetTheme() *Theme {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()
	return kb.theme
}

// PressKey sets a key to pressed state
func (kb *Keyboard) PressKey(keyID string) error {
	kb.mutex.Lock()
	defer kb.mutex.Unlock()

	kb.keyStates[keyID] = KeyStatePressed
	
	// Find the key and update its state
	for _, key := range kb.layout.Keys {
		if key.ID == keyID {
			key.State = KeyStatePressed
			// Call callback if registered
			if callback, exists := kb.callbacks[keyID]; exists {
				callback(key)
			}
			break
		}
	}

	return nil
}

// ReleaseKey sets a key to released state
func (kb *Keyboard) ReleaseKey(keyID string) error {
	kb.mutex.Lock()
	defer kb.mutex.Unlock()

	kb.keyStates[keyID] = KeyStateReleased
	
	// Find the key and update its state
	for _, key := range kb.layout.Keys {
		if key.ID == keyID {
			key.State = KeyStateReleased
			break
		}
	}

	return nil
}

// GetKeyState returns the state of a specific key
func (kb *Keyboard) GetKeyState(keyID string) KeyState {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()
	
	if state, exists := kb.keyStates[keyID]; exists {
		return state
	}
	return KeyStateReleased
}

// RegisterCallback registers a callback for key events
func (kb *Keyboard) RegisterCallback(keyID string, callback func(*Key)) {
	kb.mutex.Lock()
	defer kb.mutex.Unlock()
	kb.callbacks[keyID] = callback
}

// ListAvailableLayouts returns a list of available keyboard layouts
func (kb *Keyboard) ListAvailableLayouts() ([]string, error) {
	parser := NewLayoutParser("assets/layouts")
	layouts, err := parser.ListLayouts()
	if err != nil {
		return nil, fmt.Errorf("failed to list layouts: %w", err)
	}

	// Remove .json extension from layout names
	var layoutNames []string
	for _, layout := range layouts {
		if len(layout) > 5 && layout[len(layout)-5:] == ".json" {
			layoutNames = append(layoutNames, layout[:len(layout)-5])
		}
	}

	// Add built-in layouts
	builtInLayouts := []string{"qwerty"}
	for _, builtIn := range builtInLayouts {
		found := false
		for _, existing := range layoutNames {
			if existing == builtIn {
				found = true
				break
			}
		}
		if !found {
			layoutNames = append(layoutNames, builtIn)
		}
	}

	return layoutNames, nil
}

// SwitchLayout switches to a different keyboard layout
func (kb *Keyboard) SwitchLayout(layoutName string) error {
	return kb.LoadLayout(layoutName)
}

// GetCurrentLayoutName returns the name of the current layout
func (kb *Keyboard) GetCurrentLayoutName() string {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()
	if kb.layout != nil {
		return kb.layout.Name
	}
	return ""
}

// RefreshLayout reloads the current layout from disk
func (kb *Keyboard) RefreshLayout() error {
	currentName := kb.GetCurrentLayoutName()
	if currentName == "" {
		return fmt.Errorf("no current layout to refresh")
	}
	return kb.LoadLayout(currentName)
}

// createQWERTYLayout creates a basic QWERTY keyboard layout
func createQWERTYLayout() []*Key {
	keys := []*Key{
		// Top row
		{ID: "q", Label: "Q", Code: 16, X: 10, Y: 10, Width: 60, Height: 60},
		{ID: "w", Label: "W", Code: 17, X: 80, Y: 10, Width: 60, Height: 60},
		{ID: "e", Label: "E", Code: 18, X: 150, Y: 10, Width: 60, Height: 60},
		{ID: "r", Label: "R", Code: 19, X: 220, Y: 10, Width: 60, Height: 60},
		{ID: "t", Label: "T", Code: 20, X: 290, Y: 10, Width: 60, Height: 60},
		{ID: "y", Label: "Y", Code: 21, X: 360, Y: 10, Width: 60, Height: 60},
		{ID: "u", Label: "U", Code: 22, X: 430, Y: 10, Width: 60, Height: 60},
		{ID: "i", Label: "I", Code: 23, X: 500, Y: 10, Width: 60, Height: 60},
		{ID: "o", Label: "O", Code: 24, X: 570, Y: 10, Width: 60, Height: 60},
		{ID: "p", Label: "P", Code: 25, X: 640, Y: 10, Width: 60, Height: 60},

		// Middle row
		{ID: "a", Label: "A", Code: 30, X: 45, Y: 80, Width: 60, Height: 60},
		{ID: "s", Label: "S", Code: 31, X: 115, Y: 80, Width: 60, Height: 60},
		{ID: "d", Label: "D", Code: 32, X: 185, Y: 80, Width: 60, Height: 60},
		{ID: "f", Label: "F", Code: 33, X: 255, Y: 80, Width: 60, Height: 60},
		{ID: "g", Label: "G", Code: 34, X: 325, Y: 80, Width: 60, Height: 60},
		{ID: "h", Label: "H", Code: 35, X: 395, Y: 80, Width: 60, Height: 60},
		{ID: "j", Label: "J", Code: 36, X: 465, Y: 80, Width: 60, Height: 60},
		{ID: "k", Label: "K", Code: 37, X: 535, Y: 80, Width: 60, Height: 60},
		{ID: "l", Label: "L", Code: 38, X: 605, Y: 80, Width: 60, Height: 60},

		// Bottom row
		{ID: "z", Label: "Z", Code: 44, X: 80, Y: 150, Width: 60, Height: 60},
		{ID: "x", Label: "X", Code: 45, X: 150, Y: 150, Width: 60, Height: 60},
		{ID: "c", Label: "C", Code: 46, X: 220, Y: 150, Width: 60, Height: 60},
		{ID: "v", Label: "V", Code: 47, X: 290, Y: 150, Width: 60, Height: 60},
		{ID: "b", Label: "B", Code: 48, X: 360, Y: 150, Width: 60, Height: 60},
		{ID: "n", Label: "N", Code: 49, X: 430, Y: 150, Width: 60, Height: 60},
		{ID: "m", Label: "M", Code: 50, X: 500, Y: 150, Width: 60, Height: 60},

		// Special keys
		{ID: "space", Label: "Space", Code: 57, X: 200, Y: 220, Width: 300, Height: 60},
		{ID: "shift", Label: "Shift", Code: 42, X: 10, Y: 150, Width: 60, Height: 60, Modifier: true},
		{ID: "backspace", Label: "⌫", Code: 14, X: 710, Y: 10, Width: 80, Height: 60},
		{ID: "enter", Label: "Enter", Code: 28, X: 675, Y: 80, Width: 115, Height: 60},
	}

	return keys
}

// loadThemeFromFile loads a theme from a JSON file
func loadThemeFromFile(name string) (*Theme, error) {
	// Try different theme directories
	themeDirs := []string{
		"assets/themes",
		"themes",
		"../assets/themes",
	}
	
	for _, themeDir := range themeDirs {
		themePath := filepath.Join(themeDir, name+".json")
		if _, err := os.Stat(themePath); err == nil {
			// File exists, try to load it
			data, err := os.ReadFile(themePath)
			if err != nil {
				continue
			}
			
			var theme Theme
			if err := json.Unmarshal(data, &theme); err != nil {
				continue
			}
			
			return &theme, nil
		}
	}
	
	return nil, fmt.Errorf("theme file %s.json not found", name)
}
