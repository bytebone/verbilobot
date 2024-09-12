package commands

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var CommandList = &bot.SetMyCommandsParams{
	Commands: []models.BotCommand{
		{
			Command:     "start",
			Description: "Start the bot",
		},
		{
			Command:     "privacy",
			Description: "Show privacy policy",
		},
	},
}
