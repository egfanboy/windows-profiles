package audio

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

// AudioTools manages audio device operations with configurable tools directory
type AudioTools struct {
	toolsDir string
}

// NewAudioTools creates a new AudioTools instance with the specified tools directory
func NewAudioTools(toolsDir string) *AudioTools {
	return &AudioTools{toolsDir: toolsDir}
}

// AudioDeviceInfo represents information about an audio device
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

// hideConsoleCommand creates a command with hidden console window on Windows
func (a *AudioTools) hideConsoleCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)

	// Hide console window on Windows
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}

	return cmd
}

// GetSvclPath returns the full path to svcl.exe
func (a *AudioTools) GetSvclPath() (string, error) {
	// If custom tools directory is set (embedded tools), use it
	if a.toolsDir != "" {
		return filepath.Join(a.toolsDir, "svcl", SvclExe), nil
	}

	// Check if we're in development mode (wails dev)
	if common.IsDevelopmentMode() {
		// Development mode: use relative path from project root
		return filepath.Join("tools", "svcl", SvclExe), nil
	}

	return "", fmt.Errorf("svcl.exe not found in tools/svcl directory")
}

// CheckSvclExists verifies that svcl.exe exists
func (a *AudioTools) CheckSvclExists() error {
	path, err := a.GetSvclPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("svcl.exe not found in tools/svcl directory")
	}
	return nil
}

// GetActiveOutputDevices retrieves active output devices using svcl.exe /scomma
func (a *AudioTools) GetActiveOutputDevices() ([]AudioDeviceInfo, error) {
	// Check if svcl.exe exists
	if err := a.CheckSvclExists(); err != nil {
		return nil, err
	}

	// Get svcl.exe path
	toolPath, err := a.GetSvclPath()
	if err != nil {
		return nil, err
	}

	// Execute svcl.exe with /scomma and capture stdout
	cmd := a.hideConsoleCommand(toolPath, "/scomma")
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
func (a *AudioTools) SetPrimaryDevice(commandLineId string) error {
	// Check if svcl.exe exists
	if err := a.CheckSvclExists(); err != nil {
		return err
	}

	// Get svcl.exe path
	toolPath, err := a.GetSvclPath()
	if err != nil {
		return err
	}

	cmd := a.hideConsoleCommand(toolPath, "/SetDefault", commandLineId, "all")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set primary audio device: %w", err)
	}
	return nil
}
