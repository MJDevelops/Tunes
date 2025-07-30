package ytdlp

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetLatestRelease() error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirect")
		},
	}

	res, _ := client.Get(latestBaseUrl)
	releasePath, _ := res.Location()
	release := filepath.Base(releasePath.String())
	location, _ := url.JoinPath(baseUrl, "download", release)

	plat := strings.Join([]string{runtime.GOOS, runtime.GOARCH}, "_")
	executable := platformExecutables[plat]

	if executable == "" {
		return errors.New("unsupported")
	}

	ExecPath = filepath.Join(binPath, executable)

	out, _ := os.Create(ExecPath)
	defer out.Close()

	downloadPath, _ := url.JoinPath(location, executable)
	res, err := http.Get(downloadPath)
	if err != nil {
		return errors.New("request failed")
	}
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return errors.New("couldn't write response data to file")
	}

	os.Chmod(ExecPath, 0750)

	return nil
}
