package main

import (
	fileutils "bytebone/verbilobot/fileutils"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Loading environment variables")
	err := godotenv.Load()
	if err != nil {
		if os.Getenv("VERBILO_TELEGRAM_TOKEN") == "" || os.Getenv("VERBILO_GROQ_TOKEN") == "" {
			log.Fatalf("Environment variables not set and no .env file found: %v", err)
		}
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	log.Println("Creating bot")
	b, err := bot.New(os.Getenv("VERBILO_TELEGRAM_TOKEN"), opts...)
	if err != nil {
		log.Panicf("Error creating bot: %v", err)
	}
	u, err := b.GetMyName(ctx, &bot.GetMyNameParams{})
	if err == nil {
		log.Println("Logged in as @" + u.Name)
	} else {
		log.Print(err)
	}

	log.Println("Registering handlers")
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandlerMatchFunc(fileMatcher, fileHandler)

	log.Println("Starting bot")
	b.Start(ctx)
	log.Println("Shutting down. Goodbye!")
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "This message does not contain any files that I can process.",
	})
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Hello, %s!", update.Message.From.FirstName),
	})
}

func fileMatcher(update *models.Update) bool {
	return update.Message.VideoNote != nil || update.Message.Audio != nil || update.Message.Voice != nil
}

func fileHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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
		default:
			return nil, fmt.Errorf("message does not contain any files")
		}
	}
	f, err := getFile()
	if err != nil {
		log.Print(err)
		return
	} else {
		log.Printf("Got file: %s", f.FileUniqueID)
	}

	path, err := fileutils.Download(b, f)
	if err != nil {
		log.Print(err)
		return
	} else {
		log.Printf("Downloaded file to: %s", path)
	}

	transcodedPath, err := fileutils.Transcode(path)
	if err != nil {
		log.Print(err)
		if err := fileutils.Delete(path); err != nil {
			log.Print(err)
		} else {
			log.Printf("Deleted %s", path)
		}
		return
	} else {
		log.Printf("Transcoded file to: %s", transcodedPath)
	}

	text, err := fileutils.Transcribe(transcodedPath)
	if err != nil {
		log.Print(err)
		if err := fileutils.Delete(path, transcodedPath); err != nil {
			log.Print(err)
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
		return
	} else {
		log.Printf("Deleted %s and %s", path, transcodedPath)
	}
}
