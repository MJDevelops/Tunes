package ytdlp

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mjdevelops/tunes/internal/pkg/util"
)

// Wrapper for yt-dlp executable
type YtDlp struct {
	path    string
	Release string
}

type Thumbnail struct {
	Url        string      `json:"url"`
	Height     json.Number `json:"height"`
	Width      json.Number `json:"width"`
	Resolution string      `json:"resolution"`
}

const baseUrl string = "https://github.com/yt-dlp/yt-dlp/releases"

var ErrUnsupported error = errors.New("unsupported platform")

func DownloadLatest(binPath string) (*YtDlp, error) {
	var err error
	ytdlp := &YtDlp{}

	latestRelease, err := fetchLatestRelease()
	if err != nil {
		return nil, errors.New("unable to fetch latest release")
	}

	executable, err := downloadRelease(latestRelease, binPath)
	if err != nil {
		return nil, err
	}

	ytdlp.path = executable
	ytdlp.Release = latestRelease

	return ytdlp, nil
}

// Creates command with the given options and sets the quiet flag
func (y *YtDlp) CreateCommandQuiet(opts ...string) *exec.Cmd {
	opts = append(opts, "-q")
	return exec.Command(y.path, opts...)
}

func (y *YtDlp) Path() string {
	return y.path
}

// Downloads the provided yt-dlp release and returns the resulting output path.
func downloadRelease(release string, toPath string) (string, error) {
	executable := getPlatformExecutable()
	if executable == "" {
		return "", ErrUnsupported
	}

	if _, err := os.Stat(toPath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(toPath, 0750)
	}

	outPath := filepath.Join(toPath, executable)

	output, err := exec.Command(outPath, "--version").Output()
	if err != nil {
		log.Printf("error fetching yt-dlp version: %v\n", err)
	} else {
		outputRelease := strings.TrimSpace(string(output))
		if release == outputRelease {
			return outPath, nil
		}
	}

	out, _ := os.Create(outPath)
	defer out.Close()

	location, _ := url.JoinPath(baseUrl, "download", release, executable)
	res, err := http.Get(location)
	if err != nil {
		return "", errors.New("request failed")
	}
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", errors.New("couldn't write response data to file")
	}

	os.Chmod(outPath, 0750)

	return outPath, nil
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
	switch util.GetPlatform() {
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
