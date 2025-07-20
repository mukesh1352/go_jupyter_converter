import sys
import json
import io
import contextlib

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
    while True:
        line = sys.stdin.readline()
        if not line:
            break

        try:
            data = json.loads(line)
            if data.get("code") == "__EXIT__":
                break

            result = execute_cell(data.get("code", ""))
            print(json.dumps(result), flush=True)
        except Exception as e:
            error_output = {
                "stdout": "",
                "stderr": f"Internal error: {str(e)}",
                "success": False
            }
            print(json.dumps(error_output), flush=True)

if __name__ == "__main__":
    main()
