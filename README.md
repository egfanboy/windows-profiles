# Windows Profile Manager

A Windows desktop application built with Wails and React that allows you to create and manage monitor and audio device profiles for different display and audio configurations.

## Features

- **Monitor Detection**: Automatically detects all connected monitors and their properties using MultiMonitorTool CLI
- **Audio Device Detection**: Automatically detects all audio devices and their properties using SVCL CLI
- **Profile Management**: Save and load different monitor and audio device configurations
- **Display Control**: Activate/deactivate monitors and set primary display via CLI tools
- **Audio Control**: Set default audio devices and manage audio device states via CLI tools

## Requirements

- Go 1.25 or later
- Node.js 18+ and npm (for frontend development)
- Windows operating system (for monitor and audio management functionality)
- C compiler (for Wails compilation)
- **MultiMonitorTool**: Included in `tools/multimonitortool/` directory
- **SVCL (SoundVolumeCommandLine)**: Included in `tools/svcl/` directory

## CLI Tools Setup

This application requires two external CLI tools from NirSoft for monitor and audio management:

### Downloading the Tools

1. **SVCL (SoundVolumeCommandLine)**
   - Download from: https://www.nirsoft.net/utils/sound_volume_command_line.html
   - Download the ZIP file and extract `svcl.exe`
   - Place `svcl.exe` in: `tools/svcl/svcl.exe`

2. **MultiMonitorTool**
   - Download from: https://www.nirsoft.net/utils/multi_monitor_tool.html
   - Download the ZIP file and extract `MultiMonitorTool.exe`
   - Place `MultiMonitorTool.exe` in: `tools/multimonitortool/MultiMonitorTool.exe`

### Directory Structure

Create the following directory structure in your project:

```
tools/
├── svcl/
│   └── svcl.exe
└── multimonitortool/
    └── MultiMonitorTool.exe
```

### Important Notes

- Both tools are portable executables - no installation required
- The application expects these exact file paths for CLI integration
- Run the tools once manually to ensure they work on your system
- On some systems, you may need to run as administrator for full functionality

## Installation

### Prerequisites
1. Install Go 1.25 or later from [golang.org](https://golang.org/dl/)
2. Install Node.js 18+ from [nodejs.org](https://nodejs.org/)
3. Install Wails CLI:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
4. Install a C compiler (required for Wails):
   - **TDM-GCC**: Download from [jmeubank.github.io/tdm-gcc](https://jmeubank.github.io/tdm-gcc/)
   - OR **MinGW-w64**: Download from [sourceforge.net/projects/mingw-w64](https://sourceforge.net/projects/mingw-w64/)
5. Ensure Go, Node.js, and GCC are in your PATH

### Development Setup
1. Navigate to the project directory
2. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   cd ..
   ```
3. Run in development mode:
   ```bash
   wails dev
   ```

### Building Executable
To build a standalone executable:
```bash
wails build
```

For production build without frontend dependencies (uses pre-built frontend):
```bash
wails build -s
```

## Project Structure

```
├── main.go              # Wails application entry point
├── app.go               # Go backend with monitor and audio management logic
├── monitor_manager_windows.go   # Windows-specific monitor API stub (CLI tools used instead)
├── go.mod              # Go module dependencies
├── wails.json          # Wails configuration
├── pkg/                # CLI tool integration packages
│   ├── monitors/       # MultiMonitorTool integration
│   │   └── monitor.go  # Monitor detection and control via CLI
│   └── audio/          # SVCL integration
│       └── audio.go    # Audio device detection and control via CLI
├── tools/              # External CLI tools
│   ├── multimonitortool/  # MultiMonitorTool.exe
│   └── svcl/              # SVCL.exe
├── frontend/           # React frontend
│   ├── src/
│   │   ├── App.tsx     # Main React component
│   │   ├── ResizableTable.tsx  # Custom resizable table component
│   │   ├── App.css     # Main application styles
│   │   └── ResizableTable.css  # Table-specific styles
│   ├── dist/           # Built frontend assets
│   ├── package.json    # Frontend dependencies
│   └── vite.config.ts  # Vite configuration
└── build/              # Built executables
```