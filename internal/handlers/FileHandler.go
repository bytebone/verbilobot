package handlers

import (
	"bytebone/verbilobot/internal/admin"
	"bytebone/verbilobot/internal/fileutils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func FileMatcher(update *models.Update) bool {
	return update.Message.Video != nil || update.Message.VideoNote != nil || update.Message.Audio != nil || update.Message.Voice != nil || update.Message.Document != nil
}

func FileHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	getFile := func() (*models.File, error) {
		switch {
		case update.Message.Video != nil:
			return b.GetFile(ctx, &bot.GetFileParams{
				FileID: update.Message.Video.FileID,
			})
		case update.Message.VideoNote != nil:
			return b.GetFile(ctx, &bot.GetFileParams{
				FileID: update.Message.VideoNote.FileID,
			})
		case update.Message.Audio != nil:
			return b.GetFile(ctx, &bot.GetFileParams{
				FileID: update.Message.Audio.FileID,
			})
		case update.Message.Voice != nil:
			return b.GetFile(ctx, &bot.GetFileParams{
				FileID: update.Message.Voice.FileID,
			})
		case update.Message.Document != nil:
			if strings.HasPrefix(update.Message.Document.MimeType, "audio/") || strings.HasPrefix(update.Message.Document.MimeType, "video/") {
				return b.GetFile(ctx, &bot.GetFileParams{
					FileID: update.Message.Document.FileID,
				})
			} else {
				log.Print("Denied file of type: " + update.Message.Document.MimeType)
				return nil, fmt.Errorf("message does not contain an audio file")
			}
		default:
			return nil, fmt.Errorf("message does not contain an audio file")
		}
	}
	f, err := getFile()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "I can only process audio files. Please try again with an audio file.",
		})
		return
	} else {
		log.Printf("Got file: %s", f.FileUniqueID)
	}

	rawFile, err := fileutils.Download(b, f)
	if err != nil {
		log.Print(err)
		admin.Alert(ctx, b, fmt.Sprintf("Download error: %v", err))
		return
	} else {
		log.Printf("Downloaded file to: %s", rawFile.Name())
	}

	transcodedFile, err := fileutils.Transcode(rawFile)
	if err != nil {
		log.Print(err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "I was unable to transcribe this file. Are you sure it contains audio?",
		})
		admin.Alert(ctx, b, fmt.Sprintf("Transcoding error: %v\nFilename: %s", err, rawFile.Name()))
		if err := fileutils.Delete(rawFile.Name()); err != nil {
			log.Print(err)
			admin.Alert(ctx, b, fmt.Sprintf("Deletion error: %v", err))
		} else {
			log.Printf("Deleted %s", rawFile.Name())
		}
		return
	} else {
		log.Printf("Transcoded file to: %s", transcodedFile.Name())
	}

	text, err := fileutils.Transcribe(transcodedFile)
	if err != nil {
		log.Println("Error: ", err)
		admin.Alert(ctx, b, fmt.Sprintf("Transcription error: %v", err))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "I encountered an unknown issue during transcription. Please try again later.",
		})
		if err := fileutils.Delete(rawFile.Name(), transcodedFile.Name()); err != nil {
			log.Print(err)
			admin.Alert(ctx, b, fmt.Sprintf("Deletion error: %v", err))
		} else {
			log.Printf("Deleted %s and %s", rawFile.Name(), transcodedFile.Name())
		}
		return
	} else {
		log.Print("Transcribed text successfully")

		// Final text is sent here
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: Buttons,
		})
	}

	err = fileutils.Delete(rawFile.Name(), transcodedFile.Name())
	if err != nil {
		log.Print(err)
		admin.Alert(ctx, b, fmt.Sprintf("Deletion error: %v", err))
		return
	} else {
		log.Printf("Deleted %s and %s", rawFile.Name(), transcodedFile.Name())
	}
}
