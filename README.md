# Go Compiler

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)
![Python](https://img.shields.io/badge/python-3.x-green.svg)

A powerful tool for executing and analyzing Python scripts and Jupyter notebooks with detailed execution tracking, visualization support, and dataset path normalization.

## Features

- **Python Script Execution**: Run and analyze Python scripts with execution time tracking
- **Jupyter Notebook Processing**: Execute Jupyter notebook cells with detailed output capture
- **Visualization Support**: Automatically save matplotlib figures generated during execution
- **Dataset Path Normalization**: Automatically rewrite dataset paths for consistent access
- **Execution Logging**: Comprehensive logging of execution results with success/failure tracking
- **Configurable Environment**: User-friendly configuration system with sensible defaults

## Installation

### Prerequisites

- Go 1.24 or later
- Python 3.x with matplotlib and other required packages

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/go_compiler.git
cd go_compiler

# Build the project
go build -o go_compiler ./cmd/go_compiler
```

## Configuration

The tool uses a configuration file located at `~/.config/tool/config.json`. On first run, a default configuration will be created with these settings:

```json
{
  "dataset_root": "~/datasets",
  "output_dir": "~/output"
}
```

- `dataset_root`: Base directory for dataset files
- `output_dir`: Directory where execution outputs and generated images are saved

## Usage

### Processing a Single File

```bash
./go_compiler --file path/to/your/script.py
```

```bash
./go_compiler --file path/to/your/notebook.ipynb
```

### Processing All Files in a Directory

```bash
./go_compiler --dir path/to/directory
```

## How It Works

1. **Python Scripts**: The tool executes Python scripts directly and captures execution time and output.

2. **Jupyter Notebooks**: For notebooks, the tool:
   - Parses the notebook JSON structure
   - Extracts and processes code cells
   - Normalizes dataset paths
   - Executes each cell in a controlled Python environment
   - Captures stdout, stderr, and generated visualizations
   - Produces detailed execution reports

3. **Dataset Path Normalization**: Automatically rewrites paths in common data loading functions to use the configured dataset root directory.

## Project Structure

```
.
├── cmd/
│   └── go_compiler/
│       └── main.go           # Main entry point
├── internal/
│   ├── api/
│   │   └── python_checker.go # Python script execution
│   ├── executor/
│   │   └── run_code.py       # Python execution environment
│   ├── parser1/
│   │   └── parser_checker.go # Jupyter notebook processing
│   └── utils/
│       ├── config.go         # Configuration management
│       └── dataset_rewriter.go # Path normalization
├── go.mod                    # Go module definition
├── go.sum                    # Go dependencies checksum
└── README.md                 # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/feature`)
3. Commit your changes (`git commit -m 'Add some  feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- The Go team for the amazing language
- The Python community for the rich ecosystem of data science tools