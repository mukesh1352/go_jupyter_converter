# Go Python Executor

A sophisticated Go-based execution engine for Python files and Jupyter notebooks with advanced features including performance monitoring, dataset path rewriting, visualization support, and concurrent execution capabilities.

## Overview

Go Python Executor is a powerful command-line utility that provides a robust execution environment for Python scripts (`.py`) and Jupyter notebooks (`.ipynb`). Built with Go for performance and reliability, it offers advanced features like automatic dataset path normalization, matplotlib visualization handling, execution timeout management, and comprehensive error reporting.

## Key Features

### 🚀 **Multi-Format Support**
- **Python Scripts (`.py`)**: Direct execution with real-time performance monitoring
- **Jupyter Notebooks (`.ipynb`)**: Cell-by-cell execution with JSON parsing and validation
- **Batch Processing**: Process entire directories or individual files

### 🔧 **Advanced Execution Engine**
- **Concurrent Execution**: Thread-safe execution with timeout management (60s default)
- **Persistent Python Environment**: Maintains global state across notebook cells
- **Memory Management**: Automatic garbage collection and resource cleanup
- **Error Isolation**: Robust error handling with detailed stack traces

### 📊 **Visualization & Output Management**
- **Matplotlib Integration**: Automatic plot saving with unique filenames
- **Configurable Output Directory**: Centralized output management
- **Image Export**: PNG format with UUID-based naming
- **Stdout/Stderr Capture**: Complete output redirection and logging

### 🛠 **Smart Configuration System**
- **Auto-Configuration**: Creates default config if none exists
- **Dataset Path Rewriting**: Automatic path normalization for data files
- **Flexible Directory Structure**: User-configurable dataset and output paths
- **JSON-based Configuration**: Easy to modify and version control

### 🎨 **Enhanced User Experience**
- **Colored Terminal Output**: Enhanced readability with color-coded messages
- **Performance Metrics**: Execution time tracking for Python scripts
- **Detailed Logging**: Comprehensive execution logs and error reporting
- **Progress Indicators**: Real-time processing status updates

## Installation

### Prerequisites

- **Go**: Version 1.24 or higher
- **Python**: Version 3.x with the following packages:
  - `matplotlib`
  - `numpy` (recommended)
  - `pandas` (recommended)
  - `json` (built-in)
- **System**: Unix-like system (Linux, macOS) or Windows with proper PATH configuration

### Quick Install

```bash
# Clone the repository
git clone <repository-url>
cd go_compiler

# Download dependencies
go mod download

# Build the executable
go build -o go_executor cmd/go_compiler/main.go

# Make executable (Unix-like systems)
chmod +x go_executor
```

### Alternative Build Methods

```bash
# Build with optimizations
go build -ldflags="-s -w" -o go_executor cmd/go_compiler/main.go

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o go_executor-linux cmd/go_compiler/main.go
GOOS=windows GOARCH=amd64 go build -o go_executor.exe cmd/go_compiler/main.go
GOOS=darwin GOARCH=amd64 go build -o go_executor-macos cmd/go_compiler/main.go
```

## Usage

### Command Line Interface

```bash
# Execute a single Python file
./go_executor --file path/to/script.py

# Process all Python/Jupyter files in a directory
./go_executor --dir path/to/directory

# Real-world examples
./go_executor --file data_analysis.py
./go_executor --file machine_learning_notebook.ipynb
./go_executor --dir ./research_notebooks
./go_executor --dir ~/data_science_projects
```

### Configuration

The application automatically creates a configuration file at `~/.config/tool/config.json`:

```json
{
  "dataset_root": "/home/user/datasets",
  "output_dir": "/home/user/output"
}
```

#### Configuration Options

- **`dataset_root`**: Base directory for dataset files (used for path rewriting)
- **`output_dir`**: Directory where generated plots and outputs are saved

#### Path Rewriting Feature

The system automatically rewrites common data loading patterns:

```python
# Original code in notebook
pd.read_csv('data.csv')
open('dataset.txt')
np.load('model.npy')

# Automatically rewritten to
pd.read_csv('/home/user/datasets/data.csv')
open('/home/user/datasets/dataset.txt')
np.load('/home/user/datasets/model.npy')
```

## Architecture

### Project Structure

```
go_compiler/
├── cmd/
│   └── go_compiler/
│       ├── main.go              # Application entry point and CLI handling
│       └── go_compiler          # Compiled binary
├── internal/
│   ├── api/
│   │   └── python_checker.go    # Python script execution engine
│   ├── executor/
│   │   └── run_code.py          # Python subprocess handler
│   ├── parser1/
│   │   └── parser_checker.go    # Jupyter notebook parser and processor
│   └── utils/
│       ├── config.go            # Configuration management system
│       └── dataset_rewriter.go  # Path rewriting utilities
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
└── README.md                    # This documentation
```

### Core Components

#### 1. **Main Controller** (`cmd/go_compiler/main.go`)
- Command-line argument parsing
- File type detection and routing
- Directory traversal and batch processing
- Error handling and logging coordination

#### 2. **Python Script Engine** (`internal/api/python_checker.go`)
- Direct Python script execution via `python3` subprocess
- Real-time stdout/stderr forwarding
- Execution time measurement and reporting
- Process lifecycle management

#### 3. **Jupyter Notebook Processor** (`internal/parser1/parser_checker.go`)
- JSON parsing of `.ipynb` files
- Cell-by-cell code extraction and execution
- Persistent Python subprocess communication
- Result aggregation and formatting

#### 4. **Python Execution Runtime** (`internal/executor/run_code.py`)
- Threaded execution environment with timeout protection
- Global variable persistence across cells
- Matplotlib figure capture and export
- Comprehensive error handling and reporting

#### 5. **Configuration System** (`internal/utils/config.go`)
- Automatic configuration file creation
- JSON-based settings management
- Path resolution and validation
- Default value provisioning

## Advanced Features

### Execution Timeout Management

```python
# Automatic timeout protection (60 seconds default)
# Long-running cells are automatically terminated
# with detailed timeout error messages
```

### Matplotlib Integration

```python
import matplotlib.pyplot as plt
import numpy as np

# Plots are automatically saved
x = np.linspace(0, 10, 100)
y = np.sin(x)
plt.plot(x, y)
plt.title('Sine Wave')
plt.show()  # Automatically saved to output directory
```

### Error Handling Examples

```bash
# Syntax errors are captured and reported
./go_executor --file broken_script.py
# Output: Detailed Python traceback with line numbers

# File not found errors
./go_executor --file nonexistent.py
# Output: Clear error message with file path

# Directory processing with mixed file types
./go_executor --dir mixed_directory
# Output: Processes supported files, skips others with warnings
```

## Dependencies

### Go Dependencies
```go
require (
    github.com/fatih/color v1.18.0          // Terminal color output
    github.com/go-redis/redis/v8 v8.11.5    // Redis client (optional)
    github.com/redis/go-redis/v9 v9.11.0    // Updated Redis client
    // Additional system dependencies managed automatically
)
```

### Python Dependencies
```python
# Core requirements
matplotlib>=3.5.0    # Visualization and plot generation
numpy>=1.21.0       # Numerical computing (recommended)
pandas>=1.3.0       # Data manipulation (recommended)

# Built-in modules used
json                # JSON parsing and serialization
sys                 # System-specific parameters
io                  # Input/output operations
contextlib          # Context management utilities
threading           # Thread-based parallelism
concurrent.futures  # High-level threading interface
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Test specific package
go test ./internal/utils
```

### Development Workflow

```bash
# Format code
go fmt ./...

# Lint code (requires golint)
golint ./...

# Vet code for potential issues
go vet ./...

# Build and test
go build && go test ./...
```

### Adding New Features

1. **New File Type Support**: Extend the `CheckExtension` function in `main.go`
2. **Configuration Options**: Add fields to the `Config` struct in `utils/config.go`
3. **Python Execution Features**: Modify `internal/executor/run_code.py`
4. **Output Formatting**: Enhance the respective checker functions

## Troubleshooting

### Common Issues

#### Python Not Found
```bash
# Ensure Python 3 is installed and accessible
python3 --version
which python3
```

#### Permission Errors
```bash
# Make sure the executable has proper permissions
chmod +x go_executor

# Check directory permissions for output
ls -la ~/.config/tool/
```

#### Configuration Issues
```bash
# Reset configuration by deleting the config file
rm ~/.config/tool/config.json
# Will be recreated with defaults on next run
```

#### Memory Issues with Large Notebooks
- The system includes automatic garbage collection
- Large datasets may require increasing system memory
- Consider processing notebooks in smaller batches

## Performance Considerations

- **Concurrent Execution**: Utilizes Go's goroutines for efficient processing
- **Memory Management**: Automatic cleanup prevents memory leaks
- **Timeout Protection**: Prevents infinite loops and runaway processes
- **Efficient JSON Parsing**: Optimized notebook processing
- **Minimal Dependencies**: Lightweight runtime footprint

## Contributing

We welcome contributions! Please follow these guidelines:

### Getting Started
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes with proper testing
4. Ensure code formatting: `go fmt ./...`
5. Run tests: `go test ./...`
6. Commit changes: `git commit -m 'Add amazing feature'`
7. Push to branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

### Code Standards
- Follow Go best practices and idioms
- Include comprehensive error handling
- Add tests for new functionality
- Update documentation for API changes
- Use meaningful commit messages

### Reporting Issues
- Use the GitHub issue tracker
- Include system information (OS, Go version, Python version)
- Provide minimal reproduction steps
- Include relevant error messages and logs

## Roadmap

- [ ] **Web Interface**: Browser-based notebook execution
- [ ] **Plugin System**: Extensible architecture for custom processors
- [ ] **Distributed Execution**: Multi-node processing capabilities
- [ ] **Enhanced Visualization**: Support for additional plotting libraries
- [ ] **Database Integration**: Direct database connectivity for notebooks
- [ ] **Cloud Storage**: S3/GCS integration for datasets and outputs
- [ ] **Performance Profiling**: Built-in code performance analysis
- [ ] **Interactive Mode**: REPL-style execution environment

## License

MIT License

Copyright (c) 2024 Go Python Executor

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.