package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"monitor-profile-manager-wails/pkg/common"

	"fyne.io/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
//go:embed tools/svcl/svcl.exe
//go:embed tools/multimonitortool/MultiMonitorTool.exe
//go:embed build/windows/icon.ico
var assets embed.FS

var profilesUpdatedCh = make(chan []string)
var profileSelectedCh = make(chan string)

// getIconBytes returns the icon bytes for the system tray
func getIconBytes() []byte {
	// Try to read embedded icon
	iconData, err := assets.ReadFile("build/windows/icon.ico")
	if err != nil {
		return nil
	}
	return iconData
}

// extractTools extracts embedded tools to a temporary directory
func extractTools() (string, error) {
	// Check if we're in development mode (wails dev)
	if common.IsDevelopmentMode() {
		return "", nil
	}

	tempDir, err := os.MkdirTemp("", "monitor-profile-tools")
	if err != nil {
		return "", err
	}

	// List of embedded tools to extract
	tools := []struct {
		srcPath  string
		destPath string
	}{
		{"tools/svcl/svcl.exe", "svcl/svcl.exe"},
		{"tools/multimonitortool/MultiMonitorTool.exe", "multimonitortool/MultiMonitorTool.exe"},
	}

	for _, tool := range tools {
		// Read embedded file
		data, err := assets.ReadFile(tool.srcPath)
		if err != nil {
			return "", fmt.Errorf("failed to read embedded %s: %v", tool.srcPath, err)
		}

		// Create destination directory
		destDir := filepath.Join(tempDir, filepath.Dir(tool.destPath))
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %v", destDir, err)
		}

		// Write file with executable permissions
		destFile := filepath.Join(tempDir, tool.destPath)

		// Remove existing file to handle multiple runs
		os.Remove(destFile) // Ignore error if file doesn't exist

		if err := os.WriteFile(destFile, data, 0755); err != nil {
			return "", fmt.Errorf("failed to write %s: %v", destFile, err)
		}
	}

	return tempDir, nil
}

func main() {
	// Parse command line flags
	startMinimized := flag.Bool("minimized", false, "Start the application minimized")
	flag.Parse()

	// Set up panic recovery to prevent crashes
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Application panic recovered: %v", r)
			log.Printf("Stack trace: %s", debug.Stack())
		}
	}()

	// Create an instance of the app structure
	app := NewApp()

	// Set up system tray
	systrayStart, systrayEnd := systray.RunWithExternalLoop(func() {
		// System tray setup
		systray.SetIcon(getIconBytes())
		systray.SetTitle("Windows Profile Manager")
		systray.SetTooltip("Windows Profile Manager")

		// Add menu items
		mApplyProfile := systray.AddMenuItem("Apply Profile", "Apply a profile")
		mShow := systray.AddMenuItem("Show", "Show main window")
		mQuit := systray.AddMenuItem("Quit", "Quit application")

		var profileSubItems = []*systray.MenuItem{}

		// handle system tray menu events
		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					runtime.WindowShow(app.ctx)
				case <-mQuit.ClickedCh:
					systray.Quit()
					if app.ctx != nil {
						runtime.Quit(app.ctx)
					}
					return

				case profile := <-profileSelectedCh:
					app.ApplyProfile(profile)
				}
			}
		}()

		// handle profile list updates
		go func() {
			for profiles := range profilesUpdatedCh {

				// Hide any existing item since we are rebuilding the list
				for _, existingSubMenuItem := range profileSubItems {
					existingSubMenuItem.Hide()
				}

				// All old items are removed so reset the sub menu item array
				profileSubItems = []*systray.MenuItem{}

				for _, profile := range profiles {
					p := profile
					subMenuItem := mApplyProfile.AddSubMenuItem(p, "Apply Profile")

					profileSubItems = append(profileSubItems, subMenuItem)
					// Handler to emit event that apply was called
					go func(profileSubMi *systray.MenuItem, profileName string) {
						for range profileSubMi.ClickedCh {
							profileSelectedCh <- profileName
						}
					}(subMenuItem, p)
				}
			}
		}()
	}, func() {
		fmt.Println("System tray exiting")
	})

	// Clean up system tray
	defer systrayEnd()

	// Start system tray
	systrayStart()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Windows Profile Manager",
		Width:  1400,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			runtime.WindowHide(ctx)
			return true
		},
		StartHidden: *startMinimized,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
