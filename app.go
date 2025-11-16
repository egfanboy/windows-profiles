package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Monitor struct {
	DeviceName  string `json:"deviceName"`
	DisplayName string `json:"displayName"`
	IsPrimary   bool   `json:"isPrimary"`
	IsActive    bool   `json:"isActive"`
	Bounds      Rect   `json:"bounds"`
}

type Rect struct {
	X      int32 `json:"x"`
	Y      int32 `json:"y"`
	Width  int32 `json:"width"`
	Height int32 `json:"height"`
}

type AudioDevice struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsDefault  bool   `json:"isDefault"`
	IsEnabled  bool   `json:"isEnabled"`
	DeviceType string `json:"deviceType"` // "output" or "input"
	State      string `json:"state"`      // "active", "disabled", "notpresent", "unplugged"
	Selected   bool   `json:"selected"`   // whether this device is selected for the profile
}

type IgnoreList struct {
	AudioDevices []string `json:"audioDevices"`
}

type Profile struct {
	Name         string        `json:"name"`
	Monitors     []Monitor     `json:"monitors"`
	AudioDevices []AudioDevice `json:"audioDevices"`
}

// MonitorManager interface defines OS-specific monitor operations
type MonitorManager interface {
	EnumDisplayMonitors() ([]Monitor, error)
	SetPrimaryMonitor(deviceName string) error
	SetMonitorState(deviceName string, active bool) error
	ApplyProfile(profile Profile) error
}

// AudioManager interface defines OS-specific audio operations
type AudioManager interface {
	EnumAudioDevices() ([]AudioDevice, error)
	SetDefaultAudioDevice(deviceID string, deviceType string) error
	EnableAudioDevice(deviceID string, enable bool) error
}

// App struct
type App struct {
	ctx            context.Context
	monitors       []Monitor
	audioDevices   []AudioDevice
	profiles       []Profile
	ignoreList     IgnoreList
	monitorManager MonitorManager
	audioManager   AudioManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}
	app.monitorManager = NewOSMonitorManager()
	app.audioManager = NewOSAudioManager()
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadIgnoreList()
	a.loadMonitors()
	a.loadAudioDevices()
	a.loadProfiles()
}

// loadIgnoreList loads the audio device ignore list from disk
func (a *App) loadIgnoreList() {
	ignoreListPath := a.getIgnoreListPath()
	data, err := os.ReadFile(ignoreListPath)
	if err != nil {
		a.ignoreList = IgnoreList{AudioDevices: []string{}}
		return
	}

	var ignoreList IgnoreList
	if err := json.Unmarshal(data, &ignoreList); err != nil {
		a.ignoreList = IgnoreList{AudioDevices: []string{}}
		return
	}

	a.ignoreList = ignoreList
}

// saveIgnoreList saves the audio device ignore list to disk
func (a *App) saveIgnoreList() error {
	ignoreListPath := a.getIgnoreListPath()
	data, err := json.MarshalIndent(a.ignoreList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal ignore list: %v", err)
	}

	if err := os.WriteFile(ignoreListPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save ignore list: %v", err)
	}

	return nil
}

// getIgnoreListPath returns the path where the ignore list is stored
func (a *App) getIgnoreListPath() string {
	profilesDir := a.getProfilesDir()
	return filepath.Join(profilesDir, "ignore_list.json")
}

// loadMonitors loads monitors using the OS-specific implementation
func (a *App) loadMonitors() {
	monitors, err := a.monitorManager.EnumDisplayMonitors()
	if err != nil {
		monitors = []Monitor{}
	}
	a.monitors = monitors
}

// loadAudioDevices loads audio devices using the OS-specific implementation
func (a *App) loadAudioDevices() {
	devices, err := a.audioManager.EnumAudioDevices()
	if err != nil {
		devices = []AudioDevice{}
	}
	a.audioDevices = devices
}

// loadProfiles loads saved monitor profiles from disk
func (a *App) loadProfiles() {
	profilesDir := a.getProfilesDir()
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return
	}

	files, err := os.ReadDir(profilesDir)
	if err != nil {
		return
	}

	a.profiles = []Profile{}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" && file.Name() != "ignore_list.json" {
			profilePath := filepath.Join(profilesDir, file.Name())
			data, err := os.ReadFile(profilePath)
			if err != nil {
				continue
			}

			var profile Profile
			if err := json.Unmarshal(data, &profile); err != nil {
				continue
			}

			a.profiles = append(a.profiles, profile)
		}
	}
}

// getProfilesDir returns the directory where profiles are stored
func (a *App) getProfilesDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./profiles"
	}
	return filepath.Join(homeDir, "MonitorProfiles")
}

// GetMonitors returns the current list of monitors
func (a *App) GetMonitors() []Monitor {
	a.loadMonitors()
	return a.monitors
}

// GetAudioDevices returns the current list of audio devices
func (a *App) GetAudioDevices() []AudioDevice {
	a.loadAudioDevices()
	return a.audioDevices
}

// GetAudioDevicesWithIgnoreStatus returns audio devices with ignore status
func (a *App) GetAudioDevicesWithIgnoreStatus() map[string]interface{} {
	a.loadAudioDevices()

	var filteredDevices []AudioDevice
	var ignoredDevices []AudioDevice

	for _, device := range a.audioDevices {
		isIgnored := a.isDeviceIgnored(device.ID)
		if isIgnored {
			ignoredDevices = append(ignoredDevices, device)
		} else {
			filteredDevices = append(filteredDevices, device)
		}
	}

	return map[string]interface{}{
		"filtered": filteredDevices,
		"ignored":  ignoredDevices,
	}
}

// isDeviceIgnored checks if a device is in the ignore list
func (a *App) isDeviceIgnored(deviceID string) bool {
	for _, ignoredID := range a.ignoreList.AudioDevices {
		if ignoredID == deviceID {
			return true
		}
	}
	return false
}

// GetProfiles returns the list of saved profiles
func (a *App) GetProfiles() []Profile {
	a.loadProfiles()
	return a.profiles
}

// SaveProfile saves a monitor profile with the given name
func (a *App) SaveProfile(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	// Filter out ignored and unselected audio devices
	var selectedAudioDevices []AudioDevice
	for _, device := range a.audioDevices {
		if !a.isDeviceIgnored(device.ID) && device.Selected {
			selectedAudioDevices = append(selectedAudioDevices, device)
		}
	}

	profile := Profile{
		Name:         name,
		Monitors:     a.monitors,
		AudioDevices: selectedAudioDevices,
	}

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %v", err)
	}

	profilePath := filepath.Join(a.getProfilesDir(), name+".json")
	if err := os.WriteFile(profilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save profile: %v", err)
	}

	a.loadProfiles()
	return nil
}

// ApplyProfile applies a monitor profile by name
func (a *App) ApplyProfile(profileName string) error {
	for _, profile := range a.profiles {
		if profile.Name == profileName {
			// Apply monitors
			if err := a.monitorManager.ApplyProfile(profile); err != nil {
				return fmt.Errorf("failed to apply monitor settings: %v", err)
			}

			// Apply audio devices
			for _, device := range profile.AudioDevices {
				if device.IsDefault && device.DeviceType == "output" {
					if err := a.audioManager.SetDefaultAudioDevice(device.ID, "output"); err != nil {
						return fmt.Errorf("failed to set default audio device: %v", err)
					}
				}
				if device.IsDefault && device.DeviceType == "input" {
					if err := a.audioManager.SetDefaultAudioDevice(device.ID, "input"); err != nil {
						return fmt.Errorf("failed to set default audio device: %v", err)
					}
				}
				if err := a.audioManager.EnableAudioDevice(device.ID, device.IsEnabled); err != nil {
					return fmt.Errorf("failed to set audio device state: %v", err)
				}
			}

			return nil
		}
	}
	return fmt.Errorf("profile '%s' not found", profileName)
}

// RefreshMonitors refreshes the monitor list
func (a *App) RefreshMonitors() []Monitor {
	a.loadMonitors()
	return a.monitors
}

// RefreshAudioDevices refreshes the audio device list
func (a *App) RefreshAudioDevices() map[string]interface{} {
	a.loadAudioDevices()
	return a.GetAudioDevicesWithIgnoreStatus()
}

// IgnoreAudioDevice adds a device to the ignore list
func (a *App) IgnoreAudioDevice(deviceID string) error {
	// Check if already ignored
	if a.isDeviceIgnored(deviceID) {
		return fmt.Errorf("device is already ignored")
	}

	a.ignoreList.AudioDevices = append(a.ignoreList.AudioDevices, deviceID)
	return a.saveIgnoreList()
}

// UnignoreAudioDevice removes a device from the ignore list
func (a *App) UnignoreAudioDevice(deviceID string) error {
	for i, ignoredID := range a.ignoreList.AudioDevices {
		if ignoredID == deviceID {
			a.ignoreList.AudioDevices = append(a.ignoreList.AudioDevices[:i], a.ignoreList.AudioDevices[i+1:]...)
			return a.saveIgnoreList()
		}
	}
	return fmt.Errorf("device is not in ignore list")
}

// SetAudioDeviceSelection sets the selection state of an audio device
func (a *App) SetAudioDeviceSelection(deviceID string, selected bool) error {
	for i := range a.audioDevices {
		if a.audioDevices[i].ID == deviceID {
			a.audioDevices[i].Selected = selected
			return nil
		}
	}
	return fmt.Errorf("device not found")
}

// GetSelectedAudioDevices returns the list of selected audio devices
func (a *App) GetSelectedAudioDevices() []AudioDevice {
	var selected []AudioDevice
	for _, device := range a.audioDevices {
		if device.Selected && !a.isDeviceIgnored(device.ID) {
			selected = append(selected, device)
		}
	}
	return selected
}
