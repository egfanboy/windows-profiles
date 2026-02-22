# Monitor Profile Manager (Wails + React)

A modern Windows desktop application built with Wails v2 and React that allows you to create and manage monitor and audio device profiles for different display and audio configurations.

## Features

- **Monitor Detection**: Automatically detects all connected monitors and their properties using MultiMonitorTool CLI
- **Audio Device Detection**: Automatically detects all audio devices and their properties using SVCL CLI
- **Profile Management**: Save and load different monitor and audio device configurations
- **Display Control**: Activate/deactivate monitors and set primary display via CLI tools
- **Audio Control**: Set default audio devices and manage audio device states via CLI tools
- **Modern UI**: Clean, responsive interface using React with resizable table columns
- **Cross-Platform**: Built with Go and React, compiles on multiple platforms (monitor and audio management only works on Windows)

## Technology Stack

- **Backend**: Go with Wails v2
- **Frontend**: React 18 with TypeScript
- **Build Tool**: Vite
- **Styling**: Modern CSS with gradients and animations
- **UI Components**: Custom resizable table component
- **CLI Tools**: MultiMonitorTool for monitor management, SVCL for audio device management

## Requirements

- Go 1.23 or later
- Node.js 18+ and npm (for frontend development)
- Windows operating system (for monitor and audio management functionality)
- C compiler (for Wails compilation)
- **MultiMonitorTool**: Included in `tools/multimonitortool/` directory
- **SVCL (SoundVolumeCommandLine)**: Included in `tools/svcl/` directory

## Installation

### Prerequisites
1. Install Go 1.23 or later from [golang.org](https://golang.org/dl/)
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

## Usage

### Detecting Monitors
1. Launch the application
2. Click "Refresh Monitors" to detect all connected displays
3. View monitor information including resolution, primary status, and device name in the resizable table

### Creating Profiles
1. Configure your monitors and audio devices as desired (using Windows display settings and sound settings)
2. Click "Save Current Profile"
3. Enter a name for your profile
4. The profile will be saved to `~/MonitorProfiles/` with both monitor and audio device configurations

### Applying Profiles
1. Select a profile from the dropdown list
2. Click "Apply Selected Profile"
3. The application will configure your monitors and audio devices according to the saved profile using CLI tools

## UI Features

### Resizable Table
- **Drag column borders** to resize table columns
- **Hover effects** on resize handles for better UX
- **Minimum width constraints** to prevent columns from becoming too small
- **Persistent sizing** during the session

### Modern Design
- **Gradient background** with glassmorphism effects
- **Smooth animations** and transitions
- **Responsive layout** that adapts to different screen sizes
- **Status indicators** with color coding (active/inactive monitors)

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

## Architecture Overview

This project uses a **CLI-based architecture** for system integration:

### CLI Tool Integration
- **MultiMonitorTool**: External CLI tool for monitor detection and configuration
  - Located in `tools/multimonitortool/MultiMonitorTool.exe`
  - Provides monitor enumeration, primary display setting, and display state control
  - Integration via CSV output parsing and command-line execution

- **SVCL (SoundVolumeCommandLine)**: External CLI tool for audio device management
  - Located in `tools/svcl/svcl.exe`
  - Provides audio device enumeration and default device setting
  - Integration via CSV output parsing and command-line execution

### Benefits of CLI Approach
- **Reliability**: Stable, well-tested third-party tools for system operations
- **Maintainability**: No complex Windows API calls in Go code
- **Flexibility**: Easy to update CLI tools without code changes
- **Compatibility**: Tools handle Windows version differences automatically

## Migration from Direct API Calls

This project was successfully migrated from direct Windows API calls to CLI-based integration with the following improvements:
- **Stability**: CLI tools handle edge cases and Windows version differences
- **Simplified Codebase**: No complex Windows API interactions in Go
- **Better Error Handling**: CLI tools provide clear error messages
- **Enhanced Functionality**: Added audio device management through SVCL
- **Maintained UI**: All original React UI features preserved
- **Cross-platform Potential**: Easier to extend to other platforms with equivalent CLI tools

## License

This project is open source. Please refer to the LICENSE file for details.
To build a redistributable, production mode package, use `wails build`.
