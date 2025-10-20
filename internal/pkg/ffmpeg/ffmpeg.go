package ffmpeg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

const ffmpegBuildsRepo string = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest"

var (
	ErrExtraction  error = errors.New("error during extraction")
	ErrUnsupported error = errors.New("unsupported platform")
)

func NewFfmpeg(path string) (*Ffmpeg, error) {
	f := &Ffmpeg{}
	executable := getPlatformExecutable()
	if executable == "" {
		return nil, ErrUnsupported
	}

	f.binPath = filepath.Join(path, executable)

	return f, nil
}

func (f *Ffmpeg) GetLatest() error {
	if f.isLatest() {
		return nil
	}

	var (
		binData []byte
		path    string
		archive ArchiveType
		err     error
	)

	switch util.GetPlatform() {
	case util.PlatformDarwinX64, util.PlatformDarwinArm64:
		path = "https://evermeet.cx/ffmpeg/getrelease/zip"
		archive = ArchiveZip
	case util.PlatformWindowsX64:
		path, err = url.JoinPath(ffmpegBuildsRepo, "ffmpeg-master-latest-win64-gpl.zip")
		archive = ArchiveZip
	case util.PlatformWindowsArm64:
		path, err = url.JoinPath(ffmpegBuildsRepo, "ffmpeg-master-latest-winarm64-gpl.zip")
		archive = ArchiveZip
	case util.PlatformLinuxX64:
		path, err = url.JoinPath(ffmpegBuildsRepo, "ffmpeg-master-latest-linux64-gpl.tar.xz")
		archive = ArchiveTar
	case util.PlatformLinuxArm64:
		path, err = url.JoinPath(ffmpegBuildsRepo, "ffmpeg-master-latest-linuxarm64-gpl.tar.xz")
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
	if platform := util.GetPlatform(); strings.Contains(platform, "darwin") {
		// TODO: Implement isLatest for macOS
		return false
	} else {
		releaseInfo := strings.Split(f.Version(), "-")
		releaseDate := getLatestReleaseDate()
		if len(releaseInfo) > 0 {
			date := releaseInfo[3]
			if parsed, _ := time.Parse("20060102", date); parsed.Equal(releaseDate) {
				return true
			}
		}
		return false
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

func getLatestReleaseDate() time.Time {
	resp, err := http.Get("https://api.github.com/repos/BtbN/FFmpeg-Builds/releases/latest")
	if err != nil {
		return time.Time{}
	}
	defer resp.Body.Close()

	var release struct {
		PublishedAt string `json:"published_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return time.Time{}
	}

	t, _ := time.Parse(time.RFC3339, release.PublishedAt)

	return t
}
