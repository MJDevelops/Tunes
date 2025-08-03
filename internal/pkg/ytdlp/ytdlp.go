package ytdlp

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type YtDlp struct {
	// Path for the binary executable
	Bin string
}

const baseUrl string = "https://github.com/yt-dlp/yt-dlp/releases"
const viperYtDlpRelease string = "executables.ytdlp.release"
const viperYtDlpPath string = "executables.ytdlp.path"

var binPath string
var platform string
var platformExecutables = map[string]string{
	"darwin_amd64":  "yt-dlp_macos",
	"windows_amd64": "yt-dlp.exe",
	"linux_amd64":   "yt-dlp_linux",
	"darwin_arm64":  "yt-dlp_macos",
}

func init() {
	wd, _ := os.Getwd()
	binPath = filepath.Join(wd, "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0750)
	}
	platform = strings.Join([]string{runtime.GOOS, runtime.GOARCH}, "_")
}
