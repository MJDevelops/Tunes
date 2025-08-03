package ytdlp

import (
	"os/exec"
	"regexp"
)

func (r *YtDlp) GetThumbnails(url string) string {
	cmd := exec.Command(r.Bin, url, "--list-thumbnails", "-q")
	oBytes, _ := cmd.Output()
	return string(oBytes)
}

func (r *YtDlp) GetHighDefinitionThumbnail(url string) string {
	line := regexp.MustCompile(`.*1920\s+1080\s+(https?:\/\/\S+)`)
	out := r.GetThumbnails(url)
	return line.FindStringSubmatch(out)[1]
}
