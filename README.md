# Go Compiler - Python & Jupyter Notebook Executor

A Go-based tool for executing Python scripts and Jupyter notebooks with real-time output capture, error handling, and execution timing.

## Features

- **Python Script Execution**: Run `.py` files with full stdout/stderr capture and execution timing
- **Jupyter Notebook Support**: Parse and execute `.ipynb` files cell by cell
- **Concurrent Execution**: Thread-safe execution with timeout protection
- **Image Generation**: Automatic matplotlib plot saving and management
- **Flexible Input**: Process single files or entire directories
- **Real-time Output**: Live output streaming with JSON-based communication

## Installation

### Prerequisites

- Go 1.24 or higher
- Python 3.x with the following packages:
  - `matplotlib`
  - `pandas` (for sample files)
  - `numpy` (for sample files)
  - `torch` (for sample files)

### Build

```bash
go build -o go_compiler
```

## Usage

### Process a Single File

```bash
# Execute a Python script
./go_compiler --file path/to/script.py

# Execute a Jupyter notebook
./go_compiler --file path/to/notebook.ipynb
```

### Process a Directory

```bash
# Process all .py and .ipynb files in a directory
./go_compiler --dir path/to/directory
```

## Architecture

### Components

1. **Main Controller** (`main.go`)
   - Command-line argument parsing
   - File type detection and routing
   - Directory traversal

2. **Python Executor** (`api/python_checker.go`)
   - Direct Python script execution
   - Execution timing measurement
   - Output streaming

3. **Jupyter Parser** (`parser1/parser_checker.go`)
   - Notebook JSON parsing
   - Cell-by-cell execution
   - Persistent Python subprocess management

4. **Code Runner** (`executor/run_code.py`)
   - Python execution environment
   - Thread-safe code execution
   - Matplotlib integration
   - Timeout handling (60 seconds per cell)

### Execution Flow

#### Python Files (.py)
```
File Input → Python Executor → Direct python3 execution → Timed output
```

#### Jupyter Notebooks (.ipynb)
```
File Input → JSON Parser → Cell Extraction → Python Subprocess → Cell-by-cell execution → Aggregated output
```

## Sample Files

The `files/` directory contains example scripts:

- `1.py`: Pandas DataFrame operations with NumPy calculations
- `2.py`: PyTorch linear regression with matplotlib visualization
- `1.ipynb` & `2.ipynb`: Notebook versions of the above scripts

## Examples

### Running a Python Script

```bash
./go_compiler --file files/1.py
```

Output:
```
Processing: /path/to/files/1.py (type: .py)
This is a Python file.
Student Scores:

      Name  Math Score  Science Score  Average Score
0    Alice          85             88           86.5
1      Bob          90             76           83.0
2  Charlie          78             93           85.5
3    David          92             85           88.5

Top Student:
Name             David
Math Score          92
Science Score       85
Average Score     88.5
Name: 3, dtype: object

Execution Time: 452.751291ms
```

### Running a Jupyter Notebook

```bash
./go_compiler --file files/1.ipynb
```

Output:
```
Processing: /path/to/files/1.ipynb (type: .ipynb)
This is a Jupyter Notebook file.
Cell 0 Output:
Success: true
Stdout: [cell execution results]
Stderr: 
------
```

## Configuration

### Timeout Settings

The execution timeout is set to 60 seconds per cell/script. Modify `EXEC_TIMEOUT` in `executor/run_code.py` to adjust:

```python
EXEC_TIMEOUT = 60  # seconds
```

### Dependencies

Key Go modules (see `go.mod`):
- Redis client support
- Color output formatting
- System utilities

## Error Handling

- **File Not Found**: Graceful error reporting with absolute path resolution
- **Execution Timeout**: Automatic termination after 60 seconds
- **Python Errors**: Full traceback capture and display
- **JSON Parsing**: Robust notebook format validation

## Development

### Project Structure
```
.
├── main.go              # Entry point and CLI handling
├── api/
│   └── python_checker.go # Python script executor
├── parser1/
│   └── parser_checker.go # Jupyter notebook parser
├── executor/
│   └── run_code.py      # Python execution environment
└── files/               # Sample scripts and notebooks
```

### Adding New File Types

1. Add case in `CheckExtension()` function in `main.go`
2. Create corresponding handler in appropriate package
3. Update file filtering in `processDirectory()`

## Performance

The Go-based approach provides efficient execution with timing measurements and proper resource management for both individual scripts and batch processing of directories.

## License

This project is open source. See the repository for license details.
