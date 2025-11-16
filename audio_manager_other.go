//go:build !windows

package main

import "fmt"

// OtherAudioManager implements AudioManager for non-Windows systems
type OtherAudioManager struct{}

func NewOSAudioManager() AudioManager {
	return &OtherAudioManager{}
}

func (o *OtherAudioManager) EnumAudioDevices() ([]AudioDevice, error) {
	return []AudioDevice{}, fmt.Errorf("audio device enumeration not supported on this platform")
}

func (o *OtherAudioManager) SetDefaultAudioDevice(deviceID string, deviceType string) error {
	return fmt.Errorf("setting default audio device not supported on this platform")
}

func (o *OtherAudioManager) EnableAudioDevice(deviceID string, enable bool) error {
	return fmt.Errorf("enabling/disabling audio devices not supported on this platform")
}
