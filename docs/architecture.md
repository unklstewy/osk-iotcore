# Architecture Overview

OSK IoT Core is a high-performance on-screen keyboard designed for Linux environments with a focus on kiosk applications. The architecture emphasizes modularity, performance, and extensibility.

## Table of Contents

- [High-Level Design](#high-level-design)
- [Package Responsibilities](#package-responsibilities)
- [Component Interactions](#component-interactions)
- [Wayland Protocol Handling](#wayland-protocol-handling)
- [Rendering Architecture](#rendering-architecture)
- [Event System](#event-system)
- [Build System](#build-system)
- [Deployment Architecture](#deployment-architecture)

## High-Level Design

OSK IoT Core follows a layered architecture with clear separation of concerns:

1. **Application Layer**: Main entry point and command-line interface
2. **UI Layer**: User interface components and event handling
3. **Keyboard Layer**: Layout management and input processing
4. **Platform Layer**: Wayland/X11 protocol handling and rendering
5. **System Layer**: CGO bindings and system integration

### Design Principles

- **Modularity**: Each component has a well-defined responsibility
- **Performance**: Leverages Go's concurrency for responsive input handling
- **Portability**: Abstracted platform-specific code for multi-display server support
- **Extensibility**: Plugin architecture for themes and layouts
- **Reliability**: Comprehensive error handling and graceful degradation

## Package Responsibilities

### Core Application

- **`cmd/oskway/`**: Main application entry point and command-line interface
  - Application lifecycle management
  - Command-line argument parsing
  - Signal handling and graceful shutdown
  - Mode selection (test, production, Wayland PoC)

### Internal Packages

- **`internal/render/`**: Rendering abstraction layer
  - OpenGL rendering backend (`opengl.go`)
  - Vulkan rendering backend (`vulkan.go`)
  - Texture management (`texture.go`)
  - Hardware acceleration support
  - Cross-platform rendering abstractions

- **`internal/wayland/`**: Wayland protocol implementation
  - Client connection management (`client.go`)
  - Protocol interface definitions (`interface.go`)
  - Event handling and dispatching (`protocol.go`)
  - Mock implementation for testing (`mock.go`)
  - Surface creation and management
  - Input method protocol support

### Public Packages

- **`pkg/keyboard/`**: Core keyboard functionality
  - Layout management and switching (`keyboard.go`)
  - Key event processing and dispatching
  - Layout file parsing (`parser.go`)
  - Theme system integration
  - Multi-language support
  - Callback system for key events

- **`ui/`**: User interface components
  - Application window management (`app.go`)
  - Widget system and rendering (`widget.go`)
  - Touch and mouse event handling
  - Theme application and styling
  - Layout positioning and sizing

### Assets and Configuration

- **`assets/layouts/`**: Keyboard layout definitions
  - QWERTY layout (`qwerty.json`)
  - Dvorak layout (`dvorak.json`)
  - Extensible JSON-based format

- **`assets/themes/`**: Visual themes and styling
  - Default theme (`default.json`)
  - Dark theme (`dark.json`)
  - Color schemes and styling rules

## Component Interactions

```
┌─────────────────────────────────────────────────────────────────┐
│                        OSK IoT Core                            │
├─────────────────────────────────────────────────────────────────┤
│  cmd/oskway (Main Application)                                 │
│  ├─ Command Line Interface                                     │
│  ├─ Application Lifecycle                                      │
│  └─ Signal Handling                                            │
├─────────────────────────────────────────────────────────────────┤
│  UI Layer                                                      │
│  ├─ ui/app.go (Application Window)                             │
│  ├─ ui/widget.go (Widget System)                               │
│  └─ Event Handling                                             │
├─────────────────────────────────────────────────────────────────┤
│  Keyboard Layer                                                │
│  ├─ pkg/keyboard/keyboard.go (Core Logic)                      │
│  ├─ pkg/keyboard/parser.go (Layout Parser)                     │
│  ├─ Layout Management                                          │
│  └─ Event Processing                                           │
├─────────────────────────────────────────────────────────────────┤
│  Platform Layer                                                │
│  ├─ internal/wayland/ (Wayland Support)                        │
│  │   ├─ client.go (Connection Management)                      │
│  │   ├─ protocol.go (Event Handling)                          │
│  │   └─ interface.go (Protocol Definitions)                   │
│  └─ internal/render/ (Rendering Backends)                      │
│      ├─ opengl.go (OpenGL Backend)                             │
│      ├─ vulkan.go (Vulkan Backend)                             │
│      └─ texture.go (Texture Management)                        │
└─────────────────────────────────────────────────────────────────┘
```

## Wayland Protocol Handling

### CGO Integration

The project uses CGO with wlroots bindings for comprehensive Wayland support:

```c
// CGO bindings in cmd/oskway/main.go
#cgo pkg-config: wayland-client
#include <wayland-client.h>
#include <wayland-client-protocol.h>
```

### Protocol Support

- **Core Protocols**: `wl_display`, `wl_registry`, `wl_compositor`, `wl_surface`
- **Shell Protocols**: Layer shell for proper keyboard positioning
- **Input Protocols**: `wl_seat`, `wl_keyboard`, `wl_pointer`, `wl_touch`
- **Buffer Management**: Shared memory buffers for efficient rendering

### Event System

```go
// Event dispatcher in internal/wayland/protocol.go
type EventDispatcher struct {
    handlers map[uint32]EventHandler
    eventCh  chan *Event
}
```

### Surface Management

- **Surface Creation**: Automated surface creation with proper damage tracking
- **Buffer Management**: Efficient shared memory buffer allocation
- **Cleanup**: Proper resource cleanup to prevent memory leaks

## Rendering Architecture

### Multi-Backend Support

- **OpenGL Backend**: Hardware-accelerated rendering with OpenGL ES
- **Vulkan Backend**: Next-generation graphics API support
- **Texture Management**: Efficient texture loading and caching

### Rendering Pipeline

1. **Layout Calculation**: Determine key positions and sizes
2. **Texture Loading**: Load theme assets and generate textures
3. **Rendering**: Draw keyboard layout with proper blending
4. **Buffer Swap**: Present rendered frame to display

## Event System

### Input Event Flow

```
Wayland Display Server
         ↓
   Event Dispatcher
         ↓
   Event Handlers
         ↓
   UI Components
         ↓
   Keyboard Logic
         ↓
   Application Callbacks
```

### Event Types

- **Registry Events**: Wayland interface discovery
- **Keyboard Events**: Key press/release events
- **Pointer Events**: Mouse movement and clicks
- **Touch Events**: Touch screen interactions

## Build System

### Makefile Structure

```makefile
# Core targets
build: deps fmt lint test compile
run: build execute
test: unit-tests integration-tests
lint: golangci-lint security-checks
fmt: goimports gofmt
```

### Build Modes

- **Development Build**: Includes debug symbols and race detection
- **Production Build**: Optimized binary with stripped symbols
- **Docker Build**: Multi-stage containerized build process

### CGO Configuration

```bash
# pkg-config integration
CGO_CFLAGS="$(pkg-config --cflags wayland-client)"
CGO_LDFLAGS="$(pkg-config --libs wayland-client)"
```

## Deployment Architecture

### Target Environments

- **Ubuntu Core 24+**: Primary deployment target
- **Kiosk Mode**: Fullscreen application deployment
- **Wayland Compositors**: Sway, Mutter, KWin compatibility

### System Integration

- **systemd Integration**: Service management and autostart
- **Session Management**: XDG session integration
- **Resource Management**: Memory and CPU usage optimization

### Security Considerations

- **Privilege Separation**: Minimal required permissions
- **Input Validation**: Comprehensive input sanitization
- **Memory Safety**: Proper resource management and cleanup

## Testing Architecture

### Test Categories

- **Unit Tests**: Individual component testing
- **Integration Tests**: Component interaction testing
- **Mock Tests**: Wayland protocol mocking for CI/CD
- **Performance Tests**: Benchmarking and profiling

### Test Infrastructure

```go
// Mock Wayland client for testing
type MockWaylandClient struct {
    events chan *Event
    state  *MockState
}
```

## Performance Considerations

### Optimization Strategies

- **Goroutine Pooling**: Efficient concurrent event handling
- **Buffer Reuse**: Minimize memory allocations
- **Texture Caching**: Reduce GPU memory usage
- **Event Batching**: Optimize input event processing

### Resource Management

- **Memory**: Proper cleanup of Wayland resources
- **CPU**: Efficient event loop implementation
- **GPU**: Optimized rendering pipeline

## Future Architecture Considerations

### Planned Enhancements

- **Plugin System**: Dynamic layout and theme loading
- **IPC Integration**: Inter-process communication for external apps
- **Accessibility**: Screen reader and assistive technology support
- **Multi-Display**: Support for multiple monitor configurations

### Scalability

- **Modular Plugins**: Hot-swappable components
- **Configuration API**: Runtime configuration changes
- **Performance Monitoring**: Built-in metrics and profiling

---

*This architecture documentation is maintained alongside the codebase and should be updated when significant architectural changes are made.*
