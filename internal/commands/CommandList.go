package commands

import (
	"github.com/go-telegram/bot"
)

var CommandList = []Command{
	{
		Command:     "start",
		Description: "Start the Bot",
		HandlerType: bot.HandlerTypeMessageText,
		MatchType:   bot.MatchTypePrefix,
		HandlerFunc: Start,
	},
	{
		Command:     "privacy",
		Description: "Show Privacy Policy",
		HandlerType: bot.HandlerTypeMessageText,
		MatchType:   bot.MatchTypePrefix,
		HandlerFunc: Privacy,
	},

	// {
	// 	Command:     "alerttest",
	// 	Description: "Test the admin alert",
	// 	HandlerType: bot.HandlerTypeMessageText,
	// 	MatchType:   bot.MatchTypePrefix,
	// 	HandlerFunc: AlertTest,
	// },

	// {
	// 	Command:     "chatid",
	// 	Description: "Get current chat ID",
	// 	HandlerType: bot.HandlerTypeMessageText,
	// 	MatchType:   bot.MatchTypePrefix,
	// 	HandlerFunc: ID,
	// },
}

type Command struct {
	Command     string
	Description string
	HandlerType bot.HandlerType
	HandlerFunc bot.HandlerFunc
	MatchType   bot.MatchType
}
