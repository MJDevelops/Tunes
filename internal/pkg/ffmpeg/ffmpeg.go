package ffmpeg

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/mholt/archives"
	"github.com/mjdevelops/tunes/internal/pkg/util"
)

// ffmpeg executable wrapper
type Ffmpeg struct {
	Version  string
	Platform string
	Path     string
}

const (
	ArchiveTar int = iota
	ArchiveZip
)

const (
	ffmpegBuildsRepo string = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest"
)

var (
	ErrExtraction  error = errors.New("error during extraction")
	ErrUnsupported error = errors.New("unsupported platform")
)

var (
	filepathUnix    = filepath.Join("bin", "ffmpeg")
	filepathWindows = filepath.Join("bin", "ffmpeg.exe")
)

func NewFfmpeg() *Ffmpeg {
	return &Ffmpeg{
		Platform: util.GetPlatform(),
	}
}

func (f *Ffmpeg) DownloadLatest() error {
	// TODO: Check existing binary

	var (
		binData []byte
		path    string
		archive int
		err     error
	)

	switch f.Platform {
	case "darwin_amd64", "darwin_arm64":
		path = "https://evermeet.cx/ffmpeg/getrelease/zip"
		archive = ArchiveZip
		f.Path = filepathUnix
	case "windows_amd64":
		path, err = url.JoinPath(ffmpegBuildsRepo, "ffmpeg-master-latest-win64-gpl.zip")
		if err != nil {
			return err
		}
		archive = ArchiveZip
		f.Path = filepathWindows
	case "windows_arm64":
		path, err = url.JoinPath(ffmpegBuildsRepo, "ffmpeg-master-latest-winarm64-gpl.zip")
		if err != nil {
			return err
		}
		archive = ArchiveZip
	case "linux_amd64":
		path, err = url.JoinPath(ffmpegBuildsRepo, "ffmpeg-master-latest-linux64-gpl.tar.xz")
		if err != nil {
			return err
		}
		archive = ArchiveTar
	case "linux_arm64":
		path, err = url.JoinPath(ffmpegBuildsRepo, "ffmpeg-master-latest-linuxarm64-gpl.tar.xz")
		if err != nil {
			return err
		}
		archive = ArchiveTar
	default:
		return ErrUnsupported
	}

	binData, err = f.downloadFfmpeg(path, archive)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(f.Path, os.O_CREATE|os.O_WRONLY, 0750)
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

func (f *Ffmpeg) downloadFfmpeg(path string, archive int) ([]byte, error) {
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

func extractFfmpeg(binData []byte, archive int) ([]byte, error) {
	switch archive {
	case ArchiveZip:
		zipReader, err := zip.NewReader(bytes.NewReader(binData), int64(len(binData)))
		if err != nil || !(len(zipReader.File) > 0) {
			return nil, ErrExtraction
		}

		for _, f := range zipReader.File {
			if name := f.FileInfo().Name(); name == "ffmpeg" || name == "ffmpeg.exe" {
				file, err := f.Open()
				if err != nil {
					return nil, err
				}
				defer file.Close()

				binData, err := io.ReadAll(file)
				if err != nil {
					return nil, err
				}

				return binData, nil
			}
		}
	case ArchiveTar:
		var format archives.Tar

		err := format.Extract(context.Background(), bytes.NewReader(binData), func(ctx context.Context, info archives.FileInfo) error {
			if name := info.Name(); name == "ffmpeg" {
				f, err := info.Open()
				if err != nil {
					return err
				}
				defer f.Close()

				binData, err = io.ReadAll(f)
				if err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		return binData, nil
	}

	return nil, errors.New("invalid archive type")
}
