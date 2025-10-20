package ffmpeg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
	"github.com/mjdevelops/tunes/internal/pkg/util"
)

// ffmpeg executable wrapper
type Ffmpeg struct {
	binPath string
}

type ArchiveType int

const (
	ArchiveTar ArchiveType = iota
	ArchiveZip
)

const (
	ffmpegBuildsRepo = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest"
	evermeetFfmpeg   = "https://evermeet.cx/ffmpeg"
)

var (
	ErrExtraction   = errors.New("error during extraction")
	ErrUnsupported  = errors.New("unsupported platform")
	ErrFetchRelease = errors.New("error fetching release version")
)

func NewFfmpeg(path string) (*Ffmpeg, error) {
	f := &Ffmpeg{}
	executable := getPlatformExecutable()
	if executable == "" {
		return nil, ErrUnsupported
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(path, 0750)
	}

	f.binPath = filepath.Join(path, executable)

	return f, nil
}

func (f *Ffmpeg) GetLatest() error {
	if f.isLatest() {
		return nil
	}

	version := getLatestReleaseVersion()
	if version == "" {
		return ErrFetchRelease
	}

	var (
		binData []byte
		path    string
		archive ArchiveType
		err     error
	)

	switch util.GetPlatform() {
	case util.PlatformDarwinX64, util.PlatformDarwinArm64:
		path, err = url.JoinPath(evermeetFfmpeg, "getrelease", "zip")
		archive = ArchiveZip
	case util.PlatformWindowsX64:
		path, err = url.JoinPath(ffmpegBuildsRepo, fmt.Sprintf("ffmpeg-n%s-latest-win64-gpl-%s.zip", version, version))
		archive = ArchiveZip
	case util.PlatformWindowsArm64:
		path, err = url.JoinPath(ffmpegBuildsRepo, fmt.Sprintf("ffmpeg-n%s-latest-winarm64-gpl-%s.zip", version, version))
		archive = ArchiveZip
	case util.PlatformLinuxX64:
		path, err = url.JoinPath(ffmpegBuildsRepo, fmt.Sprintf("ffmpeg-n%s-latest-linux64-gpl-%s.tar.xz", version, version))
		archive = ArchiveTar
	case util.PlatformLinuxArm64:
		path, err = url.JoinPath(ffmpegBuildsRepo, fmt.Sprintf("ffmpeg-n%s-latest-linuxarm64-gpl-%s.tar.xz", version, version))
		archive = ArchiveTar
	default:
		return ErrUnsupported
	}

	if err != nil {
		return err
	}

	binData, err = f.downloadFfmpeg(path, archive)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(f.binPath, os.O_CREATE|os.O_WRONLY, 0750)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(binData)
	if err != nil {
		return err
	}

	return nil
}

func (f *Ffmpeg) Version() string {
	v, err := exec.Command(f.binPath, "-version").Output()
	if err != nil {
		return ""
	}

	return strings.Split(string(v), " ")[2]
}

func (f *Ffmpeg) Path() string {
	return f.binPath
}

func (f *Ffmpeg) downloadFfmpeg(path string, archive ArchiveType) ([]byte, error) {
	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	extractedBin, err := extractFfmpeg(b, archive)
	if err != nil {
		return nil, err
	}

	return extractedBin, nil
}

func (f *Ffmpeg) isLatest() bool {
	version := f.Version()
	if version == "" {
		return false
	}

	versionNumber := strings.Split(version, "-")[0]
	latestVersion := getLatestReleaseVersion()

	if platform := util.GetPlatform(); strings.Contains(platform, "darwin") {
		return versionNumber == latestVersion
	} else {
		return versionNumber[1:4] == latestVersion
	}
}

func extractFfmpeg(binData []byte, archive ArchiveType) ([]byte, error) {
	var (
		err error
		bin []byte
	)

	ctx, cancel := context.WithCancel(context.Background())
	extractor := func(ctx context.Context, info archives.FileInfo) error {
		if name := info.Name(); name == "ffmpeg" || name == "ffmpeg.exe" {
			file, err := info.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			bin, err = io.ReadAll(file)
			if err != nil {
				return err
			}

			cancel()
		}
		return nil
	}

	switch archive {
	case ArchiveZip:
		var format archives.Zip
		err = format.Extract(ctx, bytes.NewReader(binData), extractor)
	case ArchiveTar:
		var (
			compression archives.Xz
			format      archives.Tar
		)

		decompressedReader, err := compression.OpenReader(bytes.NewReader(binData))
		if err != nil {
			return nil, ErrExtraction
		}
		defer decompressedReader.Close()

		err = format.Extract(ctx, decompressedReader, extractor)
	default:
		return nil, errors.New("invalid archive type")
	}

	if !errors.Is(err, context.Canceled) && err != nil {
		return nil, err
	}

	return bin, nil
}

func getPlatformExecutable() string {
	switch util.GetOSType() {
	case util.OSUnix:
		return "ffmpeg"
	case util.OSWindows:
		return "ffmpeg.exe"
	default:
		return ""
	}
}

func getLatestReleaseVersion() string {
	evermeetInfoUrl, _ := url.JoinPath(evermeetFfmpeg, "info", "ffmpeg", "release")
	resp, err := http.Get(evermeetInfoUrl)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var release struct {
		Version string `json:"version"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return ""
	}

	return release.Version
}
