package commands

import (
	"context"

	"github.com/bytebone/verbilobot/internal/admin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func AlertTest(ctx context.Context, b *bot.Bot, update *models.Update) {
	admin.Alert(ctx, b, "This is a test alert")
}
