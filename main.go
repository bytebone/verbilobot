package main

import (
	"bytebone/verbilobot/internal/fileutils"
	"bytebone/verbilobot/internal/handlers"
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
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

	if fileutils.CheckFFmpeg() != nil {
		log.Fatal("Couldn't run ffmpeg. Make sure that it is installed and accessible from your PATH. Or use the docker container.")
	}
	log.Println("FFmpeg is present and working")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
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
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/privacy", bot.MatchTypeExact, handlers.PrivacyHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/chatid", bot.MatchTypeExact, handlers.IDProvider)
	b.RegisterHandlerMatchFunc(handlers.FileMatcher, handlers.FileHandler)

	log.Println("Starting bot")
	b.Start(ctx)
	log.Println("Shutting down. Goodbye!")
}
