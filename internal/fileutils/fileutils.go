package fileutils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func CheckFFmpeg() (err error) {
	return exec.Command("ffmpeg", "-version").Run()
}

func Download(b *bot.Bot, f *models.File) (path string, err error) {
	link := b.FileDownloadLink(f)
	resp, err := http.Get(link)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dir := filepath.Dir(filepath.Join("files", f.FilePath))
	path = filepath.Join(dir, f.FileUniqueID+filepath.Ext(f.FilePath))
	os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	out, err := os.Create(path)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	return
}

func Transcode(path string) (out string, err error) {
	out = path + ".wav"
	err = ffmpeg.Input(path).Output(out, ffmpeg.KwArgs{"ar": "16000", "ac": "1", "map": "0:a"}).OverWriteOutput().Silent(true).Run()
	return
}

func Transcribe(path string) (text string, err error) {
	apiKey := os.Getenv("VERBILO_GROQ_TOKEN")

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("file", path)
	if err != nil {
		return
	}
	_, err = io.Copy(part, f)
	if err != nil {
		return
	}

	err = writer.WriteField("model", "whisper-large-v3")
	if err != nil {
		return
	}
	err = writer.WriteField("response_format", "text")
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/audio/transcriptions", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Body = io.NopCloser(buf)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	} else if resp.StatusCode == 429 {
		return "", fmt.Errorf(resp.Status)
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf(resp.Status)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return string(respBody), nil
}

func Delete(paths ...string) (err error) {
	for _, path := range paths {
		if err = os.Remove(path); err != nil {
			return
		}
	}
	return
}
