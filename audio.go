package main

import (
	"fmt"
)

func (a *App) SetPrimaryOutputDevice(deviceId string) error {
	var aDevice *AudioDevice
	for i := range a.audioDevices {
		if a.audioDevices[i].ID == deviceId {
			aDevice = &a.audioDevices[i]
			break
		}
	}

	if aDevice == nil {
		return fmt.Errorf("audio device not found")
	}

	return a.audioTools.SetPrimaryDevice(aDevice.ID)
}
