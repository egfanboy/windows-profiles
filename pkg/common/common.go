package common

import "os"

// IsDevelopmentMode checks if the application is running in development mode
// by looking for go.mod file in the current directory
func IsDevelopmentMode() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}
