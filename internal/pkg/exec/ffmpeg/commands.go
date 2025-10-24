package ffmpeg

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

var ErrNotAFile = errors.New("not a file")

// Transcode converts a file to the format of outFile.
//
// If inFile could not be found or inFile is not a file, an err of type fs.PathError
// or ErrNotAFile is returned respectively.
func (f *Ffmpeg) Transcode(inFile string, outFile string) (err error) {
	file, err := os.Stat(inFile)
	if err != nil {
		return err
	}

	if file.IsDir() {
		return ErrNotAFile
	}

	return exec.Command(f.binPath, "-i", inFile, outFile).Run()
}

// Version returns the version of Ffmpeg at the specified path.
//
// A version == "" is returned when the version could not be determined.
func (f *Ffmpeg) Version() (version string) {
	v, err := exec.Command(f.binPath, "-version").Output()
	if err != nil {
		return ""
	}

	return strings.Split(string(v), " ")[2]
}
