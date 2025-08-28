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
	Bin string
	ctx context.Context
}

const (
	baseUrl           string = "https://github.com/yt-dlp/yt-dlp/releases"
	viperYtDlpRelease string = "executables.ytdlp.release"
	viperYtDlpPath    string = "executables.ytdlp.path"
)

var (
	binPath             string
	platform            string
	platformExecutables = map[string]string{
		"darwin_amd64":  "yt-dlp_macos",
		"windows_amd64": "yt-dlp.exe",
		"linux_amd64":   "yt-dlp_linux",
		"darwin_arm64":  "yt-dlp_macos",
	}
)

func init() {
	wd, _ := os.Getwd()
	binPath = filepath.Join(wd, "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0750)
	}
	platform = strings.Join([]string{runtime.GOOS, runtime.GOARCH}, "_")
}

func (y *YtDlp) SetContext(ctx context.Context) {
	y.ctx = ctx
}

func GetLatestRelease() (*YtDlp, error) {
	ytdlp := &YtDlp{}
	release := fetchLatestRelease()

	if release == config.GetString(viperYtDlpRelease) {
		ytdlp.Bin = config.GetString(viperYtDlpPath)
		return ytdlp, nil
	}

	executable := platformExecutables[platform]
	location, _ := url.JoinPath(baseUrl, "download", release)

	if executable == "" {
		return nil, errors.New("unsupported")
	}

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

	config.Set(viperYtDlpRelease, release)
	config.Set(viperYtDlpPath, ytdlp.Bin)
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
