package fileutils

import (
	"encoding/json"
	"strconv"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Duration(path string) (float64, error) {
	raw, err := ffmpeg.Probe(path)
	if err != nil {
		return 0, err
	}

	var parsed struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(parsed.Format.Duration, 64)
}
