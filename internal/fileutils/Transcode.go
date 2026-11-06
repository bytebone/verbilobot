package fileutils

import (
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Transcode(inFile *os.File) (outFile *os.File, err error) {
	inPath := inFile.Name()
	outPath := inPath + ".ogg"

	err = ffmpeg.Input(inPath).Output(outPath, ffmpeg.KwArgs{"c:a": "libopus", "ac": "1", "map": "0:a:0", "filter:a": "atempo=1.5"}).OverWriteOutput().Silent(true).Run()

	outFile, err = os.Open(outPath)
	return
}
