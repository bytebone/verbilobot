package handlers

import (
	"context"
	"log"
	"os"

	"github.com/bytebone/verbilobot/internal/llm"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jpoz/groq"
)

var LLMButtons = &models.InlineKeyboardMarkup{
	InlineKeyboard: [][]models.InlineKeyboardButton{
		{
			{Text: "✍🏼 Shorten Text", CallbackData: "llm_shorten"},
			{Text: "🔘 Bullet Points", CallbackData: "llm_bulletpoints"},
		},
	},
}

func LLMCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	groqClient := groq.NewClient(groq.WithAPIKey(os.Getenv("VERBILO_GROQ_TOKEN")))
	inputText := update.CallbackQuery.Message.Message.Text

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	switch update.CallbackQuery.Data {

	case "llm_shorten":
		log.Println("Shortening transcript contents")

		shortText, err := llm.ShortenText(groqClient, inputText)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   shortText,
		})
		if err != nil {
			log.Println(err)
			return
		}

	case "llm_bulletpoints":
		log.Println("Converting message to bullet points")

		bulletText, err := llm.BulletPoints(groqClient, inputText)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			// ParseMode: models.ParseModeMarkdownV1,
			Text: bulletText,
		})
		if err != nil {
			log.Println(err)
			return
		}

	default:
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "I did not understand that",
		})
		if err != nil {
			log.Println(err)
			return
		}
	}
}
