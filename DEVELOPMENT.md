# Development Setup

This document describes the development environment and tooling for the osk-iotcore project.

## Prerequisites

- Go 1.22 or later
- Make
- Docker (for containerized development)
- Git

## Quick Start

### Local Development

1. **Install dependencies:**
   ```bash
   make deps
   ```

2. **Install development tools:**
   ```bash
   make install-tools
   ```

3. **Format code:**
   ```bash
   make fmt
   ```

4. **Run linter:**
   ```bash
   make lint
   ```

5. **Run tests:**
   ```bash
   make test
   ```

6. **Build the application:**
   ```bash
   make build
   ```

7. **Run the application:**
   ```bash
   make run
   ```

### Docker Development

1. **Start development container:**
   ```bash
   docker-compose -f docker-compose.dev.yml up -d dev
   ```

2. **Enter the development container:**
   ```bash
   docker-compose -f docker-compose.dev.yml exec dev bash
   ```

3. **Inside the container, you can use all make targets:**
   ```bash
   make build
   make test
   make lint
   make fmt
   ```

## Make Targets

The project includes a comprehensive Makefile with the following targets:

### Primary Targets
- `build` - Build the application
- `run` - Build and run the application
- `test` - Run tests with race detection and coverage
- `lint` - Run golangci-lint
- `fmt` - Format code with goimports and gofmt
- `deps` - Install/update dependencies

### Additional Targets
- `install-tools` - Install development tools
- `clean` - Clean build artifacts
- `dev-build` - Build with debug symbols
- `dev-run` - Run in development mode
- `coverage` - Show test coverage summary
- `bench` - Run benchmark tests
- `security` - Run security checks
- `help` - Show help message

### Example Usage

```bash
# Complete development workflow
make deps fmt lint test build

# Quick development iteration
make dev-build && make dev-run

# Check code quality
make lint security

# Run tests with coverage
make test coverage
```

## Development Tools

### Installed Tools

- **goimports** - Automatic import formatting
- **golangci-lint** - Comprehensive linting
- **gosec** - Security scanning
- **delve** (in Docker) - Debugging
- **air** (in Docker) - Hot reloading

### Configuration Files

- `.golangci.yml` - golangci-lint configuration
- `.air.toml` - Air hot reloading configuration (in Docker)
- `Dockerfile.dev` - Development container
- `docker-compose.dev.yml` - Development environment

### IDE Integration

The project is configured to work well with:
- VS Code (with Go extension)
- GoLand
- Vim/Neovim with Go plugins

## Project Structure

```
.
├── cmd/
│   └── oskway/          # Main application
├── internal/            # Internal packages
│   ├── render/          # Rendering logic
│   └── wayland/         # Wayland client
├── pkg/                 # Public packages
│   └── keyboard/        # Keyboard handling
├── ui/                  # UI components
├── build/               # Build artifacts
├── coverage/            # Test coverage reports
├── Makefile            # Build automation
├── Dockerfile.dev      # Development container
└── docker-compose.dev.yml # Development environment
```

## Testing

### Running Tests

```bash
# Run all tests with race detection
make test

# Run tests with coverage
make coverage

# Run benchmarks
make bench

# Run specific test
go test -race ./pkg/keyboard/
```

### Test Coverage

Tests generate coverage reports in the `coverage/` directory:
- `coverage.out` - Coverage data
- `coverage.html` - HTML coverage report

## Code Quality

### Linting

The project uses golangci-lint with a comprehensive configuration:

```bash
make lint
```

### Security

Security scanning is performed with gosec:

```bash
make security
```

### Code Formatting

Code is automatically formatted with goimports and gofmt:

```bash
make fmt
```

## Docker Development

### Development Container

The `Dockerfile.dev` provides a complete development environment with:
- Go 1.22
- All development tools
- Wayland development libraries
- Debugging tools (delve)
- Hot reloading (air)

### Multi-stage Build

The Docker setup includes multiple stages:
- `base` - Basic Go environment
- `development` - Full development tools
- `testing` - Test execution
- `linting` - Code quality checks
- `builder` - Application build
- `production` - Minimal runtime

### Usage Examples

```bash
# Start development environment
docker-compose -f docker-compose.dev.yml up -d dev

# Run tests in container
docker-compose -f docker-compose.dev.yml --profile test up test

# Run linter in container
docker-compose -f docker-compose.dev.yml --profile lint up lint

# Build in container
docker-compose -f docker-compose.dev.yml --profile build up builder
```

## Hot Reloading

When using the development container, you can use air for hot reloading:

```bash
# Inside the development container
air
```

This will automatically rebuild and restart the application when source files change.

## Debugging

### With Delve

In the development container:

```bash
# Build with debug symbols
make dev-build

# Run with delve
dlv exec build/oskway-dev
```

### VS Code Integration

The development container exposes port 40000 for remote debugging with VS Code.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run the full test suite: `make deps fmt lint test build`
5. Submit a pull request

## Troubleshooting

### Common Issues

1. **Missing tools**: Run `make install-tools`
2. **Import errors**: Run `make fmt` to fix imports
3. **Linter errors**: Run `make lint` to see issues
4. **Build failures**: Check dependencies with `make deps`

### Container Issues

1. **Volume mounts**: Ensure Docker has permission to mount the project directory
2. **Port conflicts**: Change ports in docker-compose.dev.yml if needed
3. **Build cache**: Use `docker-compose build --no-cache` to rebuild from scratch

### Performance Tips

- Use `make dev-build` for faster development builds
- Use Docker volumes for Go module cache
- Run tests in parallel with `go test -parallel N`
