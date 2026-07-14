package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytebone/verbilobot/internal/admin"
	"github.com/bytebone/verbilobot/internal/fileutils"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func fileSizeMB(f *os.File) float64 {
	fi, err := f.Stat()
	if err != nil {
		log.Printf("stat %s: %v", f.Name(), err)
		return 0
	}
	return float64(fi.Size()) / 1_000_000
}

const maxFileSize = 20 * 1000 * 1000

var errorFileTooBig = errors.New("file too big")
var errorFileNoMedia = errors.New("file has no media")
var errorFileUnsupported = errors.New("file is unsupported")

func FileMatcher(update *models.Update) bool {
	return update.Message.Video != nil || update.Message.VideoNote != nil || update.Message.Audio != nil || update.Message.Voice != nil || update.Message.Document != nil
}

func FileHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	getFile := func() (*models.File, error) {
		switch {
		case update.Message.Video != nil:
			if update.Message.Video.FileSize >= maxFileSize {
				return &models.File{
					FileSize: update.Message.Video.FileSize,
					FilePath: update.Message.Video.FileName,
				}, errorFileTooBig
			}
			return b.GetFile(ctx, &bot.GetFileParams{
				FileID: update.Message.Video.FileID,
			})
		case update.Message.VideoNote != nil:
			if update.Message.VideoNote.FileSize >= maxFileSize {
				return &models.File{
					FileSize: int64(update.Message.VideoNote.FileSize),
				}, errorFileTooBig
			}
			return b.GetFile(ctx, &bot.GetFileParams{
				FileID: update.Message.VideoNote.FileID,
			})
		case update.Message.Audio != nil:
			if update.Message.Audio.FileSize >= maxFileSize {
				return &models.File{
					FileSize: update.Message.Audio.FileSize,
					FilePath: update.Message.Audio.FileName,
				}, errorFileTooBig
			}
			return b.GetFile(ctx, &bot.GetFileParams{
				FileID: update.Message.Audio.FileID,
			})
		case update.Message.Voice != nil:
			if update.Message.Voice.FileSize >= maxFileSize {
				return &models.File{
					FileSize: update.Message.Voice.FileSize,
				}, errorFileTooBig
			}
			return b.GetFile(ctx, &bot.GetFileParams{
				FileID: update.Message.Voice.FileID,
			})
		case update.Message.Document != nil:
			if update.Message.Document.FileSize >= maxFileSize {
				return &models.File{
					FileSize: update.Message.Document.FileSize,
					FilePath: update.Message.Document.FileName,
				}, errorFileTooBig
			}
			if strings.HasPrefix(update.Message.Document.MimeType, "audio/") || strings.HasPrefix(update.Message.Document.MimeType, "video/") {
				return b.GetFile(ctx, &bot.GetFileParams{
					FileID: update.Message.Document.FileID,
				})
			} else {
				// og.Println("Denied file of type: " + update.Message.Document.MimeType)
				return nil, errorFileNoMedia
			}
		default:
			// log.Println("Denied unknown message type")
			return nil, errorFileUnsupported
		}
	}
	f, err := getFile()
	if errors.Is(err, errorFileTooBig) {
		log.Println("Denied file over 20 MB")
		admin.Alert(ctx, b, fmt.Sprintf("Telegram Error: File > 20 MB received\nOriginal: %.2f MB", float64(f.FileSize)/1_000_000))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Due to Telegram limitations, the maximum filesize I can process is 20 Megabytes.",
		})
		return
	} else if errors.Is(err, errorFileNoMedia) {
		log.Println("Denied file of type: " + update.Message.Document.MimeType)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "This doesn't seem to be an audio or video that I can process.",
		})
		return
	} else if errors.Is(err, errorFileUnsupported) {
		log.Println("Denied file of unsupported type")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "This message type is not supported.",
		})
		return
	} else {
		log.Printf("Got file: %s", f.FileUniqueID)
	}

	_, err = b.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: update.Message.Chat.ID,
		Action: models.ChatActionTyping,
	})
	if err != nil {
		log.Println(err)
	}

	rawFile, err := fileutils.Download(b, f)
	if err != nil {
		log.Println(err)
		admin.Alert(ctx, b, fmt.Sprintf("Download error: %v", err))
		admin.Teable(ctx, admin.TeableResult{
			Status:     "Failure",
			Err:        err.Error(),
			ErrorStage: "download",
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "I couldn't download your file from Telegram. Please try again later.",
		})
		return
	} else {
		log.Printf("Downloaded file to: %s", rawFile.Name())
	}

	originalMB := fileSizeMB(rawFile)
	fileType := strings.TrimPrefix(filepath.Ext(rawFile.Name()), ".")

	transcodedFile, err := fileutils.Transcode(rawFile)
	if err != nil {
		log.Println(err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "I found no audio in this file. Are you sure it contains any?",
		})
		admin.Alert(ctx, b, fmt.Sprintf("Transcoding error: %v\nOriginal: %.2f MB", err, originalMB))
		admin.Teable(ctx, admin.TeableResult{
			Status:       "Failure",
			OriginalSize: originalMB,
			Err:          err.Error(),
			ErrorStage:   "transcode",
			FileType:     fileType,
		})
		if err := fileutils.Delete(rawFile); err != nil {
			log.Println(err)
			admin.Alert(ctx, b, fmt.Sprintf("Deletion error: %v", err))
		} else {
			log.Printf("Deleted %s", rawFile.Name())
		}
		return
	} else {
		log.Printf("Transcoded file to: %s", transcodedFile.Name())
	}

	compressedMB := fileSizeMB(transcodedFile)
	dur, derr := fileutils.Duration(transcodedFile.Name())
	if derr != nil {
		log.Printf("probe: %v", derr)
	}

	// messagePlaceholder, _ := b.SendMessage(ctx, &bot.SendMessageParams{
	// 	ChatID: update.Message.Chat.ID,
	// 	Text:   "Transcription in progress, please wait...",
	// })

	text, err := fileutils.Transcribe(ctx, transcodedFile)
	if err != nil {
		log.Println("Error: ", err)
		admin.Alert(ctx, b, fmt.Sprintf("Transcription error: %v\nOriginal: %.2f MB\nCompressed: %.2f MB\nDuration: %.2fs", err, originalMB, compressedMB, dur))
		admin.Teable(ctx, admin.TeableResult{
			Status:         "Failure",
			OriginalSize:   originalMB,
			CompressedSize: compressedMB,
			Duration:       dur,
			Err:            err.Error(),
			ErrorStage:     "transcription",
			FileType:       fileType,
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "I encountered an unknown issue during transcription. Please try again later.",
		})
		if err := fileutils.Delete(rawFile, transcodedFile); err != nil {
			log.Println(err)
			admin.Alert(ctx, b, fmt.Sprintf("Deletion error: %v", err))
		} else {
			log.Printf("Deleted %s and %s", rawFile.Name(), transcodedFile.Name())
		}
		return
	} else {
		log.Println("Transcribed text successfully")

		// Final text is sent here
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: Buttons,
		})
		admin.Teable(ctx, admin.TeableResult{
			Status:         "Success",
			OriginalSize:   originalMB,
			CompressedSize: compressedMB,
			Duration:       dur,
			FileType:       fileType,
		})
		// b.EditMessageText(ctx, &bot.EditMessageTextParams{
		// 	ChatID:      update.Message.Chat.ID,
		// 	MessageID:   messagePlaceholder.ID,
		// 	Text:        text,
		// 	ReplyMarkup: Buttons,
		// })
	}

	err = fileutils.Delete(rawFile, transcodedFile)
	if err != nil {
		log.Println(err)
		admin.Alert(ctx, b, fmt.Sprintf("Deletion error: %v", err))
		return
	} else {
		log.Printf("Deleted %s and %s", rawFile.Name(), transcodedFile.Name())
	}
}
