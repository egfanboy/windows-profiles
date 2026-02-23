package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const (
	PROFILE_DIR       = "profiles"
	PROFILE_FILE_NAME = "profiles.json"
)

type Profile struct {
	Name               string `json:"name"`
	MonitorProfileName string `json:"monitorProfileName"`
	AudioProfileName   string `json:"audioProfileName"`
}

// SaveProfile saves a monitor profile with the given name
func (a *App) SaveProfile(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	profile := Profile{
		Name:               name,
		MonitorProfileName: "",
		AudioProfileName:   "",
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
