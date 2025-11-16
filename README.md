# Monitor Profile Manager (Wails + React)

A modern Windows desktop application built with Wails v2 and React that allows you to create and manage monitor profiles for different display configurations.

## Features

- **Monitor Detection**: Automatically detects all connected monitors and their properties
- **Profile Management**: Save and load different monitor configurations
- **Display Control**: Activate/deactivate monitors and set primary display
- **Modern UI**: Clean, responsive interface using React with resizable table columns
- **Cross-Platform**: Built with Go and React, compiles on multiple platforms (monitor management only works on Windows)

## Technology Stack

- **Backend**: Go with Wails v2
- **Frontend**: React 18 with TypeScript
- **Build Tool**: Vite
- **Styling**: Modern CSS with gradients and animations
- **UI Components**: Custom resizable table component

## Requirements

- Go 1.23 or later
- Node.js 18+ and npm (for frontend development)
- Windows operating system (for monitor management functionality)
- C compiler (for Wails compilation)

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
1. Configure your monitors as desired (using Windows display settings)
2. Click "Save Current Profile"
3. Enter a name for your profile
4. The profile will be saved to `~/MonitorProfiles/`

### Applying Profiles
1. Select a profile from the dropdown list
2. Click "Apply Selected Profile"
3. The application will configure your monitors according to the saved profile

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
├── app.go               # Go backend with monitor management logic
├── monitor_manager_windows.go   # Windows-specific monitor API calls
├── monitor_manager_other.go     # Non-Windows fallback implementations
├── go.mod              # Go module dependencies
├── wails.json          # Wails configuration
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

## Migration from Fyne

This project was successfully migrated from Fyne to Wails with the following improvements:
- **Modern Web UI**: React-based interface with better styling
- **Resizable Tables**: Custom implementation with drag-to-resize columns
- **Better UX**: Smooth animations, hover effects, and responsive design
- **Maintained Functionality**: All original monitor management features preserved
- **Cross-platform Potential**: Easier to extend to other platforms

## License

This project is open source. Please refer to the LICENSE file for details.
To build a redistributable, production mode package, use `wails build`.
