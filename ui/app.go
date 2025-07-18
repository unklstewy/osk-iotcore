// Package ui provides the higher-level widget and event system for the
// on-screen keyboard application.
package ui

import (
	"fmt"
	"log"
	"sync"

	"github.com/iotcore/osk-iotcore/internal/render"
	"github.com/iotcore/osk-iotcore/internal/wayland"
	"github.com/iotcore/osk-iotcore/pkg/keyboard"
)

// App represents the main application
type App struct {
	keyboard   *keyboard.Keyboard
	renderer   render.Renderer
	waylandClient wayland.WaylandClient
	eventDispatcher *wayland.EventDispatcher
	running    bool
	mutex      sync.RWMutex
	widgets    []Widget
	keyboardWidget *KeyboardWidget
}

// NewApp creates a new application instance
func NewApp(kb *keyboard.Keyboard) *App {
	return &App{
		keyboard: kb,
		widgets:  make([]Widget, 0),
	}
}

// Run starts the application main loop
func (app *App) Run() error {
	if err := app.initialize(); err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}
	defer app.cleanup()

	log.Println("Starting oskway keyboard application...")
	
	app.mutex.Lock()
	app.running = true
	app.mutex.Unlock()

	// Main event loop
	for app.isRunning() {
		if err := app.processEvents(); err != nil {
			log.Printf("Error processing events: %v", err)
		}

		if err := app.render(); err != nil {
			log.Printf("Error rendering: %v", err)
		}

		if err := app.waylandClient.Dispatch(); err != nil {
			log.Printf("Error dispatching Wayland events: %v", err)
		}
	}

	return nil
}

// Stop stops the application
func (app *App) Stop() {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.running = false
}

// isRunning returns whether the application is running
func (app *App) isRunning() bool {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	return app.running
}

// initialize sets up the application
func (app *App) initialize() error {
	// Initialize Wayland client
	client, err := wayland.NewClientInterface()
	if err != nil {
		return fmt.Errorf("failed to create Wayland client: %w", err)
	}
	app.waylandClient = client

	// Initialize event dispatcher
	app.eventDispatcher = wayland.NewEventDispatcher()

	// Initialize renderer (OpenGL by default)
	app.renderer = &render.OpenGLRenderer{}
	if err := app.renderer.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize renderer: %w", err)
	}

	// Create keyboard widget
	app.keyboardWidget = NewKeyboardWidget(app.keyboard, app.renderer)
	app.widgets = append(app.widgets, app.keyboardWidget)

	// Setup event handlers
	app.setupEventHandlers()

	return nil
}

// cleanup cleans up application resources
func (app *App) cleanup() {
	if app.renderer != nil {
		app.renderer.Close()
	}
	if app.waylandClient != nil {
		app.waylandClient.Close()
	}
}

// processEvents processes pending events
func (app *App) processEvents() error {
	select {
	case event := <-app.eventDispatcher.EventChannel():
		return app.handleEvent(event)
	default:
		return nil
	}
}

// handleEvent handles a single event
func (app *App) handleEvent(event *wayland.Event) error {
	switch event.Type {
	case wayland.EventTypePointer:
		return app.handlePointerEvent(event)
	case wayland.EventTypeKeyboard:
		return app.handleKeyboardEvent(event)
	case wayland.EventTypeTouch:
		return app.handleTouchEvent(event)
	}
	return nil
}

// handlePointerEvent handles pointer events (mouse clicks)
func (app *App) handlePointerEvent(event *wayland.Event) error {
	pointerEvent, ok := event.Data.(*wayland.PointerEvent)
	if !ok {
		return fmt.Errorf("invalid pointer event data")
	}

	// Forward to keyboard widget
	return app.keyboardWidget.HandlePointerEvent(pointerEvent)
}

// handleKeyboardEvent handles keyboard events
func (app *App) handleKeyboardEvent(event *wayland.Event) error {
	keyboardEvent, ok := event.Data.(*wayland.KeyboardEvent)
	if !ok {
		return fmt.Errorf("invalid keyboard event data")
	}

	// Forward to keyboard widget
	return app.keyboardWidget.HandleKeyboardEvent(keyboardEvent)
}

// handleTouchEvent handles touch events
func (app *App) handleTouchEvent(event *wayland.Event) error {
	touchEvent, ok := event.Data.(*wayland.TouchEvent)
	if !ok {
		return fmt.Errorf("invalid touch event data")
	}

	// Forward to keyboard widget
	return app.keyboardWidget.HandleTouchEvent(touchEvent)
}

// render renders the application
func (app *App) render() error {
	// Render all widgets
	for _, widget := range app.widgets {
		if err := widget.Render(); err != nil {
			return fmt.Errorf("failed to render widget: %w", err)
		}
	}

	// Flush renderer
	if err := app.waylandClient.Flush(); err != nil {
		return fmt.Errorf("failed to flush Wayland client: %w", err)
	}

	return nil
}

// setupEventHandlers sets up event handlers for the application
func (app *App) setupEventHandlers() {
	// Register event handlers with the event dispatcher
	app.eventDispatcher.RegisterHandler(wayland.EventTypePointer, &PointerEventHandler{app: app})
	app.eventDispatcher.RegisterHandler(wayland.EventTypeKeyboard, &KeyboardEventHandler{app: app})
	app.eventDispatcher.RegisterHandler(wayland.EventTypeTouch, &TouchEventHandler{app: app})
}

// PointerEventHandler handles pointer events
type PointerEventHandler struct {
	app *App
}

func (h *PointerEventHandler) HandleEvent(event *wayland.Event) error {
	return h.app.handlePointerEvent(event)
}

// KeyboardEventHandler handles keyboard events
type KeyboardEventHandler struct {
	app *App
}

func (h *KeyboardEventHandler) HandleEvent(event *wayland.Event) error {
	return h.app.handleKeyboardEvent(event)
}

// TouchEventHandler handles touch events
type TouchEventHandler struct {
	app *App
}

func (h *TouchEventHandler) HandleEvent(event *wayland.Event) error {
	return h.app.handleTouchEvent(event)
}
