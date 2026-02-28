package main

import (
	"fmt"
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

	return a.monitorTools.SaveMonitorConfig(profilePath)
}

func (a *App) SetMonitorEnabledState(monitorId string, active bool) error {
	if active {
		return a.monitorTools.EnableMonitor(monitorId)
	}
	return a.monitorTools.DisableMonitor(monitorId)
}

// SetMonitorPrimary sets a monitor as the primary monitor
func (a *App) SetMonitorPrimary(monitorId string) error {
	// Find the monitor and update primary status
	var monitor *Monitor
	for i := range a.monitors {
		if a.monitors[i].MonitorId == monitorId {
			monitor = &a.monitors[i]
			break
		}
	}

	if monitor == nil {
		return fmt.Errorf("monitor not found")
	}

	if !monitor.IsActive {
		return fmt.Errorf("monitor is not active")
	}

	if !monitor.IsEnabled {
		return fmt.Errorf("monitor is not enabled")
	}

	if monitor.IsPrimary {
		return nil
	}

	return a.monitorTools.SetMonitorAsPrimary(monitorId)
}
