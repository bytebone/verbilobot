package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bytebone/verbilobot/internal/commands"
	"github.com/bytebone/verbilobot/internal/fileutils"
	"github.com/bytebone/verbilobot/internal/handlers"

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

	if fileutils.CheckFFmpeg() != nil {
		log.Fatal("Couldn't run ffmpeg. Make sure that it is installed and accessible from your PATH. Or use the docker container.")
	}
	log.Println("FFmpeg is present and working")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(commands.Default),
		bot.WithCallbackQueryDataHandler("button_", bot.MatchTypePrefix, handlers.ButtonCallbackHandler),
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
		log.Println(err)
	}

	log.Println("Registering commands")
	commandsForAPI := []models.BotCommand{}
	for _, cmd := range commands.CommandList {
		b.RegisterHandler(cmd.HandlerType, fmt.Sprintf("/%s", cmd.Command), cmd.MatchType, cmd.HandlerFunc)
		commandsForAPI = append(commandsForAPI, models.BotCommand{
			Command:     cmd.Command,
			Description: cmd.Description,
		})
	}
	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: commandsForAPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	b.RegisterHandlerMatchFunc(handlers.FileMatcher, handlers.FileHandler)
	registeredCommands, err := b.GetMyCommands(ctx, &bot.GetMyCommandsParams{})
	if err != nil {
		log.Println(err)
	}
	var registeredCommandNames []string
	for _, cmd := range registeredCommands {
		registeredCommandNames = append(registeredCommandNames, cmd.Command)
	}
	log.Printf("Registered commands: %s", strings.Join(registeredCommandNames, ", "))

	log.Println("Starting bot")
	b.Start(ctx)
	log.Println("Shutting down. Goodbye!")
}
