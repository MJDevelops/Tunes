package ytdlp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mjdevelops/tunes/internal/pkg/config"
)

// Wrapper for yt-dlp executable
type YtDlp struct {
	Path string
}

type Thumbnail struct {
	Url        string      `json:"url"`
	Height     json.Number `json:"height"`
	Width      json.Number `json:"width"`
	Resolution string      `json:"resolution"`
}

const baseUrl string = "https://github.com/yt-dlp/yt-dlp/releases"

func DownloadLatestRelease(c *config.ApplicationConfig) (*YtDlp, error) {
	var err error

	executable := getPlatformExecutable()
	if executable == "" {
		return nil, errors.New("unsupported")
	}

	ytdlp := &YtDlp{}
	release, err := fetchLatestRelease()
	if err != nil {
		return nil, errors.New("unable to fetch latest release")
	}
	execPath := c.Executables.YtDlp.Path
	_, err = os.Stat(execPath)

	if release == c.Executables.YtDlp.Release && !errors.Is(err, os.ErrNotExist) {
		ytdlp.Path = execPath
		return ytdlp, nil
	}

	location, _ := url.JoinPath(baseUrl, "download", release)

	wd, _ := os.Getwd()
	binPath := filepath.Join(wd, "bin")
	if _, err := os.Stat(binPath); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(binPath, 0750)
	}

	ytdlp.Path = filepath.Join(binPath, executable)

	out, _ := os.Create(ytdlp.Path)
	defer out.Close()

	downloadPath, _ := url.JoinPath(location, executable)
	res, err := http.Get(downloadPath)
	if err != nil {
		return nil, errors.New("request failed")
	}
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return nil, errors.New("couldn't write response data to file")
	}

	os.Chmod(ytdlp.Path, 0750)

	c.Executables.YtDlp.Release = release
	c.Executables.YtDlp.Path = ytdlp.Path
	c.Write()

	return ytdlp, nil
}

// Creates command with the given options and sets the quiet flag
func (y *YtDlp) CreateCommandQuiet(opts ...string) *exec.Cmd {
	opts = append(opts, "-q")
	return exec.Command(y.Path, opts...)
}

func fetchLatestRelease() (string, error) {
	var githubRes struct {
		TagName string `json:"tag_name"`
	}

	res, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&githubRes)

	return githubRes.TagName, nil
}

func getPlatformExecutable() string {
	platform := strings.Join([]string{runtime.GOOS, runtime.GOARCH}, "_")

	switch platform {
	case "darwin_amd64", "darwin_arm64":
		return "yt-dlp_macos"
	case "windows_amd64":
		return "yt-dlp.exe"
	case "windows_386":
		return "yt-dlp_x86.exe"
	case "windows_arm64":
		return "yt-dlp_arm64.exe"
	case "linux_amd64":
		return "yt-dlp_linux"
	case "linux_arm64":
		return "yt-dlp_linux_aarch64"
	default:
		return ""
	}
}
