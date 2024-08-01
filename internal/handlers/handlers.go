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

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "This message does not contain any files that I can process.",
	})
}

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeMarkdownV1,
		Text:      fmt.Sprintf("Hello, %s! Thanks for using Verbilobot.\n\nI will transcribe any audio you throw at me. *Anything that contains audio will work* - voice messages, videos, music and whatever else you can think of. And I handle a lot of languages automatically, no need to configure anything. Just send or forward me your files, and I'll do my best!\n\nYou can learn more about how I process your data by using /privacy.", update.Message.From.FirstName),
	})
}

func PrivacyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	disableEmbeds := true
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:             update.Message.Chat.ID,
		ParseMode:          models.ParseModeMarkdown,
		LinkPreviewOptions: &models.LinkPreviewOptions{IsDisabled: &disableEmbeds},
		Text:               "*VerbiloBot* is a privacy\\-first service\\.\n\nMessages that VerbiloBot cannot process are never downloaded, stored or forwarded\\.\n\nAny processable media is downloaded, transcribed by a third party provider, then immediately deleted from our servers\\.\n\nWe do not log or store any information about you\\. No usernames, no IPs, no filenames, nothing\\. In the same fashion, nothing beyond the file you send is forwarded to the transcription provider\\.\n\nThe transcription service is provided by *Groq*\\. They do not store your data and do not train any AI with it\\. If you want to learn more, refer to their [privacy policy](https://wow.groq.com/privacy-policy/)\\.",
	})
}

func IDProvider(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Your chat ID is: %d", update.Message.Chat.ID),
	})
}

func FileMatcher(update *models.Update) bool {
	return update.Message.VideoNote != nil || update.Message.Audio != nil || update.Message.Voice != nil || update.Message.Document != nil
}

func FileHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	getFile := func() (*models.File, error) {
		switch {
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

	path, err := fileutils.Download(b, f)
	if err != nil {
		log.Print(err)
		admin.Alert(ctx, b, fmt.Sprintf("Download error: %v", err))
		return
	} else {
		log.Printf("Downloaded file to: %s", path)
	}

	transcodedPath, err := fileutils.Transcode(path)
	if err != nil {
		log.Print(err)
		if err.Error() == "exit status 0xffffffea" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "This file does not contain any audio, so there's nothing to transcribe.",
			})
			if err := fileutils.Delete(path); err != nil {
				log.Print(err)
				admin.Alert(ctx, b, fmt.Sprintf("Deletion error: %v", err))
			} else {
				log.Printf("Deleted %s", path)
			}
			return
		} else {
			admin.Alert(ctx, b, fmt.Sprintf("Transcoding error: %v", err))
			if err := fileutils.Delete(path); err != nil {
				log.Print(err)
			} else {
				log.Printf("Deleted %s", path)
			}
			return
		}
	} else {
		log.Printf("Transcoded file to: %s", transcodedPath)
	}

	text, err := fileutils.Transcribe(transcodedPath)
	if err != nil {
		log.Print(err)
		admin.Alert(ctx, b, fmt.Sprintf("Transcription error: %v", err))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "I encountered an unknown issue during transcription. Please try again later.",
		})
		if err := fileutils.Delete(path, transcodedPath); err != nil {
			log.Print(err)
			admin.Alert(ctx, b, fmt.Sprintf("Deletion error: %v", err))
		} else {
			log.Printf("Deleted %s and %s", path, transcodedPath)
		}
		return
	} else {
		log.Print("Transcribed text successfully")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text,
		})
	}

	err = fileutils.Delete(path, transcodedPath)
	if err != nil {
		log.Print(err)
		admin.Alert(ctx, b, fmt.Sprintf("Deletion error: %v", err))
		return
	} else {
		log.Printf("Deleted %s and %s", path, transcodedPath)
	}
}
