//go:build windows
// +build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                    = windows.NewLazySystemDLL("user32.dll")
	procEnumDisplayDevices    = user32.NewProc("EnumDisplayDevicesW")
	procEnumDisplaySettings   = user32.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettings = user32.NewProc("ChangeDisplaySettingsExW")
)

const (
	ENUM_CURRENT_SETTINGS  = 0xFFFFFFFF
	CDS_UPDATEREGISTRY     = 0x00000001
	DISP_CHANGE_SUCCESSFUL = 0
	DISP_CHANGE_RESTART    = 1
	DISP_CHANGE_FAILED     = uintptr(0xFFFFFFFF)
	DISP_CHANGE_BADMODE    = uintptr(0xFFFFFFFE)
	DISP_CHANGE_NOTUPDATED = uintptr(0xFFFFFFFD)
	DISP_CHANGE_BADFLAGS   = uintptr(0xFFFFFFFC)
	DISP_CHANGE_BADPARAM   = uintptr(0xFFFFFFFB)
)

type DISPLAY_DEVICE struct {
	Cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

type DEVMODE struct {
	DmDeviceName       [32]uint16
	DmSpecVersion      uint16
	DmDriverVersion    uint16
	DmSize             uint16
	DmDriverExtra      uint16
	DmFields           uint32
	DmPosition         POINT
	DmBitsPerPel       uint32
	DmPelsWidth        uint32
	DmPelsHeight       uint32
	DmDisplayFlags     uint32
	DmDisplayFrequency uint32
}

type POINT struct {
	X int32
	Y int32
}

const (
	DISPLAY_DEVICE_ACTIVE         = 0x00000001
	DISPLAY_DEVICE_PRIMARY_DEVICE = 0x00000004
)

// WindowsMonitorManager implements MonitorManager for Windows
type WindowsMonitorManager struct{}

// NewOSMonitorManager creates a Windows-specific monitor manager
func NewOSMonitorManager() MonitorManager {
	return &WindowsMonitorManager{}
}

// EnumDisplayMonitors enumerates all connected monitors on Windows
func (w *WindowsMonitorManager) EnumDisplayMonitors() ([]Monitor, error) {
	var monitors []Monitor

	deviceIndex := uint32(0)
	for {
		var displayDevice DISPLAY_DEVICE
		displayDevice.Cb = uint32(unsafe.Sizeof(displayDevice))

		ret, _, err := procEnumDisplayDevices.Call(
			uintptr(0),
			uintptr(deviceIndex),
			uintptr(unsafe.Pointer(&displayDevice)),
			uintptr(0),
		)

		if ret == 0 {
			break
		}

		if err != syscall.Errno(0) {
			return nil, fmt.Errorf("EnumDisplayDevices failed: %v", err)
		}

		deviceName := windows.UTF16ToString(displayDevice.DeviceName[:])
		deviceString := windows.UTF16ToString(displayDevice.DeviceString[:])

		var devMode DEVMODE
		devMode.DmSize = uint16(unsafe.Sizeof(devMode))

		ret, _, _ = procEnumDisplaySettings.Call(
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(deviceName))),
			uintptr(ENUM_CURRENT_SETTINGS),
			uintptr(unsafe.Pointer(&devMode)),
			uintptr(0),
		)

		if ret == 0 {
			ret, _, _ = procEnumDisplaySettings.Call(
				uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(deviceName))),
				uintptr(0),
				uintptr(unsafe.Pointer(&devMode)),
				uintptr(0),
			)
		}

		monitor := Monitor{
			DeviceName:  deviceName,
			DisplayName: deviceString,
			IsPrimary:   (displayDevice.StateFlags & DISPLAY_DEVICE_PRIMARY_DEVICE) != 0,
			IsActive:    (displayDevice.StateFlags & DISPLAY_DEVICE_ACTIVE) != 0,
		}

		if ret != 0 {
			monitor.Bounds = Rect{
				X:      int32(devMode.DmPosition.X),
				Y:      int32(devMode.DmPosition.Y),
				Width:  int32(devMode.DmPelsWidth),
				Height: int32(devMode.DmPelsHeight),
			}
		}

		monitors = append(monitors, monitor)
		deviceIndex++
	}

	fmt.Println("Monitor number", len(monitors))
	return monitors, nil
}

// SetPrimaryMonitor sets a monitor as the primary display on Windows
func (w *WindowsMonitorManager) SetPrimaryMonitor(deviceName string) error {
	deviceNamePtr, err := windows.UTF16PtrFromString(deviceName)
	if err != nil {
		return fmt.Errorf("failed to convert device name: %v", err)
	}

	var devMode DEVMODE
	devMode.DmSize = uint16(unsafe.Sizeof(devMode))
	devMode.DmFields = 0x00000020 // DM_POSITION

	ret, _, _ := procEnumDisplaySettings.Call(
		uintptr(unsafe.Pointer(deviceNamePtr)),
		uintptr(ENUM_CURRENT_SETTINGS),
		uintptr(unsafe.Pointer(&devMode)),
		uintptr(0),
	)

	if ret == 0 {
		return fmt.Errorf("failed to get current display settings")
	}

	ret, _, _ = procChangeDisplaySettings.Call(
		uintptr(unsafe.Pointer(deviceNamePtr)),
		uintptr(unsafe.Pointer(&devMode)),
		uintptr(CDS_UPDATEREGISTRY|0x00000001), // CDS_SET_PRIMARY
		uintptr(0),
		uintptr(0),
	)

	switch ret {
	case DISP_CHANGE_SUCCESSFUL:
		return nil
	case DISP_CHANGE_RESTART:
		return fmt.Errorf("computer must be restarted for changes to take effect")
	case DISP_CHANGE_FAILED:
		return fmt.Errorf("failed to change display settings")
	case DISP_CHANGE_BADMODE:
		return fmt.Errorf("display mode not supported")
	case DISP_CHANGE_NOTUPDATED:
		return fmt.Errorf("unable to write settings to registry")
	case DISP_CHANGE_BADFLAGS:
		return fmt.Errorf("invalid flags specified")
	case DISP_CHANGE_BADPARAM:
		return fmt.Errorf("invalid parameter specified")
	default:
		return fmt.Errorf("unknown error: %d", ret)
	}
}

// SetMonitorState activates or deactivates a monitor on Windows
func (w *WindowsMonitorManager) SetMonitorState(deviceName string, active bool) error {
	deviceNamePtr, err := windows.UTF16PtrFromString(deviceName)
	if err != nil {
		return fmt.Errorf("failed to convert device name: %v", err)
	}

	var devMode DEVMODE
	devMode.DmSize = uint16(unsafe.Sizeof(devMode))

	if active {
		ret, _, _ := procEnumDisplaySettings.Call(
			uintptr(unsafe.Pointer(deviceNamePtr)),
			uintptr(ENUM_CURRENT_SETTINGS),
			uintptr(unsafe.Pointer(&devMode)),
			uintptr(0),
		)

		if ret == 0 {
			ret, _, _ = procEnumDisplaySettings.Call(
				uintptr(unsafe.Pointer(deviceNamePtr)),
				uintptr(0),
				uintptr(unsafe.Pointer(&devMode)),
				uintptr(0),
			)
		}

		if ret == 0 {
			return fmt.Errorf("failed to get display settings for activation")
		}

		devMode.DmFields = 0x00000020 // DM_POSITION
	} else {
		devMode.DmPosition.X = -32000
		devMode.DmPosition.Y = -32000
		devMode.DmFields = 0x00000020 // DM_POSITION
	}

	ret, _, _ := procChangeDisplaySettings.Call(
		uintptr(unsafe.Pointer(deviceNamePtr)),
		uintptr(unsafe.Pointer(&devMode)),
		uintptr(CDS_UPDATEREGISTRY),
		uintptr(0),
		uintptr(0),
	)

	switch ret {
	case DISP_CHANGE_SUCCESSFUL:
		return nil
	case DISP_CHANGE_RESTART:
		return fmt.Errorf("computer must be restarted for changes to take effect")
	case DISP_CHANGE_FAILED:
		return fmt.Errorf("failed to change display settings")
	case DISP_CHANGE_BADMODE:
		return fmt.Errorf("display mode not supported")
	case DISP_CHANGE_NOTUPDATED:
		return fmt.Errorf("unable to write settings to registry")
	case DISP_CHANGE_BADFLAGS:
		return fmt.Errorf("invalid flags specified")
	case DISP_CHANGE_BADPARAM:
		return fmt.Errorf("invalid parameter specified")
	default:
		return fmt.Errorf("unknown error: %d", ret)
	}
}

// ApplyProfile applies a monitor profile on Windows
func (w *WindowsMonitorManager) ApplyProfile(profile Profile) error {
	for _, monitor := range profile.Monitors {
		if monitor.IsActive {
			if err := w.SetMonitorState(monitor.DeviceName, true); err != nil {
				return fmt.Errorf("failed to activate monitor %s: %v", monitor.DeviceName, err)
			}

			if monitor.IsPrimary {
				if err := w.SetPrimaryMonitor(monitor.DeviceName); err != nil {
					return fmt.Errorf("failed to set primary monitor %s: %v", monitor.DeviceName, err)
				}
			}
		} else {
			if err := w.SetMonitorState(monitor.DeviceName, false); err != nil {
				return fmt.Errorf("failed to deactivate monitor %s: %v", monitor.DeviceName, err)
			}
		}
	}

	return nil
}
