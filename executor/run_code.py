# executor/run_code.py

import sys
import json
import io
import contextlib

# Shared global context to persist variables across cells
global_context = {}

def execute_cell(code):
    output = {"stdout": "", "stderr": "", "success": True}

    stdout_capture = io.StringIO()
    stderr_capture = io.StringIO()

    try:
        with contextlib.redirect_stdout(stdout_capture), contextlib.redirect_stderr(stderr_capture):
            exec(code, global_context)
    except Exception as e:
        output["success"] = False
        output["stderr"] = str(e)

    output["stdout"] = stdout_capture.getvalue()
    output["stderr"] += stderr_capture.getvalue()

    return output

def main():
    # Receive code cell as JSON via stdin
    input_data = sys.stdin.read()
    code = json.loads(input_data).get("code", "")
    
    result = execute_cell(code)
    
    print(json.dumps(result))

if __name__ == "__main__":
    main()
