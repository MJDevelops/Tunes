package ytdlp

import (
	"os/exec"
	"regexp"
)

func GetThumbnails(url string) string {
	cmd := exec.Command(ExecPath, url, "--list-thumbnails", "-q")
	oBytes, _ := cmd.Output()
	return string(oBytes)
}

func GetHighDefinitionThumbnail(url string) string {
	line := regexp.MustCompile(`.*1920\s+1080\s+(https?:\/\/\S+)`)
	out := GetThumbnails(url)
	return line.FindStringSubmatch(out)[1]
}
