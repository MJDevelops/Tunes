package ytdlp

import (
	"os/exec"
	"regexp"
)

func (r *YoutubeResource) GetThumbnails() string {
	cmd := exec.Command(ExecPath, r.URL, "--list-thumbnails", "-q")
	oBytes, _ := cmd.Output()
	return string(oBytes)
}

func (r *YoutubeResource) GetHighDefinitionThumbnail() string {
	line := regexp.MustCompile(`.*1920\s+1080\s+(https?:\/\/\S+)`)
	out := r.GetThumbnails()
	return line.FindStringSubmatch(out)[1]
}
