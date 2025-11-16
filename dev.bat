@echo off
echo Starting Monitor Profile Manager Development...

REM Add Node.js to PATH for this session
set PATH=C:\Program Files\nodejs;%PATH%

REM Verify Node.js and npm are available
echo Verifying Node.js installation...
node --version
if %ERRORLEVEL% NEQ 0 (
    echo Error: Node.js not found in PATH
    echo Please ensure Node.js is installed in C:\Program Files\nodejs
    pause
    exit /b 1
)

npm --version
if %ERRORLEVEL% NEQ 0 (
    echo Error: npm not found in PATH
    pause
    exit /b 1
)

REM Install frontend dependencies if needed
if not exist "frontend\node_modules" (
    echo Installing frontend dependencies...
    cd frontend
    npm install
    if %ERRORLEVEL% NEQ 0 (
        echo Failed to install frontend dependencies
        pause
        exit /b 1
    )
    cd ..
)

REM Start Wails development server
echo Starting Wails development server...
wails dev

pause
