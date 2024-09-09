package fileutils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Download(b *bot.Bot, f *models.File) (outFile *os.File, err error) {
	link := b.FileDownloadLink(f)
	resp, err := http.Get(link)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dir := filepath.Dir(filepath.Join("files", f.FilePath))
	path := filepath.Join(dir, f.FileUniqueID+filepath.Ext(f.FilePath))
	os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	outFile, err = os.Create(path)
	if err != nil {
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return
	}

	return
}
