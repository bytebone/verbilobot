package llm

import (
	"github.com/jpoz/groq"
)

// Rewrites the given text as bullet points, reducing the texts to the shortest possible length without cutting down on relevant information.
func BulletPoints(c *groq.Client, text string) (shortText string, err error) {
	resp, err := c.CreateChatCompletion(groq.CompletionCreateParams{
		Model: "llama-3.1-8b-instant",
		Messages: []groq.Message{
			{
				Role:    "system",
				Content: "You will be provided a voice transcript. Your task is to reduce the text to bullet points, reducing the texts to the shortest possible length without cutting down on relevant information.",
			},
			{
				Role:    "user",
				Content: text,
			},
		},
	})

	shortText = resp.Choices[0].Message.Content

	return
}
