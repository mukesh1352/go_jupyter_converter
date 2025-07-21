# Go Interpreter

A Go-based interpreter for executing Python files and Jupyter notebooks with built-in validation and performance monitoring capabilities.

## Overview

Go Interpreter is a command-line utility that interprets and executes Python files (`.py`) and Jupyter notebooks (`.ipynb`), providing real-time execution with performance monitoring and validation features. The interpreter supports both single file execution and batch directory processing.

## Features

- **Multi-format Support**: Process both Python files (`.py`) and Jupyter notebooks (`.ipynb`)
- **Flexible Processing**: Handle single files or entire directories
- **Performance Monitoring**: Built-in execution time tracking
- **Configuration Management**: Configurable settings via config files
- **Redis Integration**: Built-in Redis support for caching and data management
- **Colored Output**: Enhanced terminal output with color support

## Installation

### Prerequisites

- Go 1.24 or higher
- Python 3.x
- Redis (optional, for caching features)

### Build from Source

```bash
git clone <repository-url>
cd go_compiler
go mod download
go build -o go_interpreter cmd/go_compiler/main.go
```

## Usage

### Command Line Options

```bash
# Interpret and execute a single file
./go_interpreter --file path/to/script.py

# Interpret all supported files in a directory
./go_interpreter --dir path/to/directory

# Examples
./go_interpreter --file examples/hello.py
./go_interpreter --dir ./python_scripts
```

### Supported File Types

- **`.py`**: Python scripts - interpreted and executed directly with performance monitoring
- **`.ipynb`**: Jupyter notebooks - interpreted and processed with dataset validation

## Project Structure

```
go_compiler/
├── cmd/
│   └── go_compiler/          # Main interpreter entry point
│       └── main.go
├── internal/
│   ├── api/                  # API and interpretation handlers
│   │   └── python_checker.go
│   ├── executor/             # Code interpretation and execution utilities
│   │   └── run_code.py
│   ├── parser1/              # File parsing and validation
│   │   └── parser_checker.go
│   └── utils/                # Utility functions and configuration
│       ├── config.go
│       └── dataset_rewriter.go
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── README.md                 # This file
```

## Configuration

The application uses a configuration system that can be customized via the `internal/utils/config.go` module. Configuration options include:

- Dataset root paths for Jupyter notebook processing
- Redis connection settings
- Output formatting preferences

## Dependencies

- **Redis**: `github.com/go-redis/redis/v8` and `github.com/redis/go-redis/v9`
- **Color Output**: `github.com/fatih/color`
- **System Utilities**: Various Go standard library packages

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build for current platform
go build -o go_interpreter cmd/go_compiler/main.go

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o go_interpreter-linux cmd/go_compiler/main.go
GOOS=windows GOARCH=amd64 go build -o go_interpreter.exe cmd/go_compiler/main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues, questions, or contributions, please open an issue on the project repository.