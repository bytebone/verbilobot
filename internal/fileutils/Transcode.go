package fileutils

import (
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Transcode(inFile *os.File) (outFile *os.File, err error) {
	inPath := inFile.Name()
	outPath := inPath + ".wav"

	err = ffmpeg.Input(inPath).Output(outPath, ffmpeg.KwArgs{"ar": "16000", "ac": "1", "map": "0:a"}).OverWriteOutput().Silent(true).Run()

	outFile, err = os.Open(outPath)
	return
}
