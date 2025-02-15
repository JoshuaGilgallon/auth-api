@echo off
REM Check if MongoDB service is already running
sc query MongoDB | find "RUNNING"
if %ERRORLEVEL% == 0 (
    echo MongoDB is already running.
) else (
    echo Starting MongoDB...
    sc start MongoDB
)
