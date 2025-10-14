package util

import (
	"path/filepath"
	"runtime"
	"strings"
)

func GetPlatform() string {
	return strings.Join([]string{runtime.GOOS, runtime.GOARCH}, "_")
}

func GetFileExtension(file string) string {
	return strings.ToLower(filepath.Ext(file))
}
