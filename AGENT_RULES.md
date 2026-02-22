# Agent Rules for Monitor Profile Manager Development

## Architecture Guidelines

### CLI-Based Architecture
- **Always** use external CLI tools for system operations instead of direct Windows API calls
- **MultiMonitorTool**: Use for all monitor detection and configuration operations
- **SVCL**: Use for all audio device detection and configuration operations
- Never mix direct Windows API calls with CLI-based operations
- CLI tools are located in `tools/` subdirectory and must be checked for existence before use

### CLI Tool Integration
- Use `os/exec.Command` to execute external CLI tools
- Always check for tool existence using `os.Stat` before execution
- Parse CSV output from CLI tools using `encoding/csv` package
- Handle CLI tool errors with descriptive context using `fmt.Errorf()`
- Implement proper timeout handling for CLI tool execution
- Clean up temporary files created by CLI tools (e.g., CSV exports)

### Error Handling
- **Never** ignore error returns
- Wrap Windows API errors with descriptive context using `fmt.Errorf()`
- Provide user-friendly error messages in the GUI while logging technical details
- Implement proper error propagation from low-level API calls to high-level UI

### GUI Development with Wails
- Use React with TypeScript for frontend development
- Implement proper data refresh patterns when underlying state changes
- Use Wails context for communication between Go backend and React frontend
- Follow React component lifecycle patterns for state updates
- Use resizable table components for displaying monitor and audio device information

## Code Structure Rules

### File Organization
```
main.go              - Wails application entry point
app.go               - Go backend with monitor and audio management logic
monitor_manager_windows.go   - Windows-specific stub (CLI tools used instead)
pkg/monitors/monitor.go     - MultiMonitorTool CLI integration
pkg/audio/audio.go           - SVCL CLI integration
go.mod              - Module dependencies
wails.json          - Wails configuration
AGENT_RULES.md      - This file
README.md           - Project documentation
tools/              - External CLI tools directory
├── multimonitortool/   - MultiMonitorTool.exe
└── svcl/               - SVCL.exe
frontend/           - React frontend
```

### Data Structures
- Define structs with JSON tags for serialization
- Use consistent naming conventions (PascalCase for exported, camelCase for unexported)
- Include all necessary fields for monitor state representation
- Keep data structures flat and avoid unnecessary nesting

### Profile Management
- Store profiles as JSON files in user's home directory under `MonitorProfiles/`
- Use atomic file operations to prevent corruption
- Validate profile data before loading
- Implement proper file permissions (0644 for files, 0755 for directories)
- Include both monitor and audio device configurations in profiles
- Store CLI tool paths and configurations for reproducibility

## Security Considerations

### System-Level Operations
- Monitor and audio configuration changes require appropriate system permissions
- Validate all CLI tool inputs before executing system changes
- Implement confirmation dialogs for destructive operations
- Log all CLI tool executions and system changes for audit purposes
- Handle CLI tool path resolution and existence checks

### CLI Tool Operations
- Use `filepath.Join()` for cross-platform path construction to CLI tools
- Validate CLI tool paths to prevent directory traversal
- Implement proper error handling for CLI tool I/O operations
- Create temporary files with appropriate permissions for CSV parsing
- Clean up temporary files after CLI tool operations

## Performance Guidelines

### Monitor and Audio Enumeration
- Cache monitor and audio information to avoid repeated CLI tool calls
- Implement lazy loading for expensive CLI operations
- Use background goroutines for long-running CLI operations
- Update UI asynchronously to prevent blocking
- Handle CLI tool timeouts and process cleanup

### Memory Management
- Avoid memory leaks in long-running applications
- Properly close CLI tool processes and temporary files when no longer needed
- Use efficient data structures for monitor and audio device lists
- Implement proper cleanup on application exit
- Handle CSV parsing memory efficiently

## Testing Requirements

### Unit Testing
- Test CLI tool integration with mock implementations
- Validate profile serialization/deserialization
- Test error handling paths for CLI tool failures
- Verify UI state management
- Test CSV parsing with various CLI tool output formats

### Integration Testing
- Test monitor detection on actual hardware with MultiMonitorTool
- Test audio device detection on actual hardware with SVCL
- Verify profile application works correctly with CLI tools
- Test error scenarios (missing CLI tools, invalid CLI output, disconnected devices)
- Validate cross-platform builds (even if functionality is limited)

## User Experience Guidelines

### Interface Design
- Provide clear visual indicators for monitor states
- Use intuitive icons and labels
- Implement responsive design for different screen sizes
- Provide immediate feedback for user actions

### Error Communication
- Display user-friendly error messages
- Provide actionable error resolution suggestions
- Use consistent error dialog styling
- Log technical details for debugging

## Development Workflow

### Build Process
- Use `wails build` for compilation
- Test on Windows with actual monitor and audio configurations
- Ensure CLI tools are included in the build distribution
- Validate cross-platform builds (even if functionality is limited)
- Use `go mod tidy` to maintain clean dependencies

### Version Control
- Commit platform-specific files together
- Document breaking changes in commit messages
- Use semantic versioning for releases
- Maintain clean git history

## Code Quality Standards

### Style Guidelines
- Follow Go formatting standards (`gofmt`)
- Use meaningful variable and function names
- Add comprehensive comments for CLI tool interactions
- Document CLI tool command-line arguments and output formats
- Implement proper package documentation

### Documentation
- Document all exported functions and types
- Explain CLI tool constants and output formats
- Provide usage examples for complex CLI operations
- Maintain up-to-date README with build instructions and CLI tool requirements
- Document CLI tool versions and compatibility requirements
