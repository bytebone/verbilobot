package commands

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeMarkdownV1,
		Text:      fmt.Sprintf("Hello, %s! Thanks for using Verbilobot.\n\nI will transcribe any audio you throw at me. *Anything that contains audio will work* - voice messages, videos, music and whatever else you can think of. And I handle a lot of languages automatically, no need to configure anything. Just send or forward me your files, and I'll do my best!\n\nYou can learn more about how I process your data by using /privacy.", update.Message.From.FirstName),
	})
}
