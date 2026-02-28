package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"monitor-profile-manager-wails/pkg/common"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
//go:embed tools/svcl/svcl.exe
//go:embed tools/multimonitortool/MultiMonitorTool.exe
var assets embed.FS

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
	// Set up panic recovery to prevent crashes
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Application panic recovered: %v", r)
			log.Printf("Stack trace: %s", debug.Stack())
		}
	}()

	// Create an instance of the app structure
	app := NewApp()

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
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
