package fileutils

import "os/exec"

func CheckFFmpeg() (err error) {
	return exec.Command("ffmpeg", "-version").Run()
}
