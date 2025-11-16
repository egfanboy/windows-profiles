//go:build !windows
// +build !windows

package main

import (
	"fmt"
)

// OtherMonitorManager implements MonitorManager for non-Windows platforms
type OtherMonitorManager struct{}

// NewOSMonitorManager creates a non-Windows monitor manager
func NewOSMonitorManager() MonitorManager {
	return &OtherMonitorManager{}
}

// EnumDisplayMonitors returns an error on non-Windows platforms
func (o *OtherMonitorManager) EnumDisplayMonitors() ([]Monitor, error) {
	return nil, fmt.Errorf("monitor management is only supported on Windows")
}

// SetPrimaryMonitor returns an error on non-Windows platforms
func (o *OtherMonitorManager) SetPrimaryMonitor(deviceName string) error {
	return fmt.Errorf("monitor management is only supported on Windows")
}

// SetMonitorState returns an error on non-Windows platforms
func (o *OtherMonitorManager) SetMonitorState(deviceName string, active bool) error {
	return fmt.Errorf("monitor management is only supported on Windows")
}

// ApplyProfile returns an error on non-Windows platforms
func (o *OtherMonitorManager) ApplyProfile(profile Profile) error {
	return fmt.Errorf("monitor management is only supported on Windows")
}
