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
import os
import re
from concurrent.futures import ThreadPoolExecutor, TimeoutError as FutureTimeout

matplotlib.use('Agg')

global_env = {}
task_queue = queue.Queue()
output_queue = queue.Queue()

EXEC_TIMEOUT = 60


def load_config():
    config_path = os.path.expanduser("~/.config/tool/config.json")
    if not os.path.exists(config_path):
        raise FileNotFoundError(f"Config file not found at {config_path}")
    with open(config_path, "r") as f:
        return json.load(f)


def apply_config_paths(code, dataset_root):
    # Rewrite common data-loading paths to use dataset_root
    patterns = [
        r'(pd\.read_csv\(\s*[\'"])([^\'"]+)([\'"])',
        r'(open\(\s*[\'"])([^\'"]+)([\'"])',
        r'(np\.load\(\s*[\'"])([^\'"]+)([\'"])',
        r'(json\.load\(\s*open\(\s*[\'"])([^\'"]+)([\'"])'
    ]

    for pattern in patterns:
        code = re.sub(pattern, lambda m: f"{m[1]}{os.path.join(dataset_root, m[2].lstrip('/'))}{m[3]}", code)

    return code


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
        config = load_config()
        dataset_root = config.get("dataset_root", "")
        output_dir = config.get("output_dir", "")

        # Apply dataset path rewriting
        code = apply_config_paths(code, dataset_root)

        os.makedirs(output_dir, exist_ok=True)

        with contextlib.redirect_stdout(stdout_capture), contextlib.redirect_stderr(stderr_capture):
            exec(code, global_env)

            figs = [plt.figure(n) for n in plt.get_fignums()]
            for fig in figs:
                filename = f"plot_{uuid.uuid4().hex}.png"
                full_path = os.path.join(output_dir, filename)
                fig.savefig(full_path)
                result["images"].append(full_path)
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
