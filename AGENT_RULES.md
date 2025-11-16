# Agent Rules for Monitor Profile Manager Development

## Architecture Guidelines

### Platform-Specific Code
- **Always** use build constraints for platform-specific functionality (`//go:build windows` and `//go:build !windows`)
- Implement Windows-specific APIs in separate files with `_windows.go` suffix
- Provide fallback implementations for non-Windows platforms in `_other.go` files
- Never mix platform-specific code with cross-platform logic

### Windows API Integration
- Use `golang.org/x/sys/windows` for Windows API calls
- Prefer lazy loading of DLLs with `windows.NewLazySystemDLL()`
- Always check return values from Windows API calls
- Convert between Go strings and UTF16 pointers using `windows.StringToUTF16Ptr()` and `windows.UTF16ToString()`
- Handle all possible return codes from Windows API functions

### Error Handling
- **Never** ignore error returns
- Wrap Windows API errors with descriptive context using `fmt.Errorf()`
- Provide user-friendly error messages in the GUI while logging technical details
- Implement proper error propagation from low-level API calls to high-level UI

### GUI Development with Fyne
- Use container-based layouts (VBox, HBox, Grid) for responsive design
- Implement proper data refresh patterns when underlying state changes
- Use dialogs for user input and error display
- Follow Fyne's widget lifecycle patterns for list updates

## Code Structure Rules

### File Organization
```
main.go              - Main application entry point and UI logic
monitor_windows.go   - Windows-specific monitor management
monitor_other.go     - Non-Windows fallback implementations
go.mod              - Module dependencies
AGENT_RULES.md      - This file
README.md           - Project documentation
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

## Security Considerations

### System-Level Operations
- Monitor configuration changes require appropriate system permissions
- Validate all user inputs before applying system changes
- Implement confirmation dialogs for destructive operations
- Log all system changes for audit purposes

### File Operations
- Use `filepath.Join()` for cross-platform path construction
- Validate file paths to prevent directory traversal
- Implement proper error handling for file I/O operations
- Create directories with appropriate permissions

## Performance Guidelines

### Monitor Enumeration
- Cache monitor information to avoid repeated system calls
- Implement lazy loading for expensive operations
- Use background goroutines for long-running operations
- Update UI asynchronously to prevent blocking

### Memory Management
- Avoid memory leaks in long-running applications
- Properly close Windows handles when no longer needed
- Use efficient data structures for monitor lists
- Implement proper cleanup on application exit

## Testing Requirements

### Unit Testing
- Test Windows API integration with mock implementations
- Validate profile serialization/deserialization
- Test error handling paths
- Verify UI state management

### Integration Testing
- Test monitor detection on actual hardware
- Verify profile application works correctly
- Test error scenarios (disconnected monitors, invalid profiles)
- Validate cross-platform compatibility

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
- Use `go build` for compilation
- Test on Windows with actual monitor configurations
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
- Add comprehensive comments for Windows API interactions
- Implement proper package documentation

### Documentation
- Document all exported functions and types
- Explain Windows API constants and structures
- Provide usage examples for complex operations
- Maintain up-to-date README with build instructions
