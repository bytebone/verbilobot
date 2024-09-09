package commands

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Privacy(ctx context.Context, b *bot.Bot, update *models.Update) {
	disableEmbeds := true
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:             update.Message.Chat.ID,
		ParseMode:          models.ParseModeMarkdown,
		LinkPreviewOptions: &models.LinkPreviewOptions{IsDisabled: &disableEmbeds},
		Text:               "*VerbiloBot* is a privacy\\-first service\\.\n\nMessages that VerbiloBot cannot process are never downloaded, stored or forwarded\\.\n\nAny processable media is downloaded, transcribed by a third party provider, then immediately deleted from our servers\\.\n\nWe do not log or store any information about you\\. No usernames, no IPs, no filenames, nothing\\. In the same fashion, nothing beyond the file you send is forwarded to the transcription provider\\.\n\nThe transcription service is provided by *Groq*\\. They do not store your data and do not train any AI with it\\. If you want to learn more, refer to their [privacy policy](https://wow.groq.com/privacy-policy/)\\.",
	})
}
