package monitors

import (
	"encoding/csv"
	"fmt"
	"monitor-profile-manager-wails/pkg/common"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// Column name constants
const (
	ColActive           = "Active"
	ColDisconnected     = "Disconnected"
	ColPrimary          = "Primary"
	ColName             = "Name"
	ColMonitorID        = "Short Monitor ID"
	ColMonitorName      = "Monitor Name"
	MultiMonitorToolExe = "MultiMonitorTool.exe"
)

// MonitorTools manages monitor operations with configurable tools directory
type MonitorTools struct {
	toolsDir string
}

// NewMonitorTools creates a new MonitorTools instance with the specified tools directory
func NewMonitorTools(toolsDir string) *MonitorTools {
	return &MonitorTools{toolsDir: toolsDir}
}

// MonitorInfo represents information about a monitor
type MonitorInfo struct {
	data map[string]string
}

func (m MonitorInfo) GetActive() bool        { return evaluateValueToBoolean(m.data[ColActive]) }
func (m MonitorInfo) GetDisconnected() bool  { return evaluateValueToBoolean(m.data[ColDisconnected]) }
func (m MonitorInfo) GetPrimary() bool       { return evaluateValueToBoolean(m.data[ColPrimary]) }
func (m MonitorInfo) GetName() string        { return m.data[ColName] }
func (m MonitorInfo) GetMonitorID() string   { return m.data[ColMonitorID] }
func (m MonitorInfo) GetMonitorName() string { return m.data[ColMonitorName] }

// Helper method to get any field by name
func (m MonitorInfo) GetField(fieldName string) string {
	return m.data[fieldName]
}

func evaluateValueToBoolean(flag string) bool {
	return flag == "Yes"
}

// hideConsoleCommand creates a command with hidden console window on Windows
func (m *MonitorTools) hideConsoleCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)

	// Hide console window on Windows
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}

	return cmd
}

// GetExecutableDir returns the directory where the executable is running
func GetExecutableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// GetMultiMonitorToolPath returns the full path to MultiMonitorTool.exe
func (m *MonitorTools) GetMultiMonitorToolPath() (string, error) {
	// If custom tools directory is set (embedded tools), use it
	if m.toolsDir != "" {
		return filepath.Join(m.toolsDir, "multimonitortool", MultiMonitorToolExe), nil
	}

	// Check if we're in development mode (wails dev)
	if common.IsDevelopmentMode() {
		// Development mode: use relative path from project root
		return filepath.Join("tools", "multimonitortool", MultiMonitorToolExe), nil
	}

	return "", fmt.Errorf("MultiMonitorTool.exe not found in current directory")
}

// CheckMultiMonitorToolExists verifies that MultiMonitorTool.exe exists
func (m *MonitorTools) CheckMultiMonitorToolExists() error {
	path, err := m.GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("MultiMonitorTool.exe not found in current directory")
	}
	return nil
}

// GetMonitorList retrieves the list of monitors using MultiMonitorTool
func (m *MonitorTools) GetMonitorList() ([]MonitorInfo, error) {
	// Check if MultiMonitorTool.exe exists
	if err := m.CheckMultiMonitorToolExists(); err != nil {
		return nil, err
	}

	// Get MultiMonitorTool path
	toolPath, err := m.GetMultiMonitorToolPath()
	if err != nil {
		return nil, err
	}

	// Export monitor list to CSV
	cmd := m.hideConsoleCommand(toolPath, "/List", "/scomma", "monitors.csv")
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute MultiMonitorTool: %w", err)
	}

	// Read and parse the CSV file
	file, err := os.Open("monitors.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open monitors.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	// Close file immediately after reading
	file.Close()

	// Parse monitor records (skip header)
	var monitors []MonitorInfo

	// Get column indexes dynamically
	colIndexes := make(map[string]int)
	if len(records) > 0 {
		header := records[0]
		for i, colName := range header {
			colIndexes[strings.TrimSpace(colName)] = i
		}
	}

	// Required columns
	requiredCols := []string{ColActive, ColDisconnected, ColPrimary, ColName, ColMonitorID, ColMonitorName}

	for i, row := range records {
		if i == 0 {
			// Skip header row
			continue
		}

		// Create monitor info using dynamic column indexes
		monitorData := make(map[string]string)
		validRow := true

		for _, colName := range requiredCols {
			if idx, exists := colIndexes[colName]; exists && idx < len(row) {
				monitorData[colName] = strings.TrimSpace(row[idx])
			} else {
				validRow = false
				break
			}
		}

		if validRow {
			monitors = append(monitors, MonitorInfo{data: monitorData})
		}
	}

	// Clean up CSV file with retry logic
	for retry := range 3 {
		err = os.Remove("monitors.csv")
		if err == nil {
			break // Successfully removed
		}
		if retry < 2 {
			// Wait a bit before retrying
			time.Sleep(100 * time.Millisecond)
		} else {
			// Log warning but don't fail the operation
			fmt.Printf("Warning: failed to remove monitors.csv after retries: %v\n", err)
		}
	}

	return monitors, nil
}

// SaveMonitorConfig saves the current monitor configuration to a file
func (m *MonitorTools) SaveMonitorConfig(configPath string) error {
	// Check if MultiMonitorTool.exe exists
	if err := m.CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := m.GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	// Execute the save config command
	cmd := m.hideConsoleCommand(toolPath, "/SaveConfig", configPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to save monitor configuration: %w", err)
	}
	return nil
}

// DisableMonitor disables the specified monitor
func (m *MonitorTools) DisableMonitor(monitorId string) error {
	// Check if MultiMonitorTool.exe exists
	if err := m.CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := m.GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	cmd := m.hideConsoleCommand(toolPath, "/disable", monitorId)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to disable monitor: %w", err)
	}
	return nil
}

// EnableMonitor enables the specified monitor
func (m *MonitorTools) EnableMonitor(monitorId string) error {
	// Check if MultiMonitorTool.exe exists
	if err := m.CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := m.GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	cmd := m.hideConsoleCommand(toolPath, "/enable", monitorId)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to enable monitor: %w", err)
	}
	return nil
}

// SetMonitorAsPrimary sets the specified monitor as the primary display
func (m *MonitorTools) SetMonitorAsPrimary(monitorId string) error {
	// Check if MultiMonitorTool.exe exists
	if err := m.CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := m.GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	cmd := m.hideConsoleCommand(toolPath, "/SetPrimary", monitorId)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set monitor as primary: %w", err)
	}
	return nil
}

// ApplyMonitorConfig applies the current monitor configuration
func (m *MonitorTools) ApplyMonitorConfig(configPath string) error {
	// Check if MultiMonitorTool.exe exists
	if err := m.CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := m.GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	// Execute the save config command
	cmd := m.hideConsoleCommand(toolPath, "/LoadConfig", configPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to load monitor configuration: %w", err)
	}

	return nil
}
