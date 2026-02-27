package main

import (
	"context"
	"encoding/json"
	"fmt"
	"monitor-profile-manager-wails/pkg/audio"
	"monitor-profile-manager-wails/pkg/monitors"
	"os"
	"path/filepath"
	"sync"
)

const (
	SETTINGS_DIRECTORY = ".windows-profile-manager"
)

// Global mutex to prevent concurrent startup operations
var startupMutex sync.Mutex

// Global mutex to prevent concurrent monitor enumeration
var monitorEnumMutex sync.Mutex

// Global mutex to prevent concurrent audio device enumeration
var audioEnumMutex sync.Mutex

type Monitor struct {
	DeviceName  string `json:"deviceName"`
	DisplayName string `json:"displayName"`
	IsPrimary   bool   `json:"isPrimary"`
	IsActive    bool   `json:"isActive"`
	IsEnabled   bool   `json:"isEnabled"` // user-controlled enable/disable state
	MonitorId   string `json:"monitorId"`
	Nickname    string `json:"nickname"` // optional custom nickname
}

type AudioDevice struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
	IsEnabled bool   `json:"isEnabled"`
	Selected  bool   `json:"selected"` // whether this device is selected for the profile
	Nickname  string `json:"nickname"` // optional custom nickname
}

type IgnoreList struct {
	AudioDevices []string `json:"audioDevices"`
}

type NicknameStorage struct {
	Monitors     map[string]string `json:"monitors"`     // deviceID -> nickname
	AudioDevices map[string]string `json:"audioDevices"` // deviceID -> nickname
}

// App struct holds the application state
type App struct {
	ctx          context.Context
	monitors     []Monitor
	audioDevices []AudioDevice
	profiles     []Profile
	ignoreList   IgnoreList
	nicknames    NicknameStorage
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	// Prevent concurrent startup operations
	startupMutex.Lock()
	defer startupMutex.Unlock()

	a.ctx = ctx

	// Load all components with error handling to prevent crashes
	if err := func() error {
		a.loadIgnoreList()
		a.loadNicknames()
		a.loadMonitors()
		a.loadAudioDevices()
		a.loadProfiles()
		return nil
	}(); err != nil {
		// If startup fails, initialize with empty defaults
		a.monitors = []Monitor{}
		a.audioDevices = []AudioDevice{}
		a.profiles = []Profile{}
		a.ignoreList = IgnoreList{AudioDevices: []string{}}
		a.nicknames = NicknameStorage{
			Monitors:     make(map[string]string),
			AudioDevices: make(map[string]string),
		}
	}
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

	// Prevent concurrent monitor loading
	monitorEnumMutex.Lock()
	defer monitorEnumMutex.Unlock()

	monitors, err := monitors.GetMonitorList()
	appMonitors := make([]Monitor, 0)
	if err == nil {
		for _, monitor := range monitors {
			appMonitor := Monitor{}
			appMonitor.IsActive = monitor.GetActive()
			appMonitor.IsPrimary = monitor.GetPrimary()
			appMonitor.DeviceName = monitor.GetName()
			appMonitor.MonitorId = monitor.GetMonitorID()
			appMonitor.DisplayName = monitor.GetMonitorName()

			appMonitors = append(appMonitors, appMonitor)
		}

	}

	if len(appMonitors) > 0 {
		// Apply nicknames to monitors
		for i := range appMonitors {
			if nickname := a.GetMonitorNickname(appMonitors[i].DeviceName); nickname != "" {
				appMonitors[i].Nickname = nickname
			}
			// Ensure isEnabled is always set
			if !appMonitors[i].IsEnabled {
				appMonitors[i].IsEnabled = true
			}
		}

		a.monitors = appMonitors
	}

}

// loadAudioDevices loads audio devices using the OS-specific implementation
func (a *App) loadAudioDevices() {
	// Prevent concurrent audio device loading
	audioEnumMutex.Lock()
	defer audioEnumMutex.Unlock()

	devices := make([]AudioDevice, 0)

	svclDevices, err := audio.GetActiveOutputDevices()

	if err != nil {
		devices = []AudioDevice{}
	} else {
		for _, device := range svclDevices {
			ad := AudioDevice{}

			ad.IsDefault = device.IsPrimary()
			ad.IsEnabled = device.IsActive()
			ad.Name = device.GetName()
			ad.ID = device.GetCommandLineID()
			devices = append(devices, ad)
		}
	}

	// Apply nicknames to audio devices
	for i := range devices {
		if nickname := a.GetAudioDeviceNickname(devices[i].ID); nickname != "" {
			devices[i].Nickname = nickname
		}
	}

	a.audioDevices = devices
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

	filteredDevices := make([]AudioDevice, 0)
	ignoredDevices := make([]AudioDevice, 0)

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

// Nickname management methods

// SetMonitorNickname sets a custom nickname for a monitor
func (a *App) SetMonitorNickname(deviceName string, nickname string) error {
	if a.nicknames.Monitors == nil {
		a.nicknames.Monitors = make(map[string]string)
	}
	if nickname == "" {
		delete(a.nicknames.Monitors, deviceName)
	} else {
		a.nicknames.Monitors[deviceName] = nickname
	}
	return a.saveNicknames()
}

// SetAudioDeviceNickname sets a custom nickname for an audio device
func (a *App) SetAudioDeviceNickname(deviceID string, nickname string) error {
	if a.nicknames.AudioDevices == nil {
		a.nicknames.AudioDevices = make(map[string]string)
	}
	if nickname == "" {
		delete(a.nicknames.AudioDevices, deviceID)
	} else {
		a.nicknames.AudioDevices[deviceID] = nickname
	}
	return a.saveNicknames()
}

// GetMonitorNickname gets the custom nickname for a monitor
func (a *App) GetMonitorNickname(deviceName string) string {
	if a.nicknames.Monitors == nil {
		return ""
	}
	return a.nicknames.Monitors[deviceName]
}

// GetAudioDeviceNickname gets the custom nickname for an audio device
func (a *App) GetAudioDeviceNickname(deviceID string) string {
	if a.nicknames.AudioDevices == nil {
		return ""
	}
	return a.nicknames.AudioDevices[deviceID]
}

// saveNicknames saves the nickname storage to disk
func (a *App) saveNicknames() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(configDir, "monitor-profile-manager")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return err
	}

	nicknamesFile := filepath.Join(appDir, "nicknames.json")
	data, err := json.MarshalIndent(a.nicknames, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(nicknamesFile, data, 0644)
}

// loadNicknames loads the nickname storage from disk
func (a *App) loadNicknames() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	nicknamesFile := filepath.Join(configDir, "monitor-profile-manager", "nicknames.json")
	data, err := os.ReadFile(nicknamesFile)
	if err != nil {
		if os.IsNotExist(err) {
			// No nicknames file exists, initialize empty storage
			a.nicknames = NicknameStorage{
				Monitors:     make(map[string]string),
				AudioDevices: make(map[string]string),
			}
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &a.nicknames)
}
