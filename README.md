# 🧮 Student Scores Analyzer

This project processes student scores for Math and Science, calculates their average scores, and identifies the top-performing student.

---

## 📂 Files

- `1.py` – Python script that handles data processing and analysis.
- `main.go` – Go-based compiler that detects file types and executes Python code.

---

## 🐍 Python Execution

Run using the standard Python 3 interpreter:

```bash
time python3 1.py


## Output
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

# Runtime
python3 1.py  0.20s user 0.08s system 37% cpu 0.745 total

# Go-Based Compiler Execution
```bash
go run main.go
```

Output
File extension is: .py
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


Execution time
Execution Time: 452.751291ms


## 📊 Performance Comparison

| Method             | Total Execution Time |
|--------------------|----------------------|
| Python Interpreter | 745 ms               |
| Go-Based Compiler  | ~453 ms              |

The **Go-based compiler** approach shows improved execution time when handling the same Python code.
