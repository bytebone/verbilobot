package llm

import (
	"context"

	"github.com/conneroisu/groq-go"
)

// Rewrites the given text as bullet points, reducing the texts to the shortest possible length without cutting down on relevant information.
func BulletPoints(ctx context.Context, c *groq.Client, text string) (shortText string, err error) {
	resp, err := c.CreateChatCompletion(ctx, groq.ChatCompletionRequest{
		Model: groq.Llama318BInstant,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleSystem,
				Content: "You will be provided a voice transcript. Your task is to reduce the text to bullet points, reducing the texts to the shortest possible length without cutting down on relevant information.",
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
