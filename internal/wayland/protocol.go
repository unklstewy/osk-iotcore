package wayland

// Event types for Wayland protocol
const (
	EventTypeRegistry = iota
	EventTypeCompositor
	EventTypeShell
	EventTypeOutput
	EventTypeSeat
	EventTypeKeyboard
	EventTypePointer
	EventTypeTouch
)

// Event represents a Wayland protocol event
type Event struct {
	Type uint32
	Data interface{}
}

// RegistryEvent contains registry event data
type RegistryEvent struct {
	Name      uint32
	Interface string
	Version   uint32
}

// KeyboardEvent contains keyboard event data
type KeyboardEvent struct {
	Serial uint32
	Key    uint32
	State  uint32
	Time   uint32
}

// PointerEvent contains pointer event data
type PointerEvent struct {
	Serial uint32
	X      int32
	Y      int32
	Button uint32
	State  uint32
	Time   uint32
}

// TouchEvent contains touch event data
type TouchEvent struct {
	Serial uint32
	ID     int32
	X      int32
	Y      int32
	Time   uint32
}

// EventHandler defines the interface for handling Wayland events
type EventHandler interface {
	HandleEvent(event *Event) error
}

// EventDispatcher manages event handling for the Wayland client
type EventDispatcher struct {
	handlers map[uint32]EventHandler
	eventCh  chan *Event
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[uint32]EventHandler),
		eventCh:  make(chan *Event, 100),
	}
}

// RegisterHandler registers an event handler for a specific event type
func (ed *EventDispatcher) RegisterHandler(eventType uint32, handler EventHandler) {
	ed.handlers[eventType] = handler
}

// DispatchEvent dispatches an event to the appropriate handler
func (ed *EventDispatcher) DispatchEvent(event *Event) error {
	if handler, exists := ed.handlers[event.Type]; exists {
		return handler.HandleEvent(event)
	}
	return nil
}

// EventChannel returns the event channel for receiving events
func (ed *EventDispatcher) EventChannel() <-chan *Event {
	return ed.eventCh
}

// SendEvent sends an event to the event channel
func (ed *EventDispatcher) SendEvent(event *Event) {
	select {
	case ed.eventCh <- event:
	default:
		// Drop event if channel is full
	}
}
