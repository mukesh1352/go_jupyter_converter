import sys
import json
import io
import contextlib
import matplotlib
import matplotlib.pyplot as plt
import uuid
import threading
import queue
import traceback
import gc
from concurrent.futures import ThreadPoolExecutor, TimeoutError as FutureTimeout

matplotlib.use('Agg')  

global_env = {}
task_queue = queue.Queue()
output_queue = queue.Queue()

EXEC_TIMEOUT = 60  

def actual_execution(code):
    stdout_capture = io.StringIO()
    stderr_capture = io.StringIO()
    result = {
        "stdout": "",
        "stderr": "",
        "success": True,
        "images": []
    }

    try:
        with contextlib.redirect_stdout(stdout_capture), contextlib.redirect_stderr(stderr_capture):
            exec(code, global_env)

            figs = [plt.figure(n) for n in plt.get_fignums()]
            for fig in figs:
                filename = f"plot_{uuid.uuid4().hex}.png"
                fig.savefig(filename)
                result["images"].append(filename)
                plt.close(fig)

    except Exception:
        result["success"] = False
        result["stderr"] = traceback.format_exc()

    result["stdout"] = stdout_capture.getvalue()
    result["stderr"] += stderr_capture.getvalue()
    return result


def execute_cell(code):
    result = {
        "stdout": "",
        "stderr": "",
        "success": False,
        "images": []
    }

    try:
        with ThreadPoolExecutor(max_workers=1) as executor:
            future = executor.submit(actual_execution, code)
            result = future.result(timeout=EXEC_TIMEOUT)
    except FutureTimeout:
        result["success"] = False
        result["stderr"] = f"TimeoutError: Code execution exceeded {EXEC_TIMEOUT} seconds."
    except Exception:
        result["success"] = False
        result["stderr"] = traceback.format_exc()
    finally:
        gc.collect()

    return result


def worker():
    while True:
        item = task_queue.get()
        if item is None:
            break
        code = item.get("code", "")
        result = execute_cell(code)
        output_queue.put(result)


def main():
    worker_thread = threading.Thread(target=worker, daemon=True)
    worker_thread.start()

    while True:
        line = sys.stdin.readline()
        if not line:
            break

        try:
            data = json.loads(line)

            if data.get("code") == "__EXIT__":
                task_queue.put(None)
                break

            task_queue.put(data)
            result = output_queue.get()
            print(json.dumps(result), flush=True)

        except Exception:
            error_output = {
                "stdout": "",
                "stderr": f"Internal error: {traceback.format_exc()}",
                "success": False,
                "images": []
            }
            print(json.dumps(error_output), flush=True)


if __name__ == "__main__":
    main()
