//go:build windows
// +build windows

package main

import (
	"fmt"
)

// WindowsMonitorManager implements MonitorManager for Windows
type WindowsMonitorManager struct{}

// NewOSMonitorManager creates a Windows-specific monitor manager
func NewOSMonitorManager() MonitorManager {
	return &WindowsMonitorManager{}
}

// EnumDisplayMonitors enumerates all connected monitors on Windows
func (w *WindowsMonitorManager) EnumDisplayMonitors() ([]Monitor, error) {
	return []Monitor{}, nil
}

// SetPrimaryMonitor sets a monitor as the primary display on Windows
func (w *WindowsMonitorManager) SetPrimaryMonitor(deviceName string) error {
	return fmt.Errorf("not implemented")
}

// SetMonitorState activates or deactivates a monitor on Windows
func (w *WindowsMonitorManager) SetMonitorState(deviceName string, active bool) error {
	return fmt.Errorf("not implemented")
}

// ApplyProfile applies a monitor profile on Windows
func (w *WindowsMonitorManager) ApplyProfile(profile Profile) error {
	return fmt.Errorf("not implemented")
}
