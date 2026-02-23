package main

import (
	"fmt"
	"monitor-profile-manager-wails/pkg/monitors"
	"path/filepath"
)

// getMonitorConfigPath returns the file path for a monitor profile config file
func (a *App) getMonitorConfigPath(profileName string) string {
	profilesDir := a.getProfilesDir()
	return filepath.Join(profilesDir, profileName+"-monitor.cfg")
}

// SaveMonitorProfile saves the current monitor configuration to a profile file
// The profile will be saved in the profiles directory with the given profileName
func (a *App) saveMonitorProfile(profileName string) error {
	if profileName == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	// Create the full path with .cfg extension
	profilePath := a.getMonitorConfigPath(profileName)

	return monitors.SaveMonitorConfig(profilePath)
}
