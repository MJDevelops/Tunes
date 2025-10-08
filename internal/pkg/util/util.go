package util

import (
	"runtime"
	"strings"
)

func GetPlatform() string {
	return strings.Join([]string{runtime.GOOS, runtime.GOARCH}, "_")
}
