# Xtris Clone

A Tetris clone written in Go using the Ebiten game engine. This project prioritizes native libraries and maintainability.

**Note**: This is a clone of Tetris, originally created by Alexey Pajitnov. The name "Xtris" is used to avoid copyright issues.

## Features

- Classic Tetris gameplay mechanics
- 7 standard Tetris pieces (tetrominoes)
- Line clearing and scoring system
- Game over detection
- Clean, maintainable Go code
- Comprehensive test coverage
- Cross-platform support (Windows, macOS, Linux)

## Prerequisites

- Go 1.21 or later
- For development: golangci-lint (optional, install with `make install-tools`)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd tetris-clone
```

2. Install dependencies:
```bash
go mod download
```

## Building

### Quick Build
Build for your current platform:
```bash
make build
# or
go build -o xtris main.go
```

### Cross-Platform Builds
Build for all supported platforms:
```bash
make build-all
```

This creates executables for:
- Windows (AMD64): `xtris.exe`
- macOS (AMD64): `xtris-darwin-amd64`
- macOS (ARM64): `xtris-darwin-arm64`
- Linux (AMD64): `xtris-linux-amd64`
- Linux (ARM64): `xtris-linux-arm64`

### Platform-Specific Builds
Build for specific platforms:
```bash
make build-windows  # Windows only
make build-macos    # macOS (AMD64 + ARM64)
make build-linux    # Linux (AMD64 + ARM64)
```

### Release Builds
Prepare all builds for release:
```bash
make release
```
This creates a `release/` directory with all platform builds.

## Running

Run the game:
```bash
make run
# or
go run main.go
```

## Development

### Available Make Commands

```bash
make help              # Show all available commands
make test              # Run tests
make test-coverage     # Run tests with coverage
make test-race         # Run tests with race detection
make lint              # Run linter
make fmt               # Format code
make dev               # Full development workflow (fmt, lint, test)
make clean             # Clean build artifacts
make install-tools     # Install development tools
```

### Testing

Run tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
make coverage-report   # Show coverage summary
make coverage-html     # Generate HTML coverage report
```

### Code Quality

Format code:
```bash
make fmt
```

Run linter:
```bash
make lint
```

### Coverage

The project maintains high test coverage. View coverage:
```bash
make coverage-report
```

Generate HTML coverage report:
```bash
make coverage-html
```

Check if coverage meets minimum threshold (80%):
```bash
make coverage-check
```

## Controls

- **Arrow Keys**: Move piece left/right
- **Down Arrow**: Soft drop
- **Up Arrow**: Hard drop
- **Space**: Rotate piece
- **P**: Pause/Resume
- **R**: Restart game
- **Q**: Quit

## Project Structure

```
tetris-clone/
├── main.go              # Main game entry point
├── main_test.go         # Test suite
├── go.mod               # Go module file
├── go.sum               # Dependency checksums
├── Makefile             # Build and development tasks
├── README.md            # This file
├── .github/workflows/   # CI/CD workflows
├── scripts/             # Development scripts
└── .golangci.yml        # Linter configuration
```

## Cross-Platform Support

This project is fully cross-platform thanks to the Ebiten game engine:

- **Windows**: AMD64 architecture
- **macOS**: Both Intel (AMD64) and Apple Silicon (ARM64)
- **Linux**: Both x86_64 (AMD64) and ARM64 architectures

The CI/CD pipeline automatically tests and builds for all supported platforms.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run the development workflow: `make dev`
5. Ensure tests pass and coverage is maintained
6. Submit a pull request

## License

This project is open source. Please respect the original Tetris game's copyright.

## Acknowledgments

- Original Tetris game by Alexey Pajitnov
- Ebiten game engine for cross-platform graphics
- Go community for excellent tooling and libraries 