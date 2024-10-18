package llm

import (
	"context"

	"github.com/conneroisu/groq-go"
)

// Removes speech patterns and filler words, returning the shortest possible form of the text without removing or altering any of the meaningful content.
func ShortenText(ctx context.Context, c *groq.Client, text string) (shortText string, err error) {
	resp, err := c.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.Llama318BInstant,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleSystem,
				Content: "You will be provided a voice transcript. Your task is to remove any speech patterns and filler words that are exclusive to spoken text, and thereby reduce the text to its shortest possible form without removing or altering any of the meaningful contents whatsoever. Do not add formatting, markdown or additional text of your own around the shortened text. The text may be provided in any language, make sure that your shortened version matches the input language.",
			},
			{
				Role:    groq.ChatMessageRoleUser,
				Content: text,
			},
		},
	})

	shortText = resp.Choices[0].Message.Content

	return
}
