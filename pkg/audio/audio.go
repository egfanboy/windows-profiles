package audio

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Column name constants
const (
	ColName          = "Device Name"
	ColCommandLineID = "Command-Line Friendly ID"
	ColDeviceType    = "Device Type"
	ColDeviceState   = "Device State"
	ColDefault       = "Default"
	SvclExe          = "svcl.exe"
)

type AudioDeviceInfo struct {
	data map[string]string
}

func (a AudioDeviceInfo) GetName() string          { return a.data[ColName] }
func (a AudioDeviceInfo) GetCommandLineID() string { return a.data[ColCommandLineID] }
func (a AudioDeviceInfo) GetDeviceType() string    { return a.data[ColDeviceType] }
func (a AudioDeviceInfo) GetDeviceState() string   { return a.data[ColDeviceState] }
func (a AudioDeviceInfo) GetDefault() string       { return a.data[ColDefault] }

// Helper method to get any field by name
func (a AudioDeviceInfo) GetField(fieldName string) string {
	return a.data[fieldName]
}

// Helper method to check if device is primary (default)
func (a AudioDeviceInfo) IsPrimary() bool {
	isDefault := a.GetDefault()

	return isDefault == "Render"
}

// Helper method to check if device is active
func (a AudioDeviceInfo) IsActive() bool {
	return strings.Contains(strings.ToLower(a.data[ColDeviceState]), "active")
}

// GetExecutableDir returns the directory where the executable is running
func GetExecutableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// GetSvclPath returns the full path to svcl.exe
func GetSvclPath() (string, error) {
	// Check if we're in development mode (wails dev) by looking for go.mod
	if _, err := os.Stat("go.mod"); err == nil {
		// Development mode: use relative path from project root
		return filepath.Join("tools", "svcl", SvclExe), nil
	}
	// Production mode: tools should be in a 'tools' subdirectory next to the exe
	exeDir, err := GetExecutableDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(exeDir, "tools", "svcl", SvclExe), nil
}

// CheckSvclExists verifies that svcl.exe exists
func CheckSvclExists() error {
	path, err := GetSvclPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("svcl.exe not found in tools/svcl directory")
	}
	return nil
}

// GetActiveOutputDevices retrieves active output devices using svcl.exe /scomma
func GetActiveOutputDevices() ([]AudioDeviceInfo, error) {
	// Check if svcl.exe exists
	if err := CheckSvclExists(); err != nil {
		return nil, err
	}

	// Get svcl.exe path
	toolPath, err := GetSvclPath()
	if err != nil {
		return nil, err
	}

	// Execute svcl.exe with /scomma and capture stdout
	cmd := exec.Command(toolPath, "/scomma")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute svcl.exe: %w", err)
	}

	// Parse the CSV output from stdout
	reader := csv.NewReader(strings.NewReader(string(output)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV output: %w", err)
	}

	// Get column indexes dynamically from header
	colIndexes := make(map[string]int)
	if len(records) > 0 {
		header := records[0]
		for i, colName := range header {
			colIndexes[strings.TrimSpace(colName)] = i
		}
	}

	var devices []AudioDeviceInfo

	// Filter for active output devices
	for i, row := range records {
		if i == 0 {
			// Skip header row
			continue
		}

		// Check if row has any data
		if len(row) == 0 {
			continue
		}

		// Get values for filtering
		direction := ""
		deviceState := ""
		deviceType := ""

		if idx, exists := colIndexes["Direction"]; exists && idx < len(row) {
			direction = strings.TrimSpace(row[idx])
		}

		if idx, exists := colIndexes["Device State"]; exists && idx < len(row) {
			deviceState = strings.TrimSpace(row[idx])
		}

		if idx, exists := colIndexes["Type"]; exists && idx < len(row) {
			deviceType = strings.TrimSpace(row[idx])
		}

		// Apply filters: Direction == "Render", State == "Active", Type == "Device"
		if direction == "Render" && deviceState == "Active" && deviceType == "Device" {
			// Create device info with relevant data
			deviceData := make(map[string]string)

			// Add all available columns
			for colName, idx := range colIndexes {
				if idx < len(row) {
					deviceData[colName] = strings.TrimSpace(row[idx])
				}
			}

			devices = append(devices, AudioDeviceInfo{data: deviceData})
		}
	}

	return devices, nil
}

// SetPrimaryDevice sets the specified audio device as the primary/default device
func SetPrimaryDevice(deviceID string) error {
	// Check if svcl.exe exists
	if err := CheckSvclExists(); err != nil {
		return err
	}

	// Get svcl.exe path
	toolPath, err := GetSvclPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(toolPath, "/SetDefault", deviceID, "all")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set primary audio device: %w", err)
	}
	return nil
}
