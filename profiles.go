package main

import (
	"encoding/json"
	"fmt"
	"monitor-profile-manager-wails/pkg/audio"
	"monitor-profile-manager-wails/pkg/monitors"
	"os"
	"path"
	"path/filepath"
)

const (
	PROFILE_DIR       = "profiles"
	PROFILE_FILE_NAME = "profiles.json"
)

type AudioProfile struct {
	DefaultOutputDeviceId string `json:"defaultOutputDeviceId"`
}

type SaveProfileRequest struct {
	Name                  string `json:"name"`
	DefaultOutputDeviceId string `json:"defaultOutputDeviceId"`
}

// MultiMonitorTool has integrated profile management. Therefore we can use the name
// to derive the profile. However, svcl does not have profile management, so we need
// to save the audio information as part of the profile.

type Profile struct {
	Name  string       `json:"name"`
	Audio AudioProfile `json:"audio"`
}

// SaveProfile saves a monitor profile with the given profile data
func (a *App) SaveProfile(request SaveProfileRequest) error {
	if request.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	profile := Profile{
		Name: request.Name,
		Audio: AudioProfile{
			DefaultOutputDeviceId: request.DefaultOutputDeviceId,
		},
	}

	err := a.saveMonitorProfile(profile.Name)
	if err != nil {
		return err
	}

	a.profiles = append(a.profiles, profile)

	return a.saveProfilesToDisk()
}

func (a *App) DeleteProfile(profileName string) error {
	newProfiles := make([]Profile, 0)
	for _, profile := range a.profiles {
		if profile.Name != profileName {
			newProfiles = append(newProfiles, profile)
		}
	}

	a.profiles = newProfiles

	// Clean up the monitor .cfg file
	monitorConfigPath := a.getMonitorConfigPath(profileName)
	if err := os.Remove(monitorConfigPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove monitor config file: %v", err)
	}

	err := a.saveProfilesToDisk()
	if err != nil {
		return err
	}

	return nil
}

// GetProfiles returns the list of saved profiles
func (a *App) GetProfiles() []Profile {
	a.loadProfiles()
	return a.profiles
}

// ApplyProfile applies a monitor profile by name
func (a *App) ApplyProfile(profileName string) error {
	var profile *Profile

	for i := range a.profiles {
		p := a.profiles[i]
		if p.Name == profileName {
			profile = &p
			break
		}
	}

	if profile == nil {
		return fmt.Errorf("profile not found: %s", profileName)
	}

	// Apply monitor profile
	err := monitors.ApplyMonitorConfig(a.getMonitorConfigPath(profileName))
	if err != nil {
		return err
	}

	// Apply audio profile
	if profile.Audio.DefaultOutputDeviceId != "" {
		err = audio.SetPrimaryDevice(profile.Audio.DefaultOutputDeviceId)
		if err != nil {
			return err
		}
	}

	return nil
}

// getProfilesDir returns the directory where profiles are stored
func (a *App) getProfilesDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./profiles"
	}
	return filepath.Join(homeDir, SETTINGS_DIRECTORY)
}

// loadProfiles loads saved monitor profiles from disk
func (a *App) loadProfiles() {
	a.profiles = []Profile{}
	profilesDir := a.getProfilesDir()
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return
	}

	data, err := os.ReadFile(path.Join(profilesDir, PROFILE_FILE_NAME))
	if err != nil {
		return
	}

	var profiles []Profile

	if err := json.Unmarshal(data, &profiles); err != nil {
		return
	}

	a.profiles = profiles
}

func (a *App) saveProfilesToDisk() error {
	profilesDir := a.getProfilesDir()
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(a.profiles)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(profilesDir, PROFILE_FILE_NAME), data, 0644)
}
