# OSK IoT Core

An on-screen keyboard written in GoLang with support for Wayland and X11.

## Overview

OSK IoT Core provides an efficient and customizable on-screen keyboard tailored for Linux environments, supporting both Wayland and X11 display servers. Developed in GoLang, it offers smooth integration with modern desktop environments.

## Features

- **Customizable Layouts**: Easily design and use custom keyboard layouts.
- **Wayland and X11 Support**: Seamless operation across the most common Linux display servers.
- **Efficiency**: Designed for minimal resource consumption and high performance.
- **Open Source**: Contributions are welcome to improve and expand functionality.

## Architecture

The keyboard utilizes Go's concurrency features to handle input efficiently, supporting distinct modules for Wayland and X11 environments.

## Getting Started

### Prerequisites

- Go 1.21+ (for development)
- Wayland development libraries (for Wayland support)
- X11 development libraries (for X11 support)
- GTK+ 3.0+ development libraries

### Installation

```bash
# Clone the repository
git clone https://github.com/unklstewy/osk-iotcore.git
cd osk-iotcore

# Install dependencies
go mod download

# Build the project
go build -o osk-iotcore

# Run the on-screen keyboard
./osk-iotcore
```

## Documentation

Detailed documentation is available in the `docs/` directory:

- [API Documentation](docs/api.md)
- [Device Integration Guide](docs/device-integration.md)
- [Deployment Guide](docs/deployment.md)
- [Configuration Reference](docs/configuration.md)

## Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/unklstewy/osk-iotcore/issues)
- **Discussions**: [GitHub Discussions](https://github.com/unklstewy/osk-iotcore/discussions)

## Roadmap

- [ ] Complete Wayland protocol implementation
- [ ] Full X11 support with input method integration
- [ ] Theme and layout customization
- [ ] Multi-language support
- [ ] Configuration GUI
- [ ] Performance optimizations

## Acknowledgments

- Thanks to the Wayland and X11 communities for their excellent documentation
- Special thanks to the Go community for providing robust libraries for GUI development

---

**OSK IoT Core** - On-Screen Keyboard for Linux
