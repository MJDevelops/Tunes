package ytdlp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mjdevelops/tunes/internal/pkg/config"
)

type YtDlp struct {
	// Path for the binary executable
	Bin           string
	DownloadQueue DownloadQueue
	ctx           context.Context
}

const baseUrl string = "https://github.com/yt-dlp/yt-dlp/releases"

var (
	binPath    string
	executable string
)

func init() {
	wd, _ := os.Getwd()
	binPath = filepath.Join(wd, "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0750)
	}

	platform := strings.Join([]string{runtime.GOOS, runtime.GOARCH}, "_")

	switch platform {
	case "darwin_amd64", "darwin_arm64":
		executable = "yt-dlp_macos"
	case "windows_amd64":
		executable = "yt-dlp.exe"
	case "windows_386":
		executable = "yt-dlp_x86.exe"
	case "windows_arm64":
		executable = "yt-dlp_arm64.exe"
	case "linux_amd64":
		executable = "yt-dlp_linux"
	}
}

func (y *YtDlp) SetContext(ctx context.Context) {
	y.ctx = ctx
}

func GetLatestRelease() (*YtDlp, error) {
	if executable == "" {
		return nil, errors.New("unsupported")
	}

	ytdlp := &YtDlp{}
	release := fetchLatestRelease()

	if release == config.GetYtDlpRelease() {
		ytdlp.Bin = config.GetYtDlpPath()
		return ytdlp, nil
	}

	location, _ := url.JoinPath(baseUrl, "download", release)

	ytdlp.Bin = filepath.Join(binPath, executable)

	out, _ := os.Create(ytdlp.Bin)
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

	os.Chmod(ytdlp.Bin, 0750)

	config.SetYtDlpRelease(release)
	config.SetYtDlpPath(ytdlp.Bin)
	config.Write()

	return ytdlp, nil
}

func fetchLatestRelease() string {
	latestBaseUrl, _ := url.JoinPath(baseUrl, "latest")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirect")
		},
	}

	res, _ := client.Get(latestBaseUrl)
	releasePath, _ := res.Location()
	return filepath.Base(releasePath.String())
}
