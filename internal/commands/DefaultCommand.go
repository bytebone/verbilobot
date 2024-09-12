package commands

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Chat.Type == "private" {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "This message does not contain any files that I can process.",
		})
		if err != nil {
			log.Println(err)
			return
		}
	}
}
