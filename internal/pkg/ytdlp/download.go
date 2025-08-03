package ytdlp

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

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

func GetLatestRelease() (*YtDlp, error) {
	ytdlp := &YtDlp{}
	release := fetchLatestRelease()

	if release == viper.GetString(viperYtDlpRelease) {
		ytdlp.Bin = viper.GetString(viperYtDlpPath)
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

	viper.Set(viperYtDlpRelease, release)
	viper.Set(viperYtDlpPath, ytdlp.Bin)
	viper.WriteConfig()

	return ytdlp, nil
}
