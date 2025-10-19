package util

import (
	"path/filepath"
	"runtime"
	"strings"
)

const (
	PlatformWindowsX64   = "windows_amd64"
	PlatformWindowsArm64 = "windows_arm64"
	PlatformLinuxX64     = "linux_amd64"
	PlatformLinuxArm64   = "linux_arm64"
	PlatformDarwinArm64  = "darwin_arm64"
	PlatformDarwinX64    = "darwin_amd64"
)

const (
	OSUnix    = "unix"
	OSWindows = "windows"
)

func GetPlatform() string {
	return strings.Join([]string{runtime.GOOS, runtime.GOARCH}, "_")
}

func GetOSType() string {
	switch runtime.GOOS {
	case "windows":
		return OSWindows
	case "linux", "darwin":
		return OSUnix
	default:
		return ""
	}
}

func GetFileExtension(file string) string {
	return strings.ToLower(filepath.Ext(file))
}
