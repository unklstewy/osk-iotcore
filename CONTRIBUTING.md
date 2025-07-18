# Contributing to OSK IoT Core

We welcome contributions from the community to enhance the functionality and usability of OSK IoT Core. Please follow these guidelines when contributing to our project.

## Table of Contents

- [Code Style](#code-style)
- [Development Environment](#development-environment)
- [Pull Request Workflow](#pull-request-workflow)
- [Development Practices](#development-practices)
- [Testing](#testing)
- [Commit Guidelines](#commit-guidelines)
- [Reporting Issues](#reporting-issues)
- [Review Process](#review-process)

## Code Style

### Go Formatting

- Follow the official Go formatting guidelines and use `gofmt` automatically applied by:
  ```bash
  make fmt
  ```
- Organize imports using `goimports` with local prefix handling for this project:
  ```bash
  # This is automatically handled by make fmt
  goimports -w -local github.com/iotcore/osk-iotcore .
  ```

### Linting

- Run the linter to ensure code quality:
  ```bash
  make lint
  ```
- Fix any linting errors before submitting your PR
- The project uses `golangci-lint` with comprehensive rules defined in `.golangci.yml`

### Code Quality Standards

- Keep functions focused and under 15 lines of cyclomatic complexity
- Use meaningful variable and function names
- Add comments for exported functions and types
- Follow Go naming conventions (camelCase for private, PascalCase for public)
- Avoid naked returns in functions longer than a few lines

## Development Environment

### Prerequisites

- Go 1.22 or later
- GCC compiler with CGO support
- pkg-config
- Required system dependencies (see [README.md](README.md) for installation instructions)

### Quick Setup

```bash
# Check system dependencies
./scripts/check_deps.sh

# Install Go dependencies and development tools
make deps install-tools

# Verify everything works
make test build
```

### Docker Development

For a consistent development environment:

```bash
# Start development container
docker-compose -f docker-compose.dev.yml up -d dev

# Enter the container
docker-compose -f docker-compose.dev.yml exec dev bash
```

## Pull Request Workflow

1. **Fork the Repository**: Start by forking the repository to your GitHub account.

2. **Create a Feature Branch**: Use a descriptive name following the convention:
   ```bash
   git checkout -b feature/keyboard-shortcuts
   git checkout -b fix/wayland-memory-leak
   git checkout -b docs/architecture-updates
   ```

3. **Make Your Changes**: 
   - Follow the code style guidelines
   - Add tests for new functionality
   - Update documentation as needed

4. **Test Your Changes**:
   ```bash
   make deps fmt lint test build
   ```

5. **Commit Your Changes**: Follow the commit message guidelines below.

6. **Push to Your Fork**:
   ```bash
   git push origin feature/keyboard-shortcuts
   ```

7. **Open a Pull Request**: 
   - Use a clear title and description
   - Reference any related issues
   - Include screenshots for UI changes
   - Add appropriate labels

## Development Practices

### Code Organization

- Keep changes focused and atomic
- Separate refactoring from feature additions
- Update package documentation when adding new public APIs
- Follow the existing project structure:
  - `cmd/` - Main applications
  - `internal/` - Internal packages (not for external use)
  - `pkg/` - Public packages (reusable by other projects)
  - `ui/` - User interface components

### CGO Considerations

- When working with CGO code:
  - Test on multiple platforms if possible
  - Add proper error handling for C function calls
  - Document any C dependencies
  - Use build tags appropriately (`//go:build cgo`)

### Wayland Protocol Updates

- When updating Wayland protocol handling:
  - Test with different Wayland compositors
  - Update protocol stubs using the documented process
  - Ensure backward compatibility where possible

## Testing

### Running Tests

```bash
# Run all tests with race detection
make test

# Run tests with coverage report
make coverage

# Run benchmarks
make bench

# Run specific package tests
go test -race ./pkg/keyboard/
```

### Test Types

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test component interactions
- **Mock Tests**: Use the mock Wayland client for testing UI components

### Test Guidelines

- Write tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests for multiple scenarios
- Mock external dependencies appropriately
- Test error conditions and edge cases

## Commit Guidelines

### Commit Message Format

```
type(scope): subject

body

footer
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Test additions or modifications
- `chore`: Maintenance tasks

### Examples

```
feat(keyboard): add DVORAK layout support

Implement DVORAK keyboard layout with proper key mapping
and theme integration.

Closes #123
```

```
fix(wayland): resolve memory leak in surface cleanup

Ensure proper cleanup of Wayland surfaces and buffers
when the application exits.

Fixes #456
```

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

- **Environment**: OS, Go version, Wayland compositor
- **Steps to Reproduce**: Clear steps to reproduce the issue
- **Expected Behavior**: What you expected to happen
- **Actual Behavior**: What actually happened
- **Logs**: Relevant log output or error messages
- **System Info**: Output from `./scripts/check_deps.sh`

### Feature Requests

For new features:

- **Use Case**: Describe why this feature is needed
- **Proposed Solution**: How you envision the feature working
- **Alternatives**: Any alternative approaches considered
- **Additional Context**: Screenshots, mockups, or examples

## Review Process

### What We Look For

- **Code Quality**: Follows project standards and best practices
- **Testing**: Adequate test coverage for new code
- **Documentation**: Updated docs for new features
- **Compatibility**: Works with supported platforms and versions
- **Performance**: No significant performance regressions

### Review Timeline

- Initial review within 2-3 business days
- Follow-up reviews within 1-2 business days
- Merge after approval from at least one maintainer

### Getting Help

- **Discussions**: Use GitHub Discussions for questions
- **Issues**: Create issues for bugs and feature requests
- **Documentation**: Check the `docs/` directory for guides

## Code of Conduct

This project adheres to a code of conduct that we expect all contributors to follow. Please be respectful and professional in all interactions.

## License

By contributing to OSK IoT Core, you agree that your contributions will be licensed under the MIT License.
