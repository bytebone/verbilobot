package admin

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-telegram/bot"
)

func Alert(ctx context.Context, b *bot.Bot, content string) {
	adminChatIDString := os.Getenv("VERBILO_ADMIN_CHAT_ID")
	if adminChatIDString == "" {
		return
	}
	adminChatID, err := strconv.ParseInt(adminChatIDString, 10, 64)
	if err != nil {
		log.Printf("Admin chat ID (\"%s\") has to be a whole number, but is not", adminChatIDString)
		return
	}

	params := &bot.SendMessageParams{
		ChatID: adminChatID,
		Text:   fmt.Sprintf("⚠ ADMIN ALERT\n\n%s", content),
		// ParseMode: models.ParseModeMarkdownV1,
	}

	if adminThreadIDString := os.Getenv("VERBILO_ADMIN_THREAD_ID"); adminThreadIDString != "" {
		adminThreadID, err := strconv.ParseInt(adminThreadIDString, 10, 64)
		if err != nil {
			log.Printf("Admin thread ID (\"%s\") has to be a whole number, but is not", adminThreadIDString)
			return
		}
		params.MessageThreadID = int(adminThreadID)
	}

	if _, err = b.SendMessage(ctx, params); err != nil {
		log.Printf("Could not send message: %v", err)
	}
}
