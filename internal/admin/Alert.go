package admin

import (
	"context"
	"fmt"
	"os"

	"github.com/go-telegram/bot"
)

func Alert(ctx context.Context, b *bot.Bot, content string) {
	adminChat := os.Getenv("VERBILO_ADMIN_CHAT")

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: adminChat,
		Text:   fmt.Sprintf("⚠ ADMIN ALERT\n\n%s", content),
		// ParseMode: models.ParseModeMarkdownV1,
	})

}
