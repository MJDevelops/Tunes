package ffmpeg

import (
	"errors"
	"os"
	"os/exec"
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
