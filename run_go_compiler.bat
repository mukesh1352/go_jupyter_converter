@echo off
cd /d "C:\Users\Albert Nedumudy\Desktop\Main HQ\Work\Github_Tool_Core\docker-worker\compiler\go_jupyter_converter"

REM Run Go and capture its JSON output
go run cmd\go_compiler\main.go --file "%~1"
