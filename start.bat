@echo off
echo Starting GoSim - Interactive Go Learning Simulator
echo ==================================================

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go is not installed!
    echo Please install Go from: https://golang.org/dl/
    pause
    exit /b 1
)

echo Go is installed: 
go version
echo.

REM Install dependencies
echo Installing dependencies...
go mod download

REM Build the server
echo Building server...
go build -o gosim.exe cmd/server/main.go

if %ERRORLEVEL% EQU 0 (
    echo Build successful!
    echo.
    echo Starting server on http://localhost:8080
    echo ==================================================
    echo Press Ctrl+C to stop the server
    echo.
    
    REM Run the server
    gosim.exe
) else (
    echo Build failed. Please check for errors above.
    pause
    exit /b 1
)