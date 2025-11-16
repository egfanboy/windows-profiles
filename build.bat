@echo off
echo Building Monitor Profile Manager (Wails + React)...

REM Add Node.js to PATH for this session
set PATH=C:\Program Files\nodejs;%PATH%

REM Check if Wails is installed
where wails >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Wails CLI not found. Installing...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    if %ERRORLEVEL% NEQ 0 (
        echo Failed to install Wails CLI
        pause
        exit /b 1
    )
)

REM Verify Node.js and npm are available
echo Verifying Node.js installation...
node --version
if %ERRORLEVEL% NEQ 0 (
    echo Error: Node.js not found in PATH
    echo Please ensure Node.js is installed in C:\Program Files\nodejs
    pause
    exit /b 1
)

REM Install frontend dependencies
echo Installing frontend dependencies...
cd frontend
npm install
if %ERRORLEVEL% NEQ 0 (
    echo Failed to install frontend dependencies
    pause
    exit /b 1
)

REM Build frontend
echo Building frontend...
npm run build
if %ERRORLEVEL% NEQ 0 (
    echo Frontend build failed
    pause
    exit /b 1
)
cd ..

REM Clean and rebuild Go dependencies
echo Cleaning and updating Go dependencies...
go mod tidy
if %ERRORLEVEL% NEQ 0 (
    echo Failed to update dependencies
    pause
    exit /b 1
)

REM Build the application
echo Building application...
wails build
if %ERRORLEVEL% NEQ 0 (
    echo Build failed
    pause
    exit /b 1
)

echo Build completed successfully!
echo Executable location: build\bin\monitor-profile-manager-wails.exe
echo.
echo Features included:
echo - Monitor detection and management
echo - Audio device detection and management
echo - Ignore list for audio devices
echo - Profile saving and loading
echo - Audio device switching in profiles
pause
