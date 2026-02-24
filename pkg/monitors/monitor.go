package monitors

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// GetExecutableDir returns the directory where the executable is running
func GetExecutableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// GetMultiMonitorToolPath returns the full path to MultiMonitorTool.exe
func GetMultiMonitorToolPath() (string, error) {
	// Check if we're in development mode (wails dev) by looking for go.mod
	if _, err := os.Stat("go.mod"); err == nil {
		// Development mode: use relative path from project root
		return filepath.Join("tools", "multimonitortool", MultiMonitorToolExe), nil
	}
	// Production mode: tools should be in a 'tools' subdirectory next to the exe
	exeDir, err := GetExecutableDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(exeDir, "tools", "multimonitortool", MultiMonitorToolExe), nil
}

// CheckMultiMonitorToolExists verifies that MultiMonitorTool.exe exists
func CheckMultiMonitorToolExists() error {
	path, err := GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("MultiMonitorTool.exe not found in current directory")
	}
	return nil
}

// GetMonitorList retrieves the list of monitors using MultiMonitorTool
func GetMonitorList() ([]MonitorInfo, error) {
	// Check if MultiMonitorTool.exe exists
	if err := CheckMultiMonitorToolExists(); err != nil {
		return nil, err
	}

	// Get MultiMonitorTool path
	toolPath, err := GetMultiMonitorToolPath()
	if err != nil {
		return nil, err
	}

	// Export monitor list to CSV
	cmd := exec.Command(toolPath, "/List", "/scomma", "monitors.csv")
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

// SaveMonitorConfig saves the current monitor configuration to a config file
func SaveMonitorConfig(configPath string) error {
	// Check if MultiMonitorTool.exe exists
	if err := CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	// Execute the save config command
	cmd := exec.Command(toolPath, "/SaveConfig", configPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to save monitor configuration: %w", err)
	}
	return nil
}

// SetPrimaryMonitor sets the specified monitor as the primary monitor
func SetPrimaryMonitor(monitorId string) error {
	// Check if MultiMonitorTool.exe exists
	if err := CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(toolPath, "/SetPrimary", monitorId, "/Silent")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set primary monitor: %w", err)
	}
	return nil
}

func DisableMonitor(monitorId string) error {
	// Check if MultiMonitorTool.exe exists
	if err := CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(toolPath, "/disable", monitorId)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to disable monitor: %w", err)
	}
	return nil
}

func EnableMonitor(monitorId string) error {
	// Check if MultiMonitorTool.exe exists
	if err := CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(toolPath, "/enable", monitorId)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to enable monitor: %w", err)
	}
	return nil
}

func SetMonitorAsPrimary(monitorId string) error {
	// Check if MultiMonitorTool.exe exists
	if err := CheckMultiMonitorToolExists(); err != nil {
		return err
	}

	// Get MultiMonitorTool path
	toolPath, err := GetMultiMonitorToolPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(toolPath, "/SetPrimary", monitorId)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set monitor as primary: %w", err)
	}
	return nil
}
